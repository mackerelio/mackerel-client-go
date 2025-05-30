package mackerel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

/*
{
  "monitors": [
    {
      "id": "2cSZzK3XfmG",
      "type": "connectivity",
      "isMute": false,
      "scopes": [],
      "excludeScopes": [],
      "alertStatusOnGone": "WARNING"
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
      "maxCheckAttempts": 1,
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
      "critical": 100.0,
      "maxCheckAttempts": 1,
      "missingDurationWarning": 360,
      "missingDurationCritical": 720
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "external",
      "isMute": false,
      "name": "example.com",
      "method": "GET",
      "url": "http://www.example.com",
      "service": "SomeService",
      "maxCheckAttempts": 1,
      "responseTimeCritical": 10000,
      "responseTimeWarning": 5000,
      "responseTimeDuration": 5,
      "certificationExpirationCritical": 15,
      "certificationExpirationWarning": 30,
      "expectedStatusCode": 200,
      "requestBody": "Request Body",
      "containsString": "Example",
      "skipCertificateVerification": true,
      "followRedirect": true,
      "headers": [
        { "name": "Cache-Control", "value": "no-cache"}
      ]
    },
    {
      "id": "3CSsK3HKiHb",
      "type": "anomalyDetection",
      "isMute": false,
      "name": "My first anomaly detection",
      "trainingPeriodFrom": 1561429260,
      "scopes": [
        "MyService: MyRole"
      ],
      "maxCheckAttempts": 3,
      "warningSensitivity": "insensitive"
    },
    {
      "id": "57We5nNtpZA",
      "type": "query",
      "isMute": false,
      "name": "LabeldMetric - custom.access_counter",
      "query": "custom.access_counter",
      "operator": ">",
      "warning": 30.0,
      "critical": 300.0,
      "legend":""
    }
  ]
}
*/
// Monitor represents interface to which each monitor type must confirm to.
type Monitor interface {
	MonitorType() string
	MonitorID() string
	MonitorName() string

	isMonitor()
}

const (
	monitorTypeConnectivity     = "connectivity"
	monitorTypeHostMetric       = "host"
	monitorTypeServiceMetric    = "service"
	monitorTypeExternalHTTP     = "external"
	monitorTypeExpression       = "expression"
	monitorTypeAnomalyDetection = "anomalyDetection"
	monitorTypeQuery            = "query"
)

// Ensure each monitor type conforms to the Monitor interface.
var (
	_ Monitor = (*MonitorConnectivity)(nil)
	_ Monitor = (*MonitorHostMetric)(nil)
	_ Monitor = (*MonitorServiceMetric)(nil)
	_ Monitor = (*MonitorExternalHTTP)(nil)
	_ Monitor = (*MonitorExpression)(nil)
	_ Monitor = (*MonitorAnomalyDetection)(nil)
	_ Monitor = (*MonitorQuery)(nil)
)

// Ensure only monitor types defined in this package can be assigned to the
// Monitor interface.
func (m *MonitorConnectivity) isMonitor()     {}
func (m *MonitorHostMetric) isMonitor()       {}
func (m *MonitorServiceMetric) isMonitor()    {}
func (m *MonitorExternalHTTP) isMonitor()     {}
func (m *MonitorExpression) isMonitor()       {}
func (m *MonitorAnomalyDetection) isMonitor() {}
func (m *MonitorQuery) isMonitor()            {}

// MonitorConnectivity represents connectivity monitor.
type MonitorConnectivity struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Memo                 string `json:"memo,omitempty"`
	AlertStatusOnGone    string `json:"alertStatusOnGone,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Scopes        []string `json:"scopes,omitempty"`
	ExcludeScopes []string `json:"excludeScopes,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorConnectivity) MonitorType() string { return monitorTypeConnectivity }

// MonitorName returns monitor name.
func (m *MonitorConnectivity) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorConnectivity) MonitorID() string { return m.ID }

// MonitorHostMetric represents host metric monitor.
type MonitorHostMetric struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Memo                 string `json:"memo,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Metric           string   `json:"metric,omitempty"`
	Operator         string   `json:"operator,omitempty"`
	Warning          *float64 `json:"warning"`
	Critical         *float64 `json:"critical"`
	Duration         uint64   `json:"duration,omitempty"`
	MaxCheckAttempts uint64   `json:"maxCheckAttempts,omitempty"`

	Scopes        []string `json:"scopes,omitempty"`
	ExcludeScopes []string `json:"excludeScopes,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorHostMetric) MonitorType() string { return monitorTypeHostMetric }

