package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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
		if query.Get("orderColumn") != "" {
			t.Error("request query 'orderColumn' param should be empty but: ", query.Get("orderColumn"))
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

func TestListHTTPServerStatsWithAllParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		if query.Get("serviceName") != "full-service" {
			t.Error("request query 'serviceName' param should be full-service but: ", query.Get("serviceName"))
		}
		if query.Get("serviceNamespace") != "test-namespace" {
			t.Error("request query 'serviceNamespace' param should be test-namespace but: ", query.Get("serviceNamespace"))
		}
		if query.Get("environment") != "production" {
			t.Error("request query 'environment' param should be production but: ", query.Get("environment"))
		}
		if query.Get("version") != "v1.0.0" {
			t.Error("request query 'version' param should be v1.0.0 but: ", query.Get("version"))
		}
		if query.Get("method") != "GET" {
			t.Error("request query 'method' param should be GET but: ", query.Get("method"))
		}
		if query.Get("route") != "/api/test" {
			t.Error("request query 'route' param should be /api/test but: ", query.Get("route"))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"method":              "GET",
					"route":               "/api/test",
					"totalMillis":         100.0,
					"averageMillis":       10.0,
					"approxP95Millis":     20.0,
					"errorRatePercentage": 0.5,
					"requestCount":        10,
				},
			},
			"hasNextPage": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	serviceNamespace := "test-namespace"
	environment := "production"
	version := "v1.0.0"
	orderColumn := "AVERAGE"
	orderDirection := "ASC"
	method := "GET"
	route := "/api/test"
	page := 2
	perPage := 50

	result, err := client.ListHTTPServerStats(&ListHTTPServerStatsParam{
		ServiceName:      "full-service",
		From:             1000000000,
		To:               2000000000,
		ServiceNamespace: &serviceNamespace,
		Environment:      &environment,
		Version:          &version,
		OrderColumn:      &orderColumn,
		OrderDirection:   &orderDirection,
		Method:           &method,
		Route:            &route,
		Page:             &page,
		PerPage:          &perPage,
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

	if result.Results[0].Method != "GET" {
		t.Error("method should be GET but: ", result.Results[0].Method)
	}

	if result.Results[0].Route != "/api/test" {
		t.Error("route should be /api/test but: ", result.Results[0].Route)
	}
}
