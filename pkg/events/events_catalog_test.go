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