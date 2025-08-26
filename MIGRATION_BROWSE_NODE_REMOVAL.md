# Migration Guide: Browse Node Events Removal

## Version
This change is introduced in v1.x.x (TODO: Update with actual version)

## Summary
Browse Node events and related functionality have been deprecated and removed from the tall-affiliate-common library. This functionality has been replaced by the PA-API enrichment flow using CloudEvent format events.

## Removed Components

### Event Constants
- `Event_04C_BrowseNodeRequested` 
- `EventTypeBrowseNodeRequested`
- `EventTypeBrowseNodeResolved` 
- `EventTypeBrowseNodeFailed`
- `EventTypeBrowseNodeDetectionFailed`

### Payload Types
- `BrowseNodeRequestedPayload`
- `BrowseNodeResolvedPayload` 
- `BrowseNodeFailedPayload`

### Functions
- `PublishBrowseNodeEvent()` in ServiceEventAdapter
- `isBrowseNodeEvent()` helper function

### Constants
- `StreamBrowseNodes`
- `GroupBrowseNodes`
- `WorkerBrowseNodes`
- `TableBrowseNode`

## Migration Path

### For Event Publishers
Replace Browse Node event publishing with PA-API enrichment events:

**Before:**
```go
adapter.PublishBrowseNodeEvent(ctx, events.EventTypeBrowseNodeRequested, productID, payload)
```

**After:**
```go
enrichmentData := &events.ProductEnrichmentRequestedData{
    ASIN:       asin,
    Region:     "de",
    RequestID:  productID,
    RetryCount: 0,
}
event, _ := events.NewProductEnrichmentRequestedEvent("service-name", enrichmentData)
adapter.PublishEvent(ctx, constants.StreamProductLifecycle, event)
```

### For Event Consumers
Update event handlers to listen for PA-API enrichment events:

**Before:**
```go
case events.EventTypeBrowseNodeRequested:
    // Handle browse node request
case events.EventTypeBrowseNodeResolved:
    // Handle browse node resolved
```

**After:**
```go
case events.ProductEnrichmentRequestedV1:
    // Handle enrichment request
case events.ProductEnrichmentCompletedV1:
    // Handle enrichment completion
```

### For Stream Configuration
Use `StreamProductLifecycle` instead of `StreamBrowseNodes`:

**Before:**
```go
streamName := constants.StreamBrowseNodes
```

**After:**
```go
streamName := constants.StreamProductLifecycle
```

## Note on Browse Node Detection
The Browse Node detection logic (`pkg/patterns/detection.go`) remains available for services that need to analyze Browse Node data internally. This is now considered an internal implementation detail rather than an event-driven workflow.

## Breaking Changes
- Any code directly referencing removed constants will fail to compile
- Services listening for Browse Node events will no longer receive them
- Database tables for browse_node may need migration or cleanup

## Recommended Actions
1. Update all imports and references to removed constants
2. Migrate event publishers to use PA-API enrichment events
3. Update event consumers to handle new event types
4. Clean up any browse_node database tables if no longer needed
5. Update service documentation to reflect the new flow