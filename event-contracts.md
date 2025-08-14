# Event Contracts

This document defines the event contracts between services in the Tall Affiliate system.

## ✅ **KONSOLIDIERT: Einheitliche Event-Definitionen**

**Alle Event-Definitionen sind jetzt in `tall-affiliate-common/pkg/events` konsolidiert.**

## Overview

Services communicate with each other using events published to Redis Streams. Each event has a specific type and payload schema. All event types and payload schemas are now defined in the common package to ensure consistency across all services.

## Event Structure

**📍 Defined in `tall-affiliate-common/pkg/events/events.go`**

All events have the following unified structure:

```go
type Event struct {
    ID            string         `json:"id"`
    Type          string         `json:"type"`
    AggregateType string         `json:"aggregate_type"`
    AggregateID   string         `json:"aggregate_id"`
    Payload       any            `json:"payload"`
    Timestamp     time.Time      `json:"timestamp"`
    Metadata      map[string]any `json:"metadata,omitempty"`
}
```

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "type": "NEW_PRODUCT_DETECTED",
  "aggregate_type": "product",
  "aggregate_id": "123e4567-e89b-12d3-a456-426614174001",
  "timestamp": "2023-06-15T12:34:56Z",
  "payload": {
    // Event-specific payload (strongly typed)
  },
  "metadata": {
    "source": "api-gateway",
    "version": "1.0"
  }
}
```

- `id`: Unique event ID (UUID)
- `type`: Event type (string, defined as constants)
- `aggregate_type`: Type of aggregate (e.g., "product", "content", "browse_node")
- `aggregate_id`: ID of the aggregate (usually the product ID)
- `timestamp`: Event timestamp (RFC3339)
- `payload`: Event-specific payload (strongly typed structures)
- `metadata`: Optional metadata for additional context

## Streams and Consumer Groups

**📍 Defined in `tall-affiliate-common/pkg/constants/constants.go`**

The system uses the following Redis Streams:

- `stream:product_lifecycle`: Events related to product lifecycle
- `stream:content_generation`: Events related to content generation
- `stream:browse_nodes`: Events related to browse node resolution
- `stream:price_tracking`: Events related to price tracking

And the following consumer groups:

- `group:product_lifecycle`: Consumer group for product lifecycle events
- `group:content_generation`: Consumer group for content generation events
- `group:browse_nodes`: Consumer group for browse node events
- `group:price_tracking`: Consumer group for price tracking events

## Service Integration

**📍 Unified Adapter Pattern in `tall-affiliate-common/pkg/adapters`**

All services now use the same adapter pattern for event handling:

```go
// Create service adapter
adapter := adapters.NewServiceEventAdapter(producer, consumer)

// Publish events with automatic stream routing
adapter.PublishProductEvent(ctx, events.EventTypeNewProductDetected, productID, asin, payload)
adapter.PublishContentEvent(ctx, events.EventTypeContentGenerationRequested, productID, payload)
adapter.PublishBrowseNodeEvent(ctx, events.EventTypeBrowseNodeRequested, productID, payload)
```

### Automatic Stream Routing

The adapter automatically routes events to the correct stream based on event type:
- Product Lifecycle Events → `stream:product_lifecycle`
- Content Generation Events → `stream:content_generation`
- Browse Node Events → `stream:browse_nodes`
- Price Tracking Events → `stream:price_tracking`

## Product Lifecycle Events

### product.new_detected

Emitted when a new product is detected (e.g., from Amazon PA-API).

**Stream:** `stream:product_lifecycle`

**Payload:**

```json
{
  "asin": "B01EXAMPLE",
  "title": "Product Title",
  "brand": "Brand Name",
  "category": "Clothing",
  "gender": "male",
  "current_price": 99.99,
  "currency": "EUR",
  "detail_page_url": "https://amazon.com/dp/B01EXAMPLE",
  "image_urls": ["https://example.com/image1.jpg", "https://example.com/image2.jpg"],
  "features": ["Feature 1", "Feature 2"],
  "source": "pa-api",
  "size": "XL",
  "color": "Blue",
  "available_sizes": ["L", "XL", "XXL"],
  "height_cm": 75.5,
  "length_cm": 85.0,
  "width_cm": 60.0,
  "weight_g": 500.0,
  "product_group": "Apparel",
  "model": "Model123",
  "is_adult_product": false,
  "variation_attributes": ["Size", "Color"],
  "browse_node_id": "1981507031",
  "tall_friendly_score": 8.5,
  "is_tall_friendly": true,
  "rating": 4.5,
  "review_count": 150
}
```

**Quality Assessment Fields:**
- `rating`: Product rating (1-5 stars) - **CRITICAL for quality assessment**
- `review_count`: Number of reviews - **CRITICAL for quality assessment**

These fields are used by the Product Lifecycle Service to determine whether to create the product immediately (`PRODUCT_CREATED`), require manual review (`PRODUCT_REVIEW_REQUIRED`), or ignore it (`PRODUCT_IGNORED`).

### product.created

Emitted when a new product is created (quality score ≥ 3.0).

**Stream:** `stream:product_lifecycle`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "title": "Product Title",
  "brand": "Brand Name",
  "current_price": 99.99,
  "currency": "EUR",
  "image_url": "https://example.com/image.jpg",
  "category": "Clothing",
  "gender": "male",
  "created_at": "2023-06-15T12:34:56Z"
}
```

### product.review_required

Emitted when a product requires manual review (quality score 1.0-2.9).

