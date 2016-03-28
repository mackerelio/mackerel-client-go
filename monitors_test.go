package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindMonitors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/monitors" {
			t.Error("request URL should be /api/v0/monitors but :", req.URL.Path)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"monitors": []map[string]interface{}{
				map[string]interface{}{
					"id":            "2cSZzK3XfmG",
					"type":          "connectivity",
					"scopes":        []string{},
					"excludeScopes": []string{},
				},
				map[string]interface{}{
					"id":                              "2c5bLca8d",
					"type":                            "external",
					"name":                            "testMonitorExternal",
					"url":                             "https://www.example.com/",
					"maxCheckAttempts":                3,
					"service":                         "someService",
					"notificationInterval":            60,
					"responseTimeCritical":            5000,
					"responseTimeWarning":             10000,
					"responseTimeDuration":            5,
					"certificationExpirationCritical": 15,
					"certificationExpirationWarning":  30,
					"containsString":                  "Foo Bar Baz",
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	monitors, err := client.FindMonitors()

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	if monitors[0].Type != "connectivity" {
		t.Error("request sends json including type but: ", monitors[0])
	}

	if monitors[1].Type != "external" {
		t.Error("request sends json including type but: ", monitors[1])
	}

	if monitors[1].Service != "someService" {
		t.Error("request sends json including service but: ", monitors[1])
	}

	if monitors[1].NotificationInterval != 60 {
		t.Error("request sends json including notificationInterval but: ", monitors[1])
	}

	if monitors[1].ResponseTimeCritical != 5000 {
		t.Error("request sends json including responseTimeCritical but: ", monitors[1])
	}

	if monitors[1].ResponseTimeWarning != 10000 {
		t.Error("request sends json including responseTimeWarning but: ", monitors[1])
	}

	if monitors[1].ResponseTimeDuration != 5 {
		t.Error("request sends json including responseTimeDuration but: ", monitors[1])
	}

	if monitors[1].CertificationExpirationCritical != 15 {
		t.Error("request sends json including certificationExpirationCritical but: ", monitors[1])
	}

	if monitors[1].CertificationExpirationWarning != 30 {
		t.Error("request sends json including certificationExpirationWarning but: ", monitors[1])
	}

	if monitors[1].ContainsString != "Foo Bar Baz" {
		t.Error("request sends json including containsString but: ", monitors[1])
	}
}
