package mackerel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func TestGetRoleFullnames(t *testing.T) {
	host := &Host{
		Roles: Roles{
			"My-Service":  []string{"db-master", "db-slave"},
			"My-Service2": []string{"proxy"},
		},
	}

	fullnames := host.GetRoleFullnames()
	sort.Strings(fullnames)

	if !reflect.DeepEqual(fullnames, []string{"My-Service2:proxy", "My-Service:db-master", "My-Service:db-slave"}) {
		t.Error("RoleFullnames should be ['My-Service2:proxy', My-Service:db-master', 'My-Service:db-slave'] but: ", fullnames)
	}
}

func TestFindHost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/9rxGOHfVF8F" {
			t.Error("request URL should be /api/v0/hosts/9rxGOHfVF8F but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]map[string]interface{}{
			"host": {
				"id":     "9rxGOHfVF8F",
				"name":   "mydb001",
				"status": "working",
				"memo":   "hello",
				"roles":  map[string][]string{"My-Service": {"db-master", "db-slave"}},
				"interfaces": []map[string]interface{}{
					{
						"name":          "lo0",
						"ipAddress":     "127.0.0.1",
						"ipv4Addresses": []string{"127.0.0.1"},
						"ipv6Addresses": []string{"fe80::1"},
						"macAddress":    "02:02:02:02:02:02",
					},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	host, err := client.FindHost("9rxGOHfVF8F")

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if host.Memo != "hello" {
		t.Error("request sends json including memo but: ", host)
	}

	if reflect.DeepEqual(host.Roles["My-Service"], []string{"db-master", "db-slave"}) != true {
		t.Errorf("Wrong data for roles: %v", host.Roles)
	}

	if len(host.Interfaces) == 1 && reflect.DeepEqual(host.Interfaces[0], Interface{
		Name:          "lo0",
		IPAddress:     "127.0.0.1",
		IPv4Addresses: []string{"127.0.0.1"},
		IPv6Addresses: []string{"fe80::1"},
		MacAddress:    "02:02:02:02:02:02",
	}) != true {
		t.Errorf("Wrong data for interfaces: %v", host.Interfaces)
	}
}

func TestFindHosts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts" {
			t.Error("request URL should be /api/v0/hosts but: ", req.URL.Path)
		}

		query := req.URL.Query()
		if query.Get("service") != "My-Service" {
			t.Error("request query 'service' param should be My-Service but: ", query.Get("service"))
		}
		if !reflect.DeepEqual(query["role"], []string{"db-master"}) {
			t.Error("request query 'role' param should be db-master but: ", query.Get("role"))
		}
		if query.Get("name") != "mydb001" {
			t.Error("request query 'name' param should be mydb001 but: ", query.Get("name"))
		}
		if !reflect.DeepEqual(query["status"], []string{"working", "standby"}) {
			t.Error("request query 'statuses' param should be ['working','standby'] but: ", query["status"])
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"hosts": {
				{
					"id":     "9rxGOHfVF8F",
					"name":   "mydb001",
					"status": "working",
					"memo":   "hello",
					"roles":  map[string][]string{"My-Service": {"db-master", "db-slave"}},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	hosts, err := client.FindHosts(&FindHostsParam{
		Service:  "My-Service",
		Roles:    []string{"db-master"},
		Statuses: []string{"working", "standby"},
		Name:     "mydb001",
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if hosts[0].Memo != "hello" {
		t.Error("request sends json including memo but: ", hosts[0])
	}

	if reflect.DeepEqual(hosts[0].Roles["My-Service"], []string{"db-master", "db-slave"}) != true {
		t.Errorf("Wrong data for roles: %v", hosts[0].Roles)
	}

}

func TestCreateHost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts" {
			t.Error("request URL should be /api/v0/hosts but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var data struct {
			Name          string        `json:"name"`
			Meta          HostMeta      `json:"meta"`
			Interfaces    []Interface   `json:"interfaces"`
			RoleFullnames []string      `json:"roleFullnames"`
			Checks        []CheckConfig `json:"checks"`
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
		if !reflect.DeepEqual(data.Checks, []CheckConfig{
			{Name: "mysql", Memo: "check mysql memo"},
			{Name: "nginx", Memo: "check nginx memo"},
		}) {
			t.Error("request sends json including checks but: ", data.Checks)
		}

		respJSON, _ := json.Marshal(map[string]string{
			"id": "123456ABCD",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	hostID, err := client.CreateHost(&CreateHostParam{
		Name:          "mydb002",
		RoleFullnames: []string{"My-Service:db-master", "My-Service:db-slave"},
		Checks: []CheckConfig{
			{Name: "mysql", Memo: "check mysql memo"},
			{Name: "nginx", Memo: "check nginx memo"},
		},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if hostID != "123456ABCD" {
		t.Error("hostID should be empty but: ", hostID)
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
			Name          string        `json:"name"`
			Meta          HostMeta      `json:"meta"`
			Interfaces    []Interface   `json:"interfaces"`
			RoleFullnames []string      `json:"roleFullnames"`
			Checks        []CheckConfig `json:"checks"`
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
		if !reflect.DeepEqual(data.Checks, []CheckConfig{
			{Name: "mysql", Memo: "check mysql memo"},
			{Name: "nginx", Memo: "check nginx memo"},
		}) {
			t.Error("request sends json including checks but: ", data.Checks)
		}

		respJSON, _ := json.Marshal(map[string]string{
			"id": "123456ABCD",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	hostID, err := client.UpdateHost("123456ABCD", &UpdateHostParam{
		Name:          "mydb002",
		RoleFullnames: []string{"My-Service:db-master", "My-Service:db-slave"},
		Checks: []CheckConfig{
			{Name: "mysql", Memo: "check mysql memo"},
			{Name: "nginx", Memo: "check nginx memo"},
		},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if hostID != "123456ABCD" {
		t.Error("hostID should be empty but: ", hostID)
	}
}

func TestUpdateHostStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/9rxGOHfVF8F/status" {
			t.Error("request URL should be /api/v0/hosts/9rxGOHfVF8F/status but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
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

		respJSON, _ := json.Marshal(map[string]bool{
			"success": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	err := client.UpdateHostStatus("9rxGOHfVF8F", "maintenance")

	if err != nil {
		t.Error("err should be nil but: ", err)
	}
}

func TestUpdateHostRoleFullnames(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/9rxGOHfVF8F/role-fullnames" {
			t.Error("request URL should be /api/v0/hosts/9rxGOHfVF8F/role-fullnames but: ", req.URL.Path)
		}

		if req.Method != "PUT" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var data struct {
			RoleFullnames []string `json:"roleFullnames"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		respJSON, _ := json.Marshal(map[string]bool{
			"success": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	err := client.UpdateHostRoleFullnames("9rxGOHfVF8F", []string{"testservice:testrole", "testservice:testrole2"})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}
}

func TestRetireHost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/123456ABCD/retire" {
			t.Error("request URL should be /api/v0/hosts/123456ABCD/retire but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var data interface{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		respJSON, _ := json.Marshal(map[string]bool{
			"success": true,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	err := client.RetireHost("123456ABCD")

	if err != nil {
		t.Error("err should be nil but: ", err)
	}
}

func TestRetireHost_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/123456ABCD/retire" {
			t.Error("request URL should be /api/v0/hosts/123456ABCD/retire but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]string{
			"error": "Host Not Found.",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	err := client.RetireHost("123456ABCD")

	if err == nil {
		t.Error("err should not be nil but: ", err)
	}

	apiErr := err.(*APIError)
	if expected := 404; apiErr.StatusCode != expected {
		t.Errorf("api error StatusCode should be %d but got: %d", expected, apiErr.StatusCode)
	}
	if expected := "API request failed: Host Not Found."; apiErr.Error() != expected {
		t.Errorf("api error string should be %s but got: %s", expected, apiErr.Error())
	}
}

func TestListHostMetricNames(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts/9rxGOHfVF8F/metric-names" {
			t.Error("request URL should be /api/v0/hosts/9rxGOHfVF8F/metric-names but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]string{
			"names": {"loadavg1", "loadavg5"},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	names, err := client.ListHostMetricNames("9rxGOHfVF8F")

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if reflect.DeepEqual(names, []string{"loadavg1", "loadavg5"}) != true {
		t.Errorf("Wrong data for metric names: %v", names)
	}
}
