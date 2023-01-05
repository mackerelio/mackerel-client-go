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

func TestCreateNotificationGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/notification-groups" {
			t.Error("request URL should be /api/v0/notification-groups but: ", req.URL.Path)
		}
		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := io.ReadAll(req.Body)

		var notificationGroup NotificationGroup
		if err := json.Unmarshal(body, &notificationGroup); err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":                        "3JwREyrZGQ9",
			"name":                      notificationGroup.Name,
			"notificationLevel":         notificationGroup.NotificationLevel,
			"childNotificationGroupIds": notificationGroup.ChildNotificationGroupIDs,
			"childChannelIds":           notificationGroup.ChildChannelIDs,
			"monitors":                  notificationGroup.Monitors,
			"services":                  notificationGroup.Services,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		_, _ = fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	param := &NotificationGroup{
		Name:                      "New Notification Group",
		NotificationLevel:         NotificationLevelAll,
		ChildNotificationGroupIDs: []string{"3mUMcLB4ks9", "2w53XJsufQG"},
		ChildChannelIDs:           []string{"2nckL8bKy6E", "2w54SREy99h", "3JwPUGFJw2f"},
		Monitors: []*NotificationGroupMonitor{
			{ID: "2CRrhj1SFwG", SkipDefault: true},
			{ID: "3TdoBWxYRQd", SkipDefault: false},
		},
		Services: []*NotificationGroupService{
			{Name: "my-service-01"},
			{Name: "my-service-02"},
		},
	}

	got, err := client.CreateNotificationGroup(param)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	want := &NotificationGroup{
		ID:                        "3JwREyrZGQ9",
		Name:                      param.Name,
		NotificationLevel:         param.NotificationLevel,
		ChildNotificationGroupIDs: param.ChildNotificationGroupIDs,
		ChildChannelIDs:           param.ChildChannelIDs,
		Monitors:                  param.Monitors,
		Services:                  param.Services,
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("fail to get correct data: diff (-got +want)\n%s", diff)
	}
}

func TestFindNotificationGroups(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/notification-groups" {
			t.Error("request URL should be /api/v0/notification-groups but: ", req.URL.Path)
		}
		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"notificationGroups": {
				{
					"id":                        "3Ja3HG3bTwq",
					"name":                      "Default",
					"notificationLevel":         "all",
					"childNotificationGroupIDs": []string{},
					"childChannelIDs":           []string{"3Ja3HG3VTaA"},
				},
				{
					"id":                        "3UJaU9eREvw",
					"name":                      "Notification Group #01",
					"notificationLevel":         "all",
					"childNotificationGroupIds": []string{"3Tdq1pz9aLm"},
					"childChannelIds":           []string{},
				},
				{
					"id":                        "3Tdq1pz9aLm",
					"name":                      "Notification Group #02",
					"notificationLevel":         "critical",
					"childNotificationGroupIds": []string{},
					"childChannelIds":           []string{},
					"monitors": []map[string]interface{}{
						{"id": "3Ja3HG5Mngw", "skipDefault": false},
					},
					"services": []map[string]string{
						{"name": "my-service-01"},
					},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		_, _ = fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	got, err := client.FindNotificationGroups()
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	want := []*NotificationGroup{
		{
			ID:                        "3Ja3HG3bTwq",
			Name:                      "Default",
			NotificationLevel:         NotificationLevelAll,
			ChildNotificationGroupIDs: []string{},
			ChildChannelIDs:           []string{"3Ja3HG3VTaA"},
		},
		{
			ID:                        "3UJaU9eREvw",
			Name:                      "Notification Group #01",
			NotificationLevel:         NotificationLevelAll,
			ChildNotificationGroupIDs: []string{"3Tdq1pz9aLm"},
			ChildChannelIDs:           []string{},
		},
		{
			ID:                        "3Tdq1pz9aLm",
			Name:                      "Notification Group #02",
			NotificationLevel:         NotificationLevelCritical,
			ChildNotificationGroupIDs: []string{},
			ChildChannelIDs:           []string{},
			Monitors: []*NotificationGroupMonitor{
				{ID: "3Ja3HG5Mngw", SkipDefault: false},
			},
			Services: []*NotificationGroupService{
				{Name: "my-service-01"},
			},
		},
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("fail to get correct data: diff: (-got +want)\n%s", diff)
	}
}

func TestUpdateNotificationGroup(t *testing.T) {
	id := "xxxxxxxxxxx"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/notification-groups/%s", id) {
			t.Error("request URL should be /api/v0/notification-groups/<ID> but: ", req.URL.Path)
		}
		if req.Method != "PUT" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		body, _ := io.ReadAll(req.Body)

		var notificationGroup NotificationGroup
		if err := json.Unmarshal(body, &notificationGroup); err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":                        id,
			"name":                      notificationGroup.Name,
			"notificationLevel":         notificationGroup.NotificationLevel,
			"childNotificationGroupIds": notificationGroup.ChildNotificationGroupIDs,
			"childChannelIds":           notificationGroup.ChildChannelIDs,
			"monitors":                  notificationGroup.Monitors,
			"services":                  notificationGroup.Services,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		_, _ = fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	param := &NotificationGroup{
		Name:                      "New Notification Group",
		NotificationLevel:         NotificationLevelCritical,
		ChildNotificationGroupIDs: []string{"3mUMcLB4ks9", "2w53XJsufQG"},
		ChildChannelIDs:           []string{"2nckL8bKy6E", "2w54SREy99h", "3JwPUGFJw2f"},
		Monitors: []*NotificationGroupMonitor{
			{ID: "2CRrhj1SFwG", SkipDefault: true},
			{ID: "3TdoBWxYRQd", SkipDefault: false},
		},
		Services: []*NotificationGroupService{
			{Name: "my-service-01"},
			{Name: "my-service-02"},
		},
	}

	got, err := client.UpdateNotificationGroup(id, param)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	want := &NotificationGroup{
		ID:                        id,
		Name:                      param.Name,
		NotificationLevel:         param.NotificationLevel,
		ChildNotificationGroupIDs: param.ChildNotificationGroupIDs,
		ChildChannelIDs:           param.ChildChannelIDs,
		Monitors:                  param.Monitors,
		Services:                  param.Services,
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("fail to get correct data: diff (-got +want)\n%s", diff)
	}
}

func TestDeleteNotificationGroup(t *testing.T) {
	id := "xxxxxxxxxxx"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/notification-groups/%s", id) {
			t.Error("request URL should be /api/v0/notification-groups/<ID> but: ", req.URL.Path)
		}
		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":                        id,
			"name":                      "My Notification Group",
			"notificationLevel":         "all",
			"childNotificationGroupIds": []string{},
			"childChannelIds":           []string{},
			"monitors": []map[string]interface{}{
				{"id": "2CRrhj1SFwG", "skipDefault": true},
			},
			"services": []map[string]string{
				{"name": "my-service-01"},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		_, _ = fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	got, err := client.DeleteNotificationGroup(id)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	want := &NotificationGroup{
		ID:                        id,
		Name:                      "My Notification Group",
		NotificationLevel:         NotificationLevelAll,
		ChildNotificationGroupIDs: []string{},
		ChildChannelIDs:           []string{},
		Monitors: []*NotificationGroupMonitor{
			{ID: "2CRrhj1SFwG", SkipDefault: true},
		},
		Services: []*NotificationGroupService{
			{Name: "my-service-01"},
		},
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("fail to get correct data: diff (-got +want)\n%s", diff)
	}
}
