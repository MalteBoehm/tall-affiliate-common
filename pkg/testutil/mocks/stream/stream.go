// Package stream provides mock implementations for stream interfaces used in testing.
// These mocks are centralized here to avoid having mock implementations in production code.
package stream

import (
	"context"
	"time"

	"github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/events"
	"github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/interfaces"
	"github.com/stretchr/testify/mock"
)

// MockProducer is a mock implementation of interfaces.StreamProducer for testing
type MockProducer struct {
	mock.Mock
}

// Ensure MockProducer implements interfaces.StreamProducer
var _ interfaces.StreamProducer = (*MockProducer)(nil)

// PublishEvent mocks the PublishEvent method
func (m *MockProducer) PublishEvent(ctx context.Context, streamName string, e *events.Event) error {
	args := m.Called(ctx, streamName, e)
	return args.Error(0)
}

// MockConsumer is a mock implementation of interfaces.StreamConsumer for testing
type MockConsumer struct {
	mock.Mock
}

// Ensure MockConsumer implements interfaces.StreamConsumer
var _ interfaces.StreamConsumer = (*MockConsumer)(nil)

// ConsumeStream mocks the ConsumeStream method
func (m *MockConsumer) ConsumeStream(
	ctx context.Context,
	streamName string,
	groupName string,
	batchSize int64,
	pollInterval time.Duration,
	handler func(context.Context, *events.Event, string) error,
) error {
	args := m.Called(ctx, streamName, groupName, batchSize, pollInterval, handler)
	return args.Error(0)
}

// Helper functions for common test scenarios

// NewMockProducerWithSuccess creates a MockProducer that always returns success
func NewMockProducerWithSuccess() *MockProducer {
	m := new(MockProducer)
	m.On("PublishEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return m
}

// NewMockProducerWithError creates a MockProducer that always returns an error
func NewMockProducerWithError(err error) *MockProducer {
	m := new(MockProducer)
	m.On("PublishEvent", mock.Anything, mock.Anything, mock.Anything).Return(err)
	return m
}

// NewMockConsumerWithSuccess creates a MockConsumer that always returns success
func NewMockConsumerWithSuccess() *MockConsumer {
	m := new(MockConsumer)
	m.On("ConsumeStream", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return m
}

// NewMockConsumerWithError creates a MockConsumer that always returns an error
func NewMockConsumerWithError(err error) *MockConsumer {
	m := new(MockConsumer)
	m.On("ConsumeStream", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(err)
	return m
}