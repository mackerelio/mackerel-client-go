package mackerel

import "fmt"

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

// FindAlertGroupSettings finds alert group settings.
func (c *Client) FindAlertGroupSettings() ([]*AlertGroupSetting, error) {
	data, err := requestGet[struct {
		AlertGroupSettings []*AlertGroupSetting `json:"alertGroupSettings"`
	}](c, "/api/v0/alert-group-settings")
	if err != nil {
		return nil, err
	}
	return data.AlertGroupSettings, nil
}

// CreateAlertGroupSetting creates an alert group setting.
func (c *Client) CreateAlertGroupSetting(param *AlertGroupSetting) (*AlertGroupSetting, error) {
	return requestPost[AlertGroupSetting](c, "/api/v0/alert-group-settings", param)
}

// GetAlertGroupSetting gets an alert group setting.
func (c *Client) GetAlertGroupSetting(id string) (*AlertGroupSetting, error) {
	path := fmt.Sprintf("/api/v0/alert-group-settings/%s", id)
	return requestGet[AlertGroupSetting](c, path)
}

// UpdateAlertGroupSetting updates an alert group setting.
func (c *Client) UpdateAlertGroupSetting(id string, param *AlertGroupSetting) (*AlertGroupSetting, error) {
	path := fmt.Sprintf("/api/v0/alert-group-settings/%s", id)
	return requestPut[AlertGroupSetting](c, path, param)
}

// DeleteAlertGroupSetting deletes an alert group setting.
func (c *Client) DeleteAlertGroupSetting(id string) (*AlertGroupSetting, error) {
	path := fmt.Sprintf("/api/v0/alert-group-settings/%s", id)
	return requestDelete[AlertGroupSetting](c, path)
}
