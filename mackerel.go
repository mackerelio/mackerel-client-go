package mackerel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	defaultBaseURL    = "https://api.mackerelio.com/"
	defaultUserAgent  = "mackerel-client-go"
	apiRequestTimeout = 30 * time.Second
)

// PrioritizedLogger is the interface that groups prioritized logging methods.
type PrioritizedLogger interface {
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warningf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

// Client api client for mackerel
type Client struct {
	BaseURL           *url.URL
	APIKey            string
	Verbose           bool
	UserAgent         string
	AdditionalHeaders http.Header
	HTTPClient        *http.Client

	// Client will send logging events to both Logger and PrioritizedLogger.
	// When neither Logger or PrioritizedLogger is set, the log package's standard logger will be used.
	Logger            *log.Logger
	PrioritizedLogger PrioritizedLogger
}

// NewClient returns new mackerel.Client
func NewClient(apikey string) *Client {
	c, _ := NewClientWithOptions(apikey, defaultBaseURL, false)
	return c
}

// NewClientWithOptions returns new mackerel.Client
func NewClientWithOptions(apikey string, rawurl string, verbose bool) (*Client, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	client.Timeout = apiRequestTimeout
	return &Client{
		BaseURL:           u,
		APIKey:            apikey,
		Verbose:           verbose,
		UserAgent:         defaultUserAgent,
		AdditionalHeaders: http.Header{},
		HTTPClient:        client,
	}, nil
}

func (c *Client) urlFor(path string, params url.Values) *url.URL {
	newURL, err := url.Parse(c.BaseURL.String())
	if err != nil {
		panic("invalid base url")
	}
	newURL.Path = path
	newURL.RawQuery = params.Encode()
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

func (c *Client) tracef(format string, v ...interface{}) {
	if c.PrioritizedLogger != nil {
		c.PrioritizedLogger.Tracef(format, v...)
	}
	if c.Logger != nil {
		c.Logger.Printf(format, v...)
	}
	if c.PrioritizedLogger == nil && c.Logger == nil {
		log.Printf(format, v...)
	}
}

// Request request to mackerel and receive response
func (c *Client) Request(req *http.Request) (resp *http.Response, err error) {
	req = c.buildReq(req)

	if c.Verbose {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			c.tracef("%s", dump)
		}
	}

	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if c.Verbose {
		dump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			c.tracef("%s", dump)
		}
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		message, err := extractErrorMessage(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, &APIError{StatusCode: resp.StatusCode, Message: resp.Status}
		}
		return nil, &APIError{StatusCode: resp.StatusCode, Message: message}
	}
	return resp, nil
}

func requestGet[T any](client *Client, path string) (*T, error) {
	return requestNoBody[T](client, http.MethodGet, path, nil)
}

func requestGetWithParams[T any](client *Client, path string, params url.Values) (*T, error) {
	return requestNoBody[T](client, http.MethodGet, path, params)
}

func requestGetAndReturnHeader[T any](client *Client, path string) (*T, http.Header, error) {
	return requestInternal[T](client, http.MethodGet, path, nil, nil)
}

func requestPost[T any](client *Client, path string, payload any) (*T, error) {
	return requestJSON[T](client, http.MethodPost, path, payload)
}

func requestPut[T any](client *Client, path string, payload any) (*T, error) {
	return requestJSON[T](client, http.MethodPut, path, payload)
}

func requestDelete[T any](client *Client, path string) (*T, error) {
	return requestNoBody[T](client, http.MethodDelete, path, nil)
}

func requestJSON[T any](client *Client, method, path string, payload any) (*T, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return nil, err
	}
	data, _, err := requestInternal[T](client, method, path, nil, &body)
	return data, err
}

func requestNoBody[T any](client *Client, method, path string, params url.Values) (*T, error) {
	data, _, err := requestInternal[T](client, method, path, params, nil)
	return data, err
}

func requestInternal[T any](client *Client, method, path string, params url.Values, body io.Reader) (*T, http.Header, error) {
	u := client.urlFor(path, params)
	url := fmt.Sprintf("%s://%s%s?%s", u.Scheme, u.Host, u.Path, u.RawQuery)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	if body != nil || method != http.MethodGet {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := client.Request(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body) // nolint
		resp.Body.Close()
	}()

	var data T
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, nil, err
	}
	return &data, resp.Header, nil
}

func (c *Client) compatRequestJSON(method string, path string, payload interface{}) (*http.Response, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.urlFor(path, url.Values{}).String(), &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return c.Request(req)
}

// ToPtr returns a pointer to the given value of any type.
func ToPtr[T any](v T) *T {
	return &v
}

// Deprecated: use other prefered method.
func (c *Client) PostJSON(path string, payload interface{}) (*http.Response, error) {
	return c.compatRequestJSON(http.MethodPost, path, payload)
}

// Deprecated: use other prefered method.
func (c *Client) PutJSON(path string, payload interface{}) (*http.Response, error) {
	return c.compatRequestJSON(http.MethodPut, path, payload)
}
