package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Api-Key") != "dummy-key" {
			t.Error("X-Api-Key header should contains passed key")
		}

		if h := req.Header.Get("User-Agent"); h != userAgent {
			t.Errorf("User-Agent shoud be '%s' but %s", userAgent, h)
		}
	}))
	defer ts.Close()

	client, _ := NewClientForTest("dummy-key", ts.URL, false)

	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)
	client.Request(req)
}

func TestFindHost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/9rxGOHfVF8F" {
			t.Error("request URL should be /api/v0/hosts/9rxGOHfVF8F but :", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but :", req.Method)
		}

		respJson, _ := json.Marshal(map[string]map[string]interface{}{
			"host": map[string]interface{}{
				"id":     "9rxGOHfVF8F",
				"name":   "mydb001",
				"status": "working",
				"memo":   "hello",
				"roles":  map[string][]string{"My-Service": []string{"db-master", "db-slave"}},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer ts.Close()

	client, _ := NewClientForTest("dummy-key", ts.URL, false)
	host, err := client.FindHost("9rxGOHfVF8F")

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	if host.Memo != "hello" {
		t.Error("request sends json including memo but: ", host)
	}

	if reflect.DeepEqual(host.Roles["My-Service"], []string{"db-master", "db-slave"}) != true {
		t.Errorf("Wrong data for roles: %v", host.Roles)
	}

}
