package mackerel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestListLogs(t *testing.T) {
	var body map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/logs" {
			t.Error("request URL should be /api/v0/logs but: ", req.URL.Path)
		}
		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}
		body = nil
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		respJSON, _ := json.Marshal(map[string]any{
			"results": []map[string]any{
				{
					"cursor":           "c1",
					"timestamp":        "2026-05-01T12:34:56.789Z",
					"severity":         "ERROR",
					"severityNumber":   17,
					"body":             "something failed",
					"traceId":          "0123456789abcdef0123456789abcdef",
					"spanId":           nil,
					"serviceName":      "shop",
					"serviceNamespace": nil,
					"attributes": []map[string]any{
						{"key": "k", "value": map[string]any{"valueType": "string", "stringValue": "v"}},
					},
					"resourceAttributes": []map[string]any{},
					"scopeAttributes":    []map[string]any{},
				},
			},
			"pagination": map[string]any{
				"hasNextPage":     true,
				"hasPreviousPage": false,
				"startCursor":     "c1",
				"endCursor":       "c1",
			},
		})
		res.Header().Set("Content-Type", "application/json")
		res.Write(respJSON) // nolint
	}))
	defer ts.Close()
	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	t.Run("required only", func(t *testing.T) {
		got, err := client.ListLogs(&ListLogsParam{ServiceName: "shop", From: 100, To: 200})
		if err != nil {
			t.Fatal(err)
		}
		wantBody := map[string]any{"serviceName": "shop", "from": float64(100), "to": float64(200)}
		if !reflect.DeepEqual(wantBody, body) {
			t.Errorf("request body should be %v but: %v", wantBody, body)
		}
		want := &ListLogsResponse{
			Results: []*LogEntry{{
				Cursor:         "c1",
				Timestamp:      time.Date(2026, 5, 1, 12, 34, 56, 789000000, time.UTC),
				Severity:       LogSeverityERROR,
				SeverityNumber: 17,
				Body:           "something failed",
				TraceID:        new("0123456789abcdef0123456789abcdef"),
				ServiceName:    "shop",
				Attributes: []*LogAttribute{
					{Key: "k", Value: &LogAttributeValue{ValueType: "string", StringValue: "v"}},
				},
				ResourceAttributes: []*LogAttribute{},
				ScopeAttributes:    []*LogAttribute{},
			}},
			Pagination: &LogPagination{HasNextPage: true, StartCursor: new("c1"), EndCursor: new("c1")},
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("the response should equal to %v but: %v", want, got)
		}
	})

	t.Run("empty slices are omitted", func(t *testing.T) {
		_, err := client.ListLogs(&ListLogsParam{
			ServiceName: "shop", From: 100, To: 200,
			Keywords: []string{}, Severities: []LogSeverity{}, BodyAttributes: []LogBodyAttributeFilter{},
		})
		if err != nil {
			t.Fatal(err)
		}
		for _, k := range []string{"keywords", "severities", "bodyAttributes"} {
			if _, ok := body[k]; ok {
				t.Errorf("%s should not be sent", k)
			}
		}
	})

	t.Run("all params", func(t *testing.T) {
		_, err := client.ListLogs(&ListLogsParam{
			ServiceName:      "shop",
			ServiceNamespace: new("ns"),
			From:             100,
			To:               200,
			Keywords:         []string{"a,b", "c"},
			Severities:       []LogSeverity{LogSeverityWARN, LogSeverityERROR},
			TraceID:          new("0123456789abcdef0123456789abcdef"),
			BodyAttributes: []LogBodyAttributeFilter{
				{KeyPath: []string{"req", "method"}, Value: "GET", Operator: LogBodyAttributeOperatorEQ, Type: LogBodyAttributeValueTypeString},
			},
			Order:      &LogOrder{Column: new(LogOrderColumnTIMESTAMP), Direction: new(OrderDirectionASC)},
			PageCursor: &LogPageCursor{Direction: LogPageCursorDirectionNEXT, Cursor: "c1"},
			PerPage:    new(10),
		})
		if err != nil {
			t.Fatal(err)
		}
		wantBody := map[string]any{
			"serviceName":      "shop",
			"serviceNamespace": "ns",
			"from":             float64(100),
			"to":               float64(200),
			"keywords":         []any{"a,b", "c"},
			"severities":       []any{"WARN", "ERROR"},
			"traceId":          "0123456789abcdef0123456789abcdef",
			"bodyAttributes": []any{
				map[string]any{"keyPath": []any{"req", "method"}, "value": "GET", "operator": "EQ", "type": "string"},
			},
			"order":      map[string]any{"column": "TIMESTAMP", "direction": "ASC"},
			"pageCursor": map[string]any{"direction": "NEXT", "cursor": "c1"},
			"perPage":    float64(10),
		}
		if !reflect.DeepEqual(wantBody, body) {
			t.Errorf("request body should be %v but: %v", wantBody, body)
		}
	})
}

func TestFindSavedLogSearches(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/saved-log-searches" {
			t.Error("request URL should be /api/v0/saved-log-searches but: ", req.URL.Path)
		}
		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}
		respJSON, _ := json.Marshal(map[string]any{
			"results": []map[string]any{
				{
					"id":          "s1",
					"name":        "errors",
					"description": nil,
					"definition": map[string]any{
						"serviceName":      "shop",
						"serviceNamespace": nil,
						"keywords":         []string{"timeout"},
						"severities":       []string{"ERROR"},
						"traceId":          nil,
						"bodyAttributes": []map[string]any{
							{"keyPath": []string{"req", "status"}, "value": "500", "operator": "GTE", "type": "int"},
						},
					},
					"createdAt": 1700000000,
					"updatedAt": 1700000100,
				},
			},
		})
		res.Header().Set("Content-Type", "application/json")
		res.Write(respJSON) // nolint
	}))
	defer ts.Close()
	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	got, err := client.FindSavedLogSearches()
	if err != nil {
		t.Fatal(err)
	}
	want := &FindSavedLogSearchesResponse{Results: []*SavedLogSearch{{
		ID:   "s1",
		Name: "errors",
		Definition: &SavedLogSearchDefinition{
			ServiceName: "shop",
			Keywords:    []string{"timeout"},
			Severities:  []LogSeverity{LogSeverityERROR},
			BodyAttributes: []LogBodyAttributeFilter{
				{KeyPath: []string{"req", "status"}, Value: "500", Operator: LogBodyAttributeOperatorGTE, Type: LogBodyAttributeValueTypeInt},
			},
		},
		CreatedAt: 1700000000,
		UpdatedAt: 1700000100,
	}}}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("the response should equal to %v but: %v", want, got)
	}
}

func TestListLogs_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(429)
		fmt.Fprint(res, `{"error": {"message": "Too Many Requests"}}`) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	_, err := client.ListLogs(&ListLogsParam{ServiceName: "shop", From: 1, To: 2})
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error should be APIError but: %#v", err)
	}
	if apiErr.StatusCode != 429 || apiErr.Message != "Too Many Requests" {
		t.Errorf("unexpected APIError: %#v", apiErr)
	}
}

func TestLogAttributeValue_keepsEmptyString(t *testing.T) {
	b, err := json.Marshal(&LogAttributeValue{ValueType: "string", StringValue: ""})
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"valueType":"string","stringValue":""}`; string(b) != want {
		t.Errorf("marshaled value should be %s but: %s", want, b)
	}
}

// newLogsPagingServer serves pages keyed by "<direction>:<cursor>" ("" for no cursor),
// and records the keys of the requests.
func newLogsPagingServer(t *testing.T, pages map[string]*ListLogsResponse, requested *[]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var p ListLogsParam
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			t.Error(err)
		}
		key := ""
		if p.PageCursor != nil {
			key = string(p.PageCursor.Direction) + ":" + p.PageCursor.Cursor
		}
		*requested = append(*requested, key)
		page, ok := pages[key]
		if !ok {
			t.Errorf("unexpected page cursor: %q", key)
			res.WriteHeader(400)
			return
		}
		respJSON, _ := json.Marshal(page)
		res.Header().Set("Content-Type", "application/json")
		res.Write(respJSON) // nolint
	}))
}

func collectLogs(t *testing.T, seq func(func(*LogEntry, error) bool)) []string {
	var got []string
	for r, err := range seq {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, r.Cursor)
	}
	return got
}

func TestListLogsSeq(t *testing.T) {
	pages := map[string]*ListLogsResponse{
		"": {
			Results:    []*LogEntry{{Cursor: "c1"}, {Cursor: "c2"}},
			Pagination: &LogPagination{HasNextPage: true, StartCursor: new("c1"), EndCursor: new("c2")},
		},
		"NEXT:c2": {
			Results:    []*LogEntry{{Cursor: "c3"}},
			Pagination: &LogPagination{HasPreviousPage: true, StartCursor: new("c3"), EndCursor: new("c3")},
		},
	}
	var requested []string
	ts := newLogsPagingServer(t, pages, &requested)
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	got := collectLogs(t, client.ListLogsSeq(t.Context(), &ListLogsParam{ServiceName: "shop", From: 1, To: 2}))
	if want := []string{"c1", "c2", "c3"}; !reflect.DeepEqual(want, got) {
		t.Errorf("the logs should be %v but: %v", want, got)
	}
	if want := []string{"", "NEXT:c2"}; !reflect.DeepEqual(want, requested) {
		t.Errorf("the requested cursors should be %v but: %v", want, requested)
	}
}

func TestListLogsSeqPrevious(t *testing.T) {
	pages := map[string]*ListLogsResponse{
		"PREVIOUS:c10": {
			Results:    []*LogEntry{{Cursor: "c8"}, {Cursor: "c9"}},
			Pagination: &LogPagination{HasNextPage: true, HasPreviousPage: true, StartCursor: new("c8"), EndCursor: new("c9")},
		},
		"PREVIOUS:c8": {
			Results:    []*LogEntry{{Cursor: "c7"}},
			Pagination: &LogPagination{HasNextPage: true, StartCursor: new("c7"), EndCursor: new("c7")},
		},
	}
	var requested []string
	ts := newLogsPagingServer(t, pages, &requested)
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	got := collectLogs(t, client.ListLogsSeq(t.Context(), &ListLogsParam{
		ServiceName: "shop", From: 1, To: 2,
		PageCursor: &LogPageCursor{Direction: LogPageCursorDirectionPREVIOUS, Cursor: "c10"},
	}))
	if want := []string{"c9", "c8", "c7"}; !reflect.DeepEqual(want, got) {
		t.Errorf("the logs should be %v but: %v", want, got)
	}
	if want := []string{"PREVIOUS:c10", "PREVIOUS:c8"}; !reflect.DeepEqual(want, requested) {
		t.Errorf("the requested cursors should be %v but: %v", want, requested)
	}
}

func TestListLogsSeqStopsOnBreak(t *testing.T) {
	pages := map[string]*ListLogsResponse{
		"": {
			Results:    []*LogEntry{{Cursor: "c1"}, {Cursor: "c2"}},
			Pagination: &LogPagination{HasNextPage: true, StartCursor: new("c1"), EndCursor: new("c2")},
		},
	}
	var requested []string
	ts := newLogsPagingServer(t, pages, &requested)
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	for r, err := range client.ListLogsSeq(t.Context(), &ListLogsParam{ServiceName: "shop", From: 1, To: 2}) {
		if err != nil {
			t.Fatal(err)
		}
		if r.Cursor == "c1" {
			break
		}
	}
	if len(requested) != 1 {
		t.Errorf("the iterator should not fetch the next page after break but requested: %v", requested)
	}
}

func TestListLogsSeqRetriesOnError(t *testing.T) {
	var bodies []string
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		b, _ := io.ReadAll(req.Body)
		bodies = append(bodies, string(b))
		if len(bodies) == 1 {
			res.WriteHeader(429)
			fmt.Fprint(res, `{"error": {"message": "Too Many Requests"}}`) // nolint
			return
		}
		respJSON, _ := json.Marshal(&ListLogsResponse{
			Results:    []*LogEntry{{Cursor: "c1"}},
			Pagination: &LogPagination{StartCursor: new("c1"), EndCursor: new("c1")},
		})
		res.Header().Set("Content-Type", "application/json")
		res.Write(respJSON) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	var errs []error
	var got []string
	for r, err := range client.ListLogsSeq(t.Context(), &ListLogsParam{ServiceName: "shop", From: 1, To: 2}) {
		if err != nil {
			errs = append(errs, err)
			continue
		}
		got = append(got, r.Cursor)
	}
	if len(errs) != 1 {
		t.Fatal("the iterator should report exactly 1 error but: ", errs)
	}
	if want := []string{"c1"}; !reflect.DeepEqual(want, got) {
		t.Errorf("the logs should be %v but: %v", want, got)
	}
	if len(bodies) != 2 || bodies[0] != bodies[1] {
		t.Errorf("the iterator should retry with the same params but: %v", bodies)
	}
}

func TestListLogsSeqCanceledContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		respJSON, _ := json.Marshal(&ListLogsResponse{Pagination: &LogPagination{HasNextPage: true, EndCursor: new("c1")}})
		res.Header().Set("Content-Type", "application/json")
		res.Write(respJSON) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	var errs []error
	for _, err := range client.ListLogsSeq(ctx, &ListLogsParam{ServiceName: "shop", From: 1, To: 2}) {
		if errs = append(errs, err); len(errs) > 1 {
			t.Fatal("the iterator should stop on a canceled context but it retried: ", len(errs))
		}
	}
	if len(errs) != 1 {
		t.Fatal("the iterator should report exactly 1 error but: ", errs)
	}
	if !errors.Is(errs[0], context.Canceled) {
		t.Error("the error should be context.Canceled but: ", errs[0])
	}
}
