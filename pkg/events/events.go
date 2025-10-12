// pkg/events/events.go
package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Event types with numbered convention (Phase_Option_Name)
// Numbers indicate sequence, Letters indicate path (A=Success, B=Failure, C=Manual, D=Retry)
const (
	// Phase 0: Scraper Intake (00)
	EVENT_00A_SCRAPER_JOB_REQUESTED = "00A_SCRAPER_JOB_REQUESTED"

	// Phase 1: Product Discovery (01-02)
	EVENT_01_PRODUCT_DETECTED         = "01_PRODUCT_DETECTED"
	EVENT_02A_PRODUCT_VALIDATED       = "02A_PRODUCT_VALIDATED"
	EVENT_02B_PRODUCT_IGNORED         = "02B_PRODUCT_IGNORED"
	EVENT_02C_PRODUCT_REVIEW_REQUIRED = "02C_PRODUCT_REVIEW_REQUIRED"

	// Phase 2: Enrichment Orchestration (03-05)
	EVENT_03_ENRICHMENT_ORCHESTRATION_STARTED = "03_ENRICHMENT_ORCHESTRATION_STARTED"
	// DEPRECATED: Event_04A_DimensionEnrichmentRequested removed - use CatalogProductEnrichmentRequestedV1 instead
	EVENT_04B_COLOR_ENRICHMENT_REQUESTED = "04B_COLOR_ENRICHMENT_REQUESTED"
	// DEPRECATED: Event_04C_BrowseNodeRequested removed - use CatalogProductEnrichmentRequestedV1 instead
	EVENT_04D_VARIANTS_ENRICHMENT_REQUESTED = "04D_VARIANTS_ENRICHMENT_REQUESTED"
	EVENT_05A_ENRICHMENT_COMPLETED          = ENRICHMENT_COMPLETED_V1
	EVENT_05B_ENRICHMENT_FAILED             = ENRICHMENT_FAILED_V1
	EVENT_05C_VARIANTS_ENRICHED             = VARIANTS_ENRICHED_V1
	EVENT_05D_ENRICHMENT_RETRY              = ENRICHMENT_RETRY_V1

	// Phase 3: Quality Assessment (06-07)
	EVENT_06_QUALITY_ASSESSMENT_REQUESTED  = QUALITY_ASSESSMENT_REQUESTED_V1
	EVENT_07A_QUALITY_ASSESSMENT_COMPLETED = QUALITY_ASSESSMENT_COMPLETED_V1
	EVENT_07B_QUALITY_ASSESSMENT_FAILED    = QUALITY_ASSESSMENT_FAILED_V1

	// Phase 4: Content & Reviews Generation (08-12)
	EVENT_08A_CONTENT_GENERATION_REQUESTED = "08A_CONTENT_GENERATION_REQUESTED"
	EVENT_08B_REVIEWS_REQUESTED            = REVIEWS_REQUESTED_V1
	EVENT_09A_CONTENT_GENERATION_STARTED   = "09A_CONTENT_GENERATION_STARTED"
	EVENT_09B_REVIEWS_FETCHED              = REVIEWS_FETCHED_V1
	EVENT_10A_CONTENT_GENERATED            = "10A_CONTENT_GENERATED"
	EVENT_10B_CONTENT_GENERATION_FAILED    = "10B_CONTENT_GENERATION_FAILED"
	EVENT_10C_REVIEWS_PROCESSED            = REVIEWS_PROCESSED_V1
	EVENT_10D_CONTENT_GENERATION_RETRIED   = CONTENT_GENERATION_RETRIED_V1
	EVENT_11A_REVIEWS_VALIDATED            = REVIEWS_VALIDATED_V1
	EVENT_11B_REVIEWS_FETCH_FAILED         = REVIEWS_FETCH_FAILED_V1
	EVENT_12A_REVIEWS_STORED               = REVIEWS_STORED_V1
	EVENT_12B_REVIEWS_ERROR                = REVIEWS_ERROR_V1

	// Phase 5: Publication & Monitoring (13-14)
	EVENT_13_PRODUCT_READY_FOR_PUBLICATION = PRODUCT_READY_FOR_PUBLICATION_V1
	EVENT_14A_PRICE_MONITORING_SCHEDULED   = PRICE_MONITORING_SCHEDULED_V1
	EVENT_14B_AVAILABILITY_CHECK_SCHEDULED = AVAILABILITY_CHECK_SCHEDULED_V1
	EVENT_14C_PERIODIC_UPDATE_SCHEDULED    = PERIODIC_UPDATE_SCHEDULED_V1

	// Additional Events (15+)
	EVENT_15A_PRICE_UPDATED               = PRICE_UPDATED_V1
	EVENT_15B_PRICE_UPDATE_FAILED         = PRICE_UPDATE_FAILED_V1
	EVENT_16A_PRODUCT_UPDATED             = PRODUCT_UPDATED_V1
	EVENT_16B_PRODUCT_UPDATE_FAILED       = PRODUCT_UPDATE_FAILED_V1
	EVENT_17_PRODUCT_AVAILABILITY_CHANGED = "17_PRODUCT_AVAILABILITY_CHANGED"
	EVENT_18_PRODUCT_STATUS_CHANGED       = PRODUCT_STATUS_CHANGED_V1
	EVENT_19_PRODUCT_DELETED              = PRODUCT_DELETED_V1
)

// Deprecated numbered event aliases maintained for backward compatibility.
// These will be removed in a future release once downstream services migrate to CAPS constants.
const (
	Event_00A_ScraperJobRequested           = EVENT_00A_SCRAPER_JOB_REQUESTED
	Event_01_ProductDetected                = EVENT_01_PRODUCT_DETECTED
	Event_02A_ProductValidated              = EVENT_02A_PRODUCT_VALIDATED
	Event_02B_ProductIgnored                = EVENT_02B_PRODUCT_IGNORED
	Event_02C_ProductReviewRequired         = EVENT_02C_PRODUCT_REVIEW_REQUIRED
	Event_03_EnrichmentOrchestrationStarted = EVENT_03_ENRICHMENT_ORCHESTRATION_STARTED
	Event_04B_ColorEnrichmentRequested      = EVENT_04B_COLOR_ENRICHMENT_REQUESTED
	Event_04D_VariantsEnrichmentRequested   = EVENT_04D_VARIANTS_ENRICHMENT_REQUESTED
	Event_05A_EnrichmentCompleted           = EVENT_05A_ENRICHMENT_COMPLETED
	Event_05B_EnrichmentFailed              = EVENT_05B_ENRICHMENT_FAILED
	Event_05C_VariantsEnriched              = EVENT_05C_VARIANTS_ENRICHED
	Event_05D_EnrichmentRetry               = EVENT_05D_ENRICHMENT_RETRY
	Event_06_QualityAssessmentRequested     = EVENT_06_QUALITY_ASSESSMENT_REQUESTED
	Event_07A_QualityAssessmentCompleted    = EVENT_07A_QUALITY_ASSESSMENT_COMPLETED
	Event_07B_QualityAssessmentFailed       = EVENT_07B_QUALITY_ASSESSMENT_FAILED
	Event_08A_ContentGenerationRequested    = EVENT_08A_CONTENT_GENERATION_REQUESTED
	Event_08B_ReviewsRequested              = EVENT_08B_REVIEWS_REQUESTED
	Event_09A_ContentGenerationStarted      = EVENT_09A_CONTENT_GENERATION_STARTED
	Event_09B_ReviewsFetched                = EVENT_09B_REVIEWS_FETCHED
	Event_10A_ContentGenerated              = EVENT_10A_CONTENT_GENERATED
	Event_10B_ContentGenerationFailed       = EVENT_10B_CONTENT_GENERATION_FAILED
	Event_10C_ReviewsProcessed              = EVENT_10C_REVIEWS_PROCESSED
	Event_10D_ContentGenerationRetried      = EVENT_10D_CONTENT_GENERATION_RETRIED
	Event_11A_ReviewsValidated              = EVENT_11A_REVIEWS_VALIDATED
	Event_11B_ReviewsFetchFailed            = EVENT_11B_REVIEWS_FETCH_FAILED
	Event_12A_ReviewsStored                 = EVENT_12A_REVIEWS_STORED
	Event_12B_ReviewsError                  = EVENT_12B_REVIEWS_ERROR
	Event_13_ProductReadyForPublication     = EVENT_13_PRODUCT_READY_FOR_PUBLICATION
	Event_14A_PriceMonitoringScheduled      = EVENT_14A_PRICE_MONITORING_SCHEDULED
	Event_14B_AvailabilityCheckScheduled    = EVENT_14B_AVAILABILITY_CHECK_SCHEDULED
	Event_14C_PeriodicUpdateScheduled       = EVENT_14C_PERIODIC_UPDATE_SCHEDULED
	Event_15A_PriceUpdated                  = EVENT_15A_PRICE_UPDATED
	Event_15B_PriceUpdateFailed             = EVENT_15B_PRICE_UPDATE_FAILED
	Event_16A_ProductUpdated                = EVENT_16A_PRODUCT_UPDATED
	Event_16B_ProductUpdateFailed           = EVENT_16B_PRODUCT_UPDATE_FAILED
	Event_17_ProductAvailabilityChanged     = EVENT_17_PRODUCT_AVAILABILITY_CHANGED
	Event_18_ProductStatusChanged           = EVENT_18_PRODUCT_STATUS_CHANGED
	Event_19_ProductDeleted                 = EVENT_19_PRODUCT_DELETED
)

// Legacy event names (deprecated - for backward compatibility)
// DEPRECATED: Use numbered constants instead of plain strings:
// - "PRODUCT_CREATED" → use Event_02A_ProductValidated or Event_08A_ContentGenerationRequested
// - "PRODUCT_IGNORED" → use Event_02B_ProductIgnored
// - "CONTENT_GENERATION_REQUESTED" → use Event_08A_ContentGenerationRequested
// - "04D_VARIANTS_ENRICHMENT_REQUESTED" → use VariationEnrichmentRequested (04B_VARIATION_ENRICHMENT_REQUESTED)
const (
	EventTypeScraperJobRequested        = Event_00A_ScraperJobRequested
	EventTypeNewProductDetected         = Event_01_ProductDetected
	EventTypeProductValidated           = Event_02A_ProductValidated
	EventTypeProductIgnored             = Event_02B_ProductIgnored
	EventTypeProductReviewRequired      = Event_02C_ProductReviewRequired
	EventTypeContentGenerationRequested = Event_08A_ContentGenerationRequested
	EventTypeContentGenerationStarted   = Event_09A_ContentGenerationStarted
	EventTypeContentGenerated           = Event_10A_ContentGenerated
	EventTypeContentGenerationFailed    = Event_10B_ContentGenerationFailed
	EventTypeContentGenerationRetried   = Event_10D_ContentGenerationRetried
	EventTypeReviewsRequested           = Event_08B_ReviewsRequested
	EventTypeReviewsFetched             = Event_09B_ReviewsFetched
	EventTypeReviewsProcessed           = Event_10C_ReviewsProcessed
	EventTypeReviewsValidated           = Event_11A_ReviewsValidated
	EventTypeReviewsFetchFailed         = Event_11B_ReviewsFetchFailed
	EventTypeReviewsStored              = Event_12A_ReviewsStored
	EventTypeReviewsError               = Event_12B_ReviewsError
	// DEPRECATED: Browse Node events removed - use CatalogProductEnrichment events instead
	EventTypeCheckPrice                 = "CHECK_PRICE"
	EventTypePriceUpdated               = Event_15A_PriceUpdated
	EventTypePriceUpdateFailed          = Event_15B_PriceUpdateFailed
	EventTypeProductUpdated             = Event_16A_ProductUpdated
	EventTypeProductAvailabilityChanged = Event_17_ProductAvailabilityChanged
	EventTypeProductStatusChanged       = Event_18_ProductStatusChanged
	EventTypeProductDeleted             = Event_19_ProductDeleted

	// Deprecated - to be removed
	EventTypeProductUnavailable     = "PRODUCT_UNAVAILABLE"
	EventTypeProductCreated         = "PRODUCT_CREATED"
	EventTypeProductUpdateRequested = "PRODUCT_UPDATE_REQUESTED"
	EventTypeContentUpdateRequested = "CONTENT_UPDATE_REQUESTED"
	EventTypeContentUpdated         = "CONTENT_UPDATED"
	EventTypeContentAnalysisFailed  = "CONTENT_ANALYSIS_FAILED"
	EventTypeReviewsCollected       = "REVIEWS_COLLECTED"
	EventTypeReviewsEnriched        = "REVIEWS_ENRICHED"
	EventTypeReviewsCached          = "REVIEWS_CACHED"
	EventTypeReviewsExpired         = "REVIEWS_EXPIRED"
	EventTypeReviewsDeleted         = "REVIEWS_DELETED"
	// DEPRECATED: Browse node detection removed - handled internally by services
)

