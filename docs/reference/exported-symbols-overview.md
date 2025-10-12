# Exported Symbols Overview

This document provides a comprehensive catalog of all exported symbols in the tall-affiliate-common package, organized by package and category.

## pkg/constants

### Constants
- **Stream names**
  - `STREAM_CONTENT_GENERATION` - "stream:content_generation"
  - `STREAM_PRODUCT_LIFECYCLE` - "stream:product_lifecycle"
  - `STREAM_PRICE_TRACKING` - "stream:price_tracking"

- **Consumer group names**
  - `GROUP_PRODUCT_LIFECYCLE` - "group:product_lifecycle"
  - `GROUP_CONTENT_GENERATION` - "group:content_generation"
  - `GROUP_PRICE_TRACKING` - "group:price_tracking"

- **Worker names**
  - `WORKER_PRODUCT_LIFECYCLE` - "product_lifecycle_worker"
  - `WORKER_CONTENT_GENERATION` - "content_generation_worker"
  - `WORKER_PRICE_TRACKING` - "price_tracking_worker"

- **Database table names**
  - `TABLE_PRODUCT` - "product"
  - `TABLE_OUTBOX` - "outbox"
  - `TABLE_PRICE_HISTORY` - "price_history"

## pkg/events

### Constants

#### Event Types (Numbered Convention)
- `EVENT_00A_SCRAPER_JOB_REQUESTED` - "00A_SCRAPER_JOB_REQUESTED"
- `EVENT_01_PRODUCT_DETECTED` - "01_PRODUCT_DETECTED"
- `EVENT_02A_PRODUCT_VALIDATED` - "02A_PRODUCT_VALIDATED"
- `EVENT_02B_PRODUCT_IGNORED` - "02B_PRODUCT_IGNORED"
- `EVENT_02C_PRODUCT_REVIEW_REQUIRED` - "02C_PRODUCT_REVIEW_REQUIRED"
- `EVENT_03_ENRICHMENT_ORCHESTRATION_STARTED` - "03_ENRICHMENT_ORCHESTRATION_STARTED"
- `EVENT_04B_COLOR_ENRICHMENT_REQUESTED` - "04B_COLOR_ENRICHMENT_REQUESTED"
- `EVENT_04D_VARIANTS_ENRICHMENT_REQUESTED` - "04D_VARIANTS_ENRICHMENT_REQUESTED"
- `EVENT_05A_ENRICHMENT_COMPLETED` - "enrichment.completed.v1"
- `EVENT_05B_ENRICHMENT_FAILED` - "enrichment.failed.v1"
- `EVENT_05C_VARIANTS_ENRICHED` - "variants.enriched.v1"
- `EVENT_05D_ENRICHMENT_RETRY` - "enrichment.retry.v1"
- `EVENT_06_QUALITY_ASSESSMENT_REQUESTED` - "quality.assessment.requested.v1"
- `EVENT_07A_QUALITY_ASSESSMENT_COMPLETED` - "quality.assessment.completed.v1"
- `EVENT_07B_QUALITY_ASSESSMENT_FAILED` - "quality.assessment.failed.v1"
- `EVENT_08A_CONTENT_GENERATION_REQUESTED` - "08A_CONTENT_GENERATION_REQUESTED"
- `EVENT_08B_REVIEWS_REQUESTED` - "reviews.requested.v1"
- `EVENT_09A_CONTENT_GENERATION_STARTED` - "09A_CONTENT_GENERATION_STARTED"
- `EVENT_09B_REVIEWS_FETCHED` - "reviews.fetched.v1"
- `EVENT_10A_CONTENT_GENERATED` - "10A_CONTENT_GENERATED"
- `EVENT_10B_CONTENT_GENERATION_FAILED` - "10B_CONTENT_GENERATION_FAILED"
- `EVENT_10C_REVIEWS_PROCESSED` - "reviews.processed.v1"
- `EVENT_10D_CONTENT_GENERATION_RETRIED` - "content.generation.retried.v1"
- `EVENT_11A_REVIEWS_VALIDATED` - "reviews.validated.v1"
- `EVENT_11B_REVIEWS_FETCH_FAILED` - "reviews.fetch_failed.v1"
- `EVENT_12A_REVIEWS_STORED` - "reviews.stored.v1"
- `EVENT_12B_REVIEWS_ERROR` - "reviews.error.v1"
- `EVENT_13_PRODUCT_READY_FOR_PUBLICATION` - "product.ready_for_publication.v1"
- `EVENT_14A_PRICE_MONITORING_SCHEDULED` - "price.monitoring.scheduled.v1"
- `EVENT_14B_AVAILABILITY_CHECK_SCHEDULED` - "product.availability_check.scheduled.v1"
- `EVENT_14C_PERIODIC_UPDATE_SCHEDULED` - "product.periodic_update.scheduled.v1"
- `EVENT_15A_PRICE_UPDATED` - "price.updated.v1"
- `EVENT_15B_PRICE_UPDATE_FAILED` - "price.update_failed.v1"
- `EVENT_16A_PRODUCT_UPDATED` - "product.updated.v1"
- `EVENT_16B_PRODUCT_UPDATE_FAILED` - "product.update_failed.v1"
- `EVENT_17_PRODUCT_AVAILABILITY_CHANGED` - "17_PRODUCT_AVAILABILITY_CHANGED"
- `EVENT_18_PRODUCT_STATUS_CHANGED` - "product.status.changed.v1"
- `EVENT_19_PRODUCT_DELETED` - "product.deleted.v1"

