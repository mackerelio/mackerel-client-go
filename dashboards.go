package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
{
  "dashboards": [
	{
	  "id": "2c5bLca8d",
	  "title": "My Dashboard",
	  "bodyMarkdown": "# A test dashboard",
	  "urlPath": "2u4PP3TJqbu",
	  "createdAt": 1439346145003,
	  "updatedAt": 1439346145003
	}
  ]
}
*/

// Dashboard information
type Dashboard struct {
	ID           string `json:"id,omitempty"`
	Title        string `json:"title,omitempty"`
	BodyMarkDown string `json:"bodyMarkdown,omitempty"`
	URLPath      string `json:"urlPath,omitempty"`
	CreatedAt    int64  `json:"createdAt,omitempty"`
	UpdatedAt    int64  `json:"updatedAt,omitempty"`
}

// FindDashboards find dashboards
func (c *Client) FindDashboards() ([]*Dashboard, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/dashboards").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Dashboards []*(Dashboard) `json:"dashboards"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Dashboards, err
}

// FindDashboard find dashboard
func (c *Client) FindDashboard(dashboardID string) (*Dashboard, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/api/v0/dashboards/%s", dashboardID)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Dashboard
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, err
}

// CreateDashboard creating dashboard
func (c *Client) CreateDashboard(param *Dashboard) (*Dashboard, error) {
	resp, err := c.PostJSON("/api/v0/dashboards", param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Dashboard
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdateDashboard update dashboard
func (c *Client) UpdateDashboard(dashboardID string, param *Dashboard) (*Dashboard, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/api/v0/dashboards/%s", dashboardID), param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Dashboard
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// DeleteDashboard delete dashboard
func (c *Client) DeleteDashboard(dashboardID string) (*Dashboard, error) {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/dashboards/%s", dashboardID)).String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Dashboard
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
