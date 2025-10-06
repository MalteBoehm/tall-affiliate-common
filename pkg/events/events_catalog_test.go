package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCatalogProductEnrichmentRequestedEvent(t *testing.T) {
	data := &ProductEnrichmentRequestedData{
		ASIN:       "B07PXGQC1Q",
		Region:     "de",
		RequestID:  "req-123",
		RetryCount: 0,
	}

	event, err := NewCatalogProductEnrichmentRequestedEvent("test-service", data)
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, CatalogProductEnrichmentRequestedV1, event.Type)
	assert.Equal(t, "catalog.product", event.AggregateType)
	assert.Equal(t, "B07PXGQC1Q", event.AggregateID)

	var payload ProductEnrichmentRequestedData
	err = event.UnmarshalPayload(&payload)
	require.NoError(t, err)
	assert.Equal(t, data.ASIN, payload.ASIN)
	assert.Equal(t, data.Region, payload.Region)
	assert.Equal(t, data.RequestID, payload.RequestID)
	assert.Equal(t, data.RetryCount, payload.RetryCount)
}

func TestNewCatalogProductEnrichmentCompletedEvent(t *testing.T) {
	now := time.Now().UTC()
	data := &ProductEnrichedData{
		ASIN:      "B07PXGQC1Q",
		Region:    "de",
		RequestID: "req-123",
		ColorVariants: []ColorVariant{
			{
				ASIN:      "B07PXGQC1R",
				ColorName: "Red",
				Images:    []ImageSet{},
			},
		},
		ProcessingMS: 1500,
		EnrichedAt:   now,
	}

	event, err := NewCatalogProductEnrichmentCompletedEvent("test-service", data)
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, CatalogProductEnrichmentCompletedV1, event.Type)
	assert.Equal(t, "catalog.product", event.AggregateType)
	assert.Equal(t, "B07PXGQC1Q", event.AggregateID)

	var payload ProductEnrichedData
	err = event.UnmarshalPayload(&payload)
	require.NoError(t, err)
	assert.Equal(t, data.ASIN, payload.ASIN)
	assert.Equal(t, data.Region, payload.Region)
	assert.Equal(t, data.RequestID, payload.RequestID)
	assert.Equal(t, len(data.ColorVariants), len(payload.ColorVariants))
	assert.Equal(t, data.ProcessingMS, payload.ProcessingMS)
}

func TestNewCatalogProductEnrichmentFailedEvent(t *testing.T) {
	now := time.Now().UTC()
	data := &ProductEnrichmentFailedData{
		ASIN:         "B07PXGQC1Q",
		Region:       "de",
		ErrorCode:    "RATE_LIMIT",
		ErrorMessage: "API rate limit exceeded",
		RequestID:    "req-123",
		RetryCount:   3,
		FailedAt:     now,
	}

	event, err := NewCatalogProductEnrichmentFailedEvent("test-service", data)
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, CatalogProductEnrichmentFailedV1, event.Type)
	assert.Equal(t, "catalog.product", event.AggregateType)
	assert.Equal(t, "B07PXGQC1Q", event.AggregateID)

	var payload ProductEnrichmentFailedData
	err = event.UnmarshalPayload(&payload)
	require.NoError(t, err)
	assert.Equal(t, data.ASIN, payload.ASIN)
	assert.Equal(t, data.Region, payload.Region)
	assert.Equal(t, data.ErrorCode, payload.ErrorCode)
	assert.Equal(t, data.ErrorMessage, payload.ErrorMessage)
	assert.Equal(t, data.RequestID, payload.RequestID)
	assert.Equal(t, data.RetryCount, payload.RetryCount)
}

func TestCatalogProductEnrichmentRequestedEventValidation(t *testing.T) {
	tests := []struct {
		name    string
		data    *ProductEnrichmentRequestedData
		wantErr bool
	}{
		{
			name: "valid data",
			data: &ProductEnrichmentRequestedData{
				ASIN:      "B07PXGQC1Q",
				Region:    "de",
				RequestID: "req-123",
			},
			wantErr: false,
		},
		{
			name: "missing ASIN",
			data: &ProductEnrichmentRequestedData{
				Region:    "de",
				RequestID: "req-123",
			},
			wantErr: true,
		},
		{
			name: "missing region",
			data: &ProductEnrichmentRequestedData{
				ASIN:      "B07PXGQC1Q",
				RequestID: "req-123",
			},
			wantErr: true,
		},
		{
			name: "missing request ID",
			data: &ProductEnrichmentRequestedData{
				ASIN:   "B07PXGQC1Q",
				Region: "de",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := NewCatalogProductEnrichmentRequestedEvent("test-service", tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
			}
		})
	}
}

func TestNormalizeEventTypeCatalogEvents(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		found    bool
	}{
		{
			input:    "SCRAPER_JOB_REQUESTED",
			expected: Event_00A_ScraperJobRequested,
			found:    true,
		},
		{
			input:    Event_00A_ScraperJobRequested,
			expected: Event_00A_ScraperJobRequested,
			found:    true,
		},
		{
			input:    "catalog.product.enrichment.requested.v1",
			expected: CatalogProductEnrichmentRequestedV1,
			found:    true,
		},
		{
			input:    "catalog.product.enrichment.completed.v1",
			expected: CatalogProductEnrichmentCompletedV1,
			found:    true,
		},
		{
			input:    "catalog.product.enrichment.failed.v1",
			expected: CatalogProductEnrichmentFailedV1,
			found:    true,
		},
		{
			input:    "product.enrichment.requested.v1",
			expected: ProductEnrichmentRequestedV1,
			found:    true,
		},
		{
			input:    "product.enrichment.completed.v1",
			expected: ProductEnrichmentCompletedV1,
			found:    true,
		},
		{
			input:    "product.enrichment.failed.v1",
			expected: ProductEnrichmentFailedV1,
			found:    true,
		},
		{
			input:    "unknown.event.type",
			expected: "unknown.event.type",
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, found := NormalizeEventType(tt.input)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.found, found)
		})
	}
}

func TestNewScraperJobRequestedEvent(t *testing.T) {
	jobID := "job-123"
	searchQuery := "tall jeans"
	category := "fashion"
	maxPages := 5

	event := NewScraperJobRequestedEvent(jobID, searchQuery, category, maxPages)
	require.NotNil(t, event)
	assert.Equal(t, EventTypeScraperJobRequested, event.Type)
	assert.Equal(t, "scraper.job", event.AggregateType)
	assert.Equal(t, jobID, event.AggregateID)

	var payload ScraperJobRequestedPayload
	require.NoError(t, event.UnmarshalPayload(&payload))
	assert.Equal(t, jobID, payload.JobID)
	assert.Equal(t, searchQuery, payload.SearchQuery)
	assert.Equal(t, category, payload.Category)
	assert.Equal(t, maxPages, payload.MaxPages)
	assert.WithinDuration(t, time.Now().UTC(), payload.RequestedAt, 2*time.Second)
}

