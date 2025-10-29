package mackerel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindDashboards(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/dashboards" {
			t.Error("request URL should be /api/v0/dashboards but: ", req.URL.Path)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"dashboards": []interface{}{
				map[string]interface{}{
					"id":        "2c5bLca8d",
					"title":     "My Dashboard 1",
					"urlPath":   "2u4PP3TJqbu",
					"createdAt": 1439346145003,
					"updatedAt": 1439346145003,
					"memo":      "My Dashboard memo 1",
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
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

	if dashboards[0].Title != "My Dashboard 1" {
		t.Error("request sends json including title but: ", dashboards[0])
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

	if dashboards[0].Memo != "My Dashboard memo 1" {
		t.Error("request sends json including memo but: ", dashboards[0])
	}
}

func TestFindDashboard(t *testing.T) {

	testID := "2c5bLca8e"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/dashboards/%s", testID) {
			t.Error("request URL should be /api/v0/dashboards/<ID> but: ", req.URL.Path)
		}

		respJSON, _ := json.Marshal(
			map[string]interface{}{
				"id":        "2c5bLca8e",
				"createdAt": 1552909732,
				"updatedAt": 1552992837,
				"title":     "My Custom Dashboard(Current)",
				"urlPath":   "2u4PP3TJqbv",
				"memo":      "A test Current Dashboard",
				"widgets": []map[string]interface{}{
					{
						"type":     "markdown",
						"title":    "markdown_widget",
						"markdown": "# body",
						"layout": map[string]interface{}{
							"x":      0,
							"y":      0,
							"width":  24,
							"height": 3,
						},
					},
					{
						"type":  "graph",
						"title": "graph_widget",
						"graph": map[string]interface{}{
							"type":   "host",
							"hostId": "2u4PP3TJqbw",
							"name":   "loadavg.loadavg15",
						},
						"layout": map[string]interface{}{
							"x":      0,
							"y":      7,
							"width":  8,
							"height": 10,
						},
						"referenceLines": []map[string]interface{}{
							{
								"label": "critical",
								"value": 1.23,
							},
						},
					},
					{
						"type":         "value",
						"title":        "value_widget",
						"fractionSize": 2,
						"suffix":       "total",
						"formatRules": []map[string]interface{}{
							{
								"name":      "SLO",
								"threshold": 2.34,
								"operator":  ">",
							},
						},
						"metric": map[string]interface{}{
							"type":       "expression",
							"expression": "alias(scale(\nsum(\n  group(\n    host(2u4PP3TJqbx,loadavg.*)\n  )\n),\n1\n), 'test')",
						},
						"layout": map[string]interface{}{
							"x":      0,
							"y":      17,
							"width":  8,
							"height": 5,
						},
					},
					{
						"type":         "alertStatus",
						"title":        "alert_status_widget",
						"roleFullname": "test:dashboard",
						"layout": map[string]interface{}{
							"x":      9,
							"y":      3,
							"width":  6,
							"height": 6,
						},
					},
				},
			})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))

	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	dashboard, err := client.FindDashboard(testID)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if dashboard.ID != "2c5bLca8e" {
		t.Error("request sends json including id but: ", dashboard)
	}

	if dashboard.Title != "My Custom Dashboard(Current)" {
		t.Error("request sends json including title but: ", dashboard)
	}

	if dashboard.URLPath != "2u4PP3TJqbv" {
		t.Error("request sends json including urlpath but: ", dashboard)
	}

	if dashboard.CreatedAt != 1552909732 {
		t.Error("request sends json including createdAt but: ", dashboard)
	}

	if dashboard.UpdatedAt != 1552992837 {
		t.Error("request sends json including updatedAt but: ", dashboard)
	}

	if dashboard.Memo != "A test Current Dashboard" {
		t.Error("request sends json including memo but:", dashboard)
	}

	// Widget Test : Widget(Common) && Markdown && Layout(Common)
	if dashboard.Widgets[0].Type != "markdown" {
		t.Error("request sends json including widgets.type but:", dashboard)
	}

	if dashboard.Widgets[0].Title != "markdown_widget" {
		t.Error("request sends json including widgets.title but:", dashboard)
	}

	if dashboard.Widgets[0].Markdown != "# body" {
		t.Error("request sends json including widgets.markdown but:", dashboard)
	}

	if dashboard.Widgets[0].Layout.X != 0 {
		t.Error("request sends json including widgets.layout.x but:", dashboard)
	}

	if dashboard.Widgets[0].Layout.Y != 0 {
		t.Error("request sends json including widgets.layout.y  but:", dashboard)
	}

	if dashboard.Widgets[0].Layout.Width != 24 {
		t.Error("request sends json including widgets.layout.width  but:", dashboard)
	}

	if dashboard.Widgets[0].Layout.Height != 3 {
		t.Error("request sends json including widgets.layout.height  but:", dashboard)
	}

	// Widget Test : Graph ( && Host Type)
	if dashboard.Widgets[1].Graph.Type != "host" {
		t.Error("request sends json including widgets.graph.type but:", dashboard)
	}

	if dashboard.Widgets[1].Graph.HostID != "2u4PP3TJqbw" {
		t.Error("request sends json including widgets.graph.hostId but:", dashboard)
	}

	if dashboard.Widgets[1].Graph.Name != "loadavg.loadavg15" {
		t.Error("request sends json including widgets.graph.name but:", dashboard)
	}

	if dashboard.Widgets[1].ReferenceLines[0].Label != "critical" {
		t.Error("request sends json including widgets.graph.referenceLines.label but:", dashboard)
	}

	if dashboard.Widgets[1].ReferenceLines[0].Value != 1.23 {
		t.Error("request sends json including widgets.graph.referenceLines.value but:", dashboard)
	}

	// Widget Test : Metric ( && Expression Type)

	if dashboard.Widgets[2].Metric.Type != "expression" {
		t.Error("request sends json including widgets.metric.type but:", dashboard)
	}

	if dashboard.Widgets[2].Metric.Expression != "alias(scale(\nsum(\n  group(\n    host(2u4PP3TJqbx,loadavg.*)\n  )\n),\n1\n), 'test')" {
		t.Error("request sends json including widgets.metric.expression but:", dashboard)
	}

	if *(dashboard.Widgets[2].FractionSize) != 2 {
		t.Error("request sends json including widgets.fractionsize but:", dashboard)
	}

	if dashboard.Widgets[2].Suffix != "total" {
		t.Error("request sends json including widgets.suffix but:", dashboard)
	}

	if dashboard.Widgets[2].FormatRules[0].Name != "SLO" {
		t.Error("request sends json including widgets.formatRules.name but:", dashboard)
	}

	if dashboard.Widgets[2].FormatRules[0].Threshold != 2.34 {
		t.Error("request sends json including widgets.formatRules.threshold but:", dashboard)
	}

	if dashboard.Widgets[2].FormatRules[0].Operator != ">" {
		t.Error("request sends json including widgets.formatRules.operator but:", dashboard)
	}

	// Widget Test : AlertStatus

	if dashboard.Widgets[3].RoleFullName != "test:dashboard" {
		t.Error("request sends json including widgets.roleFullname but:", dashboard)
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

		body, _ := io.ReadAll(req.Body)

		var data struct {
			ID        string   `json:"id"`
			Title     string   `json:"title"`
			URLPath   string   `json:"urlPath"`
			CreatedAt int64    `json:"createdAt"`
			UpdatedAt int64    `json:"updatedAt"`
			Memo      string   `json:"memo"`
			Widgets   []Widget `json:"widgets"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":        "2c5bLca8d",
			"title":     "My Dashboard",
			"urlPath":   "2u4PP3TJqbu",
			"createdAt": 1439346145003,
			"updatedAt": 1439346145003,
			"memo":      "My Dashboard Memo",
			"widgets":   []map[string]interface{}{},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	dashboard, err := client.CreateDashboard(&Dashboard{
		Title:   "My Dashboard",
		URLPath: "2u4PP3TJqbu",
		Memo:    "My Dashboard Memo",
		Widgets: []Widget{},
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

	if dashboard.URLPath != "2u4PP3TJqbu" {
		t.Error("request sends json including urlpath but: ", dashboard)
	}

	if dashboard.CreatedAt != 1439346145003 {
		t.Error("request sends json including createdAt but: ", dashboard)
	}

	if dashboard.UpdatedAt != 1439346145003 {
		t.Error("request sends json including updatedAt but: ", dashboard)
	}

	if dashboard.Memo != "My Dashboard Memo" {
		t.Error("request sends json including memo but: ", dashboard)
	}

	if len(dashboard.Widgets) != 0 {
		t.Error("request sends json including widgets but: ", dashboard)
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

		body, _ := io.ReadAll(req.Body)

		var data struct {
			ID        string   `json:"id"`
			Title     string   `json:"title"`
			URLPath   string   `json:"urlPath"`
			CreatedAt int64    `json:"createdAt"`
			UpdatedAt int64    `json:"updatedAt"`
			Memo      string   `json:"memo"`
			Widgets   []Widget `json:"widgets"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":        "2c5bLca8d",
			"title":     "My Dashboard",
			"urlPath":   "2u4PP3TJqbu",
			"createdAt": 1439346145003,
			"updatedAt": 1439346145003,
			"memo":      "My Dashboard Memo",
			"widgets":   []map[string]interface{}{},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	dashboard, err := client.UpdateDashboard(testID, &Dashboard{
		Title:   "My Dashboard",
		URLPath: "2u4PP3TJqbu",
		Memo:    "My Dashboard Memo",
		Widgets: []Widget{},
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

	if dashboard.URLPath != "2u4PP3TJqbu" {
		t.Error("request sends json including urlpath but: ", dashboard)
	}

	if dashboard.CreatedAt != 1439346145003 {
		t.Error("request sends json including createdAt but: ", dashboard)
	}

	if dashboard.UpdatedAt != 1439346145003 {
		t.Error("request sends json including updatedAt but: ", dashboard)
	}

	if dashboard.Memo != "My Dashboard Memo" {
		t.Error("request sends json including memo but: ", dashboard)
	}

	if len(dashboard.Widgets) != 0 {
		t.Error("request sends json including widgets but: ", dashboard)
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
			"id":        "2c5bLca8d",
			"title":     "My Dashboard",
			"urlPath":   "2u4PP3TJqbu",
			"createdAt": 1439346145003,
			"updatedAt": 1439346145003,
			"memo":      "My Dashboard Memo",
			"widgets":   []map[string]interface{}{},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
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

	if dashboard.URLPath != "2u4PP3TJqbu" {
		t.Error("request sends json including urlpath but: ", dashboard)
	}

	if dashboard.CreatedAt != 1439346145003 {
		t.Error("request sends json including createdAt but: ", dashboard)
	}

	if dashboard.UpdatedAt != 1439346145003 {
		t.Error("request sends json including updatedAt but: ", dashboard)
	}

	if dashboard.Memo != "My Dashboard Memo" {
		t.Error("request sends json including memo but: ", dashboard)
	}

	if len(dashboard.Widgets) != 0 {
		t.Error("request sends json including widgets but: ", dashboard)
	}
}

func TestRangeMarshalJSON(t *testing.T) {
	tests := []struct {
		r    Range
		want string
	}{
		{
			r: Range{
				Type:  "absolute",
				Start: 100,
				End:   200,
			},
			want: `{"type":"absolute","start":100,"end":200}`,
		},
		{
			r: Range{
				Type:   "relative",
				Period: 100,
				Offset: 0,
			},
			want: `{"type":"relative","period":100,"offset":0}`,
		},
		{
			r:    Range{},
			want: `null`,
		},
	}
	for _, tt := range tests {
		b, err := json.Marshal(tt.r)
		if err != nil {
			t.Fatal(err)
		}
		if s := string(b); s != tt.want {
			t.Errorf("MarshalJSON(%v) = %q; want %q", tt.r, s, tt.want)
		}
	}
}

func TestMetricMarshalJSON(t *testing.T) {
	tests := []struct {
		m    Metric
		want string
	}{
		{
			m: Metric{
				Type:   "host",
				Name:   "loadavg5",
				HostID: "0000",
			},
			want: `{"type":"host","name":"loadavg5","hostId":"0000"}`,
		},
		{
			m: Metric{
				Type:        "service",
				Name:        "metric0",
				ServiceName: "service0",
			},
			want: `{"type":"service","name":"metric0","serviceName":"service0"}`,
		},
		{
			m: Metric{
				Type:       "expression",
				Expression: "max(role(my-service:db, loadavg5))",
			},
			want: `{"type":"expression","expression":"max(role(my-service:db, loadavg5))"}`,
		},
		{
			m: Metric{
				Type:   "query",
				Query:  "up{}",
				Legend: "",
			},
			want: `{"type":"query","query":"up{}","legend":""}`,
		},
		{
			m:    Metric{},
			want: "null",
		},
	}
	for _, tt := range tests {
		b, err := json.Marshal(tt.m)
		if err != nil {
			t.Fatal(err)
		}
		if s := string(b); s != tt.want {
			t.Errorf("MarshalJSON(%v) = %q; want %q", tt.m, s, tt.want)
		}
	}
}

func TestGraphMarshalJSON(t *testing.T) {
	tests := []struct {
		g    Graph
		want string
	}{
		{
			g: Graph{
				Type:   "host",
				Name:   "loadavg5",
				HostID: "0000",
			},
			want: `{"type":"host","name":"loadavg5","hostId":"0000"}`,
		},
		{
			g: Graph{
				Type:         "role",
				Name:         "loadavg5",
				RoleFullName: "my-service:db",
				IsStacked:    true,
			},
			want: `{"type":"role","name":"loadavg5","roleFullname":"my-service:db","isStacked":true}`,
		},
		{
			g: Graph{
				Type:        "service",
				Name:        "my-metric",
				ServiceName: "my-service",
			},
			want: `{"type":"service","name":"my-metric","serviceName":"my-service"}`,
		},
		{
			g: Graph{
				Type:       "expression",
				Expression: "max(role(my-service:db, loadavg5))",
			},
			want: `{"type":"expression","expression":"max(role(my-service:db, loadavg5))"}`,
		},
		{
			g: Graph{
				Type:   "query",
				Query:  "up{}",
				Legend: "{{instance}}",
			},
			want: `{"type":"query","query":"up{}","legend":"{{instance}}"}`,
		},
		{
			g:    Graph{},
			want: "null",
		},
	}
	for _, tt := range tests {
		b, err := json.Marshal(tt.g)
		if err != nil {
			t.Fatal(err)
		}
		if s := string(b); s != tt.want {
			t.Errorf("MarshalJSON(%v) = %q; want %q", tt.g, s, tt.want)
		}
	}
}
