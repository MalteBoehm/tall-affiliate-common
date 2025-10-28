package cloudevents

import (
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

// Event is a minimal CloudEvents 1.0 compatible envelope (no runtime deps)
type Event struct {
    ID              string          `json:"id"`
    Source          string          `json:"source"`
    Type            string          `json:"type"`
    SpecVersion     string          `json:"specversion"`
    Subject         string          `json:"subject"`
    Time            string          `json:"time"`
    DataContentType string          `json:"datacontenttype"`
    Data            json.RawMessage `json:"data"`
    TraceID         string          `json:"trace_id,omitempty"`
    CorrelationID   string          `json:"correlation_id,omitempty"`
    Tenant          string          `json:"tenant,omitempty"`
}

// New builds a new CloudEvent envelope
func New(source, subject, eventType string, data any, opts ...Option) (*Event, error) {
    payload, err := json.Marshal(data)
    if err != nil { return nil, err }
    e := &Event{
        ID:              uuid.New().String(),
        Source:          source,
        Type:            eventType,
        SpecVersion:     "1.0",
        Subject:         subject,
        Time:            time.Now().UTC().Format(time.RFC3339),
        DataContentType: "application/json",
        Data:            payload,
    }
    for _, o := range opts { o(e) }
    return e, nil
}

// Option functional option
type Option func(*Event)

func WithTrace(traceID string) Option       { return func(e *Event) { e.TraceID = traceID } }
func WithCorrelation(id string) Option      { return func(e *Event) { e.CorrelationID = id } }
func WithTenant(tenant string) Option       { return func(e *Event) { e.Tenant = tenant } }
func WithID(id string) Option               { return func(e *Event) { e.ID = id } }
func WithTime(t time.Time) Option           { return func(e *Event) { e.Time = t.UTC().Format(time.RFC3339) } }
func WithContentType(ct string) Option      { return func(e *Event) { e.DataContentType = ct } }

// Headers returns a JSON-serializable headers map for storage
func (e *Event) Headers() map[string]any {
    h := map[string]any{
        "specversion":     e.SpecVersion,
        "source":          e.Source,
        "subject":         e.Subject,
        "time":            e.Time,
        "datacontenttype": e.DataContentType,
    }
    if e.TraceID != "" { h["trace_id"] = e.TraceID }
    if e.CorrelationID != "" { h["correlation_id"] = e.CorrelationID }
    if e.Tenant != "" { h["tenant"] = e.Tenant }
    return h
}