// Legacy orchestration event names (mapped to new convention)
const (
	// DEPRECATED: Dimension Enrichment events removed - use CatalogProductEnrichment events instead

	// Quality Assessment
	QualityAssessmentRequested = Event_06_QualityAssessmentRequested
	QualityAssessmentCompleted = Event_07A_QualityAssessmentCompleted
	QualityAssessmentFailed    = Event_07B_QualityAssessmentFailed

	// Scheduling
	PriceMonitoringScheduled   = Event_14A_PriceMonitoringScheduled
	AvailabilityCheckScheduled = Event_14B_AvailabilityCheckScheduled
	PeriodicUpdateScheduled    = Event_14C_PeriodicUpdateScheduled

	// PA-API Color Enrichment
	ColorEnrichmentRequested     = Event_04B_ColorEnrichmentRequested
	ColorEnrichmentCompleted     = "05A_COLOR_ENRICHMENT_COMPLETED"
	ColorEnrichmentFailed        = "05B_COLOR_ENRICHMENT_FAILED"
	VariationEnrichmentRequested = "04B_VARIATION_ENRICHMENT_REQUESTED"
	VariationEnrichmentCompleted = "05A_VARIATION_ENRICHMENT_COMPLETED"
	VariationEnrichmentFailed    = "05B_VARIATION_ENRICHMENT_FAILED"

	// PA-API Variant Enrichment (single call for both color and image variants)
	VariantsEnrichmentRequested = Event_04D_VariantsEnrichmentRequested
	VariantsEnrichmentCompleted = Event_05C_VariantsEnriched
	VariantsEnrichmentFailed    = "05B_VARIANTS_ENRICHMENT_FAILED"

	// PA-API Enrichment Event Types (CloudEvent format)
	// DEPRECATED: Use catalog.* variants instead for consistency
	ProductEnrichmentRequestedV1  = "product.enrichment.requested.v1"
	ProductEnrichmentCompletedV1  = "product.enrichment.completed.v1"
	ProductEnrichmentFailedV1     = "product.enrichment.failed.v1"
	VariantsEnrichmentRequestedV1 = "product.variants.enrichment.requested.v1"
	VariantsEnrichmentCompletedV1 = "product.variants.enrichment.completed.v1"
	VariantsEnrichmentFailedV1    = "product.variants.enrichment.failed.v1"
)

// CAPS CloudEvent constants for new implementations
// These follow the CAPS naming convention while maintaining CloudEvent semantics
const (
	// Product Discovery Events
	CATALOG_PRODUCT_DETECTED_V1        = "catalog.product.detected.v1"
	CATALOG_PRODUCT_VALIDATED_V1       = "catalog.product.validated.v1"
	CATALOG_PRODUCT_IGNORED_V1         = "catalog.product.ignored.v1"
	CATALOG_PRODUCT_REVIEW_REQUIRED_V1 = "catalog.product.review_required.v1"

	// Content Generation Events
	CONTENT_GENERATION_REQUESTED_V1 = "content.generation.requested.v1"
	CONTENT_GENERATION_STARTED_V1   = "content.generation.started.v1"
	CONTENT_GENERATED_V1            = "content.generated.v1"
	CONTENT_GENERATION_FAILED_V1    = "content.generation.failed.v1"
	CONTENT_GENERATION_RETRIED_V1   = "content.generation.retried.v1"

	// Reviews Events
	REVIEWS_REQUESTED_V1    = "reviews.requested.v1"
	REVIEWS_FETCHED_V1      = "reviews.fetched.v1"
	REVIEWS_PROCESSED_V1    = "reviews.processed.v1"
	REVIEWS_VALIDATED_V1    = "reviews.validated.v1"
	REVIEWS_FETCH_FAILED_V1 = "reviews.fetch_failed.v1"
	REVIEWS_STORED_V1       = "reviews.stored.v1"
	REVIEWS_ERROR_V1        = "reviews.error.v1"

	// Product Enrichment Events (replaces 04A/04C flows)
	CATALOG_PRODUCT_ENRICHMENT_REQUESTED_V1 = "catalog.product.enrichment.requested.v1"
	CATALOG_PRODUCT_ENRICHMENT_COMPLETED_V1 = "catalog.product.enrichment.completed.v1"
	CATALOG_PRODUCT_ENRICHMENT_FAILED_V1    = "catalog.product.enrichment.failed.v1"

	// Enrichment Events (Legacy 05 series)
	ENRICHMENT_COMPLETED_V1 = "enrichment.completed.v1"
	ENRICHMENT_FAILED_V1    = "enrichment.failed.v1"
	VARIANTS_ENRICHED_V1    = "variants.enriched.v1"
	ENRICHMENT_RETRY_V1     = "enrichment.retry.v1"

	// Quality Assessment Events
	QUALITY_ASSESSMENT_REQUESTED_V1 = "quality.assessment.requested.v1"
	QUALITY_ASSESSMENT_COMPLETED_V1 = "quality.assessment.completed.v1"
	QUALITY_ASSESSMENT_FAILED_V1    = "quality.assessment.failed.v1"

	// Product Lifecycle Events
	PRODUCT_READY_FOR_PUBLICATION_V1 = "product.ready_for_publication.v1"
	PRODUCT_UPDATED_V1               = "product.updated.v1"
	PRODUCT_UPDATE_FAILED_V1         = "product.update_failed.v1"
	PRODUCT_AVAILABILITY_CHANGED_V1  = "product.availability.changed.v1"
	PRODUCT_STATUS_CHANGED_V1        = "product.status.changed.v1"
	PRODUCT_DELETED_V1               = "product.deleted.v1"

	// Price Monitoring Events
	PRICE_UPDATED_V1                = "price.updated.v1"
	PRICE_UPDATE_FAILED_V1          = "price.update_failed.v1"
	PRICE_MONITORING_SCHEDULED_V1   = "price.monitoring.scheduled.v1"
	AVAILABILITY_CHECK_SCHEDULED_V1 = "product.availability_check.scheduled.v1"
	PERIODIC_UPDATE_SCHEDULED_V1    = "product.periodic_update.scheduled.v1"

	// PA-API Enrichment Events (Legacy)
	PRODUCT_ENRICHMENT_REQUESTED_V1          = "product.enrichment.requested.v1"
	PRODUCT_ENRICHMENT_COMPLETED_V1          = "product.enrichment.completed.v1"
	PRODUCT_ENRICHMENT_FAILED_V1             = "product.enrichment.failed.v1"
	PRODUCT_VARIANTS_ENRICHMENT_REQUESTED_V1 = "product.variants.enrichment.requested.v1"
	PRODUCT_VARIANTS_ENRICHMENT_COMPLETED_V1 = "product.variants.enrichment.completed.v1"
	PRODUCT_VARIANTS_ENRICHMENT_FAILED_V1    = "product.variants.enrichment.failed.v1"
)

// Reviews-related constants (CAPS naming)
const (
	REVIEWS_SOURCE_AMAZON_API = "amazon_api"
	REVIEWS_SOURCE_MANUAL     = "manual"
)

const (
	REVIEWS_ERROR_TYPE_FETCH      = "fetch_error"
	REVIEWS_ERROR_TYPE_PROCESS    = "process_error"
	REVIEWS_ERROR_TYPE_VALIDATION = "validation_error"
	REVIEWS_ERROR_TYPE_CACHE      = "cache_error"
	REVIEWS_ERROR_TYPE_DATABASE   = "database_error"
)

