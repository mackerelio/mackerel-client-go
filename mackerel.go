package mackerel

import (
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

type Client struct {
	BaseUrl   *url.URL
	ApiKey    string
	Verbose   bool
	UserAgent string
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func NewClient(apikey string) *Client {
	u, _ := url.Parse(defaultBaseURL)
	return &Client{u, apikey, false, defaultUserAgent}
}

func NewClientForTest(apikey string, rawurl string, verbose bool) (*Client, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Client{u, apikey, verbose, defaultUserAgent}, nil
}

func (c *Client) urlFor(path string) *url.URL {
	newUrl, err := url.Parse(c.BaseUrl.String())
	if err != nil {
		panic("invalid url passed")
	}

	newUrl.Path = path

	return newUrl
}

func (c *Client) buildReq(req *http.Request) *http.Request {
	req.Header.Set("X-Api-Key", c.ApiKey)
	req.Header.Set("User-Agent", c.UserAgent)

	return req
}

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
	return resp, nil
}