// MonitorName returns monitor name.
func (m *MonitorHostMetric) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorHostMetric) MonitorID() string { return m.ID }

// MonitorServiceMetric represents service metric monitor.
type MonitorServiceMetric struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Memo                 string `json:"memo,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Service                 string   `json:"service,omitempty"`
	Metric                  string   `json:"metric,omitempty"`
	Operator                string   `json:"operator,omitempty"`
	Warning                 *float64 `json:"warning"`
	Critical                *float64 `json:"critical"`
	Duration                uint64   `json:"duration,omitempty"`
	MaxCheckAttempts        uint64   `json:"maxCheckAttempts,omitempty"`
	MissingDurationWarning  uint64   `json:"missingDurationWarning,omitempty"`
	MissingDurationCritical uint64   `json:"missingDurationCritical,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorServiceMetric) MonitorType() string { return monitorTypeServiceMetric }

// MonitorName returns monitor name.
func (m *MonitorServiceMetric) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorServiceMetric) MonitorID() string { return m.ID }

// MonitorExternalHTTP represents external HTTP monitor.
type MonitorExternalHTTP struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Memo                 string `json:"memo,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Method                          string   `json:"method,omitempty"`
	URL                             string   `json:"url,omitempty"`
	MaxCheckAttempts                uint64   `json:"maxCheckAttempts,omitempty"`
	Service                         string   `json:"service,omitempty"`
	ResponseTimeCritical            *float64 `json:"responseTimeCritical,omitempty"`
	ResponseTimeWarning             *float64 `json:"responseTimeWarning,omitempty"`
	ResponseTimeDuration            *uint64  `json:"responseTimeDuration,omitempty"`
	RequestBody                     string   `json:"requestBody,omitempty"`
	ContainsString                  string   `json:"containsString,omitempty"`
	CertificationExpirationCritical *uint64  `json:"certificationExpirationCritical,omitempty"`
	CertificationExpirationWarning  *uint64  `json:"certificationExpirationWarning,omitempty"`
	SkipCertificateVerification     bool     `json:"skipCertificateVerification,omitempty"`
	FollowRedirect                  bool     `json:"followRedirect,omitempty"`
	ExpectedStatusCode              *int     `json:"expectedStatusCode,omitempty"`
	// Empty list of headers and nil are different. You have to specify empty
	// list as headers explicitly if you want to remove all headers instead of
	// using nil.
	Headers []HeaderField `json:"headers"`
}

// HeaderField represents key-value pairs in an HTTP header for external http
// monitoring.
type HeaderField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MonitorType returns monitor type.
func (m *MonitorExternalHTTP) MonitorType() string { return monitorTypeExternalHTTP }

// MonitorName returns monitor name.
func (m *MonitorExternalHTTP) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorExternalHTTP) MonitorID() string { return m.ID }

// MonitorExpression represents expression monitor.
type MonitorExpression struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Memo                 string `json:"memo,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Expression              string   `json:"expression,omitempty"`
	Operator                string   `json:"operator,omitempty"`
	Warning                 *float64 `json:"warning"`
	Critical                *float64 `json:"critical"`
	EvaluateBackwardMinutes *uint64  `json:"evaluateBackwardMinutes,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorExpression) MonitorType() string { return monitorTypeExpression }

// MonitorName returns monitor name.
func (m *MonitorExpression) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorExpression) MonitorID() string { return m.ID }

// MonitorAnomalyDetection represents anomaly detection monitor.
type MonitorAnomalyDetection struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Memo                 string `json:"memo,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	WarningSensitivity  string `json:"warningSensitivity,omitempty"`
	CriticalSensitivity string `json:"criticalSensitivity,omitempty"`
	TrainingPeriodFrom  uint64 `json:"trainingPeriodFrom,omitempty"`
	MaxCheckAttempts    uint64 `json:"maxCheckAttempts,omitempty"`

	Scopes []string `json:"scopes"`
}