// Canonical CloudEvents types (domain.subdomain.action.v1 format)
// DEPRECATED: Use CAPS constants instead. Will be removed in a future release.
const (
	// Product Discovery Events
	CatalogProductDetectedV1       = CATALOG_PRODUCT_DETECTED_V1
	CatalogProductValidatedV1      = CATALOG_PRODUCT_VALIDATED_V1
	CatalogProductIgnoredV1        = CATALOG_PRODUCT_IGNORED_V1
	CatalogProductReviewRequiredV1 = CATALOG_PRODUCT_REVIEW_REQUIRED_V1

	// Content Generation Events
	ContentGenerationRequestedV1 = CONTENT_GENERATION_REQUESTED_V1
	ContentGenerationStartedV1   = CONTENT_GENERATION_STARTED_V1
	ContentGeneratedV1           = CONTENT_GENERATED_V1
	ContentGenerationFailedV1    = CONTENT_GENERATION_FAILED_V1
	ContentGenerationRetriedV1   = CONTENT_GENERATION_RETRIED_V1

	// Reviews Events
	ReviewsRequestedV1   = REVIEWS_REQUESTED_V1
	ReviewsFetchedV1     = REVIEWS_FETCHED_V1
	ReviewsProcessedV1   = REVIEWS_PROCESSED_V1
	ReviewsValidatedV1   = REVIEWS_VALIDATED_V1
	ReviewsFetchFailedV1 = REVIEWS_FETCH_FAILED_V1
	ReviewsStoredV1      = REVIEWS_STORED_V1
	ReviewsErrorV1       = REVIEWS_ERROR_V1

	// Product Enrichment Events (replaces 04A/04C flows)
	CatalogProductEnrichmentRequestedV1 = CATALOG_PRODUCT_ENRICHMENT_REQUESTED_V1
	CatalogProductEnrichmentCompletedV1 = CATALOG_PRODUCT_ENRICHMENT_COMPLETED_V1
	CatalogProductEnrichmentFailedV1    = CATALOG_PRODUCT_ENRICHMENT_FAILED_V1

	// Enrichment Events (Legacy 05 series)
	EnrichmentCompletedV1 = ENRICHMENT_COMPLETED_V1
	EnrichmentFailedV1    = ENRICHMENT_FAILED_V1
	VariantsEnrichedV1    = VARIANTS_ENRICHED_V1
	EnrichmentRetryV1     = ENRICHMENT_RETRY_V1

	// Quality Assessment Events
	QualityAssessmentRequestedV1 = QUALITY_ASSESSMENT_REQUESTED_V1
	QualityAssessmentCompletedV1 = QUALITY_ASSESSMENT_COMPLETED_V1
	QualityAssessmentFailedV1    = QUALITY_ASSESSMENT_FAILED_V1

	// Product Lifecycle Events
	ProductReadyForPublicationV1 = PRODUCT_READY_FOR_PUBLICATION_V1
	ProductUpdatedV1             = PRODUCT_UPDATED_V1
	ProductUpdateFailedV1        = PRODUCT_UPDATE_FAILED_V1
	ProductAvailabilityChangedV1 = PRODUCT_AVAILABILITY_CHANGED_V1
	ProductStatusChangedV1       = PRODUCT_STATUS_CHANGED_V1
	ProductDeletedV1             = PRODUCT_DELETED_V1

	// Price Monitoring Events
	PriceUpdatedV1               = PRICE_UPDATED_V1
	PriceUpdateFailedV1          = PRICE_UPDATE_FAILED_V1
	PriceMonitoringScheduledV1   = PRICE_MONITORING_SCHEDULED_V1
	AvailabilityCheckScheduledV1 = AVAILABILITY_CHECK_SCHEDULED_V1
	PeriodicUpdateScheduledV1    = PERIODIC_UPDATE_SCHEDULED_V1
)

// Event represents a domain event
type Event struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	AggregateType string         `json:"aggregate_type"`
	AggregateID   string         `json:"aggregate_id"`
	Payload       any            `json:"payload"`
	Timestamp     time.Time      `json:"timestamp"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

// ScraperJobRequestedPayload represents the payload for an intake scraper job
type ScraperJobRequestedPayload struct {
	JobID       string    `json:"job_id"`
	SearchQuery string    `json:"search_query"`
	Category    string    `json:"category"`
	MaxPages    int       `json:"max_pages"`
	RequestedAt time.Time `json:"requested_at"`
}

// ProductCreatedPayload represents the payload for a PRODUCT_CREATED event
type ProductCreatedPayload struct {
	ASIN           string   `json:"asin"`
	Title          string   `json:"title"`
	Brand          string   `json:"brand,omitempty"`
	Category       string   `json:"category,omitempty"`
	Gender         string   `json:"gender,omitempty"`
	CurrentPrice   float64  `json:"current_price,omitempty"`
	Currency       string   `json:"currency,omitempty"`
	DetailPageURL  string   `json:"detail_page_url,omitempty"`
	ImageUrls      []string `json:"image_urls,omitempty"`
	Features       []string `json:"features,omitempty"`
	BrowseNodeID   string   `json:"browse_node_id,omitempty"`
	BrowseNodeTags []string `json:"browse_node_tags,omitempty"`
}

// ProductUpdatedPayload represents the payload for a PRODUCT_UPDATED event
type ProductUpdatedPayload struct {
	ASIN      string  `json:"asin"`
	Title     string  `json:"title,omitempty"`
	Brand     string  `json:"brand,omitempty"`
	Price     float64 `json:"price,omitempty"`
	Currency  string  `json:"currency,omitempty"`
	Available bool    `json:"available"`
}

// ContentGenerationRequestedPayload represents the payload for a CONTENT_GENERATION_REQUESTED event
type ContentGenerationRequestedPayload struct {
	ASIN        string    `json:"asin"`
	ProductID   string    `json:"product_id"`
	Priority    int       `json:"priority,omitempty"` // 1-5, wobei 5 die höchste Priorität ist
	RequestedAt time.Time `json:"requested_at"`
}

// DEPRECATED: BrowseNode payload types removed
// Use ProductEnrichmentRequestedData, ProductEnrichedData, and ProductEnrichmentFailedData instead

// ContentGeneratedPayload represents the payload for a CONTENT_GENERATED event
type ContentGeneratedPayload struct {
	ASIN           string    `json:"asin"`
	ProductID      string    `json:"product_id"`
	ContentType    string    `json:"content_type"` // z.B. "description", "faq"
	ContentID      string    `json:"content_id,omitempty"`
	GeneratedAt    time.Time `json:"generated_at"`
	ReviewsContent string    `json:"reviews_content,omitempty"`
	FAQContent     string    `json:"faq_content,omitempty"`
	Gender         string    `json:"gender,omitempty"`
}

// ContentGenerationFailedPayload represents the payload for a CONTENT_GENERATION_FAILED event
type ContentGenerationFailedPayload struct {
	ASIN      string    `json:"asin"`
	ProductID string    `json:"product_id"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
}

// ReviewsCollectedPayload represents the payload for a REVIEWS_COLLECTED event
type ReviewsCollectedPayload struct {
	ASIN        string    `json:"asin"`
	ProductID   string    `json:"product_id"`
	ReviewCount int       `json:"review_count"`
	CollectedAt time.Time `json:"collected_at"`
}

// NewProductDetectedPayload represents the payload for a NEW_PRODUCT_DETECTED event
type NewProductDetectedPayload struct {
	ASIN           string                 `json:"asin"`
	Title          string                 `json:"title"`
	Brand          string                 `json:"brand,omitempty"`
	Category       string                 `json:"category,omitempty"`
	Gender         string                 `json:"gender,omitempty"`
	CurrentPrice   float64                `json:"current_price,omitempty"`
	Currency       string                 `json:"currency,omitempty"`
	DetailPageURL  string                 `json:"detail_page_url,omitempty"`
	ImageUrls      []string               `json:"image_urls,omitempty"`
	Features       []string               `json:"features,omitempty"`
	BrowseNodeID   string                 `json:"browse_node_id,omitempty"`
	BrowseNodeTags []string               `json:"browse_node_tags,omitempty"` // Browse node IDs for categorization
	AmazonData     map[string]interface{} `json:"amazon_data,omitempty"`      // Complete Amazon product data

	// ENHANCED FILTER FIELDS FOR TALL PEOPLE (CRITICAL for validation)
	// Size & Color Information
	Size           string   `json:"size,omitempty"`            // "XL", "XXL", "Tall", "Long"
	Color          string   `json:"color,omitempty"`           // Color variants
	AvailableSizes []string `json:"available_sizes,omitempty"` // Available sizes

	// Item Dimensions (CRITICAL for tall people)
	HeightCm *float64 `json:"height_cm,omitempty"` // Physical height in cm
	LengthCm *float64 `json:"length_cm,omitempty"` // Length in cm (critical for pants)
	WidthCm  *float64 `json:"width_cm,omitempty"`  // Width in cm
	WeightG  *float64 `json:"weight_g,omitempty"`  // Weight in grams

	// Product Classification
	ProductGroup   string `json:"product_group,omitempty"`    // Amazon ProductGroup
	Model          string `json:"model,omitempty"`            // Product model
	IsAdultProduct *bool  `json:"is_adult_product,omitempty"` // Adult product flag

	// Variation Attributes
	VariationAttributes []map[string]string `json:"variation_attributes,omitempty"` // Size, Color variations

	// Tall-Friendly Scoring
	TallFriendlyScore *float64 `json:"tall_friendly_score,omitempty"` // Pre-calculated score
	IsTallFriendly    *bool    `json:"is_tall_friendly,omitempty"`    // Pre-calculated flag

	// Quality Assessment Fields (CRITICAL for validation)
	Rating      *float64 `json:"rating,omitempty"`       // Product rating (1-5 stars)
	ReviewCount *int64   `json:"review_count,omitempty"` // Number of reviews

	// Dimension Data from Amazon Scraper
	Dimensions    map[string]interface{} `json:"dimensions,omitempty"`     // Size chart dimensions from scraper
	HasDimensions bool                   `json:"has_dimensions,omitempty"` // Whether dimensions were found
}

// ProductIgnoredPayload represents the payload for a PRODUCT_IGNORED event
type ProductIgnoredPayload struct {
	ASIN          string    `json:"asin"`
	ProductID     string    `json:"product_id"`
	Reason        string    `json:"reason"` // e.g., "no_size_information", "no_reviews", "no_images", "price_out_of_range"
	MissingFields []string  `json:"missing_fields,omitempty"`
	IgnoredAt     time.Time `json:"ignored_at"`
}

// ProductReviewRequiredPayload represents the payload for a PRODUCT_REVIEW_REQUIRED event
type ProductReviewRequiredPayload struct {
	ASIN            string    `json:"asin"`
	ProductID       string    `json:"product_id"`
	Reason          string    `json:"reason"`           // e.g., "borderline_quality", "size_information_unclear", "low_rating"
	SuggestedAction string    `json:"suggested_action"` // e.g., "manual_size_check", "quality_review"
	Score           float64   `json:"score,omitempty"`  // Quality score if applicable
	RequiredAt      time.Time `json:"required_at"`
}

// DEPRECATED: DimensionEnrichment payload types removed
// Use ProductEnrichmentRequestedData, ProductEnrichedData, and ProductEnrichmentFailedData instead

// QualityAssessmentRequestedPayload represents the payload for quality assessment request
type QualityAssessmentRequestedPayload struct {
	ASIN        string                 `json:"asin"`
	ProductID   string                 `json:"product_id"`
	ProductData map[string]interface{} `json:"product_data"` // All product data including dimensions
	RequestedAt time.Time              `json:"requested_at"`
}

// QualityAssessmentCompletedPayload represents the payload for completed quality assessment
type QualityAssessmentCompletedPayload struct {
	ASIN            string    `json:"asin"`
	ProductID       string    `json:"product_id"`
	QualityScore    float64   `json:"quality_score"`
	Status          string    `json:"status"` // "validated", "ignored", "review_required"
	Reason          string    `json:"reason,omitempty"`
	MissingFields   []string  `json:"missing_fields,omitempty"`
	SuggestedAction string    `json:"suggested_action,omitempty"`
	CompletedAt     time.Time `json:"completed_at"`
}

// QualityAssessmentFailedPayload represents the payload for failed quality assessment
type QualityAssessmentFailedPayload struct {
	ASIN      string    `json:"asin"`
	ProductID string    `json:"product_id"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
}

// ColorEnrichmentRequestedPayload represents the payload for PA-API color enrichment request
type ColorEnrichmentRequestedPayload struct {
	ASIN        string    `json:"asin"`
	ProductID   string    `json:"product_id"`
	Title       string    `json:"title,omitempty"`
	RequestedAt time.Time `json:"requested_at"`
}

// ColorEnrichmentCompletedPayload represents the payload for successful color enrichment
type ColorEnrichmentCompletedPayload struct {
	ASIN             string                   `json:"asin"`
	ProductID        string                   `json:"product_id"`
	ColorVariations  []map[string]interface{} `json:"color_variations"`
	ParentASIN       string                   `json:"parent_asin,omitempty"`
	EnrichmentSource string                   `json:"enrichment_source"` // "pa-api"
	CompletedAt      time.Time                `json:"completed_at"`
}

// ColorEnrichmentFailedPayload represents the payload for failed color enrichment
type ColorEnrichmentFailedPayload struct {
	ASIN      string    `json:"asin"`
	ProductID string    `json:"product_id"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
}

// VariationEnrichmentRequestedPayload for getting all product variations
type VariationEnrichmentRequestedPayload struct {
	ASIN        string    `json:"asin"`
	ParentASIN  string    `json:"parent_asin,omitempty"`
	ProductID   string    `json:"product_id"`
	RequestedAt time.Time `json:"requested_at"`
}

// VariationEnrichmentCompletedPayload for successful variation data
type VariationEnrichmentCompletedPayload struct {
	ASIN        string                   `json:"asin"`
	ProductID   string                   `json:"product_id"`
	Variations  []map[string]interface{} `json:"variations"`
	CompletedAt time.Time                `json:"completed_at"`
}

// VariationEnrichmentFailedPayload for failed variation enrichment
type VariationEnrichmentFailedPayload struct {
	ASIN      string    `json:"asin"`
	ProductID string    `json:"product_id"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
}

// ColorVariant represents a product color variation
type ColorVariant struct {
	ColorName string     `json:"color_name"`
	ASIN      string     `json:"asin"`
	Images    []ImageSet `json:"images"`
}

// ImageSet represents product images at different sizes
type ImageSet struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

// ProductEnrichmentRequestedData represents a PA-API enrichment request
type ProductEnrichmentRequestedData struct {
	ASIN       string `json:"asin"`
	Region     string `json:"region"`
	RequestID  string `json:"request_id"`
	RetryCount int    `json:"retry_count"`
}

func (p *ProductEnrichmentRequestedData) Validate() error {
	if p.ASIN == "" {
		return fmt.Errorf("ASIN is required")
	}
	if p.Region == "" {
		return fmt.Errorf("region is required")
	}
	if p.RequestID == "" {
		return fmt.Errorf("request_id is required")
	}
	return nil
}

// ProductEnrichedData represents successful PA-API enrichment
type ProductEnrichedData struct {
	ASIN          string         `json:"asin"`
	Region        string         `json:"region"`
	RequestID     string         `json:"request_id"`
	ColorVariants []ColorVariant `json:"color_variants"`
	ProcessingMS  int64          `json:"processing_ms"`
	EnrichedAt    time.Time      `json:"enriched_at"`
}

// ProductEnrichmentFailedData represents failed PA-API enrichment
type ProductEnrichmentFailedData struct {
	ASIN         string    `json:"asin"`
	Region       string    `json:"region"`
	RequestID    string    `json:"request_id"`
	ErrorCode    string    `json:"error_code"`
	ErrorMessage string    `json:"error_message"`
	FailedAt     time.Time `json:"failed_at"`
	RetryCount   int       `json:"retry_count"`
}

// PriceMonitoringScheduledPayload represents the payload for scheduled price monitoring
type PriceMonitoringScheduledPayload struct {
	ASIN        string    `json:"asin"`
	ProductID   string    `json:"product_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	NextCheckAt time.Time `json:"next_check_at"`
}

// AvailabilityCheckScheduledPayload represents the payload for scheduled availability check
type AvailabilityCheckScheduledPayload struct {
	ASIN        string    `json:"asin"`
	ProductID   string    `json:"product_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	NextCheckAt time.Time `json:"next_check_at"`
}

// PeriodicUpdateScheduledPayload represents the payload for scheduled periodic updates
type PeriodicUpdateScheduledPayload struct {
	ASIN        string    `json:"asin"`
	ProductID   string    `json:"product_id"`
	UpdateType  string    `json:"update_type"` // "full", "price_only", "availability_only"
	ScheduledAt time.Time `json:"scheduled_at"`
	NextCheckAt time.Time `json:"next_check_at"`
}

// Helper functions for event creation
func NewEvent(eventType, aggregateType, aggregateID string, payload any) (*Event, error) {
	return &Event{
		ID:            uuid.New().String(),
		Type:          eventType,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Payload:       payload,
		Timestamp:     time.Now().UTC(),
		Metadata:      make(map[string]any),
	}, nil
}

// NewScraperJobRequestedEvent creates a canonical scraper job requested event
func NewScraperJobRequestedEvent(jobID, searchQuery, category string, maxPages int) *Event {
	payload := ScraperJobRequestedPayload{
		JobID:       jobID,
		SearchQuery: searchQuery,
		Category:    category,
		MaxPages:    maxPages,
		RequestedAt: time.Now().UTC(),
	}

	event, _ := NewEvent(EventTypeScraperJobRequested, "scraper.job", jobID, payload)
	return event
}

// UnmarshalPayload unmarshals the event payload into the given value
func (e *Event) UnmarshalPayload(v any) error {
	payloadBytes, err := json.Marshal(e.Payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(payloadBytes, v)
}

// Reviews-related constants and types
// DEPRECATED: Use CAPS constants instead. Will be removed in a future release.
const (
	ReviewsSourceAmazonAPI = REVIEWS_SOURCE_AMAZON_API
	ReviewsSourceManual    = REVIEWS_SOURCE_MANUAL
)

const (
	ReviewsErrorTypeFetch      = REVIEWS_ERROR_TYPE_FETCH
	ReviewsErrorTypeProcess    = REVIEWS_ERROR_TYPE_PROCESS
	ReviewsErrorTypeValidation = REVIEWS_ERROR_TYPE_VALIDATION
	ReviewsErrorTypeCache      = REVIEWS_ERROR_TYPE_CACHE
	ReviewsErrorTypeDatabase   = REVIEWS_ERROR_TYPE_DATABASE
)

// Reviews payload structures
type ReviewsRequestedPayload struct {
	ASIN      string         `json:"asin"`
	ProductID string         `json:"product_id"`
	Source    string         `json:"source"`
	Options   map[string]any `json:"options,omitempty"`
}

type ReviewsErrorPayload struct {
	ASIN         string `json:"asin"`
	ProductID    string `json:"product_id"`
	Source       string `json:"source"`
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`
	RetryCount   int    `json:"retry_count"`
}

// Reviews helper functions
func NewReviewsRequestedEvent(asin, productID, source string, options map[string]any) *Event {
	payload := ReviewsRequestedPayload{
		ASIN:      asin,
		ProductID: productID,
		Source:    source,
		Options:   options,
	}

	event, _ := NewEvent(EventTypeReviewsRequested, "reviews", productID, payload)
	return event
}

func NewReviewsErrorEvent(asin, productID, source, errorType, errorMessage string, retryCount int) *Event {
	payload := ReviewsErrorPayload{
		ASIN:         asin,
		ProductID:    productID,
		Source:       source,
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
		RetryCount:   retryCount,
	}

	event, _ := NewEvent(EventTypeReviewsError, "reviews", productID, payload)
	return event
}

func IsReviewsEvent(eventType string) bool {
	reviewsEvents := []string{
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
	}

	for _, re := range reviewsEvents {
		if re == eventType {
			return true
		}
	}
	return false
}

