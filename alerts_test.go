package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindAlerts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/alerts" {
			t.Error("request URL should be /api/v0/alerts but: ", req.URL.Path)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"alerts": {
				{
					"id":        "2wpLU5fBXbG",
					"status":    "CRITICAL",
					"monitorId": "2cYjfibBkaj",
					"type":      "connectivity",
					"openedAt":  1445399342,
					"hostId":    "2vJ965ygiXf",
				},
				{
					"id":        "2ust8jNxFH3",
					"status":    "CRITICAL",
					"monitorId": "2cYjfibBkaj",
					"type":      "connectivity",
					"openedAt":  1441939801,
					"hostId":    "2tFrtykgMib",
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	alerts, err := client.FindAlerts()

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if alerts.Alerts[0].Type != "connectivity" {
		t.Error("request sends json including type but: ", alerts.Alerts[0])
	}

	if alerts.Alerts[1].Type != "connectivity" {
		t.Error("request sends json including type but: ", alerts.Alerts[1])
	}

	if alerts.Alerts[1].Status != "CRITICAL" {
		t.Error("request sends json including status but: ", alerts.Alerts[1])
	}

	if alerts.Alerts[1].OpenedAt != 1441939801 {
		t.Error("request sends json including openedAt but: ", alerts.Alerts[1])
	}
}

func TestFindAlertsWithNextId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/alerts" {
			t.Error("request URL should be /api/v0/alerts but: ", req.URL.Path)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"alerts": []map[string]interface{}{
				{
					"id":        "2wpLU5fBXbG",
					"status":    "CRITICAL",
					"monitorId": "2cYjfibBkaj",
					"type":      "connectivity",
					"openedAt":  1445399342,
					"hostId":    "2vJ965ygiXf",
				},
				{
					"id":        "2ust8jNxFH3",
					"status":    "CRITICAL",
					"monitorId": "2cYjfibBkaj",
					"type":      "connectivity",
					"openedAt":  1441939801,
					"hostId":    "2tFrtykgMib",
				},
			},
			"nextId": "2fsf8jRxFG1",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	alerts, err := client.FindAlerts()

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if alerts.Alerts[0].Type != "connectivity" {
		t.Error("request sends json including type but: ", alerts.Alerts[0])
	}

	if alerts.Alerts[1].Type != "connectivity" {
		t.Error("request sends json including type but: ", alerts.Alerts[1])
	}

	if alerts.Alerts[1].Status != "CRITICAL" {
		t.Error("request sends json including status but: ", alerts.Alerts[1])
	}

	if alerts.Alerts[1].OpenedAt != 1441939801 {
		t.Error("request sends json including openedAt but: ", alerts.Alerts[1])
	}

	if alerts.NextID != "2fsf8jRxFG1" {
		t.Error("request sends json including nextId but: ", alerts.NextID)
	}
}

func TestFindAlertsByNextId(t *testing.T) {
	var nextID = "2fsf8jRxFG1"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/alerts" {
			t.Error("request URL should be /api/v0/alerts but: ", req.URL.Path)
		}
		if req.URL.RawQuery != "nextId="+nextID {
			t.Error("request Query should be /api/v/alerts?nextId=", nextID, " but: ", req.URL.RawQuery)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"alerts": []map[string]interface{}{
				{
					"id":        "2fsf8jRxFG1",
					"status":    "CRITICAL",
					"monitorId": "2cYjfibBkaj",
					"type":      "connectivity",
					"openedAt":  1445399342,
					"hostId":    "2vJ965ygiXf",
				},
				{
					"id":        "2dsg6jNxEY7",
					"status":    "CRITICAL",
					"monitorId": "2cYjfibBkaj",
					"type":      "connectivity",
					"openedAt":  1441939801,
					"hostId":    "2tFrtykgMib",
				},
			},
			"nextId": "2ghy4jDhEH3",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	alerts, err := client.FindAlertsByNextID(nextID)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if alerts.Alerts[0].Type != "connectivity" {
		t.Error("request sends json including type but: ", alerts.Alerts[0])
	}

	if alerts.Alerts[1].Type != "connectivity" {
		t.Error("request sends json including type but: ", alerts.Alerts[1])
	}

	if alerts.Alerts[1].Status != "CRITICAL" {
		t.Error("request sends json including status but: ", alerts.Alerts[1])
	}

	if alerts.Alerts[1].OpenedAt != 1441939801 {
		t.Error("request sends json including openedAt but: ", alerts.Alerts[1])
	}

	if alerts.NextID != "2ghy4jDhEH3" {
		t.Error("request sends json including nextId but: ", alerts.NextID)
	}
}
func TestFindWithClosedAlerts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/alerts" {
			t.Error("request URL should be /api/v0/alerts but: ", req.URL.Path)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"alerts": {
				{
					"id":        "2wpLU5fBXbG",
					"status":    "CRITICAL",
					"monitorId": "2cYjfibBkaj",
					"type":      "connectivity",
					"openedAt":  1445399342,
					"hostId":    "2vJ965ygiXf",
				},
				{
					"id":        "2ust8jNxFH3",
					"status":    "CRITICAL",
					"monitorId": "2cYjfibBkaj",
					"type":      "service",
					"openedAt":  1441939801,
					"hostId":    "2tFrtykgMib",
				},
				{
					"id":        "2ust8jNxFH3",
					"status":    "OK",
					"monitorId": "2cYjfibBkaj",
					"type":      "host",
					"reason":    "hoge",
					"openedAt":  1441939801,
					"closedAt":  1441940101,
					"hostId":    "2tFrtykgMib",
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	alerts, err := client.FindWithClosedAlerts()

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if alerts.Alerts[0].Type != "connectivity" {
		t.Error("request sends json including type but: ", alerts.Alerts[0])
	}

	if alerts.Alerts[1].Status != "CRITICAL" {
		t.Error("request sends json including type but: ", alerts.Alerts[1])
	}

	if alerts.Alerts[2].Status != "OK" {
		t.Error("request sends json including status but: ", alerts.Alerts[1])
	}

	if alerts.Alerts[2].ClosedAt != 1441940101 {
		t.Error("request sends json including openedAt but: ", alerts.Alerts[1])
	}
}

func TestFindWithClosedAlertsByNextId(t *testing.T) {
	var nextID = "2wpLU5fBXbG"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/alerts" {
			t.Error("request URL should be /api/v0/alerts but: ", req.URL.Path)
		}
		q := req.URL.Query()
		if q.Get("nextId") != nextID {
			t.Error("request nextId should be ", nextID, "but: ", q.Get("nextId"))
		}
		if q.Get("withClosed") != "true" {
			t.Error("request withClosed should be true but: ", q.Get("withClosed"))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"alerts": []map[string]interface{}{
				{
					"id":        "2wpLU5fBXbG",
					"status":    "CRITICAL",
					"monitorId": "2cYjfibBkaj",
					"type":      "connectivity",
					"openedAt":  1445399342,
					"hostId":    "2vJ965ygiXf",
				},
				{
					"id":        "2ust8jNxFH3",
					"status":    "CRITICAL",
					"monitorId": "2cYjfibBkaj",
					"type":      "service",
					"openedAt":  1441939801,
					"hostId":    "2tFrtykgMib",
				},
				{
					"id":        "2ust8jNxFH3",
					"status":    "OK",
					"monitorId": "2cYjfibBkaj",
					"type":      "host",
					"reason":    "hoge",
					"openedAt":  1441939801,
					"closedAt":  1441940101,
					"hostId":    "2tFrtykgMib",
				},
			},
			"nextId": "2fsf8jRxFG1",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	alerts, err := client.FindWithClosedAlertsByNextID(nextID)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if alerts.Alerts[0].Type != "connectivity" {
		t.Error("request sends json including type but: ", alerts.Alerts[0])
	}

	if alerts.Alerts[1].Status != "CRITICAL" {
		t.Error("request sends json including type but: ", alerts.Alerts[1])
	}

	if alerts.Alerts[2].Status != "OK" {
		t.Error("request sends json including status but: ", alerts.Alerts[1])
	}

	if alerts.Alerts[2].ClosedAt != 1441940101 {
		t.Error("request sends json including openedAt but: ", alerts.Alerts[1])
	}

	if alerts.NextID != "2fsf8jRxFG1" {
		t.Error("request sends json including nextId but: ", alerts.NextID)
	}
}
