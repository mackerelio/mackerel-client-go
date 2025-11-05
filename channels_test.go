package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// boolPointer is a helper function to initialize a bool pointer
func boolPointer(b bool) *bool {
	return &b
}

func TestFindChannels(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/channels" {
			t.Error("request URL should be /api/v0/channels but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"channels": {
				{
					"id":      "abcdefabc",
					"name":    "email channel",
					"type":    "email",
					"emails":  []string{"test@example.com", "test2@example.com"},
					"userIds": []string{"1234", "2345"},
					"events":  []string{"alert"},
				},
				{
					"id":   "bcdefabcd",
					"name": "slack channel",
					"type": "slack",
					"url":  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
					"mentions": map[string]interface{}{
						"ok":      "ok message",
						"warning": "warning message",
					},
					"enabledGraphImage": true,
					"events":            []string{"alert"},
				},
				{
					"id":     "cdefabcde",
					"name":   "webhook channel",
					"type":   "webhook",
					"url":    "http://example.com/webhook",
					"events": []string{"alertGroup"},
				},
				{
					"id":          "defabcdef",
					"name":        "line channel",
					"type":        "line",
					"suspendedAt": 12345678,
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	channels, err := client.FindChannels()

	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if len(channels) != 4 {
		t.Error("request has 4 channels but: ", len(channels))
	}

	if channels[0].ID != "abcdefabc" {
		t.Error("request has ID but: ", channels[0].ID)
	}
	if channels[1].ID != "bcdefabcd" {
		t.Error("request has ID but: ", channels[1].ID)
	}
	if channels[2].ID != "cdefabcde" {
		t.Error("request has ID but: ", channels[2].ID)
	}
	if channels[3].ID != "defabcdef" {
		t.Error("request has ID but: ", channels[3].ID)
	}
	if channels[0].Name != "email channel" {
		t.Error("request has Name but: ", channels[0].Name)
	}
	if channels[1].Name != "slack channel" {
		t.Error("request has Name but: ", channels[1].Name)
	}
	if channels[2].Name != "webhook channel" {
		t.Error("request has Name but: ", channels[2].Name)
	}
	if channels[3].Name != "line channel" {
		t.Error("request has Name but: ", channels[3].Name)
	}
	if channels[0].Type != "email" {
		t.Error("request has Type but: ", channels[0].Type)
	}
	if channels[1].Type != "slack" {
		t.Error("request has Type but: ", channels[1].Type)
	}
	if channels[2].Type != "webhook" {
		t.Error("request has Type but: ", channels[2].Type)
	}
	if channels[3].Type != "line" {
		t.Error("request has Type but: ", channels[3].Type)
	}
	if channels[0].SuspendedAt != nil {
		t.Error("request has SuspendedAt but: ", channels[0].SuspendedAt)
	}
	if channels[1].SuspendedAt != nil {
		t.Error("request has SuspendedAt but: ", channels[1].SuspendedAt)
	}
	if channels[2].SuspendedAt != nil {
		t.Error("request has SuspendedAt but: ", channels[2].SuspendedAt)
	}
	if *channels[3].SuspendedAt != 12345678 {
		t.Error("request has SuspendedAt but: ", *channels[3].SuspendedAt)
	}
	if reflect.DeepEqual(*(channels[0].Emails), []string{"test@example.com", "test2@example.com"}) != true {
		t.Errorf("Wrong data for emails: %v", *(channels[0].Emails))
	}
	if channels[1].Emails != nil {
		t.Errorf("Wrong data for emails: %v", *(channels[1].Emails))
	}
	if channels[2].Emails != nil {
		t.Errorf("Wrong data for emails: %v", *(channels[2].Emails))
	}
	if channels[3].Emails != nil {
		t.Errorf("Wrong data for emails: %v", *(channels[3].Emails))
	}
	if reflect.DeepEqual(*(channels[0].UserIDs), []string{"1234", "2345"}) != true {
		t.Errorf("Wrong data for userIds: %v", *(channels[0].UserIDs))
	}
	if channels[1].UserIDs != nil {
		t.Errorf("Wrong data for userIds: %v", *(channels[1].UserIDs))
	}
	if channels[2].UserIDs != nil {
		t.Errorf("Wrong data for userIds: %v", *(channels[2].UserIDs))
	}
	if channels[3].UserIDs != nil {
		t.Errorf("Wrong data for userIds: %v", *(channels[3].UserIDs))
	}
	if reflect.DeepEqual(*(channels[0].Events), []string{"alert"}) != true {
		t.Errorf("Wrong data for events: %v", *(channels[0].Events))
	}
	if reflect.DeepEqual(*(channels[1].Events), []string{"alert"}) != true {
		t.Errorf("Wrong data for events: %v", *(channels[1].Events))
	}
	if reflect.DeepEqual(*(channels[2].Events), []string{"alertGroup"}) != true {
		t.Errorf("Wrong data for events: %v", *(channels[2].Events))
	}
	if channels[3].Events != nil {
		t.Errorf("Wrong data for events: %v", *(channels[3].Events))
	}
	if channels[0].URL != "" {
		t.Error("request has no URL but: ", channels[0])
	}
	if channels[1].URL != "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX" {
		t.Error("request sends json including URL but: ", channels[1])
	}
	if channels[2].URL != "http://example.com/webhook" {
		t.Error("request sends json including URL but: ", channels[2])
	}
	if channels[3].URL != "" {
		t.Error("request has no URL but: ", channels[3])
	}
	if reflect.DeepEqual(channels[0].Mentions, Mentions{}) != true {
		t.Error("request has mentions but: ", channels[0].Mentions)
	}
	if reflect.DeepEqual(channels[1].Mentions, Mentions{OK: "ok message", Warning: "warning message"}) != true {
		t.Error("request has mentions but: ", channels[1].Mentions)
	}
	if reflect.DeepEqual(channels[2].Mentions, Mentions{}) != true {
		t.Error("request has mentions but: ", channels[2].Mentions)
	}
	if reflect.DeepEqual(channels[3].Mentions, Mentions{}) != true {
		t.Error("request has mentions but: ", channels[3].Mentions)
	}
	if channels[0].EnabledGraphImage != nil {
		t.Error("request sends json including enabledGraphImage but: ", *(channels[0].EnabledGraphImage))
	}
	if !*(channels[1].EnabledGraphImage) {
		t.Error("request sends json including enabledGraphImage but: ", *(channels[1].EnabledGraphImage))
	}
	if channels[2].EnabledGraphImage != nil {
		t.Error("request sends json including enabledGraphImage but: ", *(channels[2].EnabledGraphImage))
	}
	if channels[3].EnabledGraphImage != nil {
		t.Error("request sends json including enabledGraphImage but: ", *(channels[3].EnabledGraphImage))
	}
}

func TestCreateChannel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/channels" {
			t.Error("request URL should be /api/v0/channels but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":   "abcdefabc",
			"name": "slack channel",
			"type": "slack",
			"url":  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
			"mentions": map[string]interface{}{
				"ok":      "ok message",
				"warning": "warning message",
			},
			"enabledGraphImage": true,
			"events":            []string{"alert"},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	channel, err := client.CreateChannel(&Channel{
		Name: "slack channel",
		Type: "slack",
		URL:  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
		Mentions: Mentions{
			OK:      "ok message",
			Warning: "warning message",
		},
		EnabledGraphImage: boolPointer(true),
		Events:            &[]string{"alert"},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if channel.ID != "abcdefabc" {
		t.Error("request sends json including ID but: ", channel.ID)
	}
	if channel.Name != "slack channel" {
		t.Error("request sends json including name but: ", channel.Name)
	}
	if channel.URL != "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX" {
		t.Error("request sends json including URL but: ", channel.URL)
	}
	if reflect.DeepEqual(channel.Mentions, Mentions{OK: "ok message", Warning: "warning message"}) != true {
		t.Errorf("Wrong data for mentions: %v", channel.Mentions)
	}
	if !*channel.EnabledGraphImage {
		t.Error("request sends json including enabledGraphImage but: ", *channel.EnabledGraphImage)
	}
	if reflect.DeepEqual(*(channel.Events), []string{"alert"}) != true {
		t.Errorf("Wrong data for events: %v", *(channel.Events))
	}
}

func TestDeleteChannel(t *testing.T) {
	channelID := "abcdefabc"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/channels/%s", channelID) {
			t.Error("request URL should be /api/v0/channels/<ID> but: ", req.URL.Path)
		}

		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":   channelID,
			"name": "slack channel",
			"type": "slack",
			"url":  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
			"mentions": map[string]interface{}{
				"ok":      "ok message",
				"warning": "warning message",
			},
			"enabledGraphImage": true,
			"events":            []string{"alert"},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	channel, err := client.DeleteChannel(channelID)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if channel.ID != "abcdefabc" {
		t.Error("request sends json including ID but: ", channel.ID)
	}
	if channel.Name != "slack channel" {
		t.Error("request sends json including name but: ", channel.Name)
	}
	if channel.URL != "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX" {
		t.Error("request sends json including URL but: ", channel.URL)
	}
	if reflect.DeepEqual(channel.Mentions, Mentions{OK: "ok message", Warning: "warning message"}) != true {
		t.Errorf("Wrong data for mentions: %v", channel.Mentions)
	}
	if !*channel.EnabledGraphImage {
		t.Error("request sends json including enabledGraphImage but: ", *channel.EnabledGraphImage)
	}
	if reflect.DeepEqual(*(channel.Events), []string{"alert"}) != true {
		t.Errorf("Wrong data for events: %v", *(channel.Events))
	}
}