func GetReviewsEventPriority(eventType string) int {
	priorities := map[string]int{
		EventTypeReviewsError:       1, // Highest priority
		EventTypeReviewsRequested:   2,
		EventTypeReviewsFetched:     3,
		EventTypeReviewsStored:      4,
		EventTypeReviewsProcessed:   5,
		EventTypeReviewsValidated:   6,
		EventTypeReviewsEnriched:    7,
		EventTypeReviewsCached:      8,
		EventTypeReviewsExpired:     9,
		EventTypeReviewsDeleted:     10, // Lowest priority
		EventTypeReviewsCollected:   5,
		EventTypeReviewsFetchFailed: 2,
	}

	if priority, exists := priorities[eventType]; exists {
		return priority
	}
	return 999 // Unknown event
}

// Helper functions for orchestration events

// DEPRECATED: Dimension Enrichment constructor functions removed
// Use PA-API enrichment events instead:
// - NewProductEnrichmentRequestedEvent for requesting enrichment
// - NewProductEnrichedEvent for completed enrichment
// - NewProductEnrichmentFailedEvent for failed enrichment
// See MIGRATION_DIMENSION_ENRICHMENT_REMOVAL.md for migration guide

// NewQualityAssessmentRequestedEvent creates a new quality assessment requested event
func NewQualityAssessmentRequestedEvent(asin, productID string, productData map[string]interface{}) *Event {
	payload := QualityAssessmentRequestedPayload{
		ASIN:        asin,
		ProductID:   productID,
		ProductData: productData,
		RequestedAt: time.Now().UTC(),
	}

	event, _ := NewEvent(QualityAssessmentRequested, "product", productID, payload)
	return event
}

// NewQualityAssessmentCompletedEvent creates a new quality assessment completed event
func NewQualityAssessmentCompletedEvent(asin, productID string, score float64, status, reason string) *Event {
	payload := QualityAssessmentCompletedPayload{
		ASIN:         asin,
		ProductID:    productID,
		QualityScore: score,
		Status:       status,
		Reason:       reason,
		CompletedAt:  time.Now().UTC(),
	}

	event, _ := NewEvent(QualityAssessmentCompleted, "product", productID, payload)
	return event
}

// NewQualityAssessmentFailedEvent creates a new quality assessment failed event
func NewQualityAssessmentFailedEvent(asin, productID, reason string) *Event {
	payload := QualityAssessmentFailedPayload{
		ASIN:      asin,
		ProductID: productID,
		Reason:    reason,
		FailedAt:  time.Now().UTC(),
	}

	event, _ := NewEvent(QualityAssessmentFailed, "product", productID, payload)
	return event
}

// NewColorEnrichmentRequestedEvent creates a new color enrichment requested event
func NewColorEnrichmentRequestedEvent(asin, productID, title string) *Event {
	payload := ColorEnrichmentRequestedPayload{
		ASIN:        asin,
		ProductID:   productID,
		Title:       title,
		RequestedAt: time.Now().UTC(),
	}

	event, _ := NewEvent(ColorEnrichmentRequested, "product", productID, payload)
	return event
}

// NewColorEnrichmentCompletedEvent creates a new color enrichment completed event
func NewColorEnrichmentCompletedEvent(asin, productID string, colorVariations []map[string]interface{}, parentASIN string) *Event {
	payload := ColorEnrichmentCompletedPayload{
		ASIN:             asin,
		ProductID:        productID,
		ColorVariations:  colorVariations,
		ParentASIN:       parentASIN,
		EnrichmentSource: "pa-api",
		CompletedAt:      time.Now().UTC(),
	}

	event, _ := NewEvent(ColorEnrichmentCompleted, "product", productID, payload)
	return event
}

// NewColorEnrichmentFailedEvent creates a new color enrichment failed event
func NewColorEnrichmentFailedEvent(asin, productID, reason string) *Event {
	payload := ColorEnrichmentFailedPayload{
		ASIN:      asin,
		ProductID: productID,
		Reason:    reason,
		FailedAt:  time.Now().UTC(),
	}

	event, _ := NewEvent(ColorEnrichmentFailed, "product", productID, payload)
	return event
}

// NewPriceMonitoringScheduledEvent creates a new price monitoring scheduled event
func NewPriceMonitoringScheduledEvent(asin, productID string, nextCheckAt time.Time) *Event {
	payload := PriceMonitoringScheduledPayload{
		ASIN:        asin,
		ProductID:   productID,
		ScheduledAt: time.Now().UTC(),
		NextCheckAt: nextCheckAt,
	}

	event, _ := NewEvent(PriceMonitoringScheduled, "product", productID, payload)
	return event
}

// NewAvailabilityCheckScheduledEvent creates a new availability check scheduled event
func NewAvailabilityCheckScheduledEvent(asin, productID string, nextCheckAt time.Time) *Event {
	payload := AvailabilityCheckScheduledPayload{
		ASIN:        asin,
		ProductID:   productID,
		ScheduledAt: time.Now().UTC(),
		NextCheckAt: nextCheckAt,
	}

	event, _ := NewEvent(AvailabilityCheckScheduled, "product", productID, payload)
	return event
}

// NewPeriodicUpdateScheduledEvent creates a new periodic update scheduled event
func NewPeriodicUpdateScheduledEvent(asin, productID, updateType string, nextCheckAt time.Time) *Event {
	payload := PeriodicUpdateScheduledPayload{
		ASIN:        asin,
		ProductID:   productID,
		UpdateType:  updateType,
		ScheduledAt: time.Now().UTC(),
		NextCheckAt: nextCheckAt,
	}

	event, _ := NewEvent(PeriodicUpdateScheduled, "product", productID, payload)
	return event
}

// NewProductIgnoredEvent creates a new product ignored event
func NewProductIgnoredEvent(asin, reason string) *Event {
	payload := ProductIgnoredPayload{
		ASIN:      asin,
		ProductID: asin, // Use ASIN as product ID for ignored products
		Reason:    reason,
		IgnoredAt: time.Now().UTC(),
	}

	event, _ := NewEvent("PRODUCT_IGNORED", "product", asin, payload)
	return event
}

// NewProductEnrichmentRequestedEvent creates a new PA-API enrichment request event
// DEPRECATED: Use NewCatalogProductEnrichmentRequestedEvent instead
func NewProductEnrichmentRequestedEvent(source string, data *ProductEnrichmentRequestedData) (*Event, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}
	return NewEvent(PRODUCT_ENRICHMENT_REQUESTED_V1, "product", data.ASIN, data)
}

// NewProductEnrichedEvent creates a new PA-API enrichment success event
// DEPRECATED: Use NewCatalogProductEnrichmentCompletedEvent instead
func NewProductEnrichedEvent(source string, data *ProductEnrichedData) (*Event, error) {
	return NewEvent(PRODUCT_ENRICHMENT_COMPLETED_V1, "product", data.ASIN, data)
}

