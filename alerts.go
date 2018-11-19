package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
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
}

// NextID is next id of alert
type NextID string

// Data includes alert and next id
type Data struct {
	Alerts []*Alert `json:"alerts"`
	ID     NextID   `json:"nextId,omitempty"`
}

var d Data

func (c *Client) findAlertsWithParam(v url.Values) (Data, error) {
	u := c.urlFor("/api/v0/alerts")
	u.RawQuery = v.Encode()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", u.String(), v.Encode()), nil)
	if err != nil {
		return d, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return d, err
	}

	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return d, err
	}
	return d, err
}

// FindAlerts find open alerts
func (c *Client) FindAlerts() (Data, error) {
	d, err := c.findAlertsWithParam(url.Values{})
	return d, err
}

// FindAlertsByNextID find next open alerts by next id
func (c *Client) FindAlertsByNextID(nextID string) (Data, error) {
	v := url.Values{}
	v.Set("nextId", nextID)
	d, err := c.findAlertsWithParam(v)
	return d, err
}

// FindWithClosedAlerts find open and close alerts
func (c *Client) FindWithClosedAlerts() (Data, error) {
	v := url.Values{}
	v.Set("withClosed", "true")
	d, err := c.findAlertsWithParam(v)
	return d, err
}

// FindWithClosedAlertsByNextID find open and close alerts by next id
func (c *Client) FindWithClosedAlertsByNextID(nextID string) (Data, error) {
	v := url.Values{}
	v.Set("withClosed", "true")
	v.Set("nextId", nextID)
	d, err := c.findAlertsWithParam(v)
	return d, err
}

// CloseAlert close alert
func (c *Client) CloseAlert(alertID string, reason string) (*Alert, error) {
	var reqBody struct {
		Reason string `json:"reason"`
	}
	reqBody.Reason = reason
	resp, err := c.PostJSON(fmt.Sprintf("/api/v0/alerts/%s/close", alertID), &reqBody)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data *Alert
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
