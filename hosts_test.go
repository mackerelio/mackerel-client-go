package mackerel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

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

func TestUpdateHost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/123456ABCD" {
			t.Error("request URL should be /api/v0/hosts/123456ABCD but: ", req.URL.Path)
		}

		if req.Method != "PUT" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var data struct {
			Name          string      `json:"name"`
			Meta          HostMeta    `json:"meta"`
			Interfaces    []Interface `json:"interfaces"`
			RoleFullnames []string    `json:"roleFullnames"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		if data.Name != "mydb002" {
			t.Error("request sends json including name but: ", data.Name)
		}
		if !reflect.DeepEqual(data.RoleFullnames, []string{"My-Service:db-master", "My-Service:db-slave"}) {
			t.Error("request sends json including roleFullnames but: ", data.RoleFullnames)
		}

		respJson, _ := json.Marshal(map[string]string{
			"id": "123456ABCD",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer ts.Close()

	client, _ := NewClientForTest("dummy-key", ts.URL, false)
	hostId, err := client.UpdateHost("123456ABCD", &UpdateHostParam{
		Name:          "mydb002",
		RoleFullnames: []string{"My-Service:db-master", "My-Service:db-slave"},
	})

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	if hostId != "123456ABCD" {
		t.Error("hostId shoud be empty but: ", hostId)
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

func TestUpdateHostStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/9rxGOHfVF8F/status" {
			t.Error("request URL should be /api/v0/hosts/9rxGOHfVF8F/status but :", req.URL.Path)
		}

		if req.Method != "PUT" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var data struct {
			Status string `json:"status"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		if data.Status != "maintenance" {
			t.Error("request sends json including status but: ", data.Status)
		}

		respJson, _ := json.Marshal(map[string]bool{
			"success": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer ts.Close()

	client, _ := NewClientForTest("dummy-key", ts.URL, false)
	err := client.UpdateHostStatus("9rxGOHfVF8F", "maintenance")

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}
}

func TestRetireHost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/123456ABCD/retire" {
			t.Error("request URL should be /api/v0/hosts/123456ABCD/retire but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var data interface{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		respJson, _ := json.Marshal(map[string]bool{
			"success": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer ts.Close()

	client, _ := NewClientForTest("dummy-key", ts.URL, false)
	err := client.RetireHost("123456ABCD")

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}
}