**Stream:** `stream:product_lifecycle`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "title": "Product Title",
  "brand": "Brand Name",
  "current_price": 99.99,
  "currency": "EUR",
  "reason": "borderline_quality",
  "score": 2.5,
  "details": {
    "rating_score": 1.0,
    "review_count_score": 0.5,
    "image_score": 0.5,
    "price_score": 1.0,
    "brand_score": 0.5
  },
  "created_at": "2023-06-15T12:34:56Z"
}
```

### product.ignored

Emitted when a product is ignored due to insufficient quality (score < 1.0).

**Stream:** `stream:product_lifecycle`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "title": "Product Title",
  "reason": "insufficient_quality",
  "score": 0.5,
  "created_at": "2023-06-15T12:34:56Z"
}
```

### product.updated

Emitted when a product is updated.

**Stream:** `stream:product-lifecycle`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "updated_fields": ["title", "current_price"],
  "updated_at": "2023-06-15T12:34:56Z"
}
```

### product.deleted

Emitted when a product is deleted.

**Stream:** `stream:product-lifecycle`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "deleted_at": "2023-06-15T12:34:56Z"
}
```

### ❌ product.availability_check_requested - ENTFERNT

**Warum entfernt:**
- **Redundant**: Verfügbarkeitsdaten sind bereits in allen PA-API Responses enthalten
- **Überflüssig**: SearchItems/GetItems liefern immer Offers.Listings.Availability
- **Kein separater Check nötig**: PA-API hat keinen dedizierten "Availability Check"

**Ersetzt durch:**
- Verfügbarkeitsdaten in `NEW_PRODUCT_DETECTED` Events
- Verfügbarkeitsdaten in `PRODUCT_UPDATED` Events
- `PRODUCT_AVAILABILITY_CHANGED` bei Änderungen
```

### product.availability_changed

Emitted when the availability of a product changes.

**Stream:** `stream:product-lifecycle`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "availability": "in_stock",
  "previous_availability": "out_of_stock",
  "changed_at": "2023-06-15T12:34:56Z"
}
```

### product.price_changed

Emitted when the price of a product changes.

**Stream:** `stream:product-lifecycle`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "current_price": 99.99,
  "previous_price": 109.99,
  "currency": "EUR",
  "changed_at": "2023-06-15T12:34:56Z"
}
```

### product.browse_node_updated

Emitted when the browse node ID of a product is updated.

**Stream:** `stream:product-lifecycle`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "browse_node_id": "12345",
  "previous_browse_node_id": null,
  "updated_at": "2023-06-15T12:34:56Z"
}
```

### product.content_generation_requested

Emitted when content generation is requested for a product.

**Stream:** `stream:product-lifecycle`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "priority": 3,
  "options": {
    "include_description": true,
    "include_features": true,
    "include_reviews_summary": true,
    "include_faq": true
  },
  "requested_at": "2023-06-15T12:34:56Z"
}
```

### product.content_generated

Emitted when content is generated for a product.

**Stream:** `stream:product-lifecycle`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "content_id": "123e4567-e89b-12d3-a456-426614174001",
  "generated_at": "2023-06-15T12:34:56Z"
}
```

## Content Generation Events

### content.generation_requested

Emitted when content generation is requested.

**Stream:** `stream:content-generation`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "priority": 3,
  "options": {
    "include_description": true,
    "include_features": true,
    "include_reviews_summary": true,
    "include_faq": true
  },
  "requested_at": "2023-06-15T12:34:56Z"
}
```

### content.generated

Emitted when content is generated.

**Stream:** `stream:content-generation`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "content_id": "123e4567-e89b-12d3-a456-426614174001",
  "content_type": "description",
  "generated_at": "2023-06-15T12:34:56Z"
}
```

### content.generation_failed

Emitted when content generation fails.

**Stream:** `stream:content-generation`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "reason": "Failed to generate content: API timeout",
  "failed_at": "2023-06-15T12:34:56Z"
}
```

### content.reviews_collected

Emitted when reviews are collected.

**Stream:** `stream:content-generation`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "review_count": 10,
  "collected_at": "2023-06-15T12:34:56Z"
}
```

## Browse Node Events

### browse_node.requested

Emitted when browse node resolution is requested.

**Stream:** `stream:browse-node`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "priority": 3,
  "requested_at": "2023-06-15T12:34:56Z"
}
```

### browse_node.resolved

Emitted when a browse node is resolved.

**Stream:** `stream:browse-node`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "browse_node_id": "12345",
  "resolved_at": "2023-06-15T12:34:56Z"
}
```

### browse_node.failed

Emitted when browse node resolution fails.

**Stream:** `stream:browse-node`

**Payload:**

```json
{
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "asin": "B01EXAMPLE",
  "reason": "Failed to resolve browse node: API error",
  "failed_at": "2023-06-15T12:34:56Z"
}
```

## Event Handling

Services should handle events according to the following guidelines:

1. **Idempotency**: Event handlers should be idempotent, meaning they can be called multiple times with the same event without causing issues.
2. **Error Handling**: If an event cannot be processed, it should be moved to a dead letter queue (DLQ) after a certain number of retries.
3. **Acknowledgment**: Events should be acknowledged only after they have been successfully processed.
4. **Ordering**: Events should be processed in the order they were published.
5. **Retry**: Failed event processing should be retried with exponential backoff.

## Event Schema Evolution

When evolving event schemas, follow these guidelines:

1. **Backward Compatibility**: New versions of event schemas should be backward compatible with old versions.
2. **Optional Fields**: New fields should be optional.
3. **Default Values**: New fields should have sensible default values.
4. **Versioning**: Consider adding a version field to events if backward compatibility cannot be maintained.
