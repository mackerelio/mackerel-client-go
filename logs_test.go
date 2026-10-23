package mackerel

import (
	"encoding/json"
	"fmt"
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
				TraceID:        ToPtr("0123456789abcdef0123456789abcdef"),
				ServiceName:    "shop",
				Attributes: []*LogAttribute{
					{Key: "k", Value: &LogAttributeValue{ValueType: "string", StringValue: "v"}},
				},
				ResourceAttributes: []*LogAttribute{},
				ScopeAttributes:    []*LogAttribute{},
			}},
			Pagination: &LogPagination{HasNextPage: true, StartCursor: ToPtr("c1"), EndCursor: ToPtr("c1")},
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
			ServiceNamespace: ToPtr("ns"),
			From:             100,
			To:               200,
			Keywords:         []string{"a,b", "c"},
			Severities:       []LogSeverity{LogSeverityWARN, LogSeverityERROR},
			TraceID:          ToPtr("0123456789abcdef0123456789abcdef"),
			BodyAttributes: []LogBodyAttributeFilter{
				{KeyPath: []string{"req", "method"}, Value: "GET", Operator: LogBodyAttributeOperatorEQ, Type: LogBodyAttributeValueTypeString},
			},
			Order:      &LogOrder{Column: ToPtr(LogOrderColumnTIMESTAMP), Direction: ToPtr(OrderDirectionASC)},
			PageCursor: &LogPageCursor{Direction: LogPageCursorDirectionNEXT, Cursor: "c1"},
			PerPage:    ToPtr(10),
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
