# Migration Guide: Dimension Enrichment Events Removal

## Version
This change is introduced in v1.x.x (TODO: Update with actual version)

## Summary
Dimension Enrichment events and related functionality have been deprecated and removed from the tall-affiliate-common library. This functionality has been replaced by the PA-API enrichment flow using CloudEvent format events.

## Removed Components

### Event Constants
- `DimensionEnrichmentRequested`
- `DimensionEnrichmentCompleted`
- `DimensionEnrichmentFailed`
- `Event_04A_DimensionEnrichmentRequested`

### Payload Types
- `DimensionEnrichmentRequestedPayload`
- `DimensionEnrichmentCompletedPayload`
- `DimensionEnrichmentFailedPayload`

### Functions
- `NewDimensionEnrichmentRequestedEvent()`
- `NewDimensionEnrichmentCompletedEvent()`
- `NewDimensionEnrichmentFailedEvent()`

## Migration Path

### For Event Publishers
Replace Dimension Enrichment event publishing with PA-API enrichment events:

**Before:**
```go
event := events.NewDimensionEnrichmentRequestedEvent(asin, productID, detailPageURL)
adapter.PublishEvent(ctx, constants.StreamProductLifecycle, event)
```

**After:**
```go
enrichmentData := &events.ProductEnrichmentRequestedData{
    ASIN:       asin,
    Region:     "de",
    RequestID:  productID,
    RetryCount: 0,
    Dimensions: true,  // Request dimension data
}
event, _ := events.NewProductEnrichmentRequestedEvent("service-name", enrichmentData)
adapter.PublishEvent(ctx, constants.StreamProductLifecycle, event)
```

### For Event Consumers
Update event handlers to listen for PA-API enrichment events:

**Before:**
```go
case events.DimensionEnrichmentRequested:
    // Handle dimension enrichment request
case events.DimensionEnrichmentCompleted:
    // Handle dimension enrichment completed
case events.DimensionEnrichmentFailed:
    // Handle dimension enrichment failed
```

**After:**
```go
case events.ProductEnrichmentRequestedV1:
    // Handle enrichment request
case events.ProductEnrichmentCompletedV1:
    // Handle enrichment completion
    // Extract dimensions from enrichment data
case events.ProductEnrichmentFailedV1:
    // Handle enrichment failure
```

### For Dimension Data Access
Dimensions are now included in the PA-API enrichment response:

**Before:**
```go
payload := DimensionEnrichmentCompletedPayload{
    HeightCm: heightCm,
    LengthCm: lengthCm,
    WidthCm:  widthCm,
}
```

**After:**
```go
// Dimensions are part of the enrichment response
enrichmentResponse := ProductEnrichmentCompletedData{
    ASIN:   asin,
    Region: "de",
    Data: map[string]interface{}{
        "dimensions": map[string]interface{}{
            "height_cm": heightCm,
            "length_cm": lengthCm,
            "width_cm":  widthCm,
        },
        // Other enrichment data...
    },
}
```

## Breaking Changes
- Any code directly referencing removed constants will fail to compile
- Services listening for Dimension Enrichment events will no longer receive them
- Functions creating Dimension Enrichment events no longer exist

## Recommended Actions
1. Update all imports and references to removed constants
2. Migrate event publishers to use PA-API enrichment events
3. Update event consumers to handle new event types
4. Extract dimension data from PA-API enrichment responses
5. Update service documentation to reflect the new flow

## Note on Combined Enrichment
The PA-API enrichment flow now combines multiple enrichment types (dimensions, colors, variations, etc.) into a single request/response cycle, reducing the number of events and improving efficiency.