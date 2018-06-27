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
			t.Error("request URL should be /api/v0/graph-annotations but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
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
		respJSON, _ := json.Marshal(map[string]interface{}{
			"service":     "My-Blog",
			"roles":       []string{"Role1", "Role2"},
			"from":        1485675275,
			"to":          1485675299,
			"title":       "Deployed",
			"description": "Deployed my blog",
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	annotation, err := client.CreateGraphAnnotation(&GraphAnnotation{
		Service:     "My-Blog",
		Roles:       []string{"Role1", "Role2"},
		From:        1485675275,
		To:          1485675299,
		Title:       "Deployed",
		Description: "Deployed my blog",
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if annotation.Service != "My-Blog" {
		t.Error("request sends json including Service but: ", annotation.Service)
	}

	if !reflect.DeepEqual(annotation.Roles, []string{"Role1", "Role2"}) {
		t.Error("request sends json including Roles but: ", annotation.Roles)
	}

	if annotation.From != 1485675275 {
		t.Error("request sends json including From but: ", annotation.From)
	}

	if annotation.To != 1485675299 {
		t.Error("request sends json including To but: ", annotation.To)
	}

	if annotation.Title != "Deployed" {
		t.Error("request sends json including Title but: ", annotation.Title)
	}

	if annotation.Description != "Deployed my blog" {
		t.Error("request sends json including Description but: ", annotation.Description)
	}
}

func TestFindGraphAnnotations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/graph-annotations" {
			t.Error("request URL should be /api/v0/graph-annotations but: ", req.URL.Path)
		}
		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		query := req.URL.Query()
		if query.Get("service") != "My-Blog" {
			t.Error("request query 'service' param should be My-Blog but: ", query.Get("service"))
		}
		if query.Get("from") != "1485675275" {
			t.Error("request query 'from' param should be 1485675275 but: ", query.Get("from"))
		}
		if query.Get("to") != "1485675299" {
			t.Error("request query 'from' param should be 1485675299 but: ", query.Get("from"))
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"graphAnnotations": {
				{
					"service":     "My-Blog",
					"roles":       []string{"Role1", "Role2"},
					"from":        1485675275,
					"to":          1485675299,
					"title":       "Deployed",
					"description": "Deployed my blog",
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	annotations, err := client.FindGraphAnnotations("My-Blog", 1485675275, 1485675299)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if annotations[0].Service != "My-Blog" {
		t.Error("request sends json including Service but: ", annotations[0].Service)
	}

	if !reflect.DeepEqual(annotations[0].Roles, []string{"Role1", "Role2"}) {
		t.Error("request sends json including Roles but: ", annotations[0].Roles)
	}

	if annotations[0].From != 1485675275 {
		t.Error("request sends json including From but: ", annotations[0].From)
	}

	if annotations[0].To != 1485675299 {
		t.Error("request sends json including To but: ", annotations[0].To)
	}

	if annotations[0].Title != "Deployed" {
		t.Error("request sends json including Title but: ", annotations[0].Title)
	}

	if annotations[0].Description != "Deployed my blog" {
		t.Error("request sends json including Description but: ", annotations[0].Description)
	}
}

func TestUpdateGraphAnnotations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/graph-annotations/123456789" {
			t.Error("request URL should be /api/v0/graph-annotations/123456789 but: ", req.URL.Path)
		}
		if req.Method != "PUT" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"service":     "My-Blog",
			"roles":       []string{"Role1", "Role2"},
			"from":        1485675275,
			"to":          1485675299,
			"title":       "Deployed",
			"description": "Deployed my blog",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	annotation, err := client.UpdateGraphAnnotation("123456789",
		&GraphAnnotation{
			Service:     "My-Blog",
			Roles:       []string{"Role1", "Role2"},
			From:        1485675275,
			To:          1485675299,
			Title:       "Deployed",
			Description: "Deployed my blog",
		})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if annotation.Service != "My-Blog" {
		t.Error("request sends json including Service but: ", annotation.Service)
	}

	if !reflect.DeepEqual(annotation.Roles, []string{"Role1", "Role2"}) {
		t.Error("request sends json including Roles but: ", annotation.Roles)
	}

	if annotation.From != 1485675275 {
		t.Error("request sends json including From but: ", annotation.From)
	}

	if annotation.To != 1485675299 {
		t.Error("request sends json including To but: ", annotation.To)
	}

	if annotation.Title != "Deployed" {
		t.Error("request sends json including Title but: ", annotation.Title)
	}

	if annotation.Description != "Deployed my blog" {
		t.Error("request sends json including Description but: ", annotation.Description)
	}
}

func TestDeleteGraphAnnotations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/graph-annotations/123456789" {
			t.Error("request URL should be /api/v0/graph-annotations/123456789 but: ", req.URL.Path)
		}
		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"service":     "My-Blog",
			"roles":       []string{"Role1", "Role2"},
			"from":        1485675275,
			"to":          1485675299,
			"title":       "Deployed",
			"description": "Deployed my blog",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	annotation, err := client.DeleteGraphAnnotation("123456789")

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if annotation.Service != "My-Blog" {
		t.Error("request sends json including Service but: ", annotation.Service)
	}

	if !reflect.DeepEqual(annotation.Roles, []string{"Role1", "Role2"}) {
		t.Error("request sends json including Roles but: ", annotation.Roles)
	}

	if annotation.From != 1485675275 {
		t.Error("request sends json including From but: ", annotation.From)
	}

	if annotation.To != 1485675299 {
		t.Error("request sends json including To but: ", annotation.To)
	}

	if annotation.Title != "Deployed" {
		t.Error("request sends json including Title but: ", annotation.Title)
	}

	if annotation.Description != "Deployed my blog" {
		t.Error("request sends json including Description but: ", annotation.Description)
	}
}

func TestDeleteGraphAnnotations_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/graph-annotations/123456789" {
			t.Error("request URL should be /api/v0/graph-annotations/123456789 but: ", req.URL.Path)
		}
		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]map[string]string{
			"error": {"message": "Graph annotation not found"},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	_, err := client.DeleteGraphAnnotation("123456789")

	if err == nil {
		t.Error("err should not be nil but: ", err)
	}

	apiErr := err.(*APIError)
	if expected := 404; apiErr.StatusCode != expected {
		t.Errorf("api error StatusCode should be %d but got: %d", expected, apiErr.StatusCode)
	}
	if expected := "API request failed: Graph annotation not found"; apiErr.Error() != expected {
		t.Errorf("api error string should be %s but got: %s", expected, apiErr.Error())
	}
}
