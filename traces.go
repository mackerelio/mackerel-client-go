package mackerel

import (
	"context"
	"fmt"
	"time"
)

// TraceResponse represents the response structure from the traces API
type TraceResponse struct {
	Spans []*Span `json:"spans"`
}

// SpanKind represents the kind of span
type SpanKind string

// StatusCode represents the status code of a span
type StatusCode string

// Span represents a single span in a trace
type Span struct {
	TraceID                string       `json:"traceId"`
	SpanID                 string       `json:"spanId"`
	TraceState             string       `json:"traceState"`
	ParentSpanID           string       `json:"parentSpanId,omitempty"`
	Name                   string       `json:"name"`
	Kind                   SpanKind     `json:"kind"`
	StartTime              time.Time    `json:"startTime"`
	EndTime                time.Time    `json:"endTime"`
	Attributes             []*Attribute `json:"attributes"`
	DroppedAttributesCount int          `json:"droppedAttributesCount"`
	Events                 []*Event     `json:"events"`
	DroppedEventsCount     int          `json:"droppedEventsCount"`
	Links                  []*Link      `json:"links"`
	DroppedLinksCount      int          `json:"droppedLinksCount"`
	Status                 *Status      `json:"status"`
	Resource               *Resource    `json:"resource"`
	Scope                  *Scope       `json:"scope"`
}

// Attribute represents a span attribute
type Attribute struct {
	Key   string          `json:"key"`
	Value *AttributeValue `json:"value"`
}

// AttributeValue represents a value that can be of different types
type AttributeValue struct {
	ValueType   string                     `json:"valueType"`
	StringValue string                     `json:"stringValue,omitempty"`
	BoolValue   bool                       `json:"boolValue,omitempty"`
	IntValue    int64                      `json:"intValue,omitempty"`
	DoubleValue float64                    `json:"doubleValue,omitempty"`
	ArrayValue  []*AttributeValue          `json:"arrayValue,omitempty"`
	KvlistValue map[string]*AttributeValue `json:"kvlistValue,omitempty"`
	BytesValue  []byte                     `json:"bytesValue,omitempty"`
}

// Event represents a span event
type Event struct {
	Time                   time.Time    `json:"time"`
	Name                   string       `json:"name"`
	Attributes             []*Attribute `json:"attributes"`
	DroppedAttributesCount int          `json:"droppedAttributesCount"`
}

// Link represents a link to another span
type Link struct {
	TraceID                string       `json:"traceId"`
	SpanID                 string       `json:"spanId"`
	TraceState             string       `json:"traceState"`
	Attributes             []*Attribute `json:"attributes"`
	DroppedAttributesCount int          `json:"droppedAttributesCount"`
}

// Status represents the execution state of a span
type Status struct {
	Message string     `json:"message"`
	Code    StatusCode `json:"code"`
}

// Resource represents resource information
type Resource struct {
	Attributes             []*Attribute `json:"attributes"`
	DroppedAttributesCount int          `json:"droppedAttributesCount"`
}

// Scope represents scope information
type Scope struct {
	Name                   string       `json:"name"`
	Version                string       `json:"version"`
	Attributes             []*Attribute `json:"attributes"`
	DroppedAttributesCount int          `json:"droppedAttributesCount"`
}

// Span kind constants
const (
	SpanKindUnspecified SpanKind = "unspecified"
	SpanKindInternal    SpanKind = "internal"
	SpanKindServer      SpanKind = "server"
	SpanKindClient      SpanKind = "client"
	SpanKindProducer    SpanKind = "producer"
	SpanKindConsumer    SpanKind = "consumer"
)

// Status code constants
const (
	StatusCodeUnset StatusCode = "unset"
	StatusCodeOK    StatusCode = "ok"
	StatusCodeError StatusCode = "error"
)

// AttributeValue type constants
const (
	ValueTypeString = "string"
	ValueTypeBool   = "bool"
	ValueTypeInt    = "int"
	ValueTypeDouble = "double"
	ValueTypeArray  = "array"
	ValueTypeKvlist = "kvlist"
	ValueTypeBytes  = "bytes"
	ValueTypeEmpty  = "empty"
)

// GetTrace gets detailed trace information for the specified trace ID
func (c *Client) GetTrace(traceID string) (*TraceResponse, error) {
	return requestGetContext[TraceResponse](context.Background(), c, fmt.Sprintf("/api/v0/traces/%s", traceID))
}

// GetTraceContext is like [GetTrace].
func (c *Client) GetTraceContext(ctx context.Context, traceID string) (*TraceResponse, error) {
	return requestGetContext[TraceResponse](ctx, c, fmt.Sprintf("/api/v0/traces/%s", traceID))
}