// MonitorType returns monitor type.
func (m *MonitorAnomalyDetection) MonitorType() string { return monitorTypeAnomalyDetection }

// MonitorName returns monitor name.
func (m *MonitorAnomalyDetection) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorAnomalyDetection) MonitorID() string { return m.ID }

// MonitorQuery represents query monitor.
type MonitorQuery struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Memo                 string `json:"memo,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Query                   string   `json:"query,omitempty"`
	Operator                string   `json:"operator,omitempty"`
	Warning                 *float64 `json:"warning"`
	Critical                *float64 `json:"critical"`
	Legend                  string   `json:"legend,omitempty"`
	EvaluateBackwardMinutes *uint64  `json:"evaluateBackwardMinutes,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorQuery) MonitorType() string { return monitorTypeQuery }

// MonitorName returns monitor name.
func (m *MonitorQuery) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorQuery) MonitorID() string { return m.ID }

// FindMonitors finds monitors.
func (c *Client) FindMonitors() ([]Monitor, error) {
	data, err := requestGet[struct {
		Monitors []json.RawMessage `json:"monitors"`
	}](c, "/api/v0/monitors")
	if err != nil {
		return nil, err
	}
	ms := make([]Monitor, 0, len(data.Monitors))
	for _, rawmes := range data.Monitors {
		m, err := decodeMonitor(rawmes)
		var e *unknownMonitorTypeError
		if err != nil {
			if errors.As(err, &e) {
				break
			}
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, err
}

// GetMonitor gets a monitor.
func (c *Client) GetMonitor(monitorID string) (Monitor, error) {
	data, err := requestGet[struct {
		Monitor json.RawMessage `json:"monitor"`
	}](c, fmt.Sprintf("/api/v0/monitors/%s", monitorID))
	if err != nil {
		return nil, err
	}
	m, err := decodeMonitor(data.Monitor)
	if err != nil {
		return nil, err
	}
	return m, err
}

// CreateMonitor creates a monitor.
func (c *Client) CreateMonitor(param Monitor) (Monitor, error) {
	data, err := requestPost[json.RawMessage](c, "/api/v0/monitors", param)
	if err != nil {
		return nil, err
	}
	return decodeMonitor(*data)
}

// UpdateMonitor updates a monitor.
func (c *Client) UpdateMonitor(monitorID string, param Monitor) (Monitor, error) {
	path := fmt.Sprintf("/api/v0/monitors/%s", monitorID)
	data, err := requestPut[json.RawMessage](c, path, param)
	if err != nil {
		return nil, err
	}
	return decodeMonitor(*data)
}

// DeleteMonitor updates a monitor.
func (c *Client) DeleteMonitor(monitorID string) (Monitor, error) {
	path := fmt.Sprintf("/api/v0/monitors/%s", monitorID)
	data, err := requestDelete[json.RawMessage](c, path)
	if err != nil {
		return nil, err
	}
	return decodeMonitor(*data)
}

type unknownMonitorTypeError struct {
	Type string
}

func (e *unknownMonitorTypeError) Error() string {
	return fmt.Sprintf("unknown monitor type: %s", e.Type)
}

// decodeMonitor decodes json.RawMessage and returns monitor.
func decodeMonitor(mes json.RawMessage) (Monitor, error) {
	var typeData struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(mes, &typeData); err != nil {
		return nil, err
	}
	var m Monitor
	switch typeData.Type {
	case monitorTypeConnectivity:
		m = &MonitorConnectivity{}
	case monitorTypeHostMetric:
		m = &MonitorHostMetric{}
	case monitorTypeServiceMetric:
		m = &MonitorServiceMetric{}
	case monitorTypeExternalHTTP:
		m = &MonitorExternalHTTP{}
	case monitorTypeExpression:
		m = &MonitorExpression{}
	case monitorTypeAnomalyDetection:
		m = &MonitorAnomalyDetection{}
	case monitorTypeQuery:
		m = &MonitorQuery{}
	default:
		return nil, &unknownMonitorTypeError{Type: typeData.Type}
	}
	if err := json.Unmarshal(mes, m); err != nil {
		return nil, err
	}
	return m, nil
}

func decodeMonitorReader(r io.Reader) (Monitor, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return decodeMonitor(b)
}
