package mackerel

import (
	"context"
	"fmt"
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

// FindAlertGroupSettings finds alert group settings.
func (c *Client) FindAlertGroupSettings() ([]*AlertGroupSetting, error) {
	return c.FindAlertGroupSettingsContext(context.Background())
}

// FindAlertGroupSettingsContext finds alert group settings.
func (c *Client) FindAlertGroupSettingsContext(ctx context.Context) ([]*AlertGroupSetting, error) {
	data, err := requestGetContext[struct {
		AlertGroupSettings []*AlertGroupSetting `json:"alertGroupSettings"`
	}](ctx, c, "/api/v0/alert-group-settings")
	if err != nil {
		return nil, err
	}
	return data.AlertGroupSettings, nil
}

// CreateAlertGroupSetting creates an alert group setting.
func (c *Client) CreateAlertGroupSetting(param *AlertGroupSetting) (*AlertGroupSetting, error) {
	return c.CreateAlertGroupSettingContext(context.Background(), param)
}

// CreateAlertGroupSettingContext creates an alert group setting.
func (c *Client) CreateAlertGroupSettingContext(ctx context.Context, param *AlertGroupSetting) (*AlertGroupSetting, error) {
	return requestPostContext[AlertGroupSetting](ctx, c, "/api/v0/alert-group-settings", param)
}

// GetAlertGroupSetting gets an alert group setting.
func (c *Client) GetAlertGroupSetting(id string) (*AlertGroupSetting, error) {
	return c.GetAlertGroupSettingContext(context.Background(), id)
}

// GetAlertGroupSettingContext gets an alert group setting.
func (c *Client) GetAlertGroupSettingContext(ctx context.Context, id string) (*AlertGroupSetting, error) {
	path := fmt.Sprintf("/api/v0/alert-group-settings/%s", id)
	return requestGetContext[AlertGroupSetting](ctx, c, path)
}

// UpdateAlertGroupSetting updates an alert group setting.
func (c *Client) UpdateAlertGroupSetting(id string, param *AlertGroupSetting) (*AlertGroupSetting, error) {
	return c.UpdateAlertGroupSettingContext(context.Background(), id, param)
}

// UpdateAlertGroupSettingContext updates an alert group setting.
func (c *Client) UpdateAlertGroupSettingContext(ctx context.Context, id string, param *AlertGroupSetting) (*AlertGroupSetting, error) {
	path := fmt.Sprintf("/api/v0/alert-group-settings/%s", id)
	return requestPutContext[AlertGroupSetting](ctx, c, path, param)
}

// DeleteAlertGroupSetting deletes an alert group setting.
func (c *Client) DeleteAlertGroupSetting(id string) (*AlertGroupSetting, error) {
	return c.DeleteAlertGroupSettingContext(context.Background(), id)
}

// DeleteAlertGroupSettingContext deletes an alert group setting.
func (c *Client) DeleteAlertGroupSettingContext(ctx context.Context, id string) (*AlertGroupSetting, error) {
	path := fmt.Sprintf("/api/v0/alert-group-settings/%s", id)
	return requestDeleteContext[AlertGroupSetting](ctx, c, path)
}