#### Legacy Event Types (Deprecated)
- `EventTypeScraperJobRequested` - Event_00A_ScraperJobRequested
- `EventTypeNewProductDetected` - Event_01_ProductDetected
- `EventTypeProductValidated` - Event_02A_ProductValidated
- `EventTypeProductIgnored` - Event_02B_ProductIgnored
- `EventTypeProductReviewRequired` - Event_02C_ProductReviewRequired
- `EventTypeContentGenerationRequested` - Event_08A_ContentGenerationRequested
- `EventTypeContentGenerationStarted` - Event_09A_ContentGenerationStarted
- `EventTypeContentGenerated` - Event_10A_ContentGenerated
- `EventTypeContentGenerationFailed` - Event_10B_ContentGenerationFailed
- `EventTypeContentGenerationRetried` - Event_10D_ContentGenerationRetried
- `EventTypeReviewsRequested` - Event_08B_ReviewsRequested
- `EventTypeReviewsFetched` - Event_09B_ReviewsFetched
- `EventTypeReviewsProcessed` - Event_10C_ReviewsProcessed
- `EventTypeReviewsValidated` - Event_11A_ReviewsValidated
- `EventTypeReviewsFetchFailed` - Event_11B_ReviewsFetchFailed
- `EventTypeReviewsStored` - Event_12A_ReviewsStored
- `EventTypeReviewsError` - Event_12B_ReviewsError
- `EventTypeCheckPrice` - "CHECK_PRICE"
- `EventTypePriceUpdated` - Event_15A_PriceUpdated
- `EventTypePriceUpdateFailed` - Event_15B_PriceUpdateFailed
- `EventTypeProductUpdated` - Event_16A_ProductUpdated
- `EventTypeProductAvailabilityChanged` - Event_17_ProductAvailabilityChanged
- `EventTypeProductStatusChanged` - Event_18_ProductStatusChanged
- `EventTypeProductDeleted` - Event_19_ProductDeleted
- `EventTypeProductUnavailable` - "PRODUCT_UNAVAILABLE"
- `EventTypeProductCreated` - "PRODUCT_CREATED"
- `EventTypeProductUpdateRequested` - "PRODUCT_UPDATE_REQUESTED"
- `EventTypeContentUpdateRequested` - "CONTENT_UPDATE_REQUESTED"
- `EventTypeContentUpdated` - "CONTENT_UPDATED"
- `EventTypeContentAnalysisFailed` - "CONTENT_ANALYSIS_FAILED"
- `EventTypeReviewsCollected` - "REVIEWS_COLLECTED"
- `EventTypeReviewsEnriched` - "REVIEWS_ENRICHED"
- `EventTypeReviewsCached` - "REVIEWS_CACHED"
- `EventTypeReviewsExpired` - "REVIEWS_EXPIRED"
- `EventTypeReviewsDeleted` - "REVIEWS_DELETED"

