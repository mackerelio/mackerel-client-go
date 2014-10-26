package mackerel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostHostMetricValues(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/tsdb" {
			t.Error("request URL should be /api/v0/tsdb but :", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var values []struct {
			HostId string      `json:"hostId"`
			Name   string      `json:"name"`
			Time   float64     `json:"time"`
			Value  interface{} `json:"value"`
		}

		err := json.Unmarshal(body, &values)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		if values[0].HostId != "9rxGOHfVF8F" {
			t.Error("request sends json including hostId but: ", values[0].HostId)
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

		respJson, _ := json.Marshal(map[string]bool{
			"success": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer ts.Close()

	client, _ := NewClientForTest("dummy-key", ts.URL, false)
	err := client.PostHostMetricValues([]*HostMetricValue{
		&HostMetricValue{
			HostId: "9rxGOHfVF8F",
			Name:   "custom.metric.mysql.connections",
			Time:   123456789,
			Value:  100,
		},
	})

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}
}

func TestPostServiceMetricValues(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/services/My-Service/tsdb" {
			t.Error("request URL should be /api/v0/services/My-Service/tsdb but :", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

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

		respJson, _ := json.Marshal(map[string]bool{
			"success": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer ts.Close()

	client, _ := NewClientForTest("dummy-key", ts.URL, false)
	err := client.PostServiceMetricValues("My-Service", []*ServiceMetricValue{
		&ServiceMetricValue{
			Name:  "proxy.access_log.latency",
			Time:  123456789,
			Value: 500,
		},
	})

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}
}
