package mackerel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindRoles(t *testing.T) {
	testServiceName := "testService"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		uri := fmt.Sprintf("/api/v0/services/%s/roles", testServiceName)
		if req.URL.Path != uri {
			t.Error("request URL should be ", uri, " but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"roles": {
				{
					"name": "My-Role",
					"memo": "hello",
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	roles, err := client.FindRoles(testServiceName)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if roles[0].Memo != "hello" {
		t.Error("request sends json including memo but: ", roles[0])
	}
}

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

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Error(err.Error())
		}

		var reqBody CreateRoleParam
		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			t.Error(err.Error())
		}

		if reqBody.Name != testRoleName {
			t.Error("name (in request json parameter) should be ", testRoleName, "but: ", reqBody.Name)
		}

		if reqBody.Memo != testRoleMemo {
			t.Error("memo (in request json parameter) should be ", testRoleMemo, "but: ", reqBody.Memo)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"name": testRoleName,
			"memo": testRoleMemo,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON)) // nolint
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
		fmt.Fprint(res, string(respJSON)) // nolint
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
