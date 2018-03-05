package mackerel

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetHostMetaData(t *testing.T) {
	var (
		hostID       = "9rxGOHfVF8F"
		namespace    = "testing"
		lastModified = time.Date(2018, 3, 6, 3, 0, 0, 0, time.UTC)
	)
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		u := fmt.Sprintf("/api/v0/hosts/%s/metadata/%s", hostID, namespace)
		if req.URL.Path != u {
			t.Errorf("request URL should be %v but %v:", u, req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but :", req.Method)
		}

		respJSON := `{"type":12345,"region":"jp","env":"staging","instance_type":"c4.xlarge"}`
		res.Header()["Content-Type"] = []string{"application/json"}
		res.Header()["Last-Modified"] = []string{lastModified.Format(http.TimeFormat)}
		fmt.Fprint(res, respJSON)
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	metadataResp, err := client.GetHostMetaData(hostID, namespace)
	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	metadata := metadataResp.HostMetaData
	if metadata["type"].(float64) != 12345 {
		t.Errorf("got: %v, want: %v", metadata["type"], 12345)
	}
	if metadata["region"] != "jp" {
		t.Errorf("got: %v, want: %v", metadata["region"], "jp")
	}
	if metadata["env"] != "staging" {
		t.Errorf("got: %v, want: %v", metadata["env"], "staging")
	}
	if metadata["instance_type"] != "c4.xlarge" {
		t.Errorf("got: %v, want: %v", metadata["instance_type"], "c4.xlarge")
	}
	if !metadataResp.LastModified.Equal(lastModified) {
		t.Errorf("got: %v, want: %v", metadataResp.LastModified, lastModified)
	}
}
