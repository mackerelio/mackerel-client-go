package mackerel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindDashboards(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/dashboards" {
			t.Error("request URL should be /api/v0/dashboards but: ", req.URL.Path)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"dashboards": {
				{
					"id":           "2c5bLca8d",
					"title":        "My Dashboard(Legacy)",
					"bodyMarkDown": "# A test Legacy dashboard",
					"urlPath":      "2u4PP3TJqbu",
					"createdAt":    1439346145003,
					"updatedAt":    1439346145003,
					"isLegacy":     true,
				},
				{
					"id":        "2c5bLca8e",
					"title":     "My Custom Dashboard(Current)",
					"urlPath":   "2u4PP3TJqbv",
					"createdAt": 1552909732,
					"updatedAt": 1552992837,
					"memo":      "A test Current Dashboard",
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
		t.Error("err should be nil but: ", err)
	}

	if dashboards[0].ID != "2c5bLca8d" {
		t.Error("request sends json including id but: ", dashboards[0])
	}

	if dashboards[0].Title != "My Dashboard(Legacy)" {
		t.Error("request sends json including title but: ", dashboards[0])
	}

	if dashboards[0].BodyMarkDown != "# A test Legacy dashboard" {
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

	if dashboards[0].IsLegacy != true {
		t.Error("request sends json including IsLegacy but: ", dashboards[0])
	}
}

func TestFindDashboard(t *testing.T) {

	testID := "2c5bLca8d"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/dashboards/%s", testID) {
			t.Error("request URL should be /api/v0/dashboards/<ID> but: ", req.URL.Path)
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
		t.Error("err should be nil but: ", err)
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

func TestCreateDashboard(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/dashboards" {
			t.Error("request URL should be /api/v0/dashboards but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var data struct {
			ID           string `json:"id"`
			Title        string `json:"title"`
			BodyMarkDown string `json:"bodyMarkdown"`
			URLPath      string `json:"urlPath"`
			CreatedAt    int64  `json:"createdAt"`
			UpdatedAt    int64  `json:"updatedAt"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":           "2c5bLca8d",
			"title":        "My Dashboard",
			"bodyMarkDown": "# A test dashboard",
			"urlPath":      "2u4PP3TJqbu",
			"createdAt":    1439346145003,
			"updatedAt":    1439346145003,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	dashboard, err := client.CreateDashboard(&Dashboard{
		Title:        "My Dashboard",
		BodyMarkDown: "# A test dashboard",
		URLPath:      "2u4PP3TJqbu",
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
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

func TestUpdateDashboard(t *testing.T) {

	testID := "2c5bLca8d"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/dashboards/%s", testID) {
			t.Error("request URL should be /api/v0/dashboards/<ID> but: ", req.URL.Path)
		}

		if req.Method != "PUT" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var data struct {
			ID           string `json:"id"`
			Title        string `json:"title"`
			BodyMarkDown string `json:"bodyMarkdown"`
			URLPath      string `json:"urlPath"`
			CreatedAt    int64  `json:"createdAt"`
			UpdatedAt    int64  `json:"updatedAt"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":           "2c5bLca8d",
			"title":        "My Dashboard",
			"bodyMarkDown": "# A test dashboard",
			"urlPath":      "2u4PP3TJqbu",
			"createdAt":    1439346145003,
			"updatedAt":    1439346145003,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	dashboard, err := client.UpdateDashboard(testID, &Dashboard{
		Title:        "My Dashboard",
		BodyMarkDown: "# A test dashboard",
		URLPath:      "2u4PP3TJqbu",
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
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

func TestDeleteDashboard(t *testing.T) {

	testID := "2c5bLca8d"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/dashboards/%s", testID) {
			t.Error("request URL should be /api/v0/dashboards/<ID> but: ", req.URL.Path)
		}

		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":           "2c5bLca8d",
			"title":        "My Dashboard",
			"bodyMarkDown": "# A test dashboard",
			"urlPath":      "2u4PP3TJqbu",
			"createdAt":    1439346145003,
			"updatedAt":    1439346145003,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	dashboard, err := client.DeleteDashboard(testID)

	if err != nil {
		t.Error("err should be nil but: ", err)
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
