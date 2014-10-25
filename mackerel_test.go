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

func TestFindHosts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts.json" {
			t.Error("request URL should be /api/v0/hosts.json but :", req.URL.Path)
		}

		query := req.URL.Query()
		if query.Get("service") != "My-Service" {
			t.Error("request query 'service' param should be My-Service but :", query.Get("service"))
		}
		if !reflect.DeepEqual(query["role"], []string{"db-master"}) {
			t.Error("request query 'role' param should be db-master but :", query.Get("role"))
		}
		if query.Get("name") != "mydb001" {
			t.Error("request query 'name' param should be mydb001 but :", query.Get("name"))
		}
		if !reflect.DeepEqual(query["status"], []string{"working", "standby"}) {
			t.Error("request query 'statuses' param should be ['working','standby'] but :", query["status"])
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but :", req.Method)
		}

		respJson, _ := json.Marshal(map[string][]map[string]interface{}{
			"hosts": []map[string]interface{}{
				map[string]interface{}{
					"id":     "9rxGOHfVF8F",
					"name":   "mydb001",
					"status": "working",
					"memo":   "hello",
					"roles":  map[string][]string{"My-Service": []string{"db-master", "db-slave"}},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer ts.Close()

	client, _ := NewClientForTest("dummy-key", ts.URL, false)
	hosts, err := client.FindHosts(&FindHostsParam{
		Service:  "My-Service",
		Roles:    []string{"db-master"},
		Statuses: []string{"working", "standby"},
		Name:     "mydb001",
	})

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	if hosts[0].Memo != "hello" {
		t.Error("request sends json including memo but: ", hosts[0])
	}

	if reflect.DeepEqual(hosts[0].Roles["My-Service"], []string{"db-master", "db-slave"}) != true {
		t.Errorf("Wrong data for roles: %v", hosts[0].Roles)
	}

}
