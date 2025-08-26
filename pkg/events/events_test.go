package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEvent(t *testing.T) {
	tests := []struct {
		name          string
		eventType     string
		aggregateType string
		aggregateID   string
		payload       any
		wantErr       bool
	}{
		{
			name:          "valid product event",
			eventType:     EventTypeNewProductDetected,
			aggregateType: "product",
			aggregateID:   "product-123",
			payload: NewProductDetectedPayload{
				ASIN:  "B07PXGQC1Q",
				Title: "Test Product",
			},
			wantErr: false,
		},
		{
			name:          "valid content event",
			eventType:     EventTypeContentGenerationRequested,
			aggregateType: "content",
			aggregateID:   "product-456",
			payload: ContentGenerationRequestedPayload{
				ASIN:        "B07PXGQC1Q",
				ProductID:   "product-456",
				RequestedAt: time.Now(),
			},
			wantErr: false,
		},
		// Browse node test case removed - events deprecated
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := NewEvent(tt.eventType, tt.aggregateType, tt.aggregateID, tt.payload)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, event.ID)
			assert.Equal(t, tt.eventType, event.Type)
			assert.Equal(t, tt.aggregateType, event.AggregateType)
			assert.Equal(t, tt.aggregateID, event.AggregateID)
			assert.NotZero(t, event.Timestamp)
			assert.Equal(t, tt.payload, event.Payload)
			assert.NotNil(t, event.Metadata)
		})
	}
}

func TestEventUnmarshalPayload(t *testing.T) {
	// Create an event with a known payload
	originalPayload := ContentGenerationRequestedPayload{
		ASIN:        "B07PXGQC1Q",
		ProductID:   "product-123",
		Priority:    5,
		RequestedAt: time.Now().UTC().Truncate(time.Second), // Truncate for comparison
	}

	event, err := NewEvent(
		EventTypeContentGenerationRequested,
		"content",
		"product-123",
		originalPayload,
	)
	require.NoError(t, err)

	// Unmarshal the payload
	var unmarshaledPayload ContentGenerationRequestedPayload
	err = event.UnmarshalPayload(&unmarshaledPayload)
	require.NoError(t, err)

	// Compare the payloads
	assert.Equal(t, originalPayload.ASIN, unmarshaledPayload.ASIN)
	assert.Equal(t, originalPayload.ProductID, unmarshaledPayload.ProductID)
	assert.Equal(t, originalPayload.Priority, unmarshaledPayload.Priority)
	assert.Equal(t, originalPayload.RequestedAt.Unix(), unmarshaledPayload.RequestedAt.Unix())
}

func TestEventTypes(t *testing.T) {
	// Test that all event types are properly defined
	eventTypes := []string{
		// Product Lifecycle Events
		EventTypeNewProductDetected,
		EventTypeProductValidated,
		EventTypeProductUnavailable,
		EventTypeProductDeleted,
		EventTypeProductCreated,
		EventTypeProductUpdated,
		EventTypeProductAvailabilityChanged,
		EventTypeProductStatusChanged,
		EventTypeProductUpdateRequested,
		EventTypeProductIgnored,
		EventTypeProductReviewRequired,

		// Content Generation Events
		EventTypeContentGenerationRequested,
		EventTypeContentGenerationStarted,
		EventTypeContentGenerated,
		EventTypeContentGenerationFailed,
		EventTypeContentUpdateRequested,
		EventTypeContentUpdated,
		EventTypeContentGenerationRetried,
		EventTypeContentAnalysisFailed,

		// Reviews Events
		EventTypeReviewsRequested,
		EventTypeReviewsFetched,
		EventTypeReviewsStored,
		EventTypeReviewsCollected,
		EventTypeReviewsFetchFailed,
		EventTypeReviewsProcessed,
		EventTypeReviewsValidated,
		EventTypeReviewsEnriched,
		EventTypeReviewsCached,
		EventTypeReviewsExpired,
		EventTypeReviewsError,
		EventTypeReviewsDeleted,

		// Browse Node Events removed

		// Price Tracking Events
		EventTypeCheckPrice,
		EventTypePriceUpdated,
		EventTypePriceUpdateFailed,

		// Orchestration Events
		// Dimension enrichment events removed - replaced by PA-API enrichment
		QualityAssessmentRequested,
		QualityAssessmentCompleted,
		QualityAssessmentFailed,
		PriceMonitoringScheduled,
		AvailabilityCheckScheduled,
		PeriodicUpdateScheduled,

		// PA-API Enrichment Events
		ColorEnrichmentRequested,
		ColorEnrichmentCompleted,
		ColorEnrichmentFailed,
		VariationEnrichmentRequested,
		VariationEnrichmentCompleted,
		VariationEnrichmentFailed,

		// PA-API Enrichment Events (CloudEvent format)
		ProductEnrichmentRequestedV1,
		ProductEnrichmentCompletedV1,
		ProductEnrichmentFailedV1,
	}

	// Ensure all event types are non-empty strings
	for _, eventType := range eventTypes {
		assert.NotEmpty(t, eventType, "Event type should not be empty")
		assert.IsType(t, "", eventType, "Event type should be a string")
	}

	// Ensure we have a reasonable number of event types
	assert.GreaterOrEqual(t, len(eventTypes), 30, "Should have at least 30 event types defined")
}

func TestPayloadStructures(t *testing.T) {
	t.Run("ProductCreatedPayload", func(t *testing.T) {
		payload := ProductCreatedPayload{
			ASIN:          "B07PXGQC1Q",
			Title:         "Test Product",
			Brand:         "Test Brand",
			Category:      "Clothing",
			Gender:        "Unisex",
			CurrentPrice:  29.99,
			Currency:      "USD",
			DetailPageURL: "https://amazon.com/dp/B07PXGQC1Q",
			ImageUrls:     []string{"https://example.com/image1.jpg"},
			Features:      []string{"Feature 1", "Feature 2"},
			BrowseNodeID:  "123456",
		}

		assert.Equal(t, "B07PXGQC1Q", payload.ASIN)
		assert.Equal(t, "Test Product", payload.Title)
		assert.Equal(t, 29.99, payload.CurrentPrice)
	})

	t.Run("BrowseNodeRequestedPayload", func(t *testing.T) {
		// BrowseNodeRequestedPayload test removed - payload type deprecated
	})

	t.Run("ContentGeneratedPayload", func(t *testing.T) {
		now := time.Now()
		payload := ContentGeneratedPayload{
			ASIN:        "B07PXGQC1Q",
			ProductID:   "product-123",
			ContentType: "description",
			ContentID:   "content-456",
			GeneratedAt: now,
		}

		assert.Equal(t, "B07PXGQC1Q", payload.ASIN)
		assert.Equal(t, "description", payload.ContentType)
		assert.Equal(t, now, payload.GeneratedAt)
	})
}

func TestEventMetadata(t *testing.T) {
	event, err := NewEvent(
		EventTypeNewProductDetected,
		"product",
		"product-123",
		map[string]string{"test": "value"},
	)
	require.NoError(t, err)

	// Test that metadata is initialized
	assert.NotNil(t, event.Metadata)

	// Test adding metadata
	event.Metadata["source"] = "test"
	event.Metadata["version"] = "1.0"

	assert.Equal(t, "test", event.Metadata["source"])
	assert.Equal(t, "1.0", event.Metadata["version"])
}
