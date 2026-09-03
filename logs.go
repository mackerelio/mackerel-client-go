package mackerel

import (
	"context"
	"time"
)

type FindLogsParam struct {
	ServiceName      string                   `json:"serviceName"`
	ServiceNamespace *string                  `json:"serviceNamespace,omitzero"`
	From             time.Time                `json:"from"`
	To               time.Time                `json:"to"`
	Timezone         *string                  `json:"timezone,omitzero"`
	Keywords         []string                 `json:"keywords,omitzero"`
	Severities       []LogSeverity            `json:"severities,omitzero"`
	Attributes       []LogAttributeComparison `json:"attributes,omitzero"`
	TraceID          *string                  `json:"traceId,omitzero"`
	Order            LogOrder                 `json:"order,omitzero"`
	First            *int                     `json:"first,omitzero"`
	After            *string                  `json:"after,omitzero"`
	Last             *int                     `json:"last,omitzero"`
	Before           *string                  `json:"before,omitzero"`
}

type LogSeverity string

const (
	LogSeverityUnspecified LogSeverity = "UNSPECIFIED"
	LogSeverityTrace       LogSeverity = "TRACE"
	LogSeverityDebug       LogSeverity = "DEBUG"
	LogSeverityInfo        LogSeverity = "INFO"
	LogSeverityWarn        LogSeverity = "WARN"
	LogSeverityError       LogSeverity = "ERROR"
	LogSeverityFatal       LogSeverity = "FATAL"
)

type LogAttributeComparison struct {
	Key         string                        `json:"key"`
	Value       *LogAttributeStringComparison `json:"value,omitzero"`
	ValueInt    *LogAttributeIntComparison    `json:"valueInt,omitzero"`
	ValueDouble *LogAttributeDoubleComparison `json:"valueDouble,omitzero"`
	ValueBool   *LogAttributeBoolComparison   `json:"valueBool,omitzero"`
}

type LogAttributeStringComparison struct {
	Value    string                      `json:"value"`
	Operator LogStringComparisonOperator `json:"operator"`
}

type LogAttributeIntComparison struct {
	ValueInt int64                 `json:"valueInt"`
	Operator LogComparisonOperator `json:"operator"`
}

type LogAttributeDoubleComparison struct {
	ValueDouble float64               `json:"valueDouble"`
	Operator    LogComparisonOperator `json:"operator"`
}

type LogAttributeBoolComparison struct {
	ValueBool bool                    `json:"valueBool"`
	Operator  LogEqComparisonOperator `json:"operator"`
}

type LogStringComparisonOperator string

const (
	LogStringComparisonOperatorEQ     LogStringComparisonOperator = "EQ"
	LogStringComparisonOperatorNEQ    LogStringComparisonOperator = "NEQ"
	LogStringComparisonOperatorPREFIX LogStringComparisonOperator = "PREFIX"
)

type LogComparisonOperator string

const (
	LogComparisonOperatorEQ  LogComparisonOperator = "EQ"
	LogComparisonOperatorGT  LogComparisonOperator = "GT"
	LogComparisonOperatorGTE LogComparisonOperator = "GTE"
	LogComparisonOperatorLT  LogComparisonOperator = "LT"
	LogComparisonOperatorLTE LogComparisonOperator = "LTE"
)

type LogEqComparisonOperator string

const (
	LogEqComparisonOperatorEQ  LogEqComparisonOperator = "EQ"
	LogEqComparisonOperatorNEQ LogEqComparisonOperator = "NEQ"
)

type LogOrder struct {
	Column    *LogOrderColumn `json:"column"`
	Direction *OrderDirection `json:"direction"`
}

type LogOrderColumn string

const (
	LogOrderColumnTIMESTAMP LogOrderColumn = "TIMESTAMP"
)

type FindLogsResponse struct {
	Results  []*SimpleLog `json:"results"`
	PageInfo LogPageInfo  `json:"pageInfo"`
}

type LogPageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor,omitempty"`
	EndCursor       *string `json:"endCursor,omitempty"`
}

// SimpleLog represents a single log entry returned by the log search API
type SimpleLog struct {
	Cursor             string          `json:"cursor"`
	Timestamp          time.Time       `json:"timestamp"`
	EffectiveTimestamp time.Time       `json:"effectiveTimestamp"`
	Severity           LogSeverity     `json:"severity"`
	SeverityText       string          `json:"severityText"`
	SeverityNumber     int             `json:"severityNumber"`
	Body               string          `json:"body"`
	TraceID            *string         `json:"traceId,omitempty"`
	SpanID             *string         `json:"spanId,omitempty"`
	ServiceName        string          `json:"serviceName"`
	ServiceNamespace   string          `json:"serviceNamespace"`
	Attributes         []*LogAttribute `json:"attributes"`
	ResourceAttributes []*LogAttribute `json:"resourceAttributes"`
	ScopeAttributes    []*LogAttribute `json:"scopeAttributes"`
}

// LogAttribute represents a key-value attribute attached to a log entry
type LogAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// FindLogs searches logs
func (c *Client) FindLogs(params *FindLogsParam) (*FindLogsResponse, error) {
	return requestPostContext[FindLogsResponse](context.Background(), c, "/api/v0/logs", params)
}

// FindLogsContext is like [FindLogs].
func (c *Client) FindLogsContext(ctx context.Context, params *FindLogsParam) (*FindLogsResponse, error) {
	return requestPostContext[FindLogsResponse](ctx, c, "/api/v0/logs", params)
}
