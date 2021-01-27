package mackerel

import (
	"context"
	"encoding/json"
	"net/http"
)

type OrgService struct {
	c *Client
}

// Org information
type Org struct {
	Name string `json:"name"`
}

// GetOrg get the org
func (org *OrgService) Get(ctx context.Context) (*Org, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", org.c.urlFor("/api/v0/org").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := org.c.Request(req)
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
