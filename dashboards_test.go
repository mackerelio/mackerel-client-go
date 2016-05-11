package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindDashboards(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/dashboards" {
			t.Error("request URL should be /api/v0/dashboards but :", req.URL.Path)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"dashboards": []map[string]interface{}{
				map[string]interface{}{
					"id":           "2c5bLca8d",
					"title":        "My Dashboard",
					"bodyMarkDown": "# A test dashboard",
					"urlPath":      "2u4PP3TJqbu",
					"createdAt":    1439346145003,
					"updatedAt":    1439346145003,
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	dashboards, err := client.FindDashboards()

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	if dashboards[0].ID != "2c5bLca8d" {
		t.Error("request sends json including id but: ", dashboards[0])
	}

	if dashboards[0].Title != "My Dashboard" {
		t.Error("request sends json including title but: ", dashboards[0])
	}

	if dashboards[0].BodyMarkDown != "# A test dashboard" {
		t.Error("request sends json including bodyMarkDown but: ", dashboards[0])
	}

	if dashboards[0].URLPath != "2u4PP3TJqbu" {
		t.Error("request sends json including urlpath but: ", dashboards[0])
	}

	if dashboards[0].CreatedAt != 1439346145003 {
		t.Error("request sends json including createdAt but: ", dashboards[0])
	}

	if dashboards[0].UpdatedAt != 1439346145003 {
		t.Error("request sends json including updatedAt but: ", dashboards[0])
	}
}

func TestFindDashboard(t *testing.T) {

	testID := "2c5bLca8d"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/dashboards/%s", testID) {
			t.Error("request URL should be /api/v0/dashboards/<ID> but :", req.URL.Path)
		}

		respJSON, _ := json.Marshal(
			map[string]interface{}{
				"id":           "2c5bLca8d",
				"title":        "My Dashboard",
				"bodyMarkDown": "# A test dashboard",
				"urlPath":      "2u4PP3TJqbu",
				"createdAt":    1439346145003,
				"updatedAt":    1439346145003,
			},
		)

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	dashboard, err := client.FindDashboard(testID)

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	if dashboard.ID != "2c5bLca8d" {
		t.Error("request sends json including id but: ", dashboard)
	}

	if dashboard.Title != "My Dashboard" {
		t.Error("request sends json including title but: ", dashboard)
	}

	if dashboard.BodyMarkDown != "# A test dashboard" {
		t.Error("request sends json including bodyMarkDown but: ", dashboard)
	}

	if dashboard.URLPath != "2u4PP3TJqbu" {
		t.Error("request sends json including urlpath but: ", dashboard)
	}

	if dashboard.CreatedAt != 1439346145003 {
		t.Error("request sends json including createdAt but: ", dashboard)
	}

	if dashboard.UpdatedAt != 1439346145003 {
		t.Error("request sends json including updatedAt but: ", dashboard)
	}
}
