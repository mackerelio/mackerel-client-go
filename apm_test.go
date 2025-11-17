package mackerel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestListHTTPServerStats(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/apm/http-server-stats" {
			t.Error("request URL should be /api/v0/apm/http-server-stats but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		query := req.URL.Query()
		if query.Get("serviceName") != "test-service" {
			t.Error("request query 'serviceName' param should be test-service but: ", query.Get("serviceName"))
		}
		if query.Get("from") != "1234567890" {
			t.Error("request query 'from' param should be 1234567890 but: ", query.Get("from"))
		}
		if query.Get("to") != "1234567900" {
			t.Error("request query 'to' param should be 1234567900 but: ", query.Get("to"))
		}
		if query.Get("orderColumn") != "P95" {
			t.Error("request query 'orderColumn' param should be P95 but: ", query.Get("orderColumn"))
		}
		if query.Get("orderDirection") != "DESC" {
			t.Error("request query 'orderDirection' param should be DESC but: ", query.Get("orderDirection"))
		}
		if query.Get("page") != "1" {
			t.Error("request query 'page' param should be 1 but: ", query.Get("page"))
		}
		if query.Get("perPage") != "20" {
			t.Error("request query 'perPage' param should be 20 but: ", query.Get("perPage"))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"method":              "GET",
					"route":               "/api/users",
					"totalMillis":         837.0,
					"averageMillis":       9.01,
					"approxP95Millis":     19.89,
					"errorRatePercentage": 0.0,
					"requestCount":        93,
				},
				{
					"method":              "POST",
					"route":               "/api/posts",
					"totalMillis":         1234.5,
					"averageMillis":       12.34,
					"approxP95Millis":     25.67,
					"errorRatePercentage": 1.5,
					"requestCount":        100,
				},
			},
			"hasNextPage": false,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	orderColumn := "P95"
	orderDirection := "DESC"
	page := 1
	perPage := 20

	result, err := client.ListHTTPServerStats(&ListHTTPServerStatsParam{
		ServiceName:    "test-service",
		From:           1234567890,
		To:             1234567900,
		OrderColumn:    &orderColumn,
		OrderDirection: &orderDirection,
		Page:           &page,
		PerPage:        &perPage,
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if result.HasNextPage != false {
		t.Error("hasNextPage should be false but: ", result.HasNextPage)
	}

	if len(result.Results) != 2 {
		t.Error("results length should be 2 but: ", len(result.Results))
	}

	expectedFirst := &HTTPServerStats{
		Method:              "GET",
		Route:               "/api/users",
		TotalMillis:         837.0,
		AverageMillis:       9.01,
		ApproxP95Millis:     19.89,
		ErrorRatePercentage: 0.0,
		RequestCount:        93,
	}

	if !reflect.DeepEqual(result.Results[0], expectedFirst) {
		t.Errorf("first result should be %v but: %v", expectedFirst, result.Results[0])
	}

	expectedSecond := &HTTPServerStats{
		Method:              "POST",
		Route:               "/api/posts",
		TotalMillis:         1234.5,
		AverageMillis:       12.34,
		ApproxP95Millis:     25.67,
		ErrorRatePercentage: 1.5,
		RequestCount:        100,
	}

	if !reflect.DeepEqual(result.Results[1], expectedSecond) {
		t.Errorf("second result should be %v but: %v", expectedSecond, result.Results[1])
	}
}

func TestListHTTPServerStatsWithMinimalParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/apm/http-server-stats" {
			t.Error("request URL should be /api/v0/apm/http-server-stats but: ", req.URL.Path)
		}

		query := req.URL.Query()
		if query.Get("serviceName") != "minimal-service" {
			t.Error("request query 'serviceName' param should be minimal-service but: ", query.Get("serviceName"))
		}
		if query.Get("from") != "1000000000" {
			t.Error("request query 'from' param should be 1000000000 but: ", query.Get("from"))
		}
		if query.Get("to") != "2000000000" {
			t.Error("request query 'to' param should be 2000000000 but: ", query.Get("to"))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"results":     []map[string]interface{}{},
			"hasNextPage": false,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	result, err := client.ListHTTPServerStats(&ListHTTPServerStatsParam{
		ServiceName: "minimal-service",
		From:        1000000000,
		To:          2000000000,
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if result.HasNextPage != false {
		t.Error("hasNextPage should be false but: ", result.HasNextPage)
	}

	if len(result.Results) != 0 {
		t.Error("results length should be 0 but: ", len(result.Results))
	}
}

func TestListHTTPServerStatsContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/apm/http-server-stats" {
			t.Error("request URL should be /api/v0/apm/http-server-stats but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		query := req.URL.Query()
		if query.Get("serviceName") != "test-service" {
			t.Error("request query 'serviceName' param should be test-service but: ", query.Get("serviceName"))
		}
		if query.Get("from") != "1234567890" {
			t.Error("request query 'from' param should be 1234567890 but: ", query.Get("from"))
		}
		if query.Get("to") != "1234567900" {
			t.Error("request query 'to' param should be 1234567900 but: ", query.Get("to"))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"method":              "GET",
					"route":               "/api/users",
					"totalMillis":         837.0,
					"averageMillis":       9.01,
					"approxP95Millis":     19.89,
					"errorRatePercentage": 0.0,
					"requestCount":        93,
				},
			},
			"hasNextPage": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	ctx := context.Background()
	result, err := client.ListHTTPServerStatsContext(ctx, &ListHTTPServerStatsParam{
		ServiceName: "test-service",
		From:        1234567890,
		To:          1234567900,
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if result.HasNextPage != true {
		t.Error("hasNextPage should be true but: ", result.HasNextPage)
	}

	if len(result.Results) != 1 {
		t.Error("results length should be 1 but: ", len(result.Results))
	}

	expected := &HTTPServerStats{
		Method:              "GET",
		Route:               "/api/users",
		TotalMillis:         837.0,
		AverageMillis:       9.01,
		ApproxP95Millis:     19.89,
		ErrorRatePercentage: 0.0,
		RequestCount:        93,
	}

	if !reflect.DeepEqual(result.Results[0], expected) {
		t.Errorf("result should be %v but: %v", expected, result.Results[0])
	}
}

func TestListHTTPServerStatsContextWithTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(100 * time.Millisecond)
		respJSON, _ := json.Marshal(map[string]interface{}{
			"results":     []map[string]interface{}{},
			"hasNextPage": false,
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.ListHTTPServerStatsContext(ctx, &ListHTTPServerStatsParam{
		ServiceName: "test-service",
		From:        1234567890,
		To:          1234567900,
	})

	if err == nil {
		t.Error("err should not be nil for timeout")
	}
}

func TestListHTTPServerStatsContextWithCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(100 * time.Millisecond)
		respJSON, _ := json.Marshal(map[string]interface{}{
			"results":     []map[string]interface{}{},
			"hasNextPage": false,
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	_, err := client.ListHTTPServerStatsContext(ctx, &ListHTTPServerStatsParam{
		ServiceName: "test-service",
		From:        1234567890,
		To:          1234567900,
	})

	if err == nil {
		t.Error("err should not be nil for canceled context")
	}
}
