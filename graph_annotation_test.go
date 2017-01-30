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

func TestCreateGraphAnnotation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/graph-annotations" {
			t.Error("request URL should be /api/v0/graph-annotations but :", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be GET but :", req.Method)
		}
		body, _ := ioutil.ReadAll(req.Body)

		var data struct {
			Service     string   `json:"service,omitempty"`
			Roles       []string `json:"roles,omitempty"`
			From        int64    `json:"from,omitempty"`
			To          int64    `json:"to,omitempty"`
			Title       string   `json:"title,omitempty"`
			Description string   `json:"description,omitempty"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}
		if data.Service != "My-Blog" {
			t.Errorf("request sends json including Service but: %s", data.Service)
		}
		if !reflect.DeepEqual(data.Roles, []string{"Role1", "Role2"}) {
			t.Error("request sends json including Roles but: ", data.Roles)
		}
		if data.From != 1485675275 {
			t.Errorf("request sends json including From but: %d", data.From)
		}
		if data.To != 1485675299 {
			t.Errorf("request sends json including To but: %d", data.To)
		}
		if data.Title != "Deployed" {
			t.Errorf("request sends json including Title but: %s", data.Title)
		}
		if data.Description != "Deployed my blog" {
			t.Errorf("request sends json including Description but: %s", data.Description)
		}
		respJSON, _ := json.Marshal(map[string]string{
			"result": "OK",
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	err := client.CreateGraphAnnotation(&GraphAnnotation{
		Service:     "My-Blog",
		Roles:       []string{"Role1", "Role2"},
		From:        1485675275,
		To:          1485675299,
		Title:       "Deployed",
		Description: "Deployed my blog",
	})

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}
}
