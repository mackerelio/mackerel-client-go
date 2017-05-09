package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Service represents Mackerel "service".
type Service struct {
	Name  string   `json:"name,omitempty"`
	Memo  string   `json:"memo,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

// FindServices finds services.
func (c *Client) FindServices() ([]*Service, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/api/v0/services")).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Services []*Service `json:"services"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Services, err
}