#### Legacy Orchestration Event Names
- `QualityAssessmentRequested` - Event_06_QualityAssessmentRequested
- `QualityAssessmentCompleted` - Event_07A_QualityAssessmentCompleted
- `QualityAssessmentFailed` - Event_07B_QualityAssessmentFailed
- `PriceMonitoringScheduled` - Event_14A_PriceMonitoringScheduled
- `AvailabilityCheckScheduled` - Event_14B_AvailabilityCheckScheduled
- `PeriodicUpdateScheduled` - Event_14C_PeriodicUpdateScheduled
- `ColorEnrichmentRequested` - Event_04B_ColorEnrichmentRequested
- `ColorEnrichmentCompleted` - "05A_COLOR_ENRICHMENT_COMPLETED"
- `ColorEnrichmentFailed` - "05B_COLOR_ENRICHMENT_FAILED"
- `VariationEnrichmentRequested` - "04B_VARIATION_ENRICHMENT_REQUESTED"
- `VariationEnrichmentCompleted` - "05A_VARIATION_ENRICHMENT_COMPLETED"
- `VariationEnrichmentFailed` - "05B_VARIATION_ENRICHMENT_FAILED"
- `VariantsEnrichmentRequested` - Event_04D_VariantsEnrichmentRequested
- `VariantsEnrichmentCompleted` - Event_05C_VariantsEnriched
- `VariantsEnrichmentFailed` - "05B_VARIANTS_ENRICHMENT_FAILED"

#### PA-API Enrichment Event Types (CloudEvent Format - Deprecated)
- `ProductEnrichmentRequestedV1` - "product.enrichment.requested.v1"
- `ProductEnrichmentCompletedV1` - "product.enrichment.completed.v1"
- `ProductEnrichmentFailedV1` - "product.enrichment.failed.v1"
- `VariantsEnrichmentRequestedV1` - "product.variants.enrichment.requested.v1"
- `VariantsEnrichmentCompletedV1` - "product.variants.enrichment.completed.v1"
- `VariantsEnrichmentFailedV1` - "product.variants.enrichment.failed.v1"

#### Canonical CloudEvents Types
- `CatalogProductDetectedV1` - "catalog.product.detected.v1"
- `CatalogProductValidatedV1` - "catalog.product.validated.v1"
- `CatalogProductIgnoredV1` - "catalog.product.ignored.v1"
- `CatalogProductReviewRequiredV1` - "catalog.product.review_required.v1"
- `ContentGenerationRequestedV1` - "content.generation.requested.v1"
- `ContentGenerationStartedV1` - "content.generation.started.v1"
- `ContentGeneratedV1` - "content.generated.v1"
- `ContentGenerationFailedV1` - "content.generation.failed.v1"
- `ContentGenerationRetriedV1` - "content.generation.retried.v1"
- `ReviewsRequestedV1` - "reviews.requested.v1"
- `ReviewsFetchedV1` - "reviews.fetched.v1"
- `ReviewsProcessedV1` - "reviews.processed.v1"
- `ReviewsValidatedV1` - "reviews.validated.v1"
- `ReviewsFetchFailedV1` - "reviews.fetch_failed.v1"
- `ReviewsStoredV1` - "reviews.stored.v1"
- `ReviewsErrorV1` - "reviews.error.v1"
- `CatalogProductEnrichmentRequestedV1` - "catalog.product.enrichment.requested.v1"
- `CatalogProductEnrichmentCompletedV1` - "catalog.product.enrichment.completed.v1"
- `CatalogProductEnrichmentFailedV1` - "catalog.product.enrichment.failed.v1"
- `EnrichmentCompletedV1` - "enrichment.completed.v1"
- `EnrichmentFailedV1` - "enrichment.failed.v1"
- `VariantsEnrichedV1` - "variants.enriched.v1"
- `EnrichmentRetryV1` - "enrichment.retry.v1"
- `QualityAssessmentRequestedV1` - "quality.assessment.requested.v1"
- `QualityAssessmentCompletedV1` - "quality.assessment.completed.v1"
- `QualityAssessmentFailedV1` - "quality.assessment.failed.v1"
- `ProductReadyForPublicationV1` - "product.ready_for_publication.v1"
- `ProductUpdatedV1` - "product.updated.v1"
- `ProductUpdateFailedV1` - "product.update_failed.v1"
- `ProductAvailabilityChangedV1` - "product.availability.changed.v1"
- `ProductStatusChangedV1` - "product.status.changed.v1"
- `ProductDeletedV1` - "product.deleted.v1"
- `PriceUpdatedV1` - "price.updated.v1"
- `PriceUpdateFailedV1` - "price.update_failed.v1"
- `PriceMonitoringScheduledV1` - "price.monitoring.scheduled.v1"
- `AvailabilityCheckScheduledV1` - "product.availability_check.scheduled.v1"
- `PeriodicUpdateScheduledV1` - "product.periodic_update.scheduled.v1"

