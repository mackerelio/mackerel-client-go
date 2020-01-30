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

func TestCreateGraphDefs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/graph-defs/create" {
			t.Error("request URL should be /api/v0/graph-defs/create but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}
		body, _ := ioutil.ReadAll(req.Body)

		var datas []struct {
			Name        string             `json:"name"`
			DisplayName string             `json:"displayName"`
			Unit        string             `json:"unit"`
			Metrics     []*GraphDefsMetric `json:"metrics"`
		}

		err := json.Unmarshal(body, &datas)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}
		data := datas[0]

		if data.Name != "mackerel" {
			t.Errorf("request sends json including name but: %s", data.Name)
		}
		if data.DisplayName != "HorseMackerel" {
			t.Errorf("request sends json including DisplayName but: %s", data.Name)
		}
		if !reflect.DeepEqual(
			data.Metrics[0],
			&GraphDefsMetric{
				Name:        "saba1",
				DisplayName: "aji1",
				IsStacked:   false,
			},
		) {
			t.Error("request sends json including GraphDefsMetric but: ", data.Metrics[0])
		}
		respJSON, _ := json.Marshal(map[string]string{
			"result": "OK",
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	err := client.CreateGraphDefs([]*GraphDefsParam{
		{
			Name:        "mackerel",
			DisplayName: "HorseMackerel",
			Unit:        "percentage",
			Metrics: []*GraphDefsMetric{
				{
					Name:        "saba1",
					DisplayName: "aji1",
					IsStacked:   false,
				},
				{
					Name:        "saba2",
					DisplayName: "aji2",
					IsStacked:   false,
				},
			},
		},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}
}

func TestGraphDefsOmitJSON(t *testing.T) {
	g := GraphDefsParam{
		Metrics: []*GraphDefsMetric{
			{},
		},
	}
	want := `{"name":"","metrics":[{"name":"","isStacked":false}]}`
	b, err := json.Marshal(&g)
	if err != nil {
		t.Fatal(err)
	}
	if s := string(b); s != want {
		t.Errorf("json.Marshal(%#v) = %q; want %q", g, s, want)
	}
}
