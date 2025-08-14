package adapters

import (
	"context"
	"testing"
	"time"

	"github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProducer is a mock implementation of StreamProducer
type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) PublishEvent(ctx context.Context, streamName string, event *events.Event) error {
	args := m.Called(ctx, streamName, event)
	return args.Error(0)
}

// MockConsumer is a mock implementation of StreamConsumer
type MockConsumer struct {
	mock.Mock
}

func (m *MockConsumer) ConsumeStream(
	ctx context.Context,
	streamName string,
	groupName string,
	batchSize int64,
	pollInterval time.Duration,
	handler func(context.Context, *events.Event, string) error,
) error {
	args := m.Called(ctx, streamName, groupName, batchSize, pollInterval, handler)
	return args.Error(0)
}

func TestServiceEventAdapter_PublishEvent(t *testing.T) {
	mockProducer := new(MockProducer)
	mockConsumer := new(MockConsumer)
	adapter := NewServiceEventAdapter(mockProducer, mockConsumer)

	ctx := context.Background()
	streamName := "test-stream"
	event := &events.Event{
		ID:            "test-id",
		Type:          events.EventTypeNewProductDetected,
		AggregateType: "product",
		AggregateID:   "product-123",
		Timestamp:     time.Now(),
		Payload:       map[string]string{"test": "data"},
	}

	mockProducer.On("PublishEvent", ctx, streamName, event).Return(nil)

	err := adapter.PublishEvent(ctx, streamName, event)

	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestServiceEventAdapter_ConsumeStream(t *testing.T) {
	mockProducer := new(MockProducer)
	mockConsumer := new(MockConsumer)
	adapter := NewServiceEventAdapter(mockProducer, mockConsumer)

	ctx := context.Background()
	streamName := "test-stream"
	groupName := "test-group"
	batchSize := int64(10)
	pollInterval := 5 * time.Second
	handler := func(context.Context, *events.Event, string) error { return nil }

	mockConsumer.On("ConsumeStream", ctx, streamName, groupName, batchSize, pollInterval, mock.AnythingOfType("func(context.Context, *events.Event, string) error")).Return(nil)

	err := adapter.ConsumeStream(ctx, streamName, groupName, batchSize, pollInterval, handler)

	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

func TestServiceEventAdapter_PublishProductEvent(t *testing.T) {
	mockProducer := new(MockProducer)
	mockConsumer := new(MockConsumer)
	adapter := NewServiceEventAdapter(mockProducer, mockConsumer)

	ctx := context.Background()
	eventType := events.EventTypeNewProductDetected
	productID := "product-123"
	asin := "B07PXGQC1Q"
	payload := events.NewProductDetectedPayload{
		ASIN:  asin,
		Title: "Test Product",
	}

	// Mock the expected call
	mockProducer.On("PublishEvent", ctx, "stream:product_lifecycle", mock.MatchedBy(func(event *events.Event) bool {
		return event.Type == eventType &&
			event.AggregateType == "product" &&
			event.AggregateID == productID
	})).Return(nil)

	err := adapter.PublishProductEvent(ctx, eventType, productID, asin, payload)

	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestServiceEventAdapter_PublishContentEvent(t *testing.T) {
	mockProducer := new(MockProducer)
	mockConsumer := new(MockConsumer)
	adapter := NewServiceEventAdapter(mockProducer, mockConsumer)

	ctx := context.Background()
	eventType := events.EventTypeContentGenerationRequested
	productID := "product-123"
	payload := events.ContentGenerationRequestedPayload{
		ASIN:        "B07PXGQC1Q",
		ProductID:   productID,
		RequestedAt: time.Now(),
	}

	// Mock the expected call
	mockProducer.On("PublishEvent", ctx, "stream:content_generation", mock.MatchedBy(func(event *events.Event) bool {
		return event.Type == eventType &&
			event.AggregateType == "content" &&
			event.AggregateID == productID
	})).Return(nil)

	err := adapter.PublishContentEvent(ctx, eventType, productID, payload)

	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestServiceEventAdapter_PublishBrowseNodeEvent(t *testing.T) {
	mockProducer := new(MockProducer)
	mockConsumer := new(MockConsumer)
	adapter := NewServiceEventAdapter(mockProducer, mockConsumer)

	ctx := context.Background()
	eventType := events.EventTypeBrowseNodeRequested
	productID := "product-123"
	payload := events.BrowseNodeRequestedPayload{
		ASIN:        "B07PXGQC1Q",
		ProductID:   productID,
		RequestedAt: time.Now(),
	}

	// Mock the expected call
	mockProducer.On("PublishEvent", ctx, "stream:browse_nodes", mock.MatchedBy(func(event *events.Event) bool {
		return event.Type == eventType &&
			event.AggregateType == "browse_node" &&
			event.AggregateID == productID
	})).Return(nil)

	err := adapter.PublishBrowseNodeEvent(ctx, eventType, productID, payload)

	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestDetermineTargetStream(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		expected  string
	}{
		{
			name:      "product lifecycle event",
			eventType: events.EventTypeNewProductDetected,
			expected:  "stream:product_lifecycle",
		},
		{
			name:      "content generation event",
			eventType: events.EventTypeContentGenerationRequested,
			expected:  "stream:content_generation",
		},
		{
			name:      "browse node event",
			eventType: events.EventTypeBrowseNodeRequested,
			expected:  "stream:browse_nodes",
		},
		{
			name:      "price tracking event",
			eventType: events.EventTypeCheckPrice,
			expected:  "stream:price_tracking",
		},
		{
			name:      "unknown event type",
			eventType: "UNKNOWN_EVENT",
			expected:  "stream:product_lifecycle", // Default fallback
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetermineTargetStream(tt.eventType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEventCategorization(t *testing.T) {
	t.Run("isProductLifecycleEvent", func(t *testing.T) {
		assert.True(t, isProductLifecycleEvent(events.EventTypeNewProductDetected))
		assert.True(t, isProductLifecycleEvent(events.EventTypeProductValidated))
		assert.False(t, isProductLifecycleEvent(events.EventTypeContentGenerated))
		assert.False(t, isProductLifecycleEvent("UNKNOWN_EVENT"))
	})

	t.Run("isContentGenerationEvent", func(t *testing.T) {
		assert.True(t, isContentGenerationEvent(events.EventTypeContentGenerationRequested))
		assert.True(t, isContentGenerationEvent(events.EventTypeReviewsCollected))
		assert.False(t, isContentGenerationEvent(events.EventTypeNewProductDetected))
		assert.False(t, isContentGenerationEvent("UNKNOWN_EVENT"))
	})

	t.Run("isBrowseNodeEvent", func(t *testing.T) {
		assert.True(t, isBrowseNodeEvent(events.EventTypeBrowseNodeRequested))
		assert.True(t, isBrowseNodeEvent(events.EventTypeBrowseNodeResolved))
		assert.False(t, isBrowseNodeEvent(events.EventTypeNewProductDetected))
		assert.False(t, isBrowseNodeEvent("UNKNOWN_EVENT"))
	})

	t.Run("isPriceTrackingEvent", func(t *testing.T) {
		assert.True(t, isPriceTrackingEvent(events.EventTypeCheckPrice))
		assert.True(t, isPriceTrackingEvent(events.EventTypePriceUpdated))
		assert.False(t, isPriceTrackingEvent(events.EventTypeNewProductDetected))
		assert.False(t, isPriceTrackingEvent("UNKNOWN_EVENT"))
	})
}