#### Reviews Constants
- `REVIEWS_SOURCE_AMAZON_API` - "amazon_api"
- `REVIEWS_SOURCE_MANUAL` - "manual"
- `REVIEWS_ERROR_TYPE_FETCH` - "fetch_error"
- `REVIEWS_ERROR_TYPE_PROCESS` - "process_error"
- `REVIEWS_ERROR_TYPE_VALIDATION` - "validation_error"
- `REVIEWS_ERROR_TYPE_CACHE` - "cache_error"
- `REVIEWS_ERROR_TYPE_DATABASE` - "database_error"

### Types

#### Core Event Types
- `Event` - Main event structure with ID, Type, AggregateType, AggregateID, Payload, Timestamp, Metadata

#### Payload Types
- `ScraperJobRequestedPayload` - Scraper job request payload
- `ProductCreatedPayload` - Product creation payload
- `ProductUpdatedPayload` - Product update payload
- `ContentGenerationRequestedPayload` - Content generation request payload
- `ContentGeneratedPayload` - Content generation completion payload
- `ContentGenerationFailedPayload` - Content generation failure payload
- `ReviewsCollectedPayload` - Reviews collection payload
- `NewProductDetectedPayload` - New product detection payload
- `ProductIgnoredPayload` - Product ignored payload
- `ProductReviewRequiredPayload` - Product review required payload
- `QualityAssessmentRequestedPayload` - Quality assessment request payload
- `QualityAssessmentCompletedPayload` - Quality assessment completion payload
- `QualityAssessmentFailedPayload` - Quality assessment failure payload
- `ColorEnrichmentRequestedPayload` - Color enrichment request payload
- `ColorEnrichmentCompletedPayload` - Color enrichment completion payload
- `ColorEnrichmentFailedPayload` - Color enrichment failure payload
- `VariationEnrichmentRequestedPayload` - Variation enrichment request payload
- `VariationEnrichmentCompletedPayload` - Variation enrichment completion payload
- `VariationEnrichmentFailedPayload` - Variation enrichment failure payload
- `ProductEnrichmentRequestedData` - PA-API enrichment request data
- `ProductEnrichedData` - PA-API enrichment success data
- `ProductEnrichmentFailedData` - PA-API enrichment failure data
- `PriceMonitoringScheduledPayload` - Price monitoring schedule payload
- `AvailabilityCheckScheduledPayload` - Availability check schedule payload
- `PeriodicUpdateScheduledPayload` - Periodic update schedule payload
- `ReviewsRequestedPayload` - Reviews request payload
- `ReviewsErrorPayload` - Reviews error payload

#### Helper Types
- `ColorVariant` - Product color variation
- `ImageSet` - Product images at different sizes

### Methods

