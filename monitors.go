package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
{
  "monitors": [
    {
      "id": "2cSZzK3XfmG",
      "type": "connectivity",
      "scopes": [],
      "excludeScopes": []
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "host",
      "name": "disk.aa-00.writes.delta",
      "duration": 3,
      "metric": "disk.aa-00.writes.delta",
      "operator": ">",
      "warning": 20000.0,
      "critical": 400000.0,
      "scopes": [
        "SomeService"
      ],
      "excludeScopes": [
        "SomeService: db-slave-backup"
      ],
      "notificationInterval": 60
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "service",
      "name": "SomeService - custom.access_num.4xx_count",
      "service": "SomeService",
      "duration": 1,
      "metric": "custom.access_num.4xx_count",
      "operator": ">",
      "warning": 50.0,
      "critical": 100.0
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "external",
      "name": "example.com",
      "url": "http://www.example.com"
    }
  ]
}
*/

// Monitor information
type Monitor struct {
	ID                   string   `json:"id,omitempty"`
	Name                 string   `json:"name,omitempty"`
	Type                 string   `json:"type,omitempty"`
	Metric               string   `json:"metric,omitempty"`
	Operator             string   `json:"operator,omitempty"`
	Warning              float64  `json:"warning,omitempty"`
	Critical             float64  `json:"critical,omitempty"`
	Duration             uint64   `json:"duration,omitempty"`
	URL                  string   `json:"url,omitempty"`
	Scopes               []string `json:"scopes,omitempty"`
	Service              string   `json:"service,omitempty"`
	MaxCheckAttempts     float64  `json:"maxCheckAttempts,omitempty"`
	NotificationInterval uint64   `json:"notificationInterval,omitempty"`
	ExcludeScopes        []string `json:"excludeScopes,omitempty"`
	ResponseTimeCritical float64  `json:"responseTimeCritical,omitempty"`
	ResponseTimeWarning  float64  `json:"responseTimeWarning,omitempty"`
	ResponseTimeDuration float64  `json:"responseTimeDuration,omitempty"`
	ContainsString       string   `json:"containsString,omitempty"`
}

// FindMonitors find monitors
func (c *Client) FindMonitors() ([]*Monitor, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/monitors").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Monitors []*(Monitor) `json:"monitors"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Monitors, err
}

// CreateMonitor creating monitor
func (c *Client) CreateMonitor(param *Monitor) (*Monitor, error) {
	resp, err := c.PostJSON("/api/v0/monitors", param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Monitor
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdateMonitor update monitor
func (c *Client) UpdateMonitor(monitorID string, param *Monitor) (*Monitor, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/api/v0/monitors/%s", monitorID), param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Monitor
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// DeleteMonitor update monitor
func (c *Client) DeleteMonitor(monitorID string) (*Monitor, error) {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/monitors/%s", monitorID)).String(),
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

	var data Monitor
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
