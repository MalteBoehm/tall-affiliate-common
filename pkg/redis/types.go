package redis

// StreamMessage represents a message from a Redis stream
type StreamMessage struct {
	ID   string
	Data map[string]interface{}
}