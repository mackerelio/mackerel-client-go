package mackerel

import (
	"encoding/json"
	"fmt"
)

/*
`/dashboards` Response
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
      "fractionSize": 2,
      "suffix": "total",
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
    },
    {
      "type": "alertStatus",
      "title": "alertStatus",
      "roleFullname": "test:dashboard",
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
	ID        string   `json:"id,omitempty"`
	Title     string   `json:"title"`
	URLPath   string   `json:"urlPath"`
	CreatedAt int64    `json:"createdAt,omitempty"`
	UpdatedAt int64    `json:"updatedAt,omitempty"`
	Memo      string   `json:"memo"`
	Widgets   []Widget `json:"widgets"`
}

// Widget information
type Widget struct {
	Type           string          `json:"type"`
	Title          string          `json:"title"`
	Layout         Layout          `json:"layout"`
	Metric         Metric          `json:"metric,omitempty"`
	Graph          Graph           `json:"graph,omitempty"`
	Range          Range           `json:"range,omitempty"`
	Markdown       string          `json:"markdown,omitempty"`
	ReferenceLines []ReferenceLine `json:"referenceLines,omitempty"`
	// If this field is nil, it will be treated as a two-digit display after the decimal point.
	FractionSize *int64       `json:"fractionSize,omitempty"`
	Suffix       string       `json:"suffix,omitempty"`
	FormatRules  []FormatRule `json:"formatRules,omitempty"`
	RoleFullName string       `json:"roleFullname,omitempty"`
}

// Metric information
type Metric struct {
	Type        string `json:"type"`
	Name        string `json:"name,omitempty"`
	HostID      string `json:"hostId,omitempty"`
	ServiceName string `json:"serviceName,omitempty"`
	Expression  string `json:"expression,omitempty"`
	Query       string `json:"query,omitempty"`
	Legend      string `json:"legend,omitempty"`
}
type metricQuery struct {
	Type        string `json:"type"`
	Name        string `json:"-"`
	HostID      string `json:"-"`
	ServiceName string `json:"-"`
	Expression  string `json:"-"`
	Query       string `json:"query"`
	Legend      string `json:"legend"`
}

// MarshalJSON marshals as JSON
func (m Metric) MarshalJSON() ([]byte, error) {
	type Alias Metric
	switch m.Type {
	case "":
		return []byte("null"), nil
	case "query":
		return json.Marshal(metricQuery(m))
	default:
		return json.Marshal(Alias(m))
	}
}

// FormatRule information
type FormatRule struct {
	Name      string  `json:"name"`
	Threshold float64 `json:"threshold"`
	Operator  string  `json:"operator"`
}

// Graph information
type Graph struct {
	Type         string `json:"type"`
	Name         string `json:"name,omitempty"`
	HostID       string `json:"hostId,omitempty"`
	RoleFullName string `json:"roleFullname,omitempty"`
	IsStacked    bool   `json:"isStacked,omitempty"`
	ServiceName  string `json:"serviceName,omitempty"`
	Expression   string `json:"expression,omitempty"`
	Query        string `json:"query,omitempty"`
	Legend       string `json:"legend,omitempty"`
}
type graphQuery struct {
	Type         string `json:"type"`
	Name         string `json:"-"`
	HostID       string `json:"-"`
	RoleFullName string `json:"-"`
	IsStacked    bool   `json:"-"`
	ServiceName  string `json:"-"`
	Expression   string `json:"-"`
	Query        string `json:"query"`
	Legend       string `json:"legend"`
}

// MarshalJSON marshals as JSON
func (g Graph) MarshalJSON() ([]byte, error) {
	type Alias Graph
	switch g.Type {
	case "":
		return []byte("null"), nil
	case "query":
		return json.Marshal(graphQuery(g))
	default:
		return json.Marshal(Alias(g))
	}
}

// Range information
type Range struct {
	Type   string `json:"type"`
	Period int64  `json:"period,omitempty"`
	Offset int64  `json:"offset,omitempty"`
	Start  int64  `json:"start,omitempty"`
	End    int64  `json:"end,omitempty"`
}

type rangeAbsolute struct {
	Type   string `json:"type"`
	Period int64  `json:"-"`
	Offset int64  `json:"-"`
	Start  int64  `json:"start"`
	End    int64  `json:"end"`
}

type rangeRelative struct {
	Type   string `json:"type"`
	Period int64  `json:"period"`
	Offset int64  `json:"offset"`
	Start  int64  `json:"-"`
	End    int64  `json:"-"`
}

// MarshalJSON marshals as JSON
func (r Range) MarshalJSON() ([]byte, error) {
	switch r.Type {
	case "absolute":
		return json.Marshal(rangeAbsolute(r))
	case "relative":
		return json.Marshal(rangeRelative(r))
	default:
		return []byte("null"), nil
	}
}

// ReferenceLine information
type ReferenceLine struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// Layout information
type Layout struct {
	X      int64 `json:"x"`
	Y      int64 `json:"y"`
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

// FindDashboards finds dashboards.
func (c *Client) FindDashboards() ([]*Dashboard, error) {
	data, err := requestGet[struct {
		Dashboards []*Dashboard `json:"dashboards"`
	}](c, "/api/v0/dashboards")
	if err != nil {
		return nil, err
	}
	return data.Dashboards, nil
}

// CreateDashboard creates a dashboard.
func (c *Client) CreateDashboard(param *Dashboard) (*Dashboard, error) {
	return requestPost[Dashboard](c, "/api/v0/dashboards", param)
}

// FindDashboard finds a dashboard.
func (c *Client) FindDashboard(dashboardID string) (*Dashboard, error) {
	path := fmt.Sprintf("/api/v0/dashboards/%s", dashboardID)
	return requestGet[Dashboard](c, path)
}

// UpdateDashboard updates a dashboard.
func (c *Client) UpdateDashboard(dashboardID string, param *Dashboard) (*Dashboard, error) {
	path := fmt.Sprintf("/api/v0/dashboards/%s", dashboardID)
	return requestPut[Dashboard](c, path, param)
}

// DeleteDashboard deletes a dashboard.
func (c *Client) DeleteDashboard(dashboardID string) (*Dashboard, error) {
	path := fmt.Sprintf("/api/v0/dashboards/%s", dashboardID)
	return requestDelete[Dashboard](c, path)
}
