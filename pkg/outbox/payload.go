package outbox

import (
	"encoding/json"
	"fmt"
)

const (
	// MaxPayloadSize defines the maximum allowed payload size in bytes (16KB)
	MaxPayloadSize = 16 * 1024
)

// PayloadError represents an error related to payload handling
type PayloadError struct {
	Operation string
	Err       error
}

func (e *PayloadError) Error() string {
	return fmt.Sprintf("payload %s error: %v", e.Operation, e.Err)
}

// MarshalPayload converts any value to json.RawMessage with validation
func MarshalPayload(payload any) (json.RawMessage, error) {
	if payload == nil {
		return json.RawMessage("null"), nil
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, &PayloadError{
			Operation: "marshal",
			Err:       err,
		}
	}

	// Validate payload size
	if len(data) > MaxPayloadSize {
		return nil, &PayloadError{
			Operation: "validation",
			Err:       fmt.Errorf("payload size %d bytes exceeds maximum %d bytes", len(data), MaxPayloadSize),
		}
	}

	return json.RawMessage(data), nil
}

// UnmarshalPayload converts json.RawMessage to the target type
func UnmarshalPayload(payload json.RawMessage, target any) error {
	if len(payload) == 0 {
		return &PayloadError{
			Operation: "unmarshal",
			Err:       fmt.Errorf("empty payload"),
		}
	}

	err := json.Unmarshal(payload, target)
	if err != nil {
		return &PayloadError{
			Operation: "unmarshal",
			Err:       err,
		}
	}

	return nil
}

// ValidatePayloadSize checks if the payload size is within limits
func ValidatePayloadSize(payload json.RawMessage) error {
	if len(payload) > MaxPayloadSize {
		return &PayloadError{
			Operation: "validation",
			Err:       fmt.Errorf("payload size %d bytes exceeds maximum %d bytes", len(payload), MaxPayloadSize),
		}
	}
	return nil
}

// IsValidJSON checks if the payload contains valid JSON
func IsValidJSON(payload json.RawMessage) bool {
	var js json.RawMessage
	return json.Unmarshal(payload, &js) == nil
}

// ConvertStringToRawMessage converts a JSON string to json.RawMessage with validation
func ConvertStringToRawMessage(jsonStr string) (json.RawMessage, error) {
	if jsonStr == "" {
		return json.RawMessage("null"), nil
	}

	// Validate it's valid JSON
	var temp interface{}
	if err := json.Unmarshal([]byte(jsonStr), &temp); err != nil {
		return nil, &PayloadError{
			Operation: "conversion",
			Err:       fmt.Errorf("invalid JSON string: %v", err),
		}
	}

	// Validate size
	if len(jsonStr) > MaxPayloadSize {
		return nil, &PayloadError{
			Operation: "validation",
			Err:       fmt.Errorf("payload size %d bytes exceeds maximum %d bytes", len(jsonStr), MaxPayloadSize),
		}
	}

	return json.RawMessage(jsonStr), nil
}

// ConvertBytesToRawMessage converts []byte to json.RawMessage with validation
func ConvertBytesToRawMessage(data []byte) (json.RawMessage, error) {
	if len(data) == 0 {
		return json.RawMessage("null"), nil
	}

	// Validate it's valid JSON
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, &PayloadError{
			Operation: "conversion",
			Err:       fmt.Errorf("invalid JSON bytes: %v", err),
		}
	}

	// Validate size
	if len(data) > MaxPayloadSize {
		return nil, &PayloadError{
			Operation: "validation",
			Err:       fmt.Errorf("payload size %d bytes exceeds maximum %d bytes", len(data), MaxPayloadSize),
		}
	}

	return json.RawMessage(data), nil
}
