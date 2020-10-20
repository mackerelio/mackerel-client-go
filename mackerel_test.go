package mackerel

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestRequest_Failed(t *testing.T) {
	// Overwrite http.Client to control errors
	httpClient := &http.Client{}
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("emulated http error")
	}
	serverRequested := 0
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() { serverRequested++ }()
		// return 302 will cause error by httpClient.CheckRedirect
		http.Redirect(res, req, "http://example.com/DUMMY", 302)
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	client.HTTPClient = httpClient

	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)
	_, err := client.Request(req)

	if expectedServerRequested := 1; serverRequested != expectedServerRequested {
		t.Errorf("should request %d times but %d", expectedServerRequested, serverRequested)
	}
	if err == nil {
		t.Error("error should not be nil but nil")
	}
}

func TestSetMaxRetries(t *testing.T) {
	client := NewClient("DUMMY-KEY")
	if expectedMaxRetries := 0; client.MaxRetries != expectedMaxRetries {
		t.Errorf("client.MaxRetries should be %d, but %d", expectedMaxRetries, client.MaxRetries)
	}

	client.SetMaxRetries(2)
	if expectedMaxRetries := 2; client.MaxRetries != expectedMaxRetries {
		t.Errorf("client.MaxRetries should be %d, but %d", expectedMaxRetries, client.MaxRetries)
	}
}

func TestRequest_Retry(t *testing.T) {
	// Overwrite http.Client to control errors
	httpClient := &http.Client{}
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("emulated http error")
	}
	serverRequested := 0
	howManyTimesErrorReturns := 2
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() { serverRequested++ }()
		if howManyTimesErrorReturns > 0 {
			howManyTimesErrorReturns--
			// return 302 will cause error by httpClient.CheckRedirect
			http.Redirect(res, req, "http://example.com/DUMMY", 302)
		}
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	client.SetMaxRetries(2)
	client.HTTPClient = httpClient

	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)
	_, err := client.Request(req)

	if expectedServerRequested := 3; serverRequested != expectedServerRequested {
		t.Errorf("should request %d times but %d", expectedServerRequested, serverRequested)
	}
	if err != nil {
		t.Errorf("error should be nil but %v", err)
	}
}

func TestRequest_RetryButFailed(t *testing.T) {
	// Overwrite http.Client to control errors
	httpClient := &http.Client{}
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("emulated http error")
	}
	serverRequested := 0
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() { serverRequested++ }()
		// return 302 will cause error by httpClient.CheckRedirect
		http.Redirect(res, req, "http://example.com/DUMMY", 302)
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	client.SetMaxRetries(2)
	client.HTTPClient = httpClient

	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)
	_, err := client.Request(req)

	if expectedServerRequested := 3; serverRequested != expectedServerRequested {
		t.Errorf("should request %d times but %d", expectedServerRequested, serverRequested)
	}
	if err == nil {
		t.Error("error should not be nil but nil")
	}
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

type fakeLogger struct {
	w io.Writer
}

func (p *fakeLogger) Tracef(format string, v ...interface{}) {
	fmt.Fprintf(p.w, format, v...)
}
func (p *fakeLogger) Debugf(format string, v ...interface{})   {}
func (p *fakeLogger) Infof(format string, v ...interface{})    {}
func (p *fakeLogger) Warningf(format string, v ...interface{}) {}
func (p *fakeLogger) Errorf(format string, v ...interface{})   {}

func TestPrivateTracef(t *testing.T) {
	var (
		stdbuf bytes.Buffer
		logbuf bytes.Buffer
		pbuf   bytes.Buffer
	)
	log.SetOutput(&stdbuf)
	defer log.SetOutput(os.Stderr)
	oflags := log.Flags()
	defer log.SetFlags(oflags)
	log.SetFlags(0)

	msg := "test\n"
	t.Run("Logger+PrioritizedLogger", func(t *testing.T) {
		var c Client
		c.Logger = log.New(&logbuf, "", 0)
		c.PrioritizedLogger = &fakeLogger{w: &pbuf}
		c.tracef(msg)
		if s := stdbuf.String(); s != "" {
			t.Errorf("tracef(%q): log.Printf(%q); want %q", msg, s, "")
		}
		if s := logbuf.String(); s != msg {
			t.Errorf("tracef(%q): Logger.Printf(%q); want %q", msg, s, msg)
		}
		if s := pbuf.String(); s != msg {
			t.Errorf("tracef(%q): PrioritizedLogger.Tracef(%q); want %q", msg, s, msg)
		}
	})

	stdbuf.Reset()
	logbuf.Reset()
	pbuf.Reset()
	t.Run("Logger", func(t *testing.T) {
		var c Client
		c.Logger = log.New(&logbuf, "", 0)
		c.tracef(msg)
		if s := stdbuf.String(); s != "" {
			t.Errorf("tracef(%q): log.Printf(%q); want %q", msg, s, "")
		}
		if s := logbuf.String(); s != msg {
			t.Errorf("tracef(%q): Logger.Printf(%q); want %q", msg, s, msg)
		}
		if s := pbuf.String(); s != "" {
			t.Errorf("tracef(%q): PrioritizedLogger.Tracef(%q); want %q", msg, s, "")
		}
	})

	stdbuf.Reset()
	logbuf.Reset()
	pbuf.Reset()
	t.Run("PrioritizedLogger", func(t *testing.T) {
		var c Client
		c.PrioritizedLogger = &fakeLogger{w: &pbuf}
		c.tracef(msg)
		if s := stdbuf.String(); s != "" {
			t.Errorf("tracef(%q): log.Printf(%q); want %q", msg, s, "")
		}
		if s := logbuf.String(); s != "" {
			t.Errorf("tracef(%q): Logger.Printf(%q); want %q", msg, s, "")
		}
		if s := pbuf.String(); s != msg {
			t.Errorf("tracef(%q): PrioritizedLogger.Tracef(%q); want %q", msg, s, msg)
		}
	})

	stdbuf.Reset()
	logbuf.Reset()
	pbuf.Reset()
	t.Run("default", func(t *testing.T) {
		var c Client
		c.tracef(msg)
		if s := stdbuf.String(); s != msg {
			t.Errorf("tracef(%q): log.Printf(%q); want %q", msg, s, msg)
		}
		if s := logbuf.String(); s != "" {
			t.Errorf("tracef(%q): Logger.Printf(%q); want %q", msg, s, "")
		}
		if s := pbuf.String(); s != "" {
			t.Errorf("tracef(%q): PrioritizedLogger.Tracef(%q); want %q", msg, s, "")
		}
	})
}
