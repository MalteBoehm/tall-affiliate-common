package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/events"
	"github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/interfaces"
)

// StreamProducer implements Redis stream event publishing
type StreamProducer struct {
	client *redis.Client
	logger *slog.Logger
}

// NewStreamProducer creates a new Redis stream producer
func NewStreamProducer(client *redis.Client, logger *slog.Logger) *StreamProducer {
	if logger == nil {
		logger = slog.Default()
	}
	return &StreamProducer{
		client: client,
		logger: logger.With("component", "stream-producer"),
	}
}

// PublishEvent publishes an event to a Redis stream
func (p *StreamProducer) PublishEvent(ctx context.Context, streamName string, event *events.Event) error {
	if streamName == "" {
		return fmt.Errorf("stream name cannot be empty")
	}
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Serialize event to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to stream
	messageID, err := p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		ID:     "*", // Let Redis generate the ID
		Values: map[string]interface{}{
			"event": string(eventData),
		},
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to publish event to stream: %w", err)
	}

	p.logger.Debug("Event published to stream",
		"stream", streamName,
		"eventType", event.Type,
		"eventID", event.ID,
		"messageID", messageID)

	return nil
}

// PublishBatch publishes multiple events to a stream
func (p *StreamProducer) PublishBatch(ctx context.Context, streamName string, events []*events.Event) ([]string, error) {
	if streamName == "" {
		return nil, fmt.Errorf("stream name cannot be empty")
	}

	var messageIDs []string
	pipe := p.client.Pipeline()

	for _, event := range events {
		if event == nil {
			continue // Skip nil events
		}

		// Serialize event to JSON
		eventData, err := json.Marshal(event)
		if err != nil {
			p.logger.Error("Failed to marshal event in batch", "error", err, "eventID", event.ID)
			continue
		}

		// Add to pipeline
		pipe.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			ID:     "*",
			Values: map[string]interface{}{
				"event": string(eventData),
			},
		})
	}

	// Execute pipeline
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute batch publish: %w", err)
	}

	// Extract message IDs
	for _, cmd := range cmds {
		if xAddCmd, ok := cmd.(*redis.StringCmd); ok {
			messageID, err := xAddCmd.Result()
			if err == nil {
				messageIDs = append(messageIDs, messageID)
			}
		}
	}

	p.logger.Debug("Batch published to stream",
		"stream", streamName,
		"count", len(messageIDs))

	return messageIDs, nil
}

// StreamInfo returns information about a stream
func (p *StreamProducer) StreamInfo(ctx context.Context, streamName string) (*redis.XInfoStream, error) {
	return p.client.XInfoStream(ctx, streamName).Result()
}

// TrimStream trims a stream to a specific length
func (p *StreamProducer) TrimStream(ctx context.Context, streamName string, maxLen int64) error {
	return p.client.XTrimMaxLen(ctx, streamName, maxLen).Err()
}

// Ensure StreamProducer implements the interface
var _ interfaces.StreamProducer = (*StreamProducer)(nil)