# Migration Guide: Numbered Events to CloudEvents Format

## Overview

This guide documents the migration from numbered event formats (e.g., `10A_CONTENT_GENERATED`) to CloudEvents naming convention (e.g., `content.generated.v1`). The migration is designed to be backward compatible with dual support during the transition period.

## Migration Status Summary

### ✅ Already Migrated (Dual Support Active)
- All events now have CloudEvents aliases defined
- NormalizeEventType function supports both formats
- Backward compatibility maintained

### 🔄 Active Production Events (Require Careful Migration)

#### Critical Events (HIGH PRIORITY)
- `08A_CONTENT_GENERATION_REQUESTED` → `content.generation.requested.v1` ⚠️
- `10A_CONTENT_GENERATED` → `content.generated.v1` ⚠️
- `10B_CONTENT_GENERATION_FAILED` → `content.generation.failed.v1` ⚠️

#### Scheduling Events (MEDIUM PRIORITY)
- `14A_PRICE_MONITORING_SCHEDULED` → `price.monitoring.scheduled.v1` ⚠️
- `14B_AVAILABILITY_CHECK_SCHEDULED` → `product.availability_check.scheduled.v1` ⚠️
- `14C_PERIODIC_UPDATE_SCHEDULED` → `product.periodic_update.scheduled.v1` ⚠️

### 🚫 Unused Events (Safe to Migrate Immediately)
- `13_PRODUCT_READY_FOR_PUBLICATION` → `product.ready_for_publication.v1` ✅
- `15A_PRICE_UPDATED` → `price.updated.v1` ✅
- `15B_PRICE_UPDATE_FAILED` → `price.update_failed.v1` ✅
- `16A_PRODUCT_UPDATED` → `product.updated.v1` ✅
- `16B_PRODUCT_UPDATE_FAILED` → `product.update_failed.v1` ✅
- `17_PRODUCT_AVAILABILITY_CHANGED` → `product.availability.changed.v1` ✅
- `18_PRODUCT_STATUS_CHANGED` → `product.status.changed.v1` ✅
- `19_PRODUCT_DELETED` → `product.deleted.v1` ✅

## Event Mapping Table

| Old Format | New CloudEvents Format | Status | Used By |
|------------|------------------------|---------|---------|
| `08A_CONTENT_GENERATION_REQUESTED` | `content.generation.requested.v1` | Active | content-generation-service |
| `08B_REVIEWS_REQUESTED` | `reviews.requested.v1` | Check | Unknown |
| `09A_CONTENT_GENERATION_STARTED` | `content.generation.started.v1` | Check | Unknown |
| `09B_REVIEWS_FETCHED` | `reviews.fetched.v1` | Check | Unknown |
| `10A_CONTENT_GENERATED` | `content.generated.v1` | Active | product-lifecycle-service, content-generation-service |
| `10B_CONTENT_GENERATION_FAILED` | `content.generation.failed.v1` | Active | product-lifecycle-service, content-generation-service |
| `10C_REVIEWS_PROCESSED` | `reviews.processed.v1` | Check | Unknown |
| `10D_CONTENT_GENERATION_RETRIED` | `content.generation.retried.v1` | Check | Unknown |
| `11A_REVIEWS_VALIDATED` | `reviews.validated.v1` | Check | Unknown |
| `11B_REVIEWS_FETCH_FAILED` | `reviews.fetch_failed.v1` | Check | Unknown |
| `12A_REVIEWS_STORED` | `reviews.stored.v1` | Check | Unknown |
| `12B_REVIEWS_ERROR` | `reviews.error.v1` | Check | Unknown |
| `13_PRODUCT_READY_FOR_PUBLICATION` | `product.ready_for_publication.v1` | Unused | None |
| `14A_PRICE_MONITORING_SCHEDULED` | `price.monitoring.scheduled.v1` | Active | scheduling-worker |
| `14B_AVAILABILITY_CHECK_SCHEDULED` | `product.availability_check.scheduled.v1` | Active | scheduling-worker |
| `14C_PERIODIC_UPDATE_SCHEDULED` | `product.periodic_update.scheduled.v1` | Active | scheduling-worker |
| `15A_PRICE_UPDATED` | `price.updated.v1` | Unused | None |
| `15B_PRICE_UPDATE_FAILED` | `price.update_failed.v1` | Unused | None |
| `16A_PRODUCT_UPDATED` | `product.updated.v1` | Unused | None |
| `16B_PRODUCT_UPDATE_FAILED` | `product.update_failed.v1` | Unused | None |
| `17_PRODUCT_AVAILABILITY_CHANGED` | `product.availability.changed.v1` | Check | product-lifecycle-service (handler exists) |
| `18_PRODUCT_STATUS_CHANGED` | `product.status.changed.v1` | Unused | None |
| `19_PRODUCT_DELETED` | `product.deleted.v1` | Unused | None |

