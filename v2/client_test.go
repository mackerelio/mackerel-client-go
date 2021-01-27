package mackerel

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestBuildReq(t *testing.T) {
	cl, err := NewClient("dummy-key", nil)
	if err != nil {
		t.Fatal(err)
	}
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
