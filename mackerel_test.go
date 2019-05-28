package mackerel

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Api-Key") != "dummy-key" {
			t.Error("X-Api-Key header should contains passed key")
		}

		if h := req.Header.Get("User-Agent"); h != defaultUserAgent {
			t.Errorf("User-Agent should be '%s' but %s", defaultUserAgent, h)
		}
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)
	client.Request(req)
}

func TestUrlFor(t *testing.T) {
	client, _ := NewClientWithOptions("dummy-key", "https://example.com/with/ignored/path", false)
	xURL := "https://example.com/some/super/endpoint"
	if url := client.urlFor("/some/super/endpoint").String(); url != xURL {
		t.Errorf("urlFor should be '%s' but %s", xURL, url)
	}
}

func TestBuildReq(t *testing.T) {
	cl := NewClient("dummy-key")
	xVer := "1.0.1"
	xRev := "shasha"
	cl.AdditionalHeaders = http.Header{
		"X-Agent-Version": []string{xVer},
		"X-Revision":      []string{xRev},
	}
	cl.UserAgent = "mackerel-agent"
	req, _ := http.NewRequest("GET", cl.urlFor("/").String(), nil)
	req = cl.buildReq(req)

	if req.Header.Get("X-Api-Key") != "dummy-key" {
		t.Error("X-Api-Key header should contains passed key")
	}
	if h := req.Header.Get("User-Agent"); h != cl.UserAgent {
		t.Errorf("User-Agent should be '%s' but %s", cl.UserAgent, h)
	}
	if h := req.Header.Get("X-Agent-Version"); h != xVer {
		t.Errorf("X-Agent-Version should be '%s' but %s", xVer, h)
	}
	if h := req.Header.Get("X-Revision"); h != xRev {
		t.Errorf("X-Revision should be '%s' but %s", xRev, h)
	}
}

func TestLogger(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("OK"))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, true)
	var buf bytes.Buffer
	client.Logger = log.New(&buf, "<api>", 0)
	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)
	client.Request(req)
	s := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(s, "<api>") || !strings.HasSuffix(s, "OK") {
		t.Errorf("verbose log should match /<api>.*OK/; but %s", s)
	}
}
