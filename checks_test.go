package mackerel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindCheckMonitors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/monitoring/checks" {
			t.Error("request URL should be /api/v0/monitoring/checks but: ", req.URL.Path)
		}
		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]any{
			"checks": []map[string]string{
				{"id": "checkId1", "name": "check1"},
				{"id": "checkId2", "name": "check2"},
			},
			"nextId": "checkId3",
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	cli, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	resp, err := cli.FindCheckMonitors(nil)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if resp.Checks[0].ID != "checkId1" {
		t.Error("first check ID should be checkId1 but: ", resp.Checks[0].ID)
	}
	if resp.Checks[0].Name != "check1" {
		t.Error("first check name should be check1 but: ", resp.Checks[0].Name)
	}
	if resp.NextID != "checkId3" {
		t.Error("nextId should be checkId3 but: ", resp.NextID)
	}
}

func TestFindCheckMonitorsWithParams(t *testing.T) {
	nextID := "checkId3"
	limit := 10
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/monitoring/checks" {
			t.Error("request URL should be /api/v0/monitoring/checks but: ", req.URL.Path)
		}
		if req.URL.Query().Get("nextId") != nextID {
			t.Error("nextId should be checkId3 but: ", req.URL.Query().Get("nextId"))
		}
		if req.URL.Query().Get("limit") != "10" {
			t.Error("limit should be 10 but: ", req.URL.Query().Get("limit"))
		}

		respJSON, _ := json.Marshal(map[string]any{
			"checks": []map[string]string{
				{"id": "checkId3", "name": "check3"},
			},
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	cli, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	resp, err := cli.FindCheckMonitors(&FindCheckMonitorsParam{NextID: &nextID, Limit: &limit})
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if resp.Checks[0].ID != "checkId3" {
		t.Error("check ID should be checkId3 but: ", resp.Checks[0].ID)
	}
}

func TestCheckReports_MarshalJSON(t *testing.T) {
	crs := &CheckReports{
		Reports: []*CheckReport{
			{
				Source:     NewCheckSourceHost("hogee"),
				Name:       "chchch",
				Status:     CheckStatusCritical,
				OccurredAt: 100,
				Message:    "OKOK",
			},
		},
	}
	expect := `{"reports":[{"source":{"type":"host","hostId":"hogee"},"name":"chchch","status":"CRITICAL","message":"OKOK","occurredAt":100}]}`
	bs, _ := json.Marshal(crs)
	got := string(bs)

	if got != expect {
		t.Errorf("expect: %s, but: %s", expect, got)
	}
}

func TestClient_PostCheckReports(t *testing.T) {
	crs := &CheckReports{
		Reports: []*CheckReport{
			{
				Source:     NewCheckSourceHost("hogee"),
				Name:       "chchch",
				Status:     CheckStatusCritical,
				OccurredAt: 100,
				Message:    "OKOK",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqPath := "/api/v0/monitoring/checks/report"
		if req.URL.Path != reqPath {
			t.Errorf("request URL should be %s but: %s", reqPath, req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := io.ReadAll(req.Body)

		var values struct {
			Reports []struct {
				Source     any         `json:"source"`
				Name       string      `json:"name"`
				Status     CheckStatus `json:"status"`
				Message    string      `json:"message"`
				OccurredAt int64       `json:"occurredAt"`
			} `json:"reports"`
		}

		err := json.Unmarshal(body, &values)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		r := values.Reports[0]
		if r.Name != "chchch" {
			t.Error("request sends json including hostId but: ", r.Name)
		}
		if r.OccurredAt != 100 {
			t.Error("request sends json including time but: ", r.OccurredAt)
		}
		if r.Status != CheckStatusCritical {
			t.Error("request sends json including value but: ", r.Status)
		}
		if r.Message != "OKOK" {
			t.Error("request sends json including value but: ", r.Message)
		}

		respJSON, _ := json.Marshal(map[string]bool{
			"success": true,
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	cli, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	err := cli.PostCheckReports(crs)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

}
