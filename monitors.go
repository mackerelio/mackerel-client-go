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
      "isMute": false,
      "scopes": [],
      "excludeScopes": []
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "host",
      "isMute": false,
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
      "isMute": false,
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
      "isMute": false,
      "name": "example.com",
      "url": "http://www.example.com",
      "service": "SomeService",
      "maxCheckAttempts": 1,
      "responseTimeCritical": 10000,
      "responseTimeWarning": 5000,
      "responseTimeDuration": 5,
      "certificationExpirationCritical": 15,
      "certificationExpirationWarning": 30,
      "containsString": "Example"
    }
  ]
}
*/

// MonitorI represents interface to which each monitor type must confirm to.
// TODO(haya14busa): remove trailing `I` in the name after migrating interface.
type MonitorI interface {
	// MonitorType() must return monitor type.
	MonitorType() string
}

const (
	monitorTypeConnectivity  = "connectivity"
	monitorTypeHostMeric     = "host"
	monitorTypeServiceMetric = "service"
	monitorTypeExternalHTTP  = "external"
	monitorTypeExpression    = "expression"
)

// Ensure each monitor type conforms to the Monitor interface.
var (
	_ MonitorI = (*MonitorConnectivity)(nil)
	_ MonitorI = (*MonitorHostMetric)(nil)
	_ MonitorI = (*MonitorServiceMetric)(nil)
	_ MonitorI = (*MonitorExternalHTTP)(nil)
	_ MonitorI = (*MonitorExpression)(nil)
)

// MonitorConnectivity represents connectivity monitor.
type MonitorConnectivity struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Scopes        []string `json:"scopes,omitempty"`
	ExcludeScopes []string `json:"excludeScopes,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorConnectivity) MonitorType() string {
	return monitorTypeConnectivity
}

// MonitorHostMetric represents host metric monitor.
type MonitorHostMetric struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Metric   string  `json:"metric,omitempty"`
	Operator string  `json:"operator,omitempty"`
	Warning  float64 `json:"warning,omitempty"`
	Critical float64 `json:"critical,omitempty"`
	Duration uint64  `json:"duration,omitempty"`

	Scopes        []string `json:"scopes,omitempty"`
	ExcludeScopes []string `json:"excludeScopes,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorHostMetric) MonitorType() string {
	return monitorTypeHostMeric
}

// MonitorServiceMetric represents service metric monitor.
type MonitorServiceMetric struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Service  string  `json:"service,omitempty"`
	Metric   string  `json:"metric,omitempty"`
	Operator string  `json:"operator,omitempty"`
	Warning  float64 `json:"warning,omitempty"`
	Critical float64 `json:"critical,omitempty"`
	Duration uint64  `json:"duration,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorServiceMetric) MonitorType() string {
	return monitorTypeServiceMetric
}

// MonitorExternalHTTP represents external HTTP monitor.
type MonitorExternalHTTP struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	URL                             string  `json:"url,omitempty"`
	MaxCheckAttempts                float64 `json:"maxCheckAttempts,omitempty"`
	Service                         string  `json:"service,omitempty"`
	ResponseTimeCritical            float64 `json:"responseTimeCritical,omitempty"`
	ResponseTimeWarning             float64 `json:"responseTimeWarning,omitempty"`
	ResponseTimeDuration            float64 `json:"responseTimeDuration,omitempty"`
	ContainsString                  string  `json:"containsString,omitempty"`
	CertificationExpirationCritical uint64  `json:"certificationExpirationCritical,omitempty"`
	CertificationExpirationWarning  uint64  `json:"certificationExpirationWarning,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorExternalHTTP) MonitorType() string {
	return monitorTypeExternalHTTP
}

// MonitorExpression represents expression monitor.
type MonitorExpression struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Expression string  `json:"expression,omitempty"`
	Operator   string  `json:"operator,omitempty"`
	Warning    float64 `json:"warning,omitempty"`
	Critical   float64 `json:"critical,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorExpression) MonitorType() string {
	return monitorTypeExpression
}

// Monitor information
type Monitor struct {
	ID                              string   `json:"id,omitempty"`
	Name                            string   `json:"name,omitempty"`
	Type                            string   `json:"type,omitempty"`
	IsMute                          bool     `json:"isMute,omitempty"`
	Metric                          string   `json:"metric,omitempty"`
	Operator                        string   `json:"operator,omitempty"`
	Warning                         float64  `json:"warning,omitempty"`
	Critical                        float64  `json:"critical,omitempty"`
	Duration                        uint64   `json:"duration,omitempty"`
	URL                             string   `json:"url,omitempty"`
	Scopes                          []string `json:"scopes,omitempty"`
	Service                         string   `json:"service,omitempty"`
	MaxCheckAttempts                float64  `json:"maxCheckAttempts,omitempty"`
	NotificationInterval            uint64   `json:"notificationInterval,omitempty"`
	ExcludeScopes                   []string `json:"excludeScopes,omitempty"`
	ResponseTimeCritical            float64  `json:"responseTimeCritical,omitempty"`
	ResponseTimeWarning             float64  `json:"responseTimeWarning,omitempty"`
	ResponseTimeDuration            float64  `json:"responseTimeDuration,omitempty"`
	CertificationExpirationCritical uint64   `json:"certificationExpirationCritical,omitempty"`
	CertificationExpirationWarning  uint64   `json:"certificationExpirationWarning,omitempty"`
	ContainsString                  string   `json:"containsString,omitempty"`
	Expression                      string   `json:"expression,omitempty"`
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