## Migration Steps

### Step 1: Update tall-affiliate-common (COMPLETED ✅)
```bash
# Version 0.4.0 includes:
# - CloudEvents aliases for all events
# - Updated NormalizeEventType with full mappings
# - Backward compatibility maintained
```

### Step 2: Update Service Handlers (IN PROGRESS)

For each service, update event handlers to support both formats:

```go
// Before:
case events.Event_10A_ContentGenerated:
    // handler code

// After:
case events.Event_10A_ContentGenerated, events.ContentGeneratedV1:
    // handler code
```

### Step 3: Update Event Publishers

Update services to publish using new format:

```go
// Before:
event := events.NewEvent(events.Event_10A_ContentGenerated, ...)

// After:
event := events.NewEvent(events.ContentGeneratedV1, ...)
```

### Step 4: Service Deployment Order

**Critical Path for Production Events:**

1. **Deploy tall-affiliate-common v0.4.0** to all services
2. **Update Consumers First:**
   - content-generation-service (handles 08A, 10A, 10B)
   - product-lifecycle-service (handles 10A, 10B)
   - scheduling-worker (handles 14A, 14B, 14C)
3. **Update Publishers Second:**
   - Services that publish content generation events
   - Services that publish scheduling events
4. **Monitor and Verify:**
   - Check logs for event processing
   - Verify no events are lost
   - Monitor error rates

### Step 5: Deprecate Old Constants

After successful migration (recommended 2-4 weeks):
1. Mark old constants as deprecated
2. Add linting rules to prevent usage
3. Clean up code to use only CloudEvents format

## Service-Specific Migration

### content-generation-service
```go
// handlers.go - Update switch cases
case events.Event_08A_ContentGenerationRequested, events.ContentGenerationRequestedV1:
case events.Event_10A_ContentGenerated, events.ContentGeneratedV1:
case events.Event_10B_ContentGenerationFailed, events.ContentGenerationFailedV1:
```

### product-lifecycle-service
```go
// handlers.go - Update switch cases
case events.Event_10A_ContentGenerated, events.ContentGeneratedV1:
case events.Event_10B_ContentGenerationFailed, events.ContentGenerationFailedV1:
```

### scheduling-worker
```go
// main.go - Update switch cases
case events.Event_14A_PriceMonitoringScheduled, events.PriceMonitoringScheduledV1:
case events.Event_14B_AvailabilityCheckScheduled, events.AvailabilityCheckScheduledV1:
case events.Event_14C_PeriodicUpdateScheduled, events.PeriodicUpdateScheduledV1:
```

## Testing Strategy

### Unit Tests
- Update test fixtures to use both formats
- Verify NormalizeEventType handles all cases
- Test backward compatibility

### Integration Tests
- Deploy to staging with mixed event formats
- Verify services handle both old and new formats
- Test event flow end-to-end

### Production Verification
- Monitor event processing metrics
- Check for any dropped events
- Verify latency remains stable

## Rollback Plan

If issues arise:
1. Services already support both formats via NormalizeEventType
2. Revert publishers to old format if needed
3. No consumer changes needed (they accept both)
4. Monitor and stabilize before retrying

## Timeline

- **Week 1**: Deploy tall-affiliate-common v0.4.0
- **Week 2**: Update critical services (content-generation, product-lifecycle)
- **Week 3**: Update scheduling-worker and monitoring services  
- **Week 4**: Monitor and verify stability
- **Week 5-6**: Deprecate old constants and cleanup

## Success Criteria

- [ ] All services processing events without errors
- [ ] No increase in error rates
- [ ] Event latency remains stable
- [ ] All tests passing with new format
- [ ] Documentation updated
- [ ] Team trained on new format

## Contact

For questions or issues during migration:
- Check service logs for event processing errors
- Review this guide for troubleshooting
- Coordinate deployments through standard channels