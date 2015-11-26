package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// FindAlerts find alerts
func (c *Client) FindAlerts() ([]*Alert, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/alerts").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Alerts []*(Alert) `json:"alerts"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Alerts, err
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