// NormalizeEventType maps legacy/ad-hoc event strings to canonical numbered constants.
// Returns the normalized event type and true if normalization occurred, false otherwise.
// DEPRECATED: This function will be removed in a future release. Use CAPS constants directly.
func NormalizeEventType(s string) (string, bool) {
	switch s {
	case "CONTENT_GENERATION_REQUESTED":
		return Event_08A_ContentGenerationRequested, true
	case "CONTENT_GENERATION_STARTED":
		return Event_09A_ContentGenerationStarted, true
	case "CONTENT_GENERATED":
		return Event_10A_ContentGenerated, true
	case "CONTENT_GENERATION_FAILED":
		return Event_10B_ContentGenerationFailed, true
	case "CONTENT_GENERATION_RETRIED":
		return Event_10D_ContentGenerationRetried, true
	case "SCRAPER_JOB_REQUESTED":
		return Event_00A_ScraperJobRequested, true
	case Event_00A_ScraperJobRequested:
		return Event_00A_ScraperJobRequested, true
	case "PRODUCT_VALIDATED":
		return Event_02A_ProductValidated, true
	case "PRODUCT_IGNORED":
		return Event_02B_ProductIgnored, true
	case "PRODUCT_REVIEW_REQUIRED":
		return Event_02C_ProductReviewRequired, true
	case "04D_VARIANTS_ENRICHMENT_REQUESTED":
		return VariationEnrichmentRequested, true
	case "ENRICHMENT_RETRY":
		return Event_05D_EnrichmentRetry, true
	case "QUALITY_ASSESSMENT_REQUESTED":
		return Event_06_QualityAssessmentRequested, true
	case "QUALITY_ASSESSMENT_COMPLETED":
		return Event_07A_QualityAssessmentCompleted, true
	case "QUALITY_ASSESSMENT_FAILED":
		return Event_07B_QualityAssessmentFailed, true
	// Product lifecycle events - CloudEvents format mapping
	case "product.deleted.v1":
		return Event_19_ProductDeleted, true
	case "product.status.changed.v1":
		return Event_18_ProductStatusChanged, true
	case "product.availability.changed.v1":
		return Event_17_ProductAvailabilityChanged, true
	case "product.updated.v1":
		return Event_16A_ProductUpdated, true
	case "product.update_failed.v1":
		return Event_16B_ProductUpdateFailed, true
	case "price.updated.v1":
		return Event_15A_PriceUpdated, true
	case "price.update_failed.v1":
		return Event_15B_PriceUpdateFailed, true
	case "product.ready_for_publication.v1":
		return Event_13_ProductReadyForPublication, true
	// Scheduling events
	case "price.monitoring.scheduled.v1":
		return Event_14A_PriceMonitoringScheduled, true
	case "product.availability_check.scheduled.v1":
		return Event_14B_AvailabilityCheckScheduled, true
	case "product.periodic_update.scheduled.v1":
		return Event_14C_PeriodicUpdateScheduled, true
	// Content generation events
	case "content.generation.requested.v1":
		return Event_08A_ContentGenerationRequested, true
	case "content.generation.started.v1":
		return Event_09A_ContentGenerationStarted, true
	case "content.generated.v1":
		return Event_10A_ContentGenerated, true
	case "content.generation.failed.v1":
		return Event_10B_ContentGenerationFailed, true
	case "content.generation.retried.v1":
		return Event_10D_ContentGenerationRetried, true
	// Reviews events
	case "reviews.requested.v1":
		return Event_08B_ReviewsRequested, true
	case "reviews.fetched.v1":
		return Event_09B_ReviewsFetched, true
	case "reviews.processed.v1":
		return Event_10C_ReviewsProcessed, true
	case "reviews.validated.v1":
		return Event_11A_ReviewsValidated, true
	case "reviews.fetch_failed.v1":
		return Event_11B_ReviewsFetchFailed, true
	case "reviews.stored.v1":
		return Event_12A_ReviewsStored, true
	case "reviews.error.v1":
		return Event_12B_ReviewsError, true
	// Quality assessment events
	case "quality.assessment.requested.v1":
		return Event_06_QualityAssessmentRequested, true
	case "quality.assessment.completed.v1":
		return Event_07A_QualityAssessmentCompleted, true
	case "quality.assessment.failed.v1":
		return Event_07B_QualityAssessmentFailed, true
	// Catalog enrichment events mapping
	case "catalog.product.enrichment.requested.v1":
		return CATALOG_PRODUCT_ENRICHMENT_REQUESTED_V1, true
	case "catalog.product.enrichment.completed.v1":
		return CATALOG_PRODUCT_ENRICHMENT_COMPLETED_V1, true
	case "catalog.product.enrichment.failed.v1":
		return CATALOG_PRODUCT_ENRICHMENT_FAILED_V1, true
	// Product enrichment events (keep for backward compatibility)
	case "product.enrichment.requested.v1":
		return PRODUCT_ENRICHMENT_REQUESTED_V1, true
	case "product.enrichment.completed.v1":
		return PRODUCT_ENRICHMENT_COMPLETED_V1, true
	case "product.enrichment.failed.v1":
		return PRODUCT_ENRICHMENT_FAILED_V1, true
	default:
		return s, false
	}
}

// NewProductEnrichmentFailedEvent creates a new PA-API enrichment failure event
// DEPRECATED: Use NewCatalogProductEnrichmentFailedEvent instead
func NewProductEnrichmentFailedEvent(source string, data *ProductEnrichmentFailedData) (*Event, error) {
	return NewEvent(PRODUCT_ENRICHMENT_FAILED_V1, "product", data.ASIN, data)
}

// Catalog event constructor aliases for PA-API enrichment

// NewCatalogProductEnrichmentRequestedEvent creates a new catalog PA-API enrichment request event
// This is the preferred function for new implementations
func NewCatalogProductEnrichmentRequestedEvent(source string, data *ProductEnrichmentRequestedData) (*Event, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}
	return NewEvent(CATALOG_PRODUCT_ENRICHMENT_REQUESTED_V1, "catalog.product", data.ASIN, data)
}

// NewCatalogProductEnrichmentCompletedEvent creates a new catalog PA-API enrichment success event
// This is the preferred function for new implementations
func NewCatalogProductEnrichmentCompletedEvent(source string, data *ProductEnrichedData) (*Event, error) {
	return NewEvent(CATALOG_PRODUCT_ENRICHMENT_COMPLETED_V1, "catalog.product", data.ASIN, data)
}

// NewCatalogProductEnrichmentFailedEvent creates a new catalog PA-API enrichment failure event
// This is the preferred function for new implementations
func NewCatalogProductEnrichmentFailedEvent(source string, data *ProductEnrichmentFailedData) (*Event, error) {
	return NewEvent(CATALOG_PRODUCT_ENRICHMENT_FAILED_V1, "catalog.product", data.ASIN, data)
}