func TestNormalizeEventTypeContentGeneration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		found    bool
	}{
		{
			input:    "CONTENT_GENERATION_REQUESTED",
			expected: Event_08A_ContentGenerationRequested,
			found:    true,
		},
		{
			input:    "CONTENT_GENERATION_STARTED",
			expected: Event_09A_ContentGenerationStarted,
			found:    true,
		},
		{
			input:    "CONTENT_GENERATED",
			expected: Event_10A_ContentGenerated,
			found:    true,
		},
		{
			input:    "CONTENT_GENERATION_FAILED",
			expected: Event_10B_ContentGenerationFailed,
			found:    true,
		},
		{
			input:    "CONTENT_GENERATION_RETRIED",
			expected: Event_10D_ContentGenerationRetried,
			found:    true,
		},
		{
			input:    "REVIEWS_REQUESTED",
			expected: Event_08B_ReviewsRequested,
			found:    true,
		},
		{
			input:    "REVIEWS_FETCHED",
			expected: Event_09B_ReviewsFetched,
			found:    true,
		},
		{
			input:    "REVIEWS_PROCESSED",
			expected: Event_10C_ReviewsProcessed,
			found:    true,
		},
		{
			input:    "REVIEWS_VALIDATED",
			expected: Event_11A_ReviewsValidated,
			found:    true,
		},
		{
			input:    "REVIEWS_FETCH_FAILED",
			expected: Event_11B_ReviewsFetchFailed,
			found:    true,
		},
		{
			input:    "REVIEWS_STORED",
			expected: Event_12A_ReviewsStored,
			found:    true,
		},
		{
			input:    "REVIEWS_ERROR",
			expected: Event_12B_ReviewsError,
			found:    true,
		},
		{
			input:    "content.generation.requested.v1",
			expected: ContentGenerationRequestedV1,
			found:    true,
		},
		{
			input:    "content.generation.started.v1",
			expected: ContentGenerationStartedV1,
			found:    true,
		},
		{
			input:    "content.generated.v1",
			expected: ContentGeneratedV1,
			found:    true,
		},
		{
			input:    "content.generation.failed.v1",
			expected: ContentGenerationFailedV1,
			found:    true,
		},
		{
			input:    "content.generation.retried.v1",
			expected: ContentGenerationRetriedV1,
			found:    true,
		},
		{
			input:    "reviews.requested.v1",
			expected: ReviewsRequestedV1,
			found:    true,
		},
		{
			input:    "reviews.fetched.v1",
			expected: ReviewsFetchedV1,
			found:    true,
		},
		{
			input:    "reviews.processed.v1",
			expected: ReviewsProcessedV1,
			found:    true,
		},
		{
			input:    "reviews.validated.v1",
			expected: ReviewsValidatedV1,
			found:    true,
		},
		{
			input:    "reviews.fetch_failed.v1",
			expected: ReviewsFetchFailedV1,
			found:    true,
		},
		{
			input:    "reviews.stored.v1",
			expected: ReviewsStoredV1,
			found:    true,
		},
		{
			input:    "reviews.error.v1",
			expected: ReviewsErrorV1,
			found:    true,
		},
		{
			input:    "UNKNOWN_CONTENT_EVENT",
			expected: "UNKNOWN_CONTENT_EVENT",
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, found := NormalizeEventType(tt.input)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.found, found)
		})
	}
}

func TestContentGenerationEventCreationAndNormalization(t *testing.T) {
	// Test creating events with legacy types and normalizing them
	now := time.Now().UTC()

	// Legacy content generation requested
	legacyType := "CONTENT_GENERATION_REQUESTED"
	payload := ContentGenerationRequestedPayload{
		ASIN:        "B07PXGQC1Q",
		ProductID:   "product-123",
		Priority:    4,
		RequestedAt: now,
	}

	event, err := NewEvent(legacyType, "content", "product-123", payload)
	require.NoError(t, err)

	// Normalize the event type
	normalizedType, found := NormalizeEventType(event.Type)
	require.True(t, found)
	assert.Equal(t, Event_08A_ContentGenerationRequested, normalizedType)

	// Update event with normalized type
	event.Type = normalizedType

	// Verify payload roundtrip after normalization
	var unmarshaledPayload ContentGenerationRequestedPayload
	err = event.UnmarshalPayload(&unmarshaledPayload)
	require.NoError(t, err)

	assert.Equal(t, payload.ASIN, unmarshaledPayload.ASIN)
	assert.Equal(t, payload.ProductID, unmarshaledPayload.ProductID)
	assert.Equal(t, payload.Priority, unmarshaledPayload.Priority)
	assert.WithinDuration(t, payload.RequestedAt, unmarshaledPayload.RequestedAt, time.Second)

	// Test CloudEvents format normalization
	cloudEventType := "content.generation.requested.v1"
	normalizedCloudType, foundCloud := NormalizeEventType(cloudEventType)
	require.True(t, foundCloud)
	assert.Equal(t, ContentGenerationRequestedV1, normalizedCloudType)
}
