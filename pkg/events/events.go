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
	// Phase 1: Product Discovery (01-02)
	Event_01_ProductDetected        = "01_PRODUCT_DETECTED"
	Event_02A_ProductValidated      = "02A_PRODUCT_VALIDATED"
	Event_02B_ProductIgnored        = "02B_PRODUCT_IGNORED"
	Event_02C_ProductReviewRequired = "02C_PRODUCT_REVIEW_REQUIRED"

	// Phase 2: Enrichment Orchestration (03-05)
	Event_03_EnrichmentOrchestrationStarted = "03_ENRICHMENT_ORCHESTRATION_STARTED"
	Event_04A_DimensionEnrichmentRequested  = "04A_DIMENSION_ENRICHMENT_REQUESTED"
	Event_04B_ColorEnrichmentRequested      = "04B_COLOR_ENRICHMENT_REQUESTED"
	Event_04C_BrowseNodeRequested           = "04C_BROWSE_NODE_REQUESTED"
	Event_05A_EnrichmentCompleted           = "05A_ENRICHMENT_COMPLETED"
	Event_05B_EnrichmentFailed              = "05B_ENRICHMENT_FAILED"
	Event_05D_EnrichmentRetry               = "05D_ENRICHMENT_RETRY"

	// Phase 3: Quality Assessment (06-07)
	Event_06_QualityAssessmentRequested  = "06_QUALITY_ASSESSMENT_REQUESTED"
	Event_07A_QualityAssessmentCompleted = "07A_QUALITY_ASSESSMENT_COMPLETED"
	Event_07B_QualityAssessmentFailed    = "07B_QUALITY_ASSESSMENT_FAILED"

	// Phase 4: Content & Reviews Generation (08-12)
	Event_08A_ContentGenerationRequested  = "08A_CONTENT_GENERATION_REQUESTED"
	Event_08B_ReviewsRequested            = "08B_REVIEWS_REQUESTED"
	Event_09A_ContentGenerationStarted    = "09A_CONTENT_GENERATION_STARTED"
	Event_09B_ReviewsFetched              = "09B_REVIEWS_FETCHED"
	Event_10A_ContentGenerated            = "10A_CONTENT_GENERATED"
	Event_10B_ContentGenerationFailed     = "10B_CONTENT_GENERATION_FAILED"
	Event_10C_ReviewsProcessed            = "10C_REVIEWS_PROCESSED"
	Event_10D_ContentGenerationRetried    = "10D_CONTENT_GENERATION_RETRIED"
	Event_11A_ReviewsValidated            = "11A_REVIEWS_VALIDATED"
	Event_11B_ReviewsFetchFailed          = "11B_REVIEWS_FETCH_FAILED"
	Event_12A_ReviewsStored               = "12A_REVIEWS_STORED"
	Event_12B_ReviewsError                = "12B_REVIEWS_ERROR"

	// Phase 5: Publication & Monitoring (13-14)
	Event_13_ProductReadyForPublication  = "13_PRODUCT_READY_FOR_PUBLICATION"
	Event_14A_PriceMonitoringScheduled   = "14A_PRICE_MONITORING_SCHEDULED"
	Event_14B_AvailabilityCheckScheduled = "14B_AVAILABILITY_CHECK_SCHEDULED"
	Event_14C_PeriodicUpdateScheduled    = "14C_PERIODIC_UPDATE_SCHEDULED"

	// Additional Events (15+)
	Event_15A_PriceUpdated             = "15A_PRICE_UPDATED"
	Event_15B_PriceUpdateFailed        = "15B_PRICE_UPDATE_FAILED"
	Event_16A_ProductUpdated           = "16A_PRODUCT_UPDATED"
	Event_16B_ProductUpdateFailed      = "16B_PRODUCT_UPDATE_FAILED"
	Event_17_ProductAvailabilityChanged = "17_PRODUCT_AVAILABILITY_CHANGED"
	Event_18_ProductStatusChanged       = "18_PRODUCT_STATUS_CHANGED"
	Event_19_ProductDeleted             = "19_PRODUCT_DELETED"

	// Legacy event names (deprecated - for backward compatibility)
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
	EventTypeBrowseNodeRequested        = Event_04C_BrowseNodeRequested
	EventTypeBrowseNodeResolved         = "05A_BROWSE_NODE_RESOLVED"
	EventTypeBrowseNodeFailed           = "05B_BROWSE_NODE_FAILED"
	EventTypeCheckPrice                 = "CHECK_PRICE"
	EventTypePriceUpdated               = Event_15A_PriceUpdated
	EventTypePriceUpdateFailed          = Event_15B_PriceUpdateFailed
	EventTypeProductUpdated             = Event_16A_ProductUpdated
	EventTypeProductAvailabilityChanged = Event_17_ProductAvailabilityChanged
	EventTypeProductStatusChanged       = Event_18_ProductStatusChanged
	EventTypeProductDeleted             = Event_19_ProductDeleted

	// Deprecated - to be removed
	EventTypeProductUnavailable        = "PRODUCT_UNAVAILABLE"
	EventTypeProductCreated             = "PRODUCT_CREATED"
	EventTypeProductUpdateRequested     = "PRODUCT_UPDATE_REQUESTED"
	EventTypeContentUpdateRequested     = "CONTENT_UPDATE_REQUESTED"
	EventTypeContentUpdated             = "CONTENT_UPDATED"
	EventTypeContentAnalysisFailed      = "CONTENT_ANALYSIS_FAILED"
	EventTypeReviewsCollected           = "REVIEWS_COLLECTED"
	EventTypeReviewsEnriched            = "REVIEWS_ENRICHED"
	EventTypeReviewsCached              = "REVIEWS_CACHED"
	EventTypeReviewsExpired             = "REVIEWS_EXPIRED"
	EventTypeReviewsDeleted             = "REVIEWS_DELETED"
	EventTypeBrowseNodeDetectionFailed = "BROWSE_NODE_DETECTION_FAILED"
)