#### Event Methods
- `NewEvent(eventType, aggregateType, aggregateID string, payload any) (*Event, error)` - Creates a new event
- `(e *Event) UnmarshalPayload(v any) error` - Unmarshals event payload
- `NewScraperJobRequestedEvent(jobID, searchQuery, category string, maxPages int) *Event` - Creates scraper job requested event
- `NewReviewsRequestedEvent(asin, productID, source string, options map[string]any) *Event` - Creates reviews requested event
- `NewReviewsErrorEvent(asin, productID, source, errorType, errorMessage string, retryCount int) *Event` - Creates reviews error event

#### ProductEnrichmentRequestedData Methods
- `(p *ProductEnrichmentRequestedData) Validate() error` - Validates enrichment request data

#### Helper Functions
- `IsReviewsEvent(eventType string) bool` - Checks if event type is reviews-related
- `GetReviewsEventPriority(eventType string) int` - Gets reviews event priority
- `NewQualityAssessmentRequestedEvent(asin, productID string, productData map[string]interface{}) *Event` - Creates quality assessment requested event
- `NewQualityAssessmentCompletedEvent(asin, productID string, score float64, status, reason string) *Event` - Creates quality assessment completed event
- `NewQualityAssessmentFailedEvent(asin, productID, reason string) *Event` - Creates quality assessment failed event
- `NewColorEnrichmentRequestedEvent(asin, productID, title string) *Event` - Creates color enrichment requested event
- `NewColorEnrichmentCompletedEvent(asin, productID string, colorVariations []map[string]interface{}, parentASIN string) *Event` - Creates color enrichment completed event
- `NewColorEnrichmentFailedEvent(asin, productID, reason string) *Event` - Creates color enrichment failed event
- `NewPriceMonitoringScheduledEvent(asin, productID string, nextCheckAt time.Time) *Event` - Creates price monitoring scheduled event
- `NewAvailabilityCheckScheduledEvent(asin, productID string, nextCheckAt time.Time) *Event` - Creates availability check scheduled event
- `NewPeriodicUpdateScheduledEvent(asin, productID, updateType string, nextCheckAt time.Time) *Event` - Creates periodic update scheduled event
- `NewProductIgnoredEvent(asin, reason string) *Event` - Creates product ignored event
- `NewProductEnrichmentRequestedEvent(source string, data *ProductEnrichmentRequestedData) (*Event, error)` - Creates PA-API enrichment request event
- `NewProductEnrichedEvent(source string, data *ProductEnrichedData) (*Event, error)` - Creates PA-API enrichment success event
- `NormalizeEventType(s string) (string, bool)` - Maps legacy event strings to canonical constants
- `NewProductEnrichmentFailedEvent(source string, data *ProductEnrichmentFailedData) (*Event, error)` - Creates PA-API enrichment failure event
- `NewCatalogProductEnrichmentRequestedEvent(source string, data *ProductEnrichmentRequestedData) (*Event, error)` - Creates catalog PA-API enrichment request event
- `NewCatalogProductEnrichmentCompletedEvent(source string, data *ProductEnrichedData) (*Event, error)` - Creates catalog PA-API enrichment success event
- `NewCatalogProductEnrichmentFailedEvent(source string, data *ProductEnrichmentFailedData) (*Event, error)` - Creates catalog PA-API enrichment failure event

## pkg/interfaces

### Types
- `EventHandler` - Function signature for handling events: `func(ctx context.Context, event *events.Event) error`
- `StreamProducer` - Interface for publishing events to streams
- `StreamConsumer` - Interface for consuming events from streams
- `StreamClient` - Interface for clients that can both produce and consume events

### Interface Methods

#### StreamProducer
- `PublishEvent(ctx context.Context, streamName string, event *events.Event) error`

#### StreamConsumer
- `ConsumeStream(ctx context.Context, streamName string, groupName string, batchSize int64, pollInterval time.Duration, handler func(context.Context, *events.Event, string) error) error`

#### StreamClient
- Inherits all methods from StreamProducer and StreamConsumer

