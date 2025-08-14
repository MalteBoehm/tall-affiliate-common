package redis

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/events"
)

func TestStreamProducer_PublishEvent(t *testing.T) {
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

	producer := NewStreamProducer(client, nil)

	t.Run("publishes event successfully", func(t *testing.T) {
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

		err := producer.PublishEvent(ctx, streamName, testEvent)
		assert.NoError(t, err)

		// Verify event was published
		messages, err := client.XRange(ctx, streamName, "-", "+").Result()
		require.NoError(t, err)
		assert.Len(t, messages, 1)

		// Verify event data
		eventData, ok := messages[0].Values["event"].(string)
		require.True(t, ok)

		var publishedEvent events.Event
		err = json.Unmarshal([]byte(eventData), &publishedEvent)
		require.NoError(t, err)
		assert.Equal(t, testEvent.ID, publishedEvent.ID)
		assert.Equal(t, testEvent.Type, publishedEvent.Type)
	})

	t.Run("handles nil event", func(t *testing.T) {
		err := producer.PublishEvent(ctx, streamName, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event cannot be nil")
	})

	t.Run("validates stream name", func(t *testing.T) {
		testEvent := &events.Event{
			ID:   "test-456",
			Type: events.EventTypeProductCreated,
		}

		err := producer.PublishEvent(ctx, "", testEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stream name cannot be empty")
	})
}

func TestStreamProducer_PublishBatch(t *testing.T) {
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

	producer := NewStreamProducer(client, nil)

	t.Run("publishes batch of events", func(t *testing.T) {
		events := []*events.Event{
			{
				ID:            "batch-1",
				Type:          events.EventTypeProductCreated,
				AggregateType: "product",
				AggregateID:   "prod-1",
				Timestamp:     time.Now(),
			},
			{
				ID:            "batch-2",
				Type:          events.EventTypeProductUpdated,
				AggregateType: "product",
				AggregateID:   "prod-2",
				Timestamp:     time.Now(),
			},
			{
				ID:            "batch-3",
				Type:          events.EventTypeProductDeleted,
				AggregateType: "product",
				AggregateID:   "prod-3",
				Timestamp:     time.Now(),
			},
		}

		messageIDs, err := producer.PublishBatch(ctx, streamName, events)
		assert.NoError(t, err)
		assert.Len(t, messageIDs, 3)

		// Verify all events were published
		messages, err := client.XRange(ctx, streamName, "-", "+").Result()
		require.NoError(t, err)
		assert.Len(t, messages, 3)
	})

	t.Run("handles empty batch", func(t *testing.T) {
		messageIDs, err := producer.PublishBatch(ctx, streamName, []*events.Event{})
		assert.NoError(t, err)
		assert.Empty(t, messageIDs)
	})

	t.Run("handles nil events in batch", func(t *testing.T) {
		events := []*events.Event{
			{ID: "valid-1", Type: events.EventTypeProductCreated},
			nil, // This should be skipped
			{ID: "valid-2", Type: events.EventTypeProductUpdated},
		}

		messageIDs, err := producer.PublishBatch(ctx, streamName, events)
		assert.NoError(t, err)
		assert.Len(t, messageIDs, 2) // Only valid events published
	})
}