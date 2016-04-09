package mackerel

import (
	"encoding/json"
	"net/http"
)

// Org information
type Org struct {
	Name string `json:"name"`
}

// GetOrg get the org
func (c *Client) GetOrg() (*Org, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/org").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	var data Org
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
