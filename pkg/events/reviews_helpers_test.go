package events

import (
	"testing"
)

func TestNewReviewsRequestedEvent(t *testing.T) {
	asin := "B00EXAMPLE"
	productID := "product-123"
	source := ReviewsSourceAmazonAPI
	options := map[string]any{
		"pages": 3,
		"sort":  "recent",
	}

	event := NewReviewsRequestedEvent(asin, productID, source, options)

	if event.Type != EventTypeReviewsRequested {
		t.Errorf("Expected event type %s, got %s", EventTypeReviewsRequested, event.Type)
	}

	if event.AggregateID != productID {
		t.Errorf("Expected aggregate ID %s, got %s", productID, event.AggregateID)
	}

	payload, ok := event.Payload.(ReviewsRequestedPayload)
	if !ok {
		t.Errorf("Expected ReviewsRequestedPayload, got %T", event.Payload)
	}

	if payload.ASIN != asin {
		t.Errorf("Expected ASIN %s, got %s", asin, payload.ASIN)
	}

	if payload.ProductID != productID {
		t.Errorf("Expected ProductID %s, got %s", productID, payload.ProductID)
	}

	if payload.Source != source {
		t.Errorf("Expected Source %s, got %s", source, payload.Source)
	}

	if payload.Options["pages"] != 3 {
		t.Errorf("Expected pages option 3, got %v", payload.Options["pages"])
	}
}

func TestNewReviewsErrorEvent(t *testing.T) {
	asin := "B00EXAMPLE"
	productID := "product-123"
	source := ReviewsSourceAmazonAPI
	errorType := ReviewsErrorTypeFetch
	errorMessage := "API rate limit exceeded"
	retryCount := 2

	event := NewReviewsErrorEvent(asin, productID, source, errorType, errorMessage, retryCount)

	if event.Type != EventTypeReviewsError {
		t.Errorf("Expected event type %s, got %s", EventTypeReviewsError, event.Type)
	}

	payload, ok := event.Payload.(ReviewsErrorPayload)
	if !ok {
		t.Errorf("Expected ReviewsErrorPayload, got %T", event.Payload)
	}

	if payload.ErrorType != errorType {
		t.Errorf("Expected ErrorType %s, got %s", errorType, payload.ErrorType)
	}

	if payload.ErrorMessage != errorMessage {
		t.Errorf("Expected ErrorMessage %s, got %s", errorMessage, payload.ErrorMessage)
	}

	if payload.RetryCount != retryCount {
		t.Errorf("Expected RetryCount %d, got %d", retryCount, payload.RetryCount)
	}
}

func TestIsReviewsEvent(t *testing.T) {
	tests := []struct {
		eventType string
		expected  bool
	}{
		{EventTypeReviewsRequested, true},
		{EventTypeReviewsFetched, true},
		{EventTypeReviewsError, true},
		{EventTypeProductCreated, false},
		{EventTypeContentGenerated, false},
		{"UNKNOWN_EVENT", false},
	}

	for _, test := range tests {
		result := IsReviewsEvent(test.eventType)
		if result != test.expected {
			t.Errorf("IsReviewsEvent(%s): expected %v, got %v", test.eventType, test.expected, result)
		}
	}
}

func TestGetReviewsEventPriority(t *testing.T) {
	tests := []struct {
		eventType string
		expected  int
	}{
		{EventTypeReviewsError, 1}, // Highest priority
		{EventTypeReviewsRequested, 2},
		{EventTypeReviewsFetched, 3},
		{EventTypeReviewsDeleted, 10}, // Lowest priority
		{"UNKNOWN_EVENT", 999},        // Unknown event
	}

	for _, test := range tests {
		result := GetReviewsEventPriority(test.eventType)
		if result != test.expected {
			t.Errorf("GetReviewsEventPriority(%s): expected %d, got %d", test.eventType, test.expected, result)
		}
	}
}

func TestAllReviewsEventTypes(t *testing.T) {
	// Test that all event types are properly defined
	eventTypes := []string{
		EventTypeReviewsRequested,
		EventTypeReviewsFetched,
		EventTypeReviewsStored,
		EventTypeReviewsProcessed,
		EventTypeReviewsValidated,
		EventTypeReviewsEnriched,
		EventTypeReviewsCached,
		EventTypeReviewsExpired,
		EventTypeReviewsError,
		EventTypeReviewsDeleted,
	}

	if len(eventTypes) != 10 {
		t.Errorf("Expected 10 reviews event types, got %d", len(eventTypes))
	}

	// Test that all are recognized as reviews events
	for _, eventType := range eventTypes {
		if !IsReviewsEvent(eventType) {
			t.Errorf("Event type %s should be recognized as reviews event", eventType)
		}
	}
}
