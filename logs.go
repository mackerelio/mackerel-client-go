package mackerel

import (
	"context"
	"iter"
	"slices"
	"time"
)

// LogSeverity is the severity of a log, mapped to OpenTelemetry SeverityNumber.
type LogSeverity string

// Log severities
const (
	LogSeverityUNSPECIFIED LogSeverity = "UNSPECIFIED"
	LogSeverityTRACE       LogSeverity = "TRACE"
	LogSeverityDEBUG       LogSeverity = "DEBUG"
	LogSeverityINFO        LogSeverity = "INFO"
	LogSeverityWARN        LogSeverity = "WARN"
	LogSeverityERROR       LogSeverity = "ERROR"
	LogSeverityFATAL       LogSeverity = "FATAL"
)

// LogBodyAttributeFilter is a condition on a field of a structured (JSON) log body.
type LogBodyAttributeFilter struct {
	KeyPath  []string                  `json:"keyPath"`
	Value    string                    `json:"value"`
	Operator LogBodyAttributeOperator  `json:"operator"`
	Type     LogBodyAttributeValueType `json:"type"`
}

// LogBodyAttributeOperator is an operator to compare a field of a log body.
type LogBodyAttributeOperator string

// Log body attribute operators
const (
	LogBodyAttributeOperatorEQ         LogBodyAttributeOperator = "EQ"
	LogBodyAttributeOperatorNEQ        LogBodyAttributeOperator = "NEQ"
	LogBodyAttributeOperatorGT         LogBodyAttributeOperator = "GT"
	LogBodyAttributeOperatorGTE        LogBodyAttributeOperator = "GTE"
	LogBodyAttributeOperatorLT         LogBodyAttributeOperator = "LT"
	LogBodyAttributeOperatorLTE        LogBodyAttributeOperator = "LTE"
	LogBodyAttributeOperatorSTARTSWITH LogBodyAttributeOperator = "STARTS_WITH"
)

// LogBodyAttributeValueType is the type of a field of a log body to compare as.
type LogBodyAttributeValueType string

// Log body attribute value types
const (
	LogBodyAttributeValueTypeString LogBodyAttributeValueType = "string"
	LogBodyAttributeValueTypeInt    LogBodyAttributeValueType = "int"
	LogBodyAttributeValueTypeDouble LogBodyAttributeValueType = "double"
	LogBodyAttributeValueTypeBool   LogBodyAttributeValueType = "bool"
)

// LogOrderColumn is a column to sort logs by.
type LogOrderColumn string

// Log order columns
const (
	LogOrderColumnTIMESTAMP LogOrderColumn = "TIMESTAMP"
)

// LogOrder is the sort condition of logs.
type LogOrder struct {
	Column    *LogOrderColumn `json:"column,omitzero"`
	Direction *OrderDirection `json:"direction,omitzero"`
}

// LogPageCursorDirection is the direction to read logs from a cursor.
type LogPageCursorDirection string

// Log page cursor directions
const (
	LogPageCursorDirectionNEXT     LogPageCursorDirection = "NEXT"
	LogPageCursorDirectionPREVIOUS LogPageCursorDirection = "PREVIOUS"
)

// LogPageCursor is the starting point of paging. The log at the cursor is excluded.
type LogPageCursor struct {
	Direction LogPageCursorDirection `json:"direction"`
	Cursor    string                 `json:"cursor"`
}

// ListLogsParam is the request of [Client.ListLogs].
type ListLogsParam struct {
	ServiceName      string  `json:"serviceName"`
	ServiceNamespace *string `json:"serviceNamespace,omitzero"`
	From             int64   `json:"from"`
	To               int64   `json:"to"`
	// Slices use omitempty instead of omitzero because the API rejects empty arrays.
	Keywords       []string                 `json:"keywords,omitempty"`
	Severities     []LogSeverity            `json:"severities,omitempty"`
	TraceID        *string                  `json:"traceId,omitzero"`
	BodyAttributes []LogBodyAttributeFilter `json:"bodyAttributes,omitempty"`
	Order          *LogOrder                `json:"order,omitzero"`
	PageCursor     *LogPageCursor           `json:"pageCursor,omitzero"`
	PerPage        *int                     `json:"perPage,omitzero"`
}

// ListLogsResponse is the response of [Client.ListLogs].
type ListLogsResponse struct {
	Results    []*LogEntry    `json:"results"`
	Pagination *LogPagination `json:"pagination"`
}

// LogEntry is a log record.
type LogEntry struct {
	Cursor             string          `json:"cursor"`
	Timestamp          time.Time       `json:"timestamp"`
	Severity           LogSeverity     `json:"severity"`
	SeverityNumber     int             `json:"severityNumber"`
	Body               string          `json:"body"`
	TraceID            *string         `json:"traceId"`
	SpanID             *string         `json:"spanId"`
	ServiceName        string          `json:"serviceName"`
	ServiceNamespace   *string         `json:"serviceNamespace"`
	Attributes         []*LogAttribute `json:"attributes"`
	ResourceAttributes []*LogAttribute `json:"resourceAttributes"`
	ScopeAttributes    []*LogAttribute `json:"scopeAttributes"`
}

// LogAttribute is an attribute of a log record, its resource or its scope.
type LogAttribute struct {
	Key   string             `json:"key"`
	Value *LogAttributeValue `json:"value"`
}

// LogAttributeValue is the value of [LogAttribute]. Currently ValueType is always "string".
type LogAttributeValue struct {
	ValueType   string `json:"valueType"`
	StringValue string `json:"stringValue"`
}

// LogPagination is the paging information of [ListLogsResponse].
type LogPagination struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}

// SavedLogSearchDefinition is the search condition of a saved log search.
type SavedLogSearchDefinition struct {
	ServiceName      string                   `json:"serviceName"`
	ServiceNamespace *string                  `json:"serviceNamespace"`
	Keywords         []string                 `json:"keywords"`
	Severities       []LogSeverity            `json:"severities"`
	TraceID          *string                  `json:"traceId"`
	BodyAttributes   []LogBodyAttributeFilter `json:"bodyAttributes"`
}

// SavedLogSearch is a saved log search.
type SavedLogSearch struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Description *string                   `json:"description"`
	Definition  *SavedLogSearchDefinition `json:"definition"`
	CreatedAt   int64                     `json:"createdAt"`
	UpdatedAt   int64                     `json:"updatedAt"`
}

// FindSavedLogSearchesResponse is the response of [Client.FindSavedLogSearches].
type FindSavedLogSearchesResponse struct {
	Results []*SavedLogSearch `json:"results"`
}

// ListLogs searches logs
func (c *Client) ListLogs(params *ListLogsParam) (*ListLogsResponse, error) {
	return c.ListLogsContext(context.Background(), params)
}

// ListLogsContext is like [ListLogs].
func (c *Client) ListLogsContext(ctx context.Context, params *ListLogsParam) (*ListLogsResponse, error) {
	return requestPostContext[ListLogsResponse](ctx, c, "/api/v0/logs", params)
}

// ListLogsSeq returns an iterator over logs that are selected by params.
// When it reaches the end of the logs in a response, it automatically fetches
// the subsequent page from the API using the same params except [ListLogsParam.PageCursor],
// then the iterator continues its results.
func (c *Client) ListLogsSeq(ctx context.Context, params *ListLogsParam) iter.Seq2[*LogEntry, error] {
	return func(yield func(*LogEntry, error) bool) {
		direction := LogPageCursorDirectionNEXT
		if params.PageCursor != nil {
			direction = params.PageCursor.Direction
		}
		params := *params
		for {
			res, err := c.ListLogsContext(ctx, &params)
			if err != nil {
				if !yield(nil, err) {
					return
				}
				if ctx.Err() != nil {
					return
				}
				continue
			}
			results := res.Results
			if direction == LogPageCursorDirectionPREVIOUS {
				// The API returns a PREVIOUS page in the sort order, so reverse it to keep moving away from the cursor.
				results = slices.Clone(results)
				slices.Reverse(results)
			}
			for _, r := range results {
				if !yield(r, nil) {
					return
				}
			}
			if res.Pagination == nil {
				break
			}
			hasMore, cursor := res.Pagination.HasNextPage, res.Pagination.EndCursor
			if direction == LogPageCursorDirectionPREVIOUS {
				hasMore, cursor = res.Pagination.HasPreviousPage, res.Pagination.StartCursor
			}
			if !hasMore || cursor == nil {
				break
			}
			params.PageCursor = &LogPageCursor{Direction: direction, Cursor: *cursor}
		}
	}
}

// FindSavedLogSearches lists saved log searches
func (c *Client) FindSavedLogSearches() (*FindSavedLogSearchesResponse, error) {
	return c.FindSavedLogSearchesContext(context.Background())
}

// FindSavedLogSearchesContext is like [FindSavedLogSearches].
func (c *Client) FindSavedLogSearchesContext(ctx context.Context) (*FindSavedLogSearchesResponse, error) {
	return requestGetContext[FindSavedLogSearchesResponse](ctx, c, "/api/v0/saved-log-searches")
}
