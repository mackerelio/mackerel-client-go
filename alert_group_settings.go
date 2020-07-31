package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AlertGroupSetting represents a Mackerel alert group setting.
// ref. https://mackerel.io/api-docs/entry/alert-group-settings
type AlertGroupSetting struct {
	ID                   string   `json:"id,omitempty"`
	Name                 string   `json:"name"`
	Memo                 string   `json:"memo,omitempty"`
	ServiceScopes        []string `json:"serviceScopes,omitempty"`
	RoleScopes           []string `json:"roleScopes,omitempty"`
	MonitorScopes        []string `json:"monitorScopes,omitempty"`
	NotificationInterval uint64   `json:"notificationInterval,omitempty"`
}

// FindAlertGroupSettings finds alert group settings
func (c *Client) FindAlertGroupSettings() ([]*AlertGroupSetting, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/alert-group-settings").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		AlertGroupSettings []*AlertGroupSetting `json:"alertGroupSettings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.AlertGroupSettings, nil
}

// CreateAlertGroupSetting creates a alert group setting
func (c *Client) CreateAlertGroupSetting(param *AlertGroupSetting) (*AlertGroupSetting, error) {
	resp, err := c.PostJSON("/api/v0/alert-group-settings", param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var alertGroupSetting AlertGroupSetting
	if err := json.NewDecoder(resp.Body).Decode(&alertGroupSetting); err != nil {
		return nil, err
	}

	return &alertGroupSetting, nil
}

// GetAlertGroupSetting gets alert group setting specified by ID
func (c *Client) GetAlertGroupSetting(id string) (*AlertGroupSetting, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/api/v0/alert-group-settings/%s", id)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var alertGroupSetting AlertGroupSetting
	if err := json.NewDecoder(resp.Body).Decode(&alertGroupSetting); err != nil {
		return nil, err
	}

	return &alertGroupSetting, nil
}

// UpdateAlertGroupSetting updates a alert group setting
func (c *Client) UpdateAlertGroupSetting(id string, param *AlertGroupSetting) (*AlertGroupSetting, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/api/v0/alert-group-settings/%s", id), param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var alertGroupSetting AlertGroupSetting
	if err := json.NewDecoder(resp.Body).Decode(&alertGroupSetting); err != nil {
		return nil, err
	}

	return &alertGroupSetting, nil
}

// DeleteAlertGroupSetting deletes a alert group setting specified by ID.
func (c *Client) DeleteAlertGroupSetting(id string) (*AlertGroupSetting, error) {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/alert-group-settings/%s", id)).String(),
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

	var alertGroupSetting AlertGroupSetting
	if err := json.NewDecoder(resp.Body).Decode(&alertGroupSetting); err != nil {
		return nil, err
	}

	return &alertGroupSetting, nil
}
