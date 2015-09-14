package mackerel

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	defaultBaseURL    = "https://mackerel.io/api/v0"
	defaultUserAgent  = "mackerel-client-go"
	apiRequestTimeout = 30 * time.Second
)

// Client api client for mackerel
type Client struct {
	BaseURL           *url.URL
	APIKey            string
	Verbose           bool
	UserAgent         string
	AdditionalHeaders http.Header
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// NewClient returns new mackerel.Client
func NewClient(apikey string) *Client {
	u, _ := url.Parse(defaultBaseURL)
	return &Client{u, apikey, false, defaultUserAgent, http.Header{}}
}

// NewClientWithOptions returns new mackerel.Client
func NewClientWithOptions(apikey string, rawurl string, verbose bool) (*Client, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Client{u, apikey, verbose, defaultUserAgent, http.Header{}}, nil
}

func (c *Client) urlFor(path string) *url.URL {
	newURL, err := url.Parse(c.BaseURL.String())
	if err != nil {
		panic("invalid url passed")
	}

	newURL.Path = path

	return newURL
}

func (c *Client) buildReq(req *http.Request) *http.Request {
	for header, values := range c.AdditionalHeaders {
		for _, v := range values {
			req.Header.Add(header, v)
		}
	}
	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("User-Agent", c.UserAgent)
	return req
}

// Request request to mackerel and receive response
func (c *Client) Request(req *http.Request) (resp *http.Response, err error) {
	req = c.buildReq(req)

	if c.Verbose {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			log.Printf("%s", dump)
		}
	}

	client := &http.Client{} // same as http.DefaultClient
	client.Timeout = apiRequestTimeout
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.Verbose {
		dump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			log.Printf("%s", dump)
		}
	}
	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		return resp, fmt.Errorf("API result failed: %s", resp.Status)
	}
	return resp, nil
}

// CloseReponse clos body stream if response exists.
func (c *Client) CloseReponse(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}
