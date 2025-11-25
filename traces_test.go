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

func TestListTraces(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/traces" {
			t.Error("request URL should be /api/v0/traces but: ", req.URL.Path)
		}
		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		// Sample response based on API documentation
		respJSON, _ := json.Marshal(map[string]any{
			"results": []map[string]any{
				{
					"traceId":              "550e8400e29b41d4a716446655440000",
					"serviceName":          "shoppingcart",
					"serviceNamespace":     "shop",
					"environment":          "production",
					"title":                "GET /api/users",
					"traceStartAt":         1718802000,
					"traceLatencyMillis":   1234,
					"serviceStartAt":       1718802100,
					"serviceLatencyMillis": 567,
				},
			},
			"hasNextPage": true,
		})

		res.Header().Set("Content-Type", "application/json")
		res.Write(respJSON) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	got, err := client.ListTraces(&ListTracesParam{
		ServiceName: "shoppingcart",
		From:        1718801900,
		To:          1718802200,
	})
	if err != nil {
		t.Fatal(err)
	}

	want := &ListTracesResponse{
		Results: []*ListTracesResult{
			{
				TraceID:              "550e8400e29b41d4a716446655440000",
				ServiceName:          "shoppingcart",
				ServiceNamespace:     "shop",
				Environment:          "production",
				Title:                "GET /api/users",
				TraceStartAt:         1718802000,
				TraceLatencyMillis:   1234,
				ServiceStartAt:       1718802100,
				ServiceLatencyMillis: 567,
			},
		},
		HasNextPage: true,
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("the response should equal to %v", want)
	}
}

func TestListTracesSeq(t *testing.T) {
	pages := []*ListTracesResponse{
		{
			Results: []*ListTracesResult{
				{
					TraceID:              "550e8400e29b41d4a716446655440000",
					ServiceName:          "shoppingcart",
					Title:                "GET /api/users",
					TraceStartAt:         1718802000,
					TraceLatencyMillis:   1234,
					ServiceStartAt:       1718802100,
					ServiceLatencyMillis: 567,
				},
			},
			HasNextPage: true,
		},
		{
			Results: []*ListTracesResult{
				{
					TraceID:              "550e8400e29b41d4a716446655440000",
					ServiceName:          "authserver",
					Title:                "GET /api/users",
					TraceStartAt:         1718802010,
					TraceLatencyMillis:   1234,
					ServiceStartAt:       1718802120,
					ServiceLatencyMillis: 567,
				},
			},
			HasNextPage: false,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var r ListTracesParam
		if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
			t.Fatal(err)
		}
		page := 1
		if r.Page != nil {
			page = *r.Page
		}
		respJSON, _ := json.Marshal(pages[page-1])
		res.Header().Set("Content-Type", "application/json")
		res.Write(respJSON) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	traces := client.ListTracesSeq(t.Context(), &ListTracesParam{
		ServiceName: "shoppingcart",
		From:        1718801900,
		To:          1718802200,
	})
	var got []*ListTracesResult
	for r, err := range traces {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, r)
	}
	var want []*ListTracesResult
	for _, p := range pages {
		want = append(want, p.Results...)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("the response should equal to %v", want)
	}
}

func TestGetTrace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/traces/0123456789abcdef0123456789abcdef" {
			t.Error("request URL should be /api/v0/traces/0123456789abcdef0123456789abcdef but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		// Sample response based on API documentation
		respJSON, _ := json.Marshal(map[string]interface{}{
			"spans": []map[string]interface{}{
				{
					"traceId":    "0123456789abcdef0123456789abcdef",
					"spanId":     "0123456789abcdef",
					"traceState": "congo=xx,key=val",
					"name":       "test-span",
					"kind":       "internal",
					"startTime":  "2025-07-09T14:03:02.123Z",
					"endTime":    "2025-07-09T14:03:02.456Z",
					"attributes": []map[string]interface{}{
						{
							"key": "http.route",
							"value": map[string]interface{}{
								"valueType":   "string",
								"stringValue": "/",
							},
						},
					},
					"droppedAttributesCount": 0,
					"events": []map[string]interface{}{
						{
							"time":                   "2025-07-09T14:03:02.789Z",
							"name":                   "event1",
							"attributes":             []interface{}{},
							"droppedAttributesCount": 0,
						},
					},
					"droppedEventsCount": 0,
					"links": []map[string]interface{}{
						{
							"traceId":                "abcdef0123456789abcdef0123456789",
							"spanId":                 "abcdefabcdef0102",
							"traceState":             "",
							"attributes":             []interface{}{},
							"droppedAttributesCount": 0,
						},
					},
					"droppedLinksCount": 0,
					"status": map[string]interface{}{
						"message": "status message",
						"code":    "ok",
					},
					"resource": map[string]interface{}{
						"attributes":             []interface{}{},
						"droppedAttributesCount": 0,
					},
					"scope": map[string]interface{}{
						"name":                   "my-library",
						"version":                "1.0.0",
						"attributes":             []interface{}{},
						"droppedAttributesCount": 0,
					},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	trace, err := client.GetTrace("0123456789abcdef0123456789abcdef")

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if trace == nil {
		t.Error("trace should not be nil")
		return
	}

	if len(trace.Spans) != 1 {
		t.Error("trace should have 1 span but: ", len(trace.Spans))
		return
	}

	span := trace.Spans[0]
	if span.TraceID != "0123456789abcdef0123456789abcdef" {
		t.Error("span.TraceID should be 0123456789abcdef0123456789abcdef but: ", span.TraceID)
	}

	if span.SpanID != "0123456789abcdef" {
		t.Error("span.SpanID should be 0123456789abcdef but: ", span.SpanID)
	}

	if span.Name != "test-span" {
		t.Error("span.Name should be test-span but: ", span.Name)
	}

	if span.Kind != "internal" {
		t.Error("span.Kind should be internal but: ", span.Kind)
	}

	if span.TraceState != "congo=xx,key=val" {
		t.Error("span.TraceState should be congo=xx,key=val but: ", span.TraceState)
	}

	expectedStartTime, _ := time.Parse(time.RFC3339, "2025-07-09T14:03:02.123Z")
	if !span.StartTime.Equal(expectedStartTime) {
		t.Error("span.StartTime should be 2025-07-09T14:03:02.123Z but: ", span.StartTime)
	}

	expectedEndTime, _ := time.Parse(time.RFC3339, "2025-07-09T14:03:02.456Z")
	if !span.EndTime.Equal(expectedEndTime) {
		t.Error("span.EndTime should be 2025-07-09T14:03:02.456Z but: ", span.EndTime)
	}

	if len(span.Attributes) != 1 {
		t.Error("span should have 1 attribute but: ", len(span.Attributes))
	} else {
		attr := span.Attributes[0]
		if attr.Key != "http.route" {
			t.Error("attribute key should be http.route but: ", attr.Key)
		}
		if attr.Value.ValueType != "string" {
			t.Error("attribute value type should be string but: ", attr.Value.ValueType)
		}
		if attr.Value.StringValue != "/" {
			t.Error("attribute string value should be / but: ", attr.Value.StringValue)
		}
	}

	if len(span.Events) != 1 {
		t.Error("span should have 1 event but: ", len(span.Events))
	} else {
		event := span.Events[0]
		if event.Name != "event1" {
			t.Error("event name should be event1 but: ", event.Name)
		}
		expectedTime, _ := time.Parse(time.RFC3339, "2025-07-09T14:03:02.789Z")
		if !event.Time.Equal(expectedTime) {
			t.Error("event time should be 2025-07-09T14:03:02.789Z but: ", event.Time)
		}
	}

	if len(span.Links) != 1 {
		t.Error("span should have 1 link but: ", len(span.Links))
	} else {
		link := span.Links[0]
		if link.TraceID != "abcdef0123456789abcdef0123456789" {
			t.Error("link traceId should be abcdef0123456789abcdef0123456789 but: ", link.TraceID)
		}
		if link.SpanID != "abcdefabcdef0102" {
			t.Error("link spanId should be abcdefabcdef0102 but: ", link.SpanID)
		}
	}

	if span.Status == nil {
		t.Error("span.Status should not be nil")
	} else {
		if span.Status.Code != "ok" {
			t.Error("status code should be ok but: ", span.Status.Code)
		}
		if span.Status.Message != "status message" {
			t.Error("status message should be 'status message' but: ", span.Status.Message)
		}
	}

	if span.Resource == nil {
		t.Error("span.Resource should not be nil")
	}

	if span.Scope == nil {
		t.Error("span.Scope should not be nil")
	} else {
		if span.Scope.Name != "my-library" {
			t.Error("scope name should be my-library but: ", span.Scope.Name)
		}
		if span.Scope.Version != "1.0.0" {
			t.Error("scope version should be 1.0.0 but: ", span.Scope.Version)
		}
	}
}

func TestGetTrace_AnyValueTypes(t *testing.T) {
	// Test different AnyValue types
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		respJSON, _ := json.Marshal(map[string]interface{}{
			"spans": []map[string]interface{}{
				{
					"traceId":    "0123456789abcdef0123456789abcdef",
					"spanId":     "0123456789abcdef",
					"traceState": "",
					"name":       "test-span",
					"kind":       "internal",
					"startTime":  "2025-07-09T14:03:02.999Z",
					"endTime":    "2025-07-09T14:03:03.001Z",
					"attributes": []map[string]interface{}{
						{
							"key": "string.attr",
							"value": map[string]interface{}{
								"valueType":   "string",
								"stringValue": "test",
							},
						},
						{
							"key": "bool.attr",
							"value": map[string]interface{}{
								"valueType": "bool",
								"boolValue": true,
							},
						},
						{
							"key": "int.attr",
							"value": map[string]interface{}{
								"valueType": "int",
								"intValue":  42,
							},
						},
						{
							"key": "double.attr",
							"value": map[string]interface{}{
								"valueType":   "double",
								"doubleValue": 3.14,
							},
						},
						{
							"key": "array.attr",
							"value": map[string]interface{}{
								"valueType": "array",
								"arrayValue": []map[string]interface{}{
									{
										"valueType": "int",
										"intValue":  10,
									},
									{
										"valueType": "int",
										"intValue":  20,
									},
								},
							},
						},
						{
							"key": "kvlist.attr",
							"value": map[string]interface{}{
								"valueType": "kvlist",
								"kvlistValue": map[string]interface{}{
									"en": map[string]interface{}{
										"valueType":   "string",
										"stringValue": "success",
									},
								},
							},
						},
						{
							"key": "empty.attr",
							"value": map[string]interface{}{
								"valueType": "empty",
							},
						},
					},
					"droppedAttributesCount": 0,
					"events":                 []interface{}{},
					"droppedEventsCount":     0,
					"links":                  []interface{}{},
					"droppedLinksCount":      0,
					"status": map[string]interface{}{
						"message": "",
						"code":    "unset",
					},
					"resource": map[string]interface{}{
						"attributes":             []interface{}{},
						"droppedAttributesCount": 0,
					},
					"scope": map[string]interface{}{
						"name":                   "",
						"version":                "",
						"attributes":             []interface{}{},
						"droppedAttributesCount": 0,
					},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	trace, err := client.GetTrace("0123456789abcdef0123456789abcdef")

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if trace == nil || len(trace.Spans) == 0 {
		t.Error("trace should have spans")
		return
	}

	span := trace.Spans[0]
	if len(span.Attributes) != 7 {
		t.Error("span should have 7 attributes but: ", len(span.Attributes))
		return
	}

	// Test string value
	stringAttr := span.Attributes[0]
	if stringAttr.Value.ValueType != "string" || stringAttr.Value.StringValue != "test" {
		t.Error("string attribute should be 'test' but got:", stringAttr.Value.StringValue)
	}

	// Test bool value
	boolAttr := span.Attributes[1]
	if boolAttr.Value.ValueType != "bool" || !boolAttr.Value.BoolValue {
		t.Error("bool attribute should be true but got:", boolAttr.Value.BoolValue)
	}

	// Test int value
	intAttr := span.Attributes[2]
	if intAttr.Value.ValueType != "int" || intAttr.Value.IntValue != 42 {
		t.Error("int attribute should be 42 but got:", intAttr.Value.IntValue)
	}

	// Test double value
	doubleAttr := span.Attributes[3]
	if doubleAttr.Value.ValueType != "double" || doubleAttr.Value.DoubleValue != 3.14 {
		t.Error("double attribute should be 3.14 but got:", doubleAttr.Value.DoubleValue)
	}

	// Test array value
	arrayAttr := span.Attributes[4]
	if arrayAttr.Value.ValueType != "array" || len(arrayAttr.Value.ArrayValue) != 2 {
		t.Error("array attribute should have 2 elements but got:", len(arrayAttr.Value.ArrayValue))
	}

	// Test kvlist value
	kvlistAttr := span.Attributes[5]
	if kvlistAttr.Value.ValueType != "kvlist" || len(kvlistAttr.Value.KvlistValue) != 1 {
		t.Error("kvlist attribute should have 1 key but got:", len(kvlistAttr.Value.KvlistValue))
	}

	// Test empty value
	emptyAttr := span.Attributes[6]
	if emptyAttr.Value.ValueType != "empty" {
		t.Error("empty attribute should have valueType 'empty' but got:", emptyAttr.Value.ValueType)
	}
}

func TestGetTrace_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(404)
		fmt.Fprint(res, `{"error": {"message": "trace not found"}}`) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	_, err := client.GetTrace("nonexistent")

	if err == nil {
		t.Error("should return error for 404 response")
	}
}
