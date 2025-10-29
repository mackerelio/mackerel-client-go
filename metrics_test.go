package mackerel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPostHostMetricValues(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/tsdb" {
			t.Error("request URL should be /api/v0/tsdb but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := io.ReadAll(req.Body)

		var values []struct {
			HostID string      `json:"hostId"`
			Name   string      `json:"name"`
			Time   float64     `json:"time"`
			Value  interface{} `json:"value"`
		}

		err := json.Unmarshal(body, &values)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		if values[0].HostID != "9rxGOHfVF8F" {
			t.Error("request sends json including hostId but: ", values[0].HostID)
		}
		if values[0].Name != "custom.metric.mysql.connections" {
			t.Error("request sends json including name but: ", values[0].Name)
		}
		if values[0].Time != 123456789 {
			t.Error("request sends json including time but: ", values[0].Time)
		}
		if values[0].Value.(float64) != 100 {
			t.Error("request sends json including value but: ", values[0].Value)
		}

		respJSON, _ := json.Marshal(map[string]bool{
			"success": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	err := client.PostHostMetricValues([]*HostMetricValue{
		{
			HostID: "9rxGOHfVF8F",
			MetricValue: &MetricValue{
				Name:  "custom.metric.mysql.connections",
				Time:  123456789,
				Value: 100,
			},
		},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}
}

func TestPostServiceMetricValues(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/services/My-Service/tsdb" {
			t.Error("request URL should be /api/v0/services/My-Service/tsdb but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := io.ReadAll(req.Body)

		var values []struct {
			Name  string      `json:"name"`
			Time  float64     `json:"time"`
			Value interface{} `json:"value"`
		}

		err := json.Unmarshal(body, &values)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		if values[0].Name != "proxy.access_log.latency" {
			t.Error("request sends json including name but: ", values[0].Name)
		}
		if values[0].Time != 123456789 {
			t.Error("request sends json including time but: ", values[0].Time)
		}
		if values[0].Value.(float64) != 500 {
			t.Error("request sends json including value but: ", values[0].Value)
		}

		respJSON, _ := json.Marshal(map[string]bool{
			"success": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	err := client.PostServiceMetricValues("My-Service", []*MetricValue{
		{
			Name:  "proxy.access_log.latency",
			Time:  123456789,
			Value: 500,
		},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}
}

func TestFetchLatestMetricValues(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/tsdb/latest" {
			t.Error("request URL should be /api/v0/tsdb/latest but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		query := req.URL.Query()
		if !reflect.DeepEqual(query["hostId"], []string{"123456ABCD", "654321ABCD"}) {
			t.Error("request query 'hostId' param should be ['123456ABCD', '654321ABCD'] but: ", query["hostId"])
		}
		if !reflect.DeepEqual(query["name"], []string{"mysql.connections.Connections", "mysql.connections.Thread_created"}) {
			t.Error("request query 'name' param should be ['mysql.connections.Connections', 'mysql.connections.Thread_created'] but: ", query["name"])
		}

		respJSON, _ := json.Marshal(map[string]map[string]map[string]*MetricValue{
			"tsdbLatest": {
				"123456ABCD": {
					"mysql.connections.Connections": {
						Name:  "mysql.connections.Connections",
						Time:  123456789,
						Value: 200,
					},
					"mysql.connections.Thread_created": {
						Name:  "mysql.connections.Thread_created",
						Time:  123456789,
						Value: 220,
					},
				},
				"654321ABCD": {
					"mysql.connections.Connections": {
						Name:  "mysql.connections.Connections",
						Time:  123456789,
						Value: 300,
					},
					"mysql.connections.Thread_created": {
						Name:  "mysql.connections.Thread_created",
						Time:  123456789,
						Value: 310,
					},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	hostIDs := []string{"123456ABCD", "654321ABCD"}
	metricNames := []string{"mysql.connections.Connections", "mysql.connections.Thread_created"}
	latestMetricValues, err := client.FetchLatestMetricValues(hostIDs, metricNames)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if latestMetricValues["123456ABCD"]["mysql.connections.Connections"].Value.(float64) != 200 {
		t.Error("123456ABCD host mysql.connections.Connections should be 200 but: ", latestMetricValues["123456ABCD"]["mysql.connections.Connections"].Value)
	}

	if latestMetricValues["654321ABCD"]["mysql.connections.Connections"].Value.(float64) != 300 {
		t.Error("654321ABCD host mysql.connections.Connections should be 300 but: ", latestMetricValues["654321ABCD"]["mysql.connections.Connections"].Value)
	}
}

func TestFetchHostMetricValues(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/123456ABCD/metrics" {
			t.Error("request URL should be /api/v0/hosts/123456ABCD/metrics but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		query := req.URL.Query()
		if !reflect.DeepEqual(query["name"], []string{"mysql.connections.Connections"}) {
			t.Error("request query 'name' param should be ['mysql.connections.Connections'] but: ", query["name"])
		}

		respJSON, _ := json.Marshal(map[string][]MetricValue{
			"metrics": {
				{
					Time:  1450000800,
					Value: 200,
				},
				{
					Time:  1450000860,
					Value: 220,
				},
				{
					Time:  1450000920,
					Value: 240,
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	hostID := "123456ABCD"
	metricName := "mysql.connections.Connections"
	metricValues, err := client.FetchHostMetricValues(hostID, metricName, 1450000700, 1450001000)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if metricValues[0].Value.(float64) != 200 {
		t.Error("123456ABCD host mysql.connections.Connections should be 200 but: ", metricValues[0].Value)
	}

	if metricValues[1].Value.(float64) != 220 {
		t.Error("123456ABCD host mysql.connections.Connections should be 220 but: ", metricValues[1].Value)
	}

	if metricValues[2].Value.(float64) != 240 {
		t.Error("123456ABCD host mysql.connections.Connections should be 240 but: ", metricValues[2].Value)
	}
}

func TestFetchServiceMetricValues(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/services/123456ABCD/metrics" {
			t.Error("request URL should be /api/v0/services/123456ABCD/metrics but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		query := req.URL.Query()
		if !reflect.DeepEqual(query["name"], []string{"custom.access_latency.avg"}) {
			t.Error("request query 'name' param should be ['custom.access_latency.avg'] but: ", query["name"])
		}

		respJSON, _ := json.Marshal(map[string][]MetricValue{
			"metrics": {
				{
					Time:  1450000800,
					Value: 0.12,
				},
				{
					Time:  1450000860,
					Value: 0.14,
				},
				{
					Time:  1450000920,
					Value: 0.16,
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	serviceName := "123456ABCD"
	metricName := "custom.access_latency.avg"
	metricValues, err := client.FetchServiceMetricValues(serviceName, metricName, 1450000700, 1450001000)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if metricValues[0].Value.(float64) != 0.12 {
		t.Error("123456ABCD host custom.access_latency.avg should be 0.12 but: ", metricValues[0].Value)
	}

	if metricValues[1].Value.(float64) != 0.14 {
		t.Error("123456ABCD host custom.access_latency.avg should be 0.14 but: ", metricValues[1].Value)
	}

	if metricValues[2].Value.(float64) != 0.16 {
		t.Error("123456ABCD host custom.access_latency.avg should be 0.16 but: ", metricValues[2].Value)
	}
}
