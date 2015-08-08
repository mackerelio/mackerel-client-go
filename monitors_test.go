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
					"id":       "2cSZzK3XfmG",
					"type":     "passive",
					"name":     "connectivity",
					"duration": 5,
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

	if monitors[0].Type != "passive" {
		t.Error("request sends json including type but: ", monitors[0])
	}

	if monitors[0].Name != "connectivity" {
		t.Error("request sends json including name but: ", monitors[0])
	}

}