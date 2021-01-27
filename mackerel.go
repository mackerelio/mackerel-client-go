package mackerel

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	v2 "github.com/mackerelio/mackerel-client-go/v2"
)

const (
	defaultBaseURL    = "https://api.mackerelio.com/"
	defaultUserAgent  = "mackerel-client-go"
	apiRequestTimeout = 30 * time.Second
)

type PrioritizedLogger = v2.PrioritizedLogger

type Client struct {
	*v2.Client
}

// NewClient returns new mackerel.Client
func NewClient(apikey string) *Client {
	c, _ := NewClientWithOptions(apikey, defaultBaseURL, false)
	return c
}

// NewClientWithOptions returns new mackerel.Client
func NewClientWithOptions(apikey string, rawurl string, verbose bool) (*Client, error) {
	c, err := v2.NewClient(apikey, &v2.NewClientOptions{
		BaseURL: rawurl,
		Verbose: verbose,
	})
	if err != nil {
		return nil, err
	}
	return &Client{Client: c}, nil
}

// PostJSON shortcut method for posting json
func (c *Client) PostJSON(path string, payload interface{}) (*http.Response, error) {
	ctx := context.Background()
	return c.Client.PostJSON(ctx, path, payload)
}

// PutJSON shortcut method for putting json
func (c *Client) PutJSON(path string, payload interface{}) (*http.Response, error) {
	ctx := context.Background()
	return c.Client.PutJSON(ctx, path, payload)
}

func closeResponse(resp *http.Response) {
	if resp != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}

func (c *Client) urlFor(path string) *url.URL {
	u, err := url.Parse(c.BaseURL.String())
	if err != nil {
		panic("invalid url passed")
	}
	u.Path = path
	return u
}
