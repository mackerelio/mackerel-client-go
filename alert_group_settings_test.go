package mackerel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestFindAlertGroupSettings(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/alert-group-settings" {
			t.Error("request URL should be /api/v0/alert-group-settings but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"alertGroupSettings": {
				{
					"id":   "xxxxxxxxxxx",
					"name": "alert group setting #1",
				},
				{
					"id":                   "yyyyyyyyyyy",
					"name":                 "alert group setting #2",
					"memo":                 "lorem ipsum...",
					"serviceScopes":        []string{"my-service"},
					"roleScopes":           []string{"my-service: db"},
					"monitorScopes":        []string{"connectivity"},
					"notificationInterval": 60,
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	got, err := client.FindAlertGroupSettings()
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	want := []*AlertGroupSetting{
		{
			ID:   "xxxxxxxxxxx",
			Name: "alert group setting #1",
		},
		{
			ID:                   "yyyyyyyyyyy",
			Name:                 "alert group setting #2",
			Memo:                 "lorem ipsum...",
			ServiceScopes:        []string{"my-service"},
			RoleScopes:           []string{"my-service: db"},
			MonitorScopes:        []string{"connectivity"},
			NotificationInterval: 60,
		},
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("fail to get correct data: diff: (-got +want)\n%s", diff)
	}
}

func TestCreateAlertGroupSetting(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/alert-group-settings" {
			t.Error("request URL should be /api/v0/alert-group-settings but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := io.ReadAll(req.Body)
		var alertGroupSetting AlertGroupSetting
		if err := json.Unmarshal(body, &alertGroupSetting); err != nil {
			t.Fatal("request body should be decoded as json ", string(body))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":                   "xxxxxxxxxxx",
			"name":                 alertGroupSetting.Name,
			"memo":                 alertGroupSetting.Memo,
			"serviceScopes":        alertGroupSetting.ServiceScopes,
			"roleScopes":           alertGroupSetting.RoleScopes,
			"monitorScopes":        alertGroupSetting.MonitorScopes,
			"notificationInterval": alertGroupSetting.NotificationInterval,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	param := &AlertGroupSetting{
		Name: "alert group setting",
		Memo: "lorem ipsum...",
	}
	got, err := client.CreateAlertGroupSetting(param)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	want := &AlertGroupSetting{
		ID:   "xxxxxxxxxxx",
		Name: param.Name,
		Memo: param.Memo,
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("fail to get correct data: diff: (-got +want)\n%s", diff)
	}
}

func TestGetAlertGroupSetting(t *testing.T) {
	id := "xxxxxxxxxxx"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/alert-group-settings/%s", id) {
			t.Error("request URL should be /api/v0/alert-group-settings/<ID> but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":                   "xxxxxxxxxxx",
			"name":                 "alert group setting",
			"memo":                 "lorem ipsum...",
			"serviceScopes":        []string{"my-service"},
			"roleScopes":           []string{"my-service: db"},
			"monitorScopes":        []string{"connectivity"},
			"notificationInterval": 60,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	got, err := client.GetAlertGroupSetting(id)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	want := &AlertGroupSetting{
		ID:                   id,
		Name:                 "alert group setting",
		Memo:                 "lorem ipsum...",
		ServiceScopes:        []string{"my-service"},
		RoleScopes:           []string{"my-service: db"},
		MonitorScopes:        []string{"connectivity"},
		NotificationInterval: 60,
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("fail to get correct data: diff (-got +want)\n%s", diff)
	}
}

func TestUpdateAlertGroupSetting(t *testing.T) {
	id := "xxxxxxxxxxx"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/alert-group-settings/%s", id) {
			t.Error("request URL should be /api/v0/alert-group-settings/<ID> but: ", req.URL.Path)
		}

		if req.Method != "PUT" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		body, _ := io.ReadAll(req.Body)
		var alertGroupSetting AlertGroupSetting
		if err := json.Unmarshal(body, &alertGroupSetting); err != nil {
			t.Fatal("request body should be decoded as json ", string(body))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":                   id,
			"name":                 alertGroupSetting.Name,
			"memo":                 "lorem ipsum...",
			"serviceScopes":        []string{"my-service"},
			"roleScopes":           []string{"my-service: db"},
			"monitorScopes":        []string{"connectivity"},
			"notificationInterval": 60,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	param := &AlertGroupSetting{
		Name: "alert group notification updated",
	}
	got, err := client.UpdateAlertGroupSetting(id, param)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	want := &AlertGroupSetting{
		ID:                   id,
		Name:                 param.Name,
		Memo:                 "lorem ipsum...",
		ServiceScopes:        []string{"my-service"},
		RoleScopes:           []string{"my-service: db"},
		MonitorScopes:        []string{"connectivity"},
		NotificationInterval: 60,
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("fail to get correct data: diff: (-got +want)\n%s", diff)
	}
}

func TestDeleteAlertGroupSetting(t *testing.T) {
	id := "xxxxxxxxxxx"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/alert-group-settings/%s", id) {
			t.Error("request URL should be /api/v0/alert-group-settings/<ID> but: ", req.URL.Path)
		}

		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":                   id,
			"name":                 "alert group setting",
			"memo":                 "lorem ipsum...",
			"serviceScopes":        []string{"my-service"},
			"roleScopes":           []string{"my-service: db"},
			"monitorScopes":        []string{"connectivity"},
			"notificationInterval": 60,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	got, err := client.DeleteAlertGroupSetting(id)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	want := &AlertGroupSetting{
		ID:                   id,
		Name:                 "alert group setting",
		Memo:                 "lorem ipsum...",
		ServiceScopes:        []string{"my-service"},
		RoleScopes:           []string{"my-service: db"},
		MonitorScopes:        []string{"connectivity"},
		NotificationInterval: 60,
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("fail to get correct data: diff (-got +want)\n%s", diff)
	}
}
