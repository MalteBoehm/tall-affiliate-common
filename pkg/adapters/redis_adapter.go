// pkg/adapters/redis_adapter.go
package adapters

import (
	"context"
	"time"

	"github.com/MalteBoehm/tall-affiliate-common/pkg/events"
	"github.com/MalteBoehm/tall-affiliate-common/pkg/interfaces"
)

// RedisProducerAdapter adapts any Redis producer to the common StreamProducer interface
type RedisProducerAdapter struct {
	producer RedisProducer
}

// RedisProducer interface that any Redis producer implementation must satisfy
type RedisProducer interface {
	PublishEvent(ctx context.Context, streamName string, event *events.Event) error
}

// NewRedisProducerAdapter creates a new RedisProducerAdapter
func NewRedisProducerAdapter(producer RedisProducer) *RedisProducerAdapter {
	return &RedisProducerAdapter{
		producer: producer,
	}
}

// PublishEvent publishes an event to a stream
func (p *RedisProducerAdapter) PublishEvent(ctx context.Context, streamName string, event *events.Event) error {
	return p.producer.PublishEvent(ctx, streamName, event)
}

// Ensure RedisProducerAdapter implements interfaces.StreamProducer
var _ interfaces.StreamProducer = (*RedisProducerAdapter)(nil)

// RedisConsumerAdapter adapts any Redis consumer to the common StreamConsumer interface
type RedisConsumerAdapter struct {
	consumer RedisConsumer
}

// RedisConsumer interface that any Redis consumer implementation must satisfy
type RedisConsumer interface {
	ConsumeStream(
		ctx context.Context,
		streamName string,
		groupName string,
		batchSize int64,
		pollInterval time.Duration,
		handler func(context.Context, *events.Event, string) error,
	) error
}

// NewRedisConsumerAdapter creates a new RedisConsumerAdapter
func NewRedisConsumerAdapter(consumer RedisConsumer) *RedisConsumerAdapter {
	return &RedisConsumerAdapter{
		consumer: consumer,
	}
}

// ConsumeStream consumes events from a stream
func (c *RedisConsumerAdapter) ConsumeStream(
	ctx context.Context,
	streamName string,
	groupName string,
	batchSize int64,
	pollInterval time.Duration,
	handler func(context.Context, *events.Event, string) error,
) error {
	return c.consumer.ConsumeStream(ctx, streamName, groupName, batchSize, pollInterval, handler)
}

// Ensure RedisConsumerAdapter implements interfaces.StreamConsumer
var _ interfaces.StreamConsumer = (*RedisConsumerAdapter)(nil)
