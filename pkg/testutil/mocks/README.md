# Test Mocks

This directory contains centralized mock implementations for testing purposes. 

## Important: No Mocks in Production Code!

All mock implementations should be placed here in the `testutil` package, NOT in the production code directories (`/internal/`).

## Available Mocks

### Stream Mocks (`stream/`)

Mock implementations for stream interfaces:
- `MockProducer` - Mocks `interfaces.StreamProducer`
- `MockConsumer` - Mocks `interfaces.StreamConsumer`

Usage example:
```go
import "github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/testutil/mocks/stream"

func TestMyFunction(t *testing.T) {
    // Create a mock producer that always succeeds
    producer := stream.NewMockProducerWithSuccess()
    
    // Or create one with custom behavior
    producer := new(stream.MockProducer)
    producer.On("PublishEvent", mock.Anything, "my-stream", mock.Anything).Return(nil)
    
    // Use in your test
    service := NewService(producer)
    // ... test logic
    
    // Verify expectations
    producer.AssertExpectations(t)
}
```

## Guidelines for Adding New Mocks

1. **Only add mocks here if:**
   - The interface is defined in a common package
   - The mock will be used by multiple services
   - The mock doesn't require package-private access

2. **Keep mocks in their package if:**
   - They need access to package-private types
   - They are only used by that package's tests
   - They contain package-specific mock data

3. **Naming Convention:**
   - Use `Mock` prefix: `MockProducer`, `MockClient`
   - Place in subdirectory matching the interface domain: `stream/`, `database/`, etc.

4. **Documentation:**
   - Add godoc comments explaining what interface is mocked
   - Include usage examples in comments
   - Update this README with new mock types

## Migration Status

- ✅ Stream mocks migrated from `product-lifecycle-service/internal/worker/stream_mocks.go`
- ❌ Oxylabs mocks remain in their packages (tight coupling)
- ❌ Worker mocks remain in their packages (heavy internal usage)
- ❌ Notifier mocks remain in their packages (unused, internal interface)