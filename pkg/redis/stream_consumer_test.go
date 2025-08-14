package redis

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/MalteBoehm/tall-affiliate-common/pkg/events"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamConsumer_CreateConsumerGroup(t *testing.T) {
	// Skip if no Redis available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available")
	}

	// Clean up test data
	streamName := "test:stream:" + time.Now().Format("20060102150405")
	defer client.Del(ctx, streamName)

	consumer := NewStreamConsumer(client, streamName, "test-group", "test-consumer", nil)

	t.Run("creates new consumer group", func(t *testing.T) {
		err := consumer.CreateConsumerGroup(ctx)
		assert.NoError(t, err)

		// Verify group exists
		groups, err := client.XInfoGroups(ctx, streamName).Result()
		require.NoError(t, err)
		assert.Len(t, groups, 1)
		assert.Equal(t, "test-group", groups[0].Name)
	})

	t.Run("handles existing consumer group", func(t *testing.T) {
		// Try to create again
		err := consumer.CreateConsumerGroup(ctx)
		assert.NoError(t, err) // Should not error on existing group
	})
}

func TestStreamConsumer_ConsumeStream(t *testing.T) {
	// Skip if no Redis available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available")
	}

	// Clean up test data
	streamName := "test:stream:" + time.Now().Format("20060102150405")
	defer client.Del(ctx, streamName)

	consumer := NewStreamConsumer(client, streamName, "test-group", "test-consumer", nil)

	// Create consumer group
	err := consumer.CreateConsumerGroup(ctx)
	require.NoError(t, err)

	t.Run("consumes events from stream", func(t *testing.T) {
		// Publish test event
		testEvent := &events.Event{
			ID:            "test-123",
			Type:          events.EventTypeProductCreated,
			AggregateType: "product",
			AggregateID:   "prod-123",
			Payload: map[string]interface{}{
				"asin":  "B001234567",
				"title": "Test Product",
			},
			Timestamp: time.Now(),
		}

		eventData, err := json.Marshal(testEvent)
		require.NoError(t, err)

		// Add to stream
		_, err = client.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{
				"event": string(eventData),
			},
		}).Result()
		require.NoError(t, err)

		// Consume with timeout context
		consumeCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		consumed := false
		handler := func(ctx context.Context, event *events.Event, messageID string) error {
			assert.Equal(t, testEvent.ID, event.ID)
			assert.Equal(t, testEvent.Type, event.Type)
			consumed = true
			cancel() // Stop consuming after first message
			return nil
		}

		err = consumer.ConsumeStream(consumeCtx, streamName, "test-group", 10, 100*time.Millisecond, handler)
		assert.Equal(t, context.Canceled, err) // Expected due to cancel()
		assert.True(t, consumed, "Event should have been consumed")
	})

	t.Run("handles malformed events", func(t *testing.T) {
		// Add malformed event
		_, err = client.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{
				"event": "invalid-json",
			},
		}).Result()
		require.NoError(t, err)

		consumeCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		handler := func(ctx context.Context, event *events.Event, messageID string) error {
			t.Fatal("Handler should not be called for malformed events")
			return nil
		}

		err = consumer.ConsumeStream(consumeCtx, streamName, "test-group", 10, 100*time.Millisecond, handler)
		assert.Equal(t, context.DeadlineExceeded, err)
	})
}

func TestStreamConsumer_AcknowledgeMessage(t *testing.T) {
	// Skip if no Redis available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available")
	}

	// Clean up test data
	streamName := "test:stream:" + time.Now().Format("20060102150405")
	groupName := "test-group"
	defer client.Del(ctx, streamName)

	consumer := NewStreamConsumer(client, streamName, groupName, "test-consumer", nil)

	// Create consumer group
	err := consumer.CreateConsumerGroup(ctx)
	require.NoError(t, err)

	// Add test message
	messageID, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"data": "test",
		},
	}).Result()
	require.NoError(t, err)

	t.Run("acknowledges message successfully", func(t *testing.T) {
		err := consumer.AcknowledgeMessage(ctx, messageID)
		assert.NoError(t, err)

		// Verify message is acknowledged
		pending, err := client.XPending(ctx, streamName, groupName).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), pending.Count) // No pending messages
	})
}
