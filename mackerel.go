package mackerel

import (
	// "bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	baseURL           = "https://mackerel.io/api/v0"
	userAgent         = "mackerel-client-go"
	apiRequestTimeout = 30 * time.Second
)

type Client struct {
	BaseUrl *url.URL
	ApiKey  string
	Verbose bool
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func NewClient(apikey string) *Client {
	u, _ := url.Parse(baseURL)
	return &Client{u, apikey, false}
}

func NewClientForTest(apikey string, rawurl string, verbose bool) (*Client, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Client{u, apikey, verbose}, nil
}

func (c *Client) urlFor(path string) *url.URL {
	newUrl, err := url.Parse(c.BaseUrl.String())
	if err != nil {
		panic("invalid url passed")
	}

	newUrl.Path = path

	return newUrl
}

func (c *Client) Request(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add("X-Api-Key", c.ApiKey)
	req.Header.Set("User-Agent", userAgent)

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

func (c *Client) FindHost(id string) (*Host, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/api/v0/hosts/%s", id)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Host *Host `json:"host"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Host, err
}
