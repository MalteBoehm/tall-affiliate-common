// pkg/constants/constants.go
package constants

// Stream names (CAPS naming convention)
const (
	STREAM_CONTENT_GENERATION = "stream:content_generation"
	STREAM_PRODUCT_LIFECYCLE  = "stream:product_lifecycle"
	// StreamBrowseNodes removed - use STREAM_PRODUCT_LIFECYCLE for enrichment
	STREAM_PRICE_TRACKING = "stream:price_tracking"
)

// Consumer group names (CAPS naming convention)
const (
	GROUP_PRODUCT_LIFECYCLE  = "group:product_lifecycle"
	GROUP_CONTENT_GENERATION = "group:content_generation"
	// GroupBrowseNodes removed
	GROUP_PRICE_TRACKING = "group:price_tracking"
)

// Worker names (CAPS naming convention)
const (
	WORKER_PRODUCT_LIFECYCLE  = "product_lifecycle_worker"
	WORKER_CONTENT_GENERATION = "content_generation_worker"
	// WorkerBrowseNodes removed
	WORKER_PRICE_TRACKING = "price_tracking_worker"
)

// Database table names (CAPS naming convention)
const (
	TABLE_PRODUCT       = "product"
	TABLE_OUTBOX        = "outbox"
	TABLE_PRICE_HISTORY = "price_history"
	// TableBrowseNode removed
)

// Legacy constants for backward compatibility
// DEPRECATED: Use CAPS constants instead. Will be removed in a future release.
const (
	StreamContentGeneration = STREAM_CONTENT_GENERATION
	StreamProductLifecycle  = STREAM_PRODUCT_LIFECYCLE
	StreamPriceTracking     = STREAM_PRICE_TRACKING

	GroupProductLifecycle  = GROUP_PRODUCT_LIFECYCLE
	GroupContentGeneration = GROUP_CONTENT_GENERATION
	GroupPriceTracking     = GROUP_PRICE_TRACKING

	WorkerProductLifecycle  = WORKER_PRODUCT_LIFECYCLE
	WorkerContentGeneration = WORKER_CONTENT_GENERATION
	WorkerPriceTracking     = WORKER_PRICE_TRACKING

	TableProduct      = TABLE_PRODUCT
	TableOutbox       = TABLE_OUTBOX
	TablePriceHistory = TABLE_PRICE_HISTORY
)
