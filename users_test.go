package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/users" {
			t.Error("request URL should be but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"users": {
				{
					"id":                      "ABCDEFGHIJK",
					"screenName":              "myname",
					"email":                   "test@example.com",
					"authority":               "viewer",
					"isInRegistrationProcess": true,
					"isMFAEnabled":            true,
					"authenticationMethods":   []string{"password"},
					"joinedAt":                1560000000,
				},
			},
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	users, err := client.FindUsers()

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if users[0].ID != "ABCDEFGHIJK" {
		t.Error("request sends json including id but: ", users[0].ID)
	}

	if users[0].ScreenName != "myname" {
		t.Error("request sends json including screenName but: ", users[0].ScreenName)
	}

	if users[0].Email != "test@example.com" {
		t.Error("request sends json including email but: ", users[0].Email)
	}

	if users[0].Authority != "viewer" {
		t.Error("request sends json including authority but: ", users[0].Authority)
	}

	if users[0].IsInRegistrationProcess != true {
		t.Error("request sends json including isInRegistrationProcess but: ", users[0].IsInRegistrationProcess)
	}

	if users[0].IsMFAEnabled != true {
		t.Error("request sends json including isMFAEnabled but: ", users[0].IsMFAEnabled)
	}

	if reflect.DeepEqual(users[0].AuthenticationMethods, []string{"password"}) != true {
		t.Errorf("Wrong data for users: %v", users[0].AuthenticationMethods)
	}

	if users[0].JoinedAt != 1560000000 {
		t.Error("request sends json including joinedAt but: ", users[0].JoinedAt)
	}
}

func TestDeleteUser(t *testing.T) {

	testUserID := "ABCDEFGHIJK"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/users/%s", testUserID) {
			t.Error("request URL should be /api/v0/users/<userId> but: ", req.URL.Path)
		}

		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":                      "ABCDEFGHIJK",
			"screenName":              "myname",
			"email":                   "test@example.com",
			"authority":               "viewer",
			"isInRegistrationProcess": true,
			"isMFAEnabled":            true,
			"authenticationMethods":   []string{"password"},
			"joinedAt":                1560000000,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	user, err := client.DeleteUser(testUserID)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if user.ID != "ABCDEFGHIJK" {
		t.Error("request sends json including id but: ", user.ID)
	}

	if user.ScreenName != "myname" {
		t.Error("request sends json including screenName but: ", user.ScreenName)
	}

	if user.Email != "test@example.com" {
		t.Error("request sends json including email but: ", user.Email)
	}

	if user.Authority != "viewer" {
		t.Error("request sends json including authority but: ", user.Authority)
	}

	if user.IsInRegistrationProcess != true {
		t.Error("request sends json including isInRegistrationProcess but: ", user.IsInRegistrationProcess)
	}

	if user.IsMFAEnabled != true {
		t.Error("request sends json including isMFAEnabled but: ", user.IsMFAEnabled)
	}

	if reflect.DeepEqual(user.AuthenticationMethods, []string{"password"}) != true {
		t.Errorf("Wrong data for users: %v", user.AuthenticationMethods)
	}

	if user.JoinedAt != 1560000000 {
		t.Error("request sends json including joinedAt but: ", user.JoinedAt)
	}
}
