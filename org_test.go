package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrg(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/org" {
			t.Error("request URL should be /api/v0/org but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}
		respJSON, _ := json.Marshal(&Org{Name: "hoge"})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	org, err := client.GetOrg()
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if org.Name != "hoge" {
		t.Error("request sends json including Name but: ", org)
	}
}
