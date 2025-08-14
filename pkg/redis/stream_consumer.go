package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/MalteBoehm/tall-affiliate-common/pkg/events"
	"github.com/MalteBoehm/tall-affiliate-common/pkg/interfaces"
	"github.com/redis/go-redis/v9"
)

// StreamConsumer implements Redis stream consumption with consumer groups
type StreamConsumer struct {
	client       *redis.Client
	streamName   string
	groupName    string
	consumerName string
	logger       *slog.Logger
}

// NewStreamConsumer creates a new Redis stream consumer
func NewStreamConsumer(client *redis.Client, streamName, groupName, consumerName string, logger *slog.Logger) *StreamConsumer {
	if logger == nil {
		logger = slog.Default()
	}
	return &StreamConsumer{
		client:       client,
		streamName:   streamName,
		groupName:    groupName,
		consumerName: consumerName,
		logger:       logger.With("component", "stream-consumer", "stream", streamName),
	}
}

// CreateConsumerGroup creates the consumer group if it doesn't exist
func (c *StreamConsumer) CreateConsumerGroup(ctx context.Context) error {
	err := c.client.XGroupCreateMkStream(ctx, c.streamName, c.groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	return nil
}

// ConsumeStream implements the StreamConsumer interface
func (c *StreamConsumer) ConsumeStream(
	ctx context.Context,
	streamName string,
	groupName string,
	batchSize int64,
	pollInterval time.Duration,
	handler func(context.Context, *events.Event, string) error,
) error {
	c.logger.Info("Starting to consume stream",
		"group", groupName,
		"consumer", c.consumerName,
		"batchSize", batchSize,
		"pollInterval", pollInterval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read messages from the stream
			messages, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: c.consumerName,
				Streams:  []string{streamName, ">"},
				Count:    batchSize,
				Block:    pollInterval,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					// No messages available, continue polling
					continue
				}
				c.logger.Error("Failed to read from stream", "error", err)
				time.Sleep(pollInterval)
				continue
			}

			// Process messages
			for _, stream := range messages {
				for _, msg := range stream.Messages {
					if err := c.processMessage(ctx, msg, handler); err != nil {
						c.logger.Error("Failed to process message",
							"messageID", msg.ID,
							"error", err)
						// Continue processing other messages
						continue
					}
				}
			}
		}
	}
}

// processMessage processes a single message from the stream
func (c *StreamConsumer) processMessage(
	ctx context.Context,
	msg redis.XMessage,
	handler func(context.Context, *events.Event, string) error,
) error {
	// Extract event data from message
	eventData, ok := msg.Values["event"].(string)
	if !ok {
		return fmt.Errorf("message does not contain event data")
	}

	// Parse the event
	var event events.Event
	if err := json.Unmarshal([]byte(eventData), &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Execute handler
	c.logger.Debug("Processing event",
		"eventType", event.Type,
		"eventID", event.ID,
		"aggregateID", event.AggregateID,
		"messageID", msg.ID)

	if err := handler(ctx, &event, msg.ID); err != nil {
		return fmt.Errorf("handler failed: %w", err)
	}

	// Acknowledge the message
	if err := c.AcknowledgeMessage(ctx, msg.ID); err != nil {
		return fmt.Errorf("failed to acknowledge message: %w", err)
	}

	return nil
}

// AcknowledgeMessage acknowledges a message in the consumer group
func (c *StreamConsumer) AcknowledgeMessage(ctx context.Context, messageID string) error {
	return c.client.XAck(ctx, c.streamName, c.groupName, messageID).Err()
}

// PendingMessages returns information about pending messages
func (c *StreamConsumer) PendingMessages(ctx context.Context) (*redis.XPending, error) {
	return c.client.XPending(ctx, c.streamName, c.groupName).Result()
}

// ClaimStaleMessages claims messages that have been idle for too long
func (c *StreamConsumer) ClaimStaleMessages(ctx context.Context, minIdleTime time.Duration, count int64) ([]redis.XMessage, error) {
	// Get pending messages
	pending, err := c.client.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: c.streamName,
		Group:  c.groupName,
		Start:  "-",
		End:    "+",
		Count:  count,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get pending messages: %w", err)
	}

	var claimed []redis.XMessage
	for _, msg := range pending {
		if msg.Idle >= minIdleTime {
			// Claim the message
			msgs, err := c.client.XClaim(ctx, &redis.XClaimArgs{
				Stream:   c.streamName,
				Group:    c.groupName,
				Consumer: c.consumerName,
				MinIdle:  minIdleTime,
				Messages: []string{msg.ID},
			}).Result()
			if err != nil {
				c.logger.Error("Failed to claim message", "messageID", msg.ID, "error", err)
				continue
			}
			claimed = append(claimed, msgs...)
		}
	}

	return claimed, nil
}

// ReadMessages reads messages from the stream without a handler callback
// This is a compatibility method for workers that prefer batch processing
func (c *StreamConsumer) ReadMessages(ctx context.Context, count int64, blockTimeout time.Duration) ([]StreamMessage, error) {
	messages, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    c.groupName,
		Consumer: c.consumerName,
		Streams:  []string{c.streamName, ">"},
		Count:    count,
		Block:    blockTimeout,
	}).Result()

	if err != nil {
		if err == redis.Nil {
			// No messages available
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read messages: %w", err)
	}

	// Convert to StreamMessage format
	var result []StreamMessage
	for _, stream := range messages {
		for _, msg := range stream.Messages {
			result = append(result, StreamMessage{
				ID:   msg.ID,
				Data: msg.Values,
			})
		}
	}

	return result, nil
}

// AckMessage is an alias for AcknowledgeMessage for compatibility
func (c *StreamConsumer) AckMessage(ctx context.Context, messageID string) error {
	return c.AcknowledgeMessage(ctx, messageID)
}

// Ensure StreamConsumer implements the interface
var _ interfaces.StreamConsumer = (*StreamConsumer)(nil)