// Legacy orchestration event names (mapped to new convention)
const (
	// Dimension Enrichment
	DimensionEnrichmentRequested = Event_04A_DimensionEnrichmentRequested
	DimensionEnrichmentCompleted = "05A_DIMENSION_ENRICHMENT_COMPLETED"
	DimensionEnrichmentFailed    = "05B_DIMENSION_ENRICHMENT_FAILED"

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

	// PA-API Enrichment Event Types (CloudEvent format)
	ProductEnrichmentRequestedV1 = "product.enrichment.requested.v1"
	ProductEnrichmentCompletedV1 = "product.enrichment.completed.v1"
	ProductEnrichmentFailedV1    = "product.enrichment.failed.v1"
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

// BrowseNodeRequestedPayload represents the payload for a BROWSE_NODE_REQUESTED event
type BrowseNodeRequestedPayload struct {
	ASIN        string    `json:"asin"`
	ProductID   string    `json:"product_id"`
	RequestedAt time.Time `json:"requested_at"`
}

// BrowseNodeResolvedPayload represents the payload for a BROWSE_NODE_RESOLVED event
type BrowseNodeResolvedPayload struct {
	ASIN         string    `json:"asin"`
	ProductID    string    `json:"product_id"`
	BrowseNodeID string    `json:"browse_node_id"`
	ResolvedAt   time.Time `json:"resolved_at"`
}

// BrowseNodeFailedPayload represents the payload for a BROWSE_NODE_FAILED event
type BrowseNodeFailedPayload struct {
	ASIN      string    `json:"asin"`
	ProductID string    `json:"product_id"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
}

// ContentGeneratedPayload represents the payload for a CONTENT_GENERATED event
type ContentGeneratedPayload struct {
	ASIN        string    `json:"asin"`
	ProductID   string    `json:"product_id"`
	ContentType string    `json:"content_type"` // z.B. "description", "faq"
	ContentID   string    `json:"content_id,omitempty"`
	GeneratedAt time.Time `json:"generated_at"`
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

// DimensionEnrichmentRequestedPayload represents the payload for dimension enrichment request
type DimensionEnrichmentRequestedPayload struct {
	ASIN          string    `json:"asin"`
	ProductID     string    `json:"product_id"`
	DetailPageURL string    `json:"detail_page_url"`
	RequestedAt   time.Time `json:"requested_at"`
}

// DimensionEnrichmentCompletedPayload represents the payload for successful dimension enrichment
type DimensionEnrichmentCompletedPayload struct {
	ASIN        string    `json:"asin"`
	ProductID   string    `json:"product_id"`
	HeightCm    *float64  `json:"height_cm,omitempty"`
	LengthCm    *float64  `json:"length_cm,omitempty"`
	WidthCm     *float64  `json:"width_cm,omitempty"`
	Source      string    `json:"source"` // "amazon-scraper"
	CompletedAt time.Time `json:"completed_at"`
}

// DimensionEnrichmentFailedPayload represents the payload for failed dimension enrichment
type DimensionEnrichmentFailedPayload struct {
	ASIN      string    `json:"asin"`
	ProductID string    `json:"product_id"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
}

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

// UnmarshalPayload unmarshals the event payload into the given value
func (e *Event) UnmarshalPayload(v any) error {
	payloadBytes, err := json.Marshal(e.Payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(payloadBytes, v)
}

// Reviews-related constants and types
const (
	ReviewsSourceAmazonAPI = "amazon_api"
	ReviewsSourceManual    = "manual"
)

const (
	ReviewsErrorTypeFetch      = "fetch_error"
	ReviewsErrorTypeProcess    = "process_error"
	ReviewsErrorTypeValidation = "validation_error"
	ReviewsErrorTypeCache      = "cache_error"
	ReviewsErrorTypeDatabase   = "database_error"
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

// NewDimensionEnrichmentRequestedEvent creates a new dimension enrichment requested event
func NewDimensionEnrichmentRequestedEvent(asin, productID, detailPageURL string) *Event {
	payload := DimensionEnrichmentRequestedPayload{
		ASIN:          asin,
		ProductID:     productID,
		DetailPageURL: detailPageURL,
		RequestedAt:   time.Now().UTC(),
	}

	event, _ := NewEvent(DimensionEnrichmentRequested, "product", productID, payload)
	return event
}

// NewDimensionEnrichmentCompletedEvent creates a new dimension enrichment completed event
func NewDimensionEnrichmentCompletedEvent(asin, productID string, heightCm, lengthCm, widthCm *float64) *Event {
	payload := DimensionEnrichmentCompletedPayload{
		ASIN:        asin,
		ProductID:   productID,
		HeightCm:    heightCm,
		LengthCm:    lengthCm,
		WidthCm:     widthCm,
		Source:      "amazon-scraper",
		CompletedAt: time.Now().UTC(),
	}

	event, _ := NewEvent(DimensionEnrichmentCompleted, "product", productID, payload)
	return event
}

// NewDimensionEnrichmentFailedEvent creates a new dimension enrichment failed event
func NewDimensionEnrichmentFailedEvent(asin, productID, reason string) *Event {
	payload := DimensionEnrichmentFailedPayload{
		ASIN:      asin,
		ProductID: productID,
		Reason:    reason,
		FailedAt:  time.Now().UTC(),
	}

	event, _ := NewEvent(DimensionEnrichmentFailed, "product", productID, payload)
	return event
}

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
func NewProductEnrichmentRequestedEvent(source string, data *ProductEnrichmentRequestedData) (*Event, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}
	return NewEvent(ProductEnrichmentRequestedV1, "product", data.ASIN, data)
}

// NewProductEnrichedEvent creates a new PA-API enrichment success event
func NewProductEnrichedEvent(source string, data *ProductEnrichedData) (*Event, error) {
	return NewEvent(ProductEnrichmentCompletedV1, "product", data.ASIN, data)
}

// NewProductEnrichmentFailedEvent creates a new PA-API enrichment failure event
func NewProductEnrichmentFailedEvent(source string, data *ProductEnrichmentFailedData) (*Event, error) {
	return NewEvent(ProductEnrichmentFailedV1, "product", data.ASIN, data)
}