## pkg/adapters

### Types

#### Redis Adapter Types
- `RedisProducerAdapter` - Adapts Redis producer to StreamProducer interface
- `RedisProducer` - Interface for Redis producer implementations
- `RedisConsumerAdapter` - Adapts Redis consumer to StreamConsumer interface
- `RedisConsumer` - Interface for Redis consumer implementations

#### Service Adapter Types
- `ServiceEventAdapter` - Unified adapter for all services

### Methods

#### RedisProducerAdapter
- `NewRedisProducerAdapter(producer RedisProducer) *RedisProducerAdapter` - Creates new Redis producer adapter
- `(p *RedisProducerAdapter) PublishEvent(ctx context.Context, streamName string, event *events.Event) error` - Publishes event

#### RedisConsumerAdapter
- `NewRedisConsumerAdapter(consumer RedisConsumer) *RedisConsumerAdapter` - Creates new Redis consumer adapter
- `(c *RedisConsumerAdapter) ConsumeStream(ctx context.Context, streamName string, groupName string, batchSize int64, pollInterval time.Duration, handler func(context.Context, *events.Event, string) error) error` - Consumes events

#### ServiceEventAdapter
- `NewServiceEventAdapter(producer interfaces.StreamProducer, consumer interfaces.StreamConsumer) *ServiceEventAdapter` - Creates unified service adapter
- `(s *ServiceEventAdapter) PublishEvent(ctx context.Context, streamName string, event *events.Event) error` - Publishes event
- `(s *ServiceEventAdapter) ConsumeStream(ctx context.Context, streamName string, groupName string, batchSize int64, pollInterval time.Duration, handler func(context.Context, *events.Event, string) error) error` - Consumes events
- `(s *ServiceEventAdapter) PublishProductEvent(ctx context.Context, eventType, productID, asin string, payload any) error` - Publishes product event
- `(s *ServiceEventAdapter) PublishContentEvent(ctx context.Context, eventType, productID string, payload any) error` - Publishes content event

#### Helper Functions
- `DetermineTargetStream(eventType string) string` - Determines target stream based on event type
- `isProductLifecycleEvent(eventType string) bool` - Checks if event is product lifecycle-related
- `isContentGenerationEvent(eventType string) bool` - Checks if event is content generation-related
- `isPriceTrackingEvent(eventType string) bool` - Checks if event is price tracking-related

## pkg/models

### Types

#### Core Model Types
- `Product` - Main Amazon product structure
- `ProductImages` - Product images structure
- `ProductImage` - Single product image structure
- `VariationAttribute` - Product variation attribute structure

#### Type Definitions
- `ProductStatus` - Product status type (string)
- `ContentStatus` - Content status type (string)

### Constants

#### ProductStatus Constants
- `PRODUCT_STATUS_PENDING` - "pending"
- `PRODUCT_STATUS_ACTIVE` - "active"
- `PRODUCT_STATUS_UNAVAILABLE` - "unavailable" (deprecated)
- `PRODUCT_STATUS_TERMINATED` - "terminated" (deprecated)
- `PRODUCT_STATUS_UNSUPPORTED_CATEGORY` - "unsupported_category"
- `PRODUCT_STATUS_PROMOTED` - "promoted"
- `PRODUCT_STATUS_INACTIVE` - "inactive"
- `PRODUCT_STATUS_DRAFT` - "draft"
- `PRODUCT_STATUS_DELETED` - "deleted"
- `PRODUCT_STATUS_OUT_OF_STOCK` - "out_of_stock"
- `PRODUCT_STATUS_DISCONTINUED` - "discontinued"

#### ContentStatus Constants
- `CONTENT_STATUS_PENDING` - "pending"
- `CONTENT_STATUS_REQUESTED` - "requested"
- `CONTENT_STATUS_COMPLETE` - "completed"
- `CONTENT_STATUS_FAILED` - "failed"

#### Gender Constants
- `GENDER_MALE` - "male"
- `GENDER_FEMALE` - "female"
- `GENDER_UNISEX` - "unisex"