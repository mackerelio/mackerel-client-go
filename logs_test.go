package mackerel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestFindLogs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/logs" {
			t.Error("request URL should be /api/v0/logs but: ", req.URL.Path)
		}
		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]any{
			"results": []map[string]any{
				{
					"cursor":             "eyJ0cyI6MTIzfQ",
					"timestamp":          "2026-09-01T00:12:34.567Z",
					"effectiveTimestamp": "2026-09-01T00:12:34.567Z",
					"severity":           "ERROR",
					"severityText":       "Error",
					"severityNumber":     17,
					"body":               "connection timeout",
					"traceId":            "550e8400e29b41d4a716446655440000",
					"spanId":             "051581bf3cb55c13",
					"serviceName":        "shoppingcart",
					"serviceNamespace":   "shop",
					"attributes": []map[string]any{
						{"key": "http.method", "value": "GET"},
					},
					"resourceAttributes": []map[string]any{
						{"key": "host.name", "value": "web01"},
					},
					"scopeAttributes": []map[string]any{},
				},
			},
			"pageInfo": map[string]any{
				"hasNextPage":     true,
				"hasPreviousPage": false,
				"startCursor":     "eyJ0cyI6MTIzfQ",
				"endCursor":       "eyJ0cyI6NDU2fQ",
			},
		})

		res.Header().Set("Content-Type", "application/json")
		res.Write(respJSON) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	from := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 9, 1, 1, 0, 0, 0, time.UTC)
	got, err := client.FindLogs(&FindLogsParam{
		ServiceName: "shoppingcart",
		From:        from,
		To:          to,
	})
	if err != nil {
		t.Fatal(err)
	}

	traceID := "550e8400e29b41d4a716446655440000"
	spanID := "051581bf3cb55c13"
	startCursor := "eyJ0cyI6MTIzfQ"
	endCursor := "eyJ0cyI6NDU2fQ"
	timestamp := time.Date(2026, 9, 1, 0, 12, 34, 567000000, time.UTC)

	want := &FindLogsResponse{
		Results: []*SimpleLog{
			{
				Cursor:             "eyJ0cyI6MTIzfQ",
				Timestamp:          timestamp,
				EffectiveTimestamp: timestamp,
				Severity:           LogSeverityError,
				SeverityText:       "Error",
				SeverityNumber:     17,
				Body:               "connection timeout",
				TraceID:            &traceID,
				SpanID:             &spanID,
				ServiceName:        "shoppingcart",
				ServiceNamespace:   "shop",
				Attributes: []*LogAttribute{
					{Key: "http.method", Value: "GET"},
				},
				ResourceAttributes: []*LogAttribute{
					{Key: "host.name", Value: "web01"},
				},
				ScopeAttributes: []*LogAttribute{},
			},
		},
		PageInfo: LogPageInfo{
			HasNextPage:     true,
			HasPreviousPage: false,
			StartCursor:     &startCursor,
			EndCursor:       &endCursor,
		},
	}

	if len(got.Results) != 1 {
		t.Fatalf("expected 1 result but got %d", len(got.Results))
	}
	if !got.Results[0].Timestamp.Equal(want.Results[0].Timestamp) {
		t.Errorf("Timestamp should equal to %v but got %v", want.Results[0].Timestamp, got.Results[0].Timestamp)
	}
	if !got.Results[0].EffectiveTimestamp.Equal(want.Results[0].EffectiveTimestamp) {
		t.Errorf("EffectiveTimestamp should equal to %v but got %v", want.Results[0].EffectiveTimestamp, got.Results[0].EffectiveTimestamp)
	}
	// Already verified via Equal() above; align them before the DeepEqual comparison below.
	got.Results[0].Timestamp = want.Results[0].Timestamp
	got.Results[0].EffectiveTimestamp = want.Results[0].EffectiveTimestamp
	if !reflect.DeepEqual(want, got) {
		t.Errorf("the response should equal to %v but got %v", want, got)
	}
}

func TestFindLogs_RequestBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}

		if body["serviceName"] != "shoppingcart" {
			t.Errorf("serviceName should be shoppingcart but: %v", body["serviceName"])
		}
		if body["traceId"] != "550e8400e29b41d4a716446655440000" {
			t.Errorf("traceId should be set but: %v", body["traceId"])
		}
		if got, want := body["keywords"], []any{"timeout", "connection"}; !reflect.DeepEqual(got, want) {
			t.Errorf("keywords should be %v but: %v", want, got)
		}
		if got, want := body["severities"], []any{"ERROR", "FATAL"}; !reflect.DeepEqual(got, want) {
			t.Errorf("severities should be %v but: %v", want, got)
		}
		if got, want := body["order"], map[string]any{"column": "TIMESTAMP", "direction": "ASC"}; !reflect.DeepEqual(got, want) {
			t.Errorf("order should be %v but: %v", want, got)
		}

		attrs, ok := body["attributes"].([]any)
		if !ok || len(attrs) != 4 {
			t.Fatalf("attributes should have 4 elements but: %v", body["attributes"])
		}

		attr0 := attrs[0].(map[string]any)
		if attr0["key"] != "http.method" {
			t.Errorf("attributes[0].key should be http.method but: %v", attr0["key"])
		}
		value, ok := attr0["value"].(map[string]any)
		if !ok {
			t.Fatalf("attributes[0].value should be set but: %v", attr0["value"])
		}
		if value["value"] != "GET" || value["operator"] != "EQ" {
			t.Errorf("attributes[0].value should be {GET, EQ} but: %v", value)
		}

		attr1 := attrs[1].(map[string]any)
		valueInt, ok := attr1["valueInt"].(map[string]any)
		if !ok {
			t.Fatalf("attributes[1].valueInt should be set but: %v", attr1["valueInt"])
		}
		if valueInt["operator"] != "GTE" {
			t.Errorf("attributes[1].valueInt.operator should be GTE but: %v", valueInt["operator"])
		}

		attr2 := attrs[2].(map[string]any)
		valueDouble, ok := attr2["valueDouble"].(map[string]any)
		if !ok {
			t.Fatalf("attributes[2].valueDouble should be set but: %v", attr2["valueDouble"])
		}
		if valueDouble["valueDouble"] != 1.5 || valueDouble["operator"] != "LT" {
			t.Errorf("attributes[2].valueDouble should be {1.5, LT} but: %v", valueDouble)
		}

		attr3 := attrs[3].(map[string]any)
		valueBool, ok := attr3["valueBool"].(map[string]any)
		if !ok {
			t.Fatalf("attributes[3].valueBool should be set but: %v", attr3["valueBool"])
		}
		if valueBool["valueBool"] != true || valueBool["operator"] != "EQ" {
			t.Errorf("attributes[3].valueBool should be {true, EQ} but: %v", valueBool)
		}

		if first, ok := body["first"].(float64); !ok || first != 50 {
			t.Errorf("first should be 50 but: %v", body["first"])
		}
		if _, ok := body["last"]; ok {
			t.Errorf("last should not be set but: %v", body["last"])
		}
		if _, ok := body["before"]; ok {
			t.Errorf("before should not be set but: %v", body["before"])
		}

		respJSON, _ := json.Marshal(map[string]any{
			"results":  []map[string]any{},
			"pageInfo": map[string]any{"hasNextPage": false, "hasPreviousPage": false},
		})
		res.Header().Set("Content-Type", "application/json")
		res.Write(respJSON) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	traceID := "550e8400e29b41d4a716446655440000"
	first := 50
	column := LogOrderColumnTIMESTAMP
	direction := OrderDirectionASC
	_, err := client.FindLogs(&FindLogsParam{
		ServiceName: "shoppingcart",
		From:        time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		To:          time.Date(2026, 9, 1, 1, 0, 0, 0, time.UTC),
		TraceID:     &traceID,
		First:       &first,
		Keywords:    []string{"timeout", "connection"},
		Severities:  []LogSeverity{LogSeverityError, LogSeverityFatal},
		Order:       LogOrder{Column: &column, Direction: &direction},
		Attributes: []LogAttributeComparison{
			{
				Key: "http.method",
				Value: &LogAttributeStringComparison{
					Value:    "GET",
					Operator: LogStringComparisonOperatorEQ,
				},
			},
			{
				Key: "http.status_code",
				ValueInt: &LogAttributeIntComparison{
					ValueInt: 500,
					Operator: LogComparisonOperatorGTE,
				},
			},
			{
				Key: "latency",
				ValueDouble: &LogAttributeDoubleComparison{
					ValueDouble: 1.5,
					Operator:    LogComparisonOperatorLT,
				},
			},
			{
				Key: "is_slow",
				ValueBool: &LogAttributeBoolComparison{
					ValueBool: true,
					Operator:  LogEqComparisonOperatorEQ,
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindLogs_RequestBody_BackwardPagination(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}

		if _, ok := body["first"]; ok {
			t.Errorf("first should not be set but: %v", body["first"])
		}
		if _, ok := body["after"]; ok {
			t.Errorf("after should not be set but: %v", body["after"])
		}
		if last, ok := body["last"].(float64); !ok || last != 20 {
			t.Errorf("last should be 20 but: %v", body["last"])
		}
		if body["before"] != "eyJ0cyI6MTIzfQ" {
			t.Errorf("before should be eyJ0cyI6MTIzfQ but: %v", body["before"])
		}

		respJSON, _ := json.Marshal(map[string]any{
			"results":  []map[string]any{},
			"pageInfo": map[string]any{"hasNextPage": false, "hasPreviousPage": false},
		})
		res.Header().Set("Content-Type", "application/json")
		res.Write(respJSON) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	last := 20
	before := "eyJ0cyI6MTIzfQ"
	_, err := client.FindLogs(&FindLogsParam{
		ServiceName: "shoppingcart",
		From:        time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		To:          time.Date(2026, 9, 1, 1, 0, 0, 0, time.UTC),
		Last:        &last,
		Before:      &before,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindLogs_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(`{"error":{"message":"invalid request"}}`)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	_, err := client.FindLogs(&FindLogsParam{
		ServiceName: "shoppingcart",
		From:        time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		To:          time.Date(2026, 9, 1, 1, 0, 0, 0, time.UTC),
	})
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}
