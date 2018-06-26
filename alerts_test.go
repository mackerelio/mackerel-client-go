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
			t.Error("request URL should be /api/v0/alerts but :", req.URL.Path)
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

	if alerts[0].Type != "connectivity" {
		t.Error("request sends json including type but: ", alerts[0])
	}

	if alerts[1].Type != "connectivity" {
		t.Error("request sends json including type but: ", alerts[1])
	}

	if alerts[1].Status != "CRITICAL" {
		t.Error("request sends json including status but: ", alerts[1])
	}

	if alerts[1].OpenedAt != 1441939801 {
		t.Error("request sends json including openedAt but: ", alerts[1])
	}
}
