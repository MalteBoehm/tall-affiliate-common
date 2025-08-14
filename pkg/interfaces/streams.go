// pkg/interfaces/streams.go
package interfaces

import (
	"context"
	"time"

	"github.com/MalteBoehm/tall-affiliate-common/pkg/events"
)

// EventHandler defines the function signature for handling events
type EventHandler func(ctx context.Context, event *events.Event) error

// StreamProducer defines the interface for publishing events to a stream
type StreamProducer interface {
	// PublishEvent publishes an event to a stream
	PublishEvent(ctx context.Context, streamName string, event *events.Event) error
}

// StreamConsumer defines the interface for consuming events from a stream
type StreamConsumer interface {
	// ConsumeStream consumes events from a stream
	ConsumeStream(
		ctx context.Context,
		streamName string,
		groupName string,
		batchSize int64,
		pollInterval time.Duration,
		handler func(context.Context, *events.Event, string) error,
	) error
}

// StreamClient defines the interface for a client that can both produce and consume events
type StreamClient interface {
	StreamProducer
	StreamConsumer
}
