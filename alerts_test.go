package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func TestGetAlert(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		url := fmt.Sprintf("/api/v0/alerts/%s", "2wpLU5fBXbG")
		if req.URL.Path != url {
			t.Error("request URL should be /api/v0/alerts/<ID> but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":        "2wpLU5fBXbG",
			"status":    "CRITICAL",
			"monitorId": "2cYjfibBkaj",
			"type":      "connectivity",
			"openedAt":  1445399342,
			"hostId":    "2vJ965ygiXf",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	alert, err := client.GetAlert("2wpLU5fBXbG")
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if alert.ID != "2wpLU5fBXbG" {
		t.Error("alert id should be \"2wpLU5fBXbG\" but: ", alert.ID)
	}

	if reflect.DeepEqual(alert, &Alert{
		ID:        "2wpLU5fBXbG",
		Status:    "CRITICAL",
		MonitorID: "2cYjfibBkaj",
		Type:      "connectivity",
		HostID:    "2vJ965ygiXf",
		OpenedAt:  1445399342,
	}) != true {
		t.Errorf("Wrong data for alert: %v", alert)
	}
}

func TestFindAlertLogs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		url := fmt.Sprintf("/api/v0/alerts/%s/logs", "2wpLU5fBXbG")
		if req.URL.Path != url {
			t.Error("request URL should be /api/v0/alerts/<ID>/logs but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"logs": []map[string]interface{}{
				{
					"id":          "5m7fewuu5tS",
					"createdAt":   1735290407,
					"status":      "WARNING",
					"trigger":     "monitoring",
					"monitorId":   "5m72DB7s7sU",
					"targetValue": (*float64)(nil),
					"statusDetail": map[string]interface{}{
						"type": "check",
						"detail": map[string]interface{}{
							"message": "Uptime WARNING: 0 day(s) 0 hour(s) 6 minute(s) (398 second(s))",
							"memo":    "",
						},
					},
				}, {
					"id":           "5m7fewuu5tS",
					"createdAt":    1735290407,
					"status":       "WARNING",
					"trigger":      "monitoring",
					"monitorId":    "5m72DB7s7sU",
					"targetValue":  (*float64)(nil),
					"statusDetail": nil,
				},
			},
			"nextId": "2fsf8jRxFG1",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	logs, err := client.FindAlertLogs("2wpLU5fBXbG")
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if len(logs.AlertLogs) != 2 {
		t.Error("logs should have 1 elements but: ", len(logs.AlertLogs))
	}

	if logs.NextID != "2fsf8jRxFG1" {
		t.Error("request sends json including nextId but: ", logs.NextID)
	}

	if logs.AlertLogs[0].ID != "5m7fewuu5tS" {
		t.Error("alert id should be \"5m7fewuu5tS\" but: ", logs.AlertLogs[0].ID)
	}

	if logs.AlertLogs[0].CreatedAt != 1735290407 {
		t.Error("createdAt should be 1735290407 but: ", logs.AlertLogs[0].CreatedAt)
	}

	if logs.AlertLogs[0].Trigger != "monitoring" {
		t.Error("trigger should be \"monitoring\" but: ", logs.AlertLogs[0].Trigger)
	}

	if *logs.AlertLogs[0].MonitorID != "5m72DB7s7sU" {
		t.Error("monitorId should be \"5m72DB7s7sU\" but: ", *logs.AlertLogs[0].MonitorID)
	}

	if logs.AlertLogs[0].TargetValue != nil {
		t.Error("targetValue should be nil but: ", logs.AlertLogs[0].TargetValue)
	}

	if logs.AlertLogs[0].Status != "WARNING" {
		t.Error("alert status should be \"WARNING\" but: ", logs.AlertLogs[0].Status)
	}

	if logs.AlertLogs[0].StatusDetail.Type != "check" {
		t.Error("statusDetail type should be \"check\" but: ", logs.AlertLogs[0].StatusDetail.Type)
	}

	if logs.AlertLogs[0].StatusDetail.Detail.Message != "Uptime WARNING: 0 day(s) 0 hour(s) 6 minute(s) (398 second(s))" {
		t.Error("statusDetail message should be \"Uptime WARNING: 0 day(s) 0 hour(s) 6 minute(s) (398 second(s))\" but: ", logs.AlertLogs[0].StatusDetail.Detail.Message)
	}

	if logs.AlertLogs[0].StatusDetail.Detail.Memo != "" {
		t.Error("statusDetail memo should be empty but: ", logs.AlertLogs[0].StatusDetail.Detail.Memo)
	}

	if logs.AlertLogs[1].StatusDetail != nil {
		t.Error("statusDetail should be nil but: ", logs.AlertLogs[1].StatusDetail)
	}

}
