package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRole(t *testing.T) {
	testServiceName := "testService"
	testRoleName := "testRole"
	testRoleMemo := "this role is test"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		uri := fmt.Sprintf("/api/v0/services/%s/roles", testServiceName)
		if req.URL.Path != uri {
			t.Error("request URL should be ", uri, " but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"name": testRoleName,
			"memo": testRoleMemo,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	role, err := client.CreateRole(
		testServiceName,
		&CreateRoleParam{
			Name: testRoleName,
			Memo: testRoleMemo,
		})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if role.Name != testRoleName {
		t.Error("request sends json including name but: ", role.Name)
	}

	if role.Memo != testRoleMemo {
		t.Error("request sends json including memo but: ", role.Memo)
	}
}

func TestDeleteRole(t *testing.T) {
	testServiceName := "testService"
	testRoleName := "testRole"
	testRoleMemo := "this role is test"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		uri := fmt.Sprintf("/api/v0/services/%s/roles/%s", testServiceName, testRoleName)
		if req.URL.Path != uri {
			t.Error("request URL should be ", uri, " but: ", req.URL.Path)
		}

		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"name": testRoleName,
			"memo": testRoleMemo,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	role, err := client.DeleteRole(testServiceName, testRoleName)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if role.Name != testRoleName {
		t.Error("request sends json including name but: ", role.Name)
	}

	if role.Memo != testRoleMemo {
		t.Error("request sends json including memo but: ", role.Memo)
	}
}
