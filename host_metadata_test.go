package mackerel

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
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
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON := `{"type":12345,"region":"jp","env":"staging","instance_type":"c4.xlarge"}`
		res.Header()["Content-Type"] = []string{"application/json"}
		res.Header()["Last-Modified"] = []string{lastModified.Format(http.TimeFormat)}
		fmt.Fprint(res, respJSON) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	metadataResp, err := client.GetHostMetaData(hostID, namespace)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	metadata := metadataResp.HostMetaData
	if metadata.(map[string]interface{})["type"].(float64) != 12345 {
		t.Errorf("got: %v, want: %v", metadata.(map[string]interface{})["type"], 12345)
	}
	if metadata.(map[string]interface{})["region"] != "jp" {
		t.Errorf("got: %v, want: %v", metadata.(map[string]interface{})["region"], "jp")
	}
	if metadata.(map[string]interface{})["env"] != "staging" {
		t.Errorf("got: %v, want: %v", metadata.(map[string]interface{})["env"], "staging")
	}
	if metadata.(map[string]interface{})["instance_type"] != "c4.xlarge" {
		t.Errorf("got: %v, want: %v", metadata.(map[string]interface{})["instance_type"], "c4.xlarge")
	}
	if !metadataResp.LastModified.Equal(lastModified) {
		t.Errorf("got: %v, want: %v", metadataResp.LastModified, lastModified)
	}
}

func TestGetHostMetaDataNameSpaces(t *testing.T) {
	var (
		hostID = "9rxGOHfVF8F"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		u := fmt.Sprintf("/api/v0/hosts/%s/metadata", hostID)
		if req.URL.Path != u {
			t.Errorf("request URL should be %v but %v:", u, req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON := `{"metadata":[{"namespace":"testing1"}, {"namespace":"testing2"}]}`
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, respJSON) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	namespaces, err := client.GetHostMetaDataNameSpaces(hostID)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if !reflect.DeepEqual(namespaces, []string{"testing1", "testing2"}) {
		t.Errorf("got: %v, want: %v", namespaces, []string{"testing1", "testing2"})
	}
}

func TestPutHostMetaData(t *testing.T) {
	var (
		hostID    = "9rxGOHfVF8F"
		namespace = "testing"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		u := fmt.Sprintf("/api/v0/hosts/%s/metadata/%s", hostID, namespace)
		if req.URL.Path != u {
			t.Errorf("request URL should be %v but %v:", u, req.URL.Path)
		}

		if req.Method != "PUT" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		body, _ := io.ReadAll(req.Body)
		reqJSON := `{"env":"staging","instance_type":"c4.xlarge","region":"jp","type":12345}` + "\n"
		if string(body) != reqJSON {
			t.Errorf("request body should be %v but %v", reqJSON, string(body))
		}

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, `{"success":true}`) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	metadata := map[string]interface{}{
		"type":          12345,
		"region":        "jp",
		"env":           "staging",
		"instance_type": "c4.xlarge",
	}
	err := client.PutHostMetaData(hostID, namespace, &metadata)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
}

func TestDeleteHostMetaData(t *testing.T) {
	var (
		hostID    = "9rxGOHfVF8F"
		namespace = "testing"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		u := fmt.Sprintf("/api/v0/hosts/%s/metadata/%s", hostID, namespace)
		if req.URL.Path != u {
			t.Errorf("request URL should be %v but %v:", u, req.URL.Path)
		}

		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, `{"success":true}`) // nolint
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	err := client.DeleteHostMetaData(hostID, namespace)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
}
