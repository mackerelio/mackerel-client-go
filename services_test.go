package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFindServices(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/services" {
			t.Error("request URL should be /api/v0/services but :", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but :", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"services": {
				{
					"name":  "My-Service",
					"memo":  "hello",
					"roles": []string{"db-master", "db-slave"},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	services, err := client.FindServices()

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	if services[0].Memo != "hello" {
		t.Error("request sends json including memo but: ", services[0])
	}

	if reflect.DeepEqual(services[0].Roles, []string{"db-master", "db-slave"}) != true {
		t.Errorf("Wrong data for roles: %v", services[0].Roles)
	}

}
