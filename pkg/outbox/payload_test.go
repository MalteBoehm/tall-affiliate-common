package outbox

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalPayload(t *testing.T) {
	tests := []struct {
		name        string
		payload     any
		expectError bool
		errorType   string
	}{
		{
			name:        "nil payload",
			payload:     nil,
			expectError: false,
		},
		{
			name: "simple struct",
			payload: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "John",
				Age:  30,
			},
			expectError: false,
		},
		{
			name:        "large payload",
			payload:     map[string]string{"data": strings.Repeat("x", MaxPayloadSize)},
			expectError: true,
			errorType:   "validation",
		},
		{
			name:        "unmarshalable payload",
			payload:     make(chan int),
			expectError: true,
			errorType:   "marshal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MarshalPayload(tt.payload)

			if tt.expectError {
				require.Error(t, err)
				var payloadErr *PayloadError
				require.ErrorAs(t, err, &payloadErr)
				assert.Equal(t, tt.errorType, payloadErr.Operation)
			} else {
				require.NoError(t, err)
				assert.True(t, IsValidJSON(result))
			}
		})
	}
}

func TestUnmarshalPayload(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name        string
		payload     json.RawMessage
		target      any
		expectError bool
		errorType   string
	}{
		{
			name:        "valid payload",
			payload:     json.RawMessage(`{"name":"John","age":30}`),
			target:      &TestStruct{},
			expectError: false,
		},
		{
			name:        "empty payload",
			payload:     json.RawMessage(""),
			target:      &TestStruct{},
			expectError: true,
			errorType:   "unmarshal",
		},
		{
			name:        "invalid JSON",
			payload:     json.RawMessage(`{"name":"John","age":}`),
			target:      &TestStruct{},
			expectError: true,
			errorType:   "unmarshal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalPayload(tt.payload, tt.target)

			if tt.expectError {
				require.Error(t, err)
				var payloadErr *PayloadError
				require.ErrorAs(t, err, &payloadErr)
				assert.Equal(t, tt.errorType, payloadErr.Operation)
			} else {
				require.NoError(t, err)
				result := tt.target.(*TestStruct)
				assert.Equal(t, "John", result.Name)
				assert.Equal(t, 30, result.Age)
			}
		})
	}
}

func TestValidatePayloadSize(t *testing.T) {
	tests := []struct {
		name        string
		payload     json.RawMessage
		expectError bool
	}{
		{
			name:        "small payload",
			payload:     json.RawMessage(`{"test":"data"}`),
			expectError: false,
		},
		{
			name:        "large payload",
			payload:     json.RawMessage(strings.Repeat("x", MaxPayloadSize+1)),
			expectError: true,
		},
		{
			name:        "max size payload",
			payload:     json.RawMessage(strings.Repeat("x", MaxPayloadSize)),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePayloadSize(tt.payload)

			if tt.expectError {
				require.Error(t, err)
				var payloadErr *PayloadError
				require.ErrorAs(t, err, &payloadErr)
				assert.Equal(t, "validation", payloadErr.Operation)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestIsValidJSON(t *testing.T) {
	tests := []struct {
		name     string
		payload  json.RawMessage
		expected bool
	}{
		{
			name:     "valid JSON object",
			payload:  json.RawMessage(`{"test":"data"}`),
			expected: true,
		},
		{
			name:     "valid JSON array",
			payload:  json.RawMessage(`[1,2,3]`),
			expected: true,
		},
		{
			name:     "valid JSON string",
			payload:  json.RawMessage(`"test"`),
			expected: true,
		},
		{
			name:     "invalid JSON",
			payload:  json.RawMessage(`{"test":}`),
			expected: false,
		},
		{
			name:     "empty payload",
			payload:  json.RawMessage(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidJSON(tt.payload)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertStringToRawMessage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorType   string
	}{
		{
			name:        "valid JSON string",
			input:       `{"test":"data"}`,
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: false,
		},
		{
			name:        "invalid JSON",
			input:       `{"test":}`,
			expectError: true,
			errorType:   "conversion",
		},
		{
			name:        "large string",
			input:       strings.Repeat("x", MaxPayloadSize+1),
			expectError: true,
			errorType:   "validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertStringToRawMessage(tt.input)

			if tt.expectError {
				require.Error(t, err)
				var payloadErr *PayloadError
				require.ErrorAs(t, err, &payloadErr)
				assert.Equal(t, tt.errorType, payloadErr.Operation)
			} else {
				require.NoError(t, err)
				if tt.input == "" {
					assert.Equal(t, json.RawMessage("null"), result)
				} else {
					assert.True(t, IsValidJSON(result))
				}
			}
		})
	}
}

func TestConvertBytesToRawMessage(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectError bool
		errorType   string
	}{
		{
			name:        "valid JSON bytes",
			input:       []byte(`{"test":"data"}`),
			expectError: false,
		},
		{
			name:        "empty bytes",
			input:       []byte{},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			input:       []byte(`{"test":}`),
			expectError: true,
			errorType:   "conversion",
		},
		{
			name:        "large bytes",
			input:       []byte(strings.Repeat("x", MaxPayloadSize+1)),
			expectError: true,
			errorType:   "validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertBytesToRawMessage(tt.input)

			if tt.expectError {
				require.Error(t, err)
				var payloadErr *PayloadError
				require.ErrorAs(t, err, &payloadErr)
				assert.Equal(t, tt.errorType, payloadErr.Operation)
			} else {
				require.NoError(t, err)
				if len(tt.input) == 0 {
					assert.Equal(t, json.RawMessage("null"), result)
				} else {
					assert.True(t, IsValidJSON(result))
				}
			}
		})
	}
}
