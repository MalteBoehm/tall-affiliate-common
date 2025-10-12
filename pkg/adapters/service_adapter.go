// pkg/adapters/service_adapter.go
package adapters

import (
	"context"
	"time"

	"github.com/MalteBoehm/tall-affiliate-common/pkg/events"
	"github.com/MalteBoehm/tall-affiliate-common/pkg/interfaces"
)

// ServiceEventAdapter provides a unified adapter for all services
// This ensures all services use the same event handling pattern
type ServiceEventAdapter struct {
	producer interfaces.StreamProducer
	consumer interfaces.StreamConsumer
}

// NewServiceEventAdapter creates a new unified service adapter
func NewServiceEventAdapter(producer interfaces.StreamProducer, consumer interfaces.StreamConsumer) *ServiceEventAdapter {
	return &ServiceEventAdapter{
		producer: producer,
		consumer: consumer,
	}
}

// PublishEvent publishes an event using the common interface
func (s *ServiceEventAdapter) PublishEvent(ctx context.Context, streamName string, event *events.Event) error {
	return s.producer.PublishEvent(ctx, streamName, event)
}

// ConsumeStream consumes events using the common interface
func (s *ServiceEventAdapter) ConsumeStream(
	ctx context.Context,
	streamName string,
	groupName string,
	batchSize int64,
	pollInterval time.Duration,
	handler func(context.Context, *events.Event, string) error,
) error {
	return s.consumer.ConsumeStream(ctx, streamName, groupName, batchSize, pollInterval, handler)
}

// Helper functions for common event operations

// PublishProductEvent publishes a product-related event
func (s *ServiceEventAdapter) PublishProductEvent(ctx context.Context, eventType, productID, asin string, payload any) error {
	event, err := events.NewEvent(eventType, "product", productID, payload)
	if err != nil {
		return err
	}

	// Determine target stream based on event type
	streamName := DetermineTargetStream(eventType)
	return s.PublishEvent(ctx, streamName, event)
}

// PublishContentEvent publishes a content-related event
func (s *ServiceEventAdapter) PublishContentEvent(ctx context.Context, eventType, productID string, payload any) error {
	event, err := events.NewEvent(eventType, "content", productID, payload)
	if err != nil {
		return err
	}

	streamName := DetermineTargetStream(eventType)
	return s.PublishEvent(ctx, streamName, event)
}

// DEPRECATED: PublishBrowseNodeEvent removed - use PublishEnrichmentEvent with ProductEnrichmentRequestedData instead

// DetermineTargetStream determines the target stream based on event type
func DetermineTargetStream(eventType string) string {
	switch {
	case isProductLifecycleEvent(eventType):
		return "stream:product_lifecycle"
	case isContentGenerationEvent(eventType):
		return "stream:content_generation"
	// DEPRECATED: Browse node events removed
	case isPriceTrackingEvent(eventType):
		return "stream:price_tracking"
	default:
		return "stream:product_lifecycle" // Default fallback
	}
}

// Helper functions to categorize events
// DEPRECATED: These functions will be removed in a future release. Use CAPS constants directly.

func isProductLifecycleEvent(eventType string) bool {
	productEvents := []string{
		events.EventTypeNewProductDetected,
		events.EventTypeProductValidated,
		events.EventTypeProductUnavailable,
		events.EventTypeProductDeleted,
		events.EventTypeProductCreated,
		events.EventTypeProductUpdated,
		events.EventTypeProductAvailabilityChanged,
		events.EventTypeProductStatusChanged,
		events.EventTypeProductIgnored,
		// Also check CAPS constants
		events.CATALOG_PRODUCT_DETECTED_V1,
		events.CATALOG_PRODUCT_VALIDATED_V1,
		events.CATALOG_PRODUCT_IGNORED_V1,
		events.CATALOG_PRODUCT_REVIEW_REQUIRED_V1,
		events.PRODUCT_READY_FOR_PUBLICATION_V1,
		events.PRODUCT_UPDATED_V1,
		events.PRODUCT_UPDATE_FAILED_V1,
		events.PRODUCT_AVAILABILITY_CHANGED_V1,
		events.PRODUCT_STATUS_CHANGED_V1,
		events.PRODUCT_DELETED_V1,
	}

	for _, pe := range productEvents {
		if pe == eventType {
			return true
		}
	}
	return false
}

func isContentGenerationEvent(eventType string) bool {
	contentEvents := []string{
		events.EventTypeContentGenerationRequested,
		events.EventTypeContentGenerated,
		events.EventTypeContentGenerationFailed,
		events.EventTypeContentUpdateRequested,
		events.EventTypeContentUpdated,
		events.EventTypeReviewsCollected,
		events.EventTypeReviewsFetchFailed,
		events.EventTypeReviewsProcessed,
		// Also check CAPS constants
		events.CONTENT_GENERATION_REQUESTED_V1,
		events.CONTENT_GENERATION_STARTED_V1,
		events.CONTENT_GENERATED_V1,
		events.CONTENT_GENERATION_FAILED_V1,
		events.CONTENT_GENERATION_RETRIED_V1,
		events.REVIEWS_REQUESTED_V1,
		events.REVIEWS_FETCHED_V1,
		events.REVIEWS_PROCESSED_V1,
		events.REVIEWS_VALIDATED_V1,
		events.REVIEWS_FETCH_FAILED_V1,
		events.REVIEWS_STORED_V1,
		events.REVIEWS_ERROR_V1,
	}

	for _, ce := range contentEvents {
		if ce == eventType {
			return true
		}
	}
	return false
}

// DEPRECATED: isBrowseNodeEvent removed - browse node events no longer supported

func isPriceTrackingEvent(eventType string) bool {
	priceEvents := []string{
		events.EventTypeCheckPrice,
		events.EventTypePriceUpdated,
		events.EventTypePriceUpdateFailed,
		// Also check CAPS constants
		events.PRICE_UPDATED_V1,
		events.PRICE_UPDATE_FAILED_V1,
		events.PRICE_MONITORING_SCHEDULED_V1,
	}

	for _, pe := range priceEvents {
		if pe == eventType {
			return true
		}
	}
	return false
}

// Ensure ServiceEventAdapter implements both interfaces
var _ interfaces.StreamProducer = (*ServiceEventAdapter)(nil)
var _ interfaces.StreamConsumer = (*ServiceEventAdapter)(nil)
