package mackerel

import (
	"fmt"
	"net/url"
)

/*
{
  "alerts": [
    {
      "id": "2wpLU5fBXbG",
      "status": "CRITICAL",
      "monitorId": "2cYjfibBkaj",
      "type": "connectivity",
      "openedAt": 1445399342,
      "hostId": "2vJ965ygiXf"
    },
    {
      "id": "2ust8jNxFH3",
      "status": "CRITICAL",
      "monitorId": "2cYjfibBkaj",
      "type": "connectivity",
      "openedAt": 1441939801,
      "hostId": "2tFrtykgMib"
    }
  ]
}
*/

// Alert information
type Alert struct {
	ID        string  `json:"id,omitempty"`
	Status    string  `json:"status,omitempty"`
	MonitorID string  `json:"monitorId,omitempty"`
	Type      string  `json:"type,omitempty"`
	HostID    string  `json:"hostId,omitempty"`
	Value     float64 `json:"value,omitempty"`
	Message   string  `json:"message,omitempty"`
	Reason    string  `json:"reason,omitempty"`
	OpenedAt  int64   `json:"openedAt,omitempty"`
	ClosedAt  int64   `json:"closedAt,omitempty"`
	Memo      string  `json:"memo,omitempty"`
}

// AlertsResp includes alert and next id
type AlertsResp struct {
	Alerts []*Alert `json:"alerts"`
	NextID string   `json:"nextId,omitempty"`
}

// UpdateAlertParam is for UpdateAlert
type UpdateAlertParam struct {
	Memo string `json:"memo,omitempty"`
}

// UpdateAlertResponse is for UpdateAlert
type UpdateAlertResponse struct {
	Memo string `json:"memo,omitempty"`
}

// AlertLog is the log of alert
// See https://mackerel.io/api-docs/entry/alerts#logs
type AlertLog struct {
	ID           string   `json:"id"`
	CreatedAt    int64    `json:"createdAt"`
	Status       string   `json:"status"`
	Trigger      string   `json:"trigger"`
	MonitorID    *string  `json:"monitorId"`
	TargetValue  *float64 `json:"targetValue"`
	StatusDetail *struct {
		Type   string `json:"type"`
		Detail struct {
			Message string `json:"message"`
			Memo    string `json:"memo"`
		} `json:"detail"`
	} `json:"statusDetail,omitempty"`
}

// FindAlertLogsParam is the parameters for FindAlertLogs
type FindAlertLogsParam struct {
	NextId *string
	Limit  *int
}

// FindAlertLogsResp is for FindAlertLogs
type FindAlertLogsResp struct {
	AlertLogs []*AlertLog `json:"logs"`
	NextID    string      `json:"nextId,omitempty"`
}

func (c *Client) findAlertsWithParams(params url.Values) (*AlertsResp, error) {
	return requestGetWithParams[AlertsResp](c, "/api/v0/alerts", params)
}

// FindAlerts finds open alerts.
func (c *Client) FindAlerts() (*AlertsResp, error) {
	return c.findAlertsWithParams(nil)
}

// FindAlertsByNextID finds next open alerts by next id.
func (c *Client) FindAlertsByNextID(nextID string) (*AlertsResp, error) {
	params := url.Values{}
	params.Set("nextId", nextID)
	return c.findAlertsWithParams(params)
}

// FindWithClosedAlerts finds open and close alerts.
func (c *Client) FindWithClosedAlerts() (*AlertsResp, error) {
	params := url.Values{}
	params.Set("withClosed", "true")
	return c.findAlertsWithParams(params)
}

// FindWithClosedAlertsByNextID finds open and close alerts by next id.
func (c *Client) FindWithClosedAlertsByNextID(nextID string) (*AlertsResp, error) {
	params := url.Values{}
	params.Set("nextId", nextID)
	params.Set("withClosed", "true")
	return c.findAlertsWithParams(params)
}

// GetAlert gets an alert.
func (c *Client) GetAlert(alertID string) (*Alert, error) {
	path := fmt.Sprintf("/api/v0/alerts/%s", alertID)
	return requestGet[Alert](c, path)
}

// CloseAlert closes an alert.
func (c *Client) CloseAlert(alertID string, reason string) (*Alert, error) {
	path := fmt.Sprintf("/api/v0/alerts/%s/close", alertID)
	return requestPost[Alert](c, path, map[string]string{"reason": reason})
}

// UpdateAlert updates an alert.
func (c *Client) UpdateAlert(alertID string, param UpdateAlertParam) (*UpdateAlertResponse, error) {
	path := fmt.Sprintf("/api/v0/alerts/%s", alertID)
	return requestPut[UpdateAlertResponse](c, path, param)
}

func (p FindAlertLogsParam) toValues() url.Values {
	values := url.Values{}
	if p.NextId != nil {
		values.Set("nextId", *p.NextId)
	}
	if p.Limit != nil {
		values.Set("limit", fmt.Sprintf("%d", *p.Limit))
	}
	return values
}

// FindAlertLogs gets alert logs.
func (c *Client) FindAlertLogs(alertId string, params *FindAlertLogsParam) (*FindAlertLogsResp, error) {
	path := fmt.Sprintf("/api/v0/alerts/%s/logs", alertId)
	if params == nil {
		return requestGet[FindAlertLogsResp](c, path)
	}
	return requestGetWithParams[FindAlertLogsResp](c, path, params.toValues())
}
