package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
`/dashboards` Response
Legacy
{
  "dashboards": [
    {
       "id": "2c5bLca8d",
       "title": "My Dashboard",
       "bodyMarkdown": "# A test dashboard",
       "urlPath": "2u4PP3TJqbu",
       "createdAt": 1439346145003,
       "updatedAt": 1439346145003,
       "isLegacy": true
    }
  ]
}

Current
{
  "dashboards": [
    {
      "id":        "2c5bLca8e",
      "title":     "My Custom Dashboard(Current)",
      "urlPath":   "2u4PP3TJqbv",
      "createdAt": 1552909732,
      "updatedAt": 1552992837,
      "memo":      "A test Current Dashboard"
    }
  ]
}
*/

/*
`/dashboards/${ID}` Response`
Legacy
{
  "id": "2c5bLca8d",
  "title": "My Dashboard",
  "bodyMarkdown": "# A test dashboard",
  "urlPath": "2u4PP3TJqbu",
  "createdAt": 1439346145003,
  "updatedAt": 1439346145003,
  "isLegacy": true
}
Current
{
  "id": "2c5bLca8e",
  "createdAt": 1552909732,
  "updatedAt": 1552992837,
  "title": "My Custom Dashboard(Current),
  "urlPath": "2u4PP3TJqbv",
  "memo": "A test Current Dashboard",
  "widgets": [
    {
      "type": "markdown",
      "title": "markdown",
      "markdown": "# body",
      "layout": {
        "x": 0,
        "y": 0,
        "width": 24,
        "height": 3
      }
    },
    {
      "type": "graph",
      "title": "graph",
      "graph": {
        "type": "host",
        "hostId": "2u4PP3TJqbw",
        "name": "loadavg.loadavg15"
      },
      "layout": {
        "x": 0,
        "y": 7,
        "width": 8,
        "height": 10
      }
    },
    {
      "type": "value",
      "title": "value",
      "metric": {
        "type": "expression",
        "expression": "alias(scale(\nsum(\n  group(\n    host(2u4PP3TJqbx,loadavg.*)\n  )\n),\n1\n), 'test')"
      },
      "layout": {
        "x": 0,
        "y": 17,
        "width": 8,
        "height": 5
      }
    }
  ]
}
*/

// Dashboard information
type Dashboard struct {
	// Common to legacy dashboard and current dashboard
	ID        string `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	URLPath   string `json:"urlPath,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`

	// current dashboard
	Memo    string   `json:"memo"`
	Widgets []Widget `json:"widgets,omitempty"`

	// legacy dashboard
	IsLegacy     bool   `json:"isLegacy,omitempty"`
	BodyMarkDown string `json:"bodyMarkdown,omitempty"`
}

// Widget information
type Widget struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title"`
	Layout   Layout `json:"layout,omitempty"`
	Metric   Metric `json:"metric,omitempty"`
	Graph    Graph  `json:"graph,omitempty"`
	Range    Range  `json:"range,omitempty"`
	Markdown string `json:"markdown,omitempty"`
}

// Metric information
type Metric struct {
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	HostID      string `json:"hostId,omitempty"`
	ServiceName string `json:"serviceName,omitempty"`
	Expression  string `json:"expression,omitempty"`
}

// MarshalJSON marshals as JSON
func (m Metric) MarshalJSON() ([]byte, error) {
	type Alias Metric
	if m.Type == "" {
		return []byte("null"), nil
	}
	return json.Marshal(Alias(m))
}

// Graph information
type Graph struct {
	Type         string `json:"type,omitempty"`
	Name         string `json:"name,omitempty"`
	HostID       string `json:"hostId,omitempty"`
	RoleFullName string `json:"roleFullname,omitempty"`
	IsStacked    bool   `json:"isStacked,omitempty"`
	ServiceName  string `json:"serviceName,omitempty"`
	Expression   string `json:"expression,omitempty"`
}

// MarshalJSON marshals as JSON
func (g Graph) MarshalJSON() ([]byte, error) {
	type Alias Graph
	if g.Type == "" {
		return []byte("null"), nil
	}
	return json.Marshal(Alias(g))
}

// Range information
type Range struct {
	Type   string `json:"type,omitempty"`
	Period int64  `json:"period,omitempty"`
	Offset int64  `json:"offset,omitempty"`
	Start  int64  `json:"start,omitempty"`
	End    int64  `json:"end,omitempty"`
}

// MarshalJSON marshals as JSON
func (r Range) MarshalJSON() ([]byte, error) {
	type Alias Range
	if r.Type == "" {
		return []byte("null"), nil
	}
	return json.Marshal(Alias(r))
}

// Layout information
type Layout struct {
	X      int64 `json:"x"`
	Y      int64 `json:"y"`
	Width  int64 `json:"width,omitempty"`
	Height int64 `json:"height,omitempty"`
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
		Dashboards []*Dashboard `json:"dashboards"`
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
