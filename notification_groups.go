package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// NotificationLevel represents a level of notification.
type NotificationLevel string

// NotificationLevels
const (
	NotificationLevelAll      NotificationLevel = "all"
	NotificationLevelCritical NotificationLevel = "critical"
)

// NotificationGroup represents a Mackerel notification group.
// ref. https://mackerel.io/api-docs/entry/notification-groups
type NotificationGroup struct {
	ID                        string                      `json:"id,omitempty"`
	Name                      string                      `json:"name"`
	NotificationLevel         NotificationLevel           `json:"notificationLevel"`
	ChildNotificationGroupIDs []string                    `json:"childNotificationGroupIds"`
	ChildChannelIDs           []string                    `json:"childChannelIds"`
	Monitors                  []*NotificationGroupMonitor `json:"monitors,omitempty"`
	Services                  []*NotificationGroupService `json:"services,omitempty"`
}

// NotificationGroupMonitor represents a notification target monitor rule.
type NotificationGroupMonitor struct {
	ID          string `json:"id"`
	SkipDefault bool   `json:"skipDefault"`
}

// NotificationGroupService represents a notification target service.
type NotificationGroupService struct {
	Name string `json:"name"`
}

// CreateNotificationGroup creates a notification group.
func (c *Client) CreateNotificationGroup(param *NotificationGroup) (*NotificationGroup, error) {
	resp, err := c.PostJSON("/api/v0/notification-groups", param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var notificationGroup NotificationGroup
	if err := json.NewDecoder(resp.Body).Decode(&notificationGroup); err != nil {
		return nil, err
	}

	return &notificationGroup, nil
}

// FindNotificationGroups finds notification groups
func (c *Client) FindNotificationGroups() ([]*NotificationGroup, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/notification-groups").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		NotificationGroups []*NotificationGroup `json:"notificationGroups"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.NotificationGroups, nil
}

// UpdateNotificationGroup updates a notification group
func (c *Client) UpdateNotificationGroup(id string, param *NotificationGroup) (*NotificationGroup, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/api/v0/notification-groups/%s", id), param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var notificationGroup NotificationGroup
	if err := json.NewDecoder(resp.Body).Decode(&notificationGroup); err != nil {
		return nil, err
	}

	return &notificationGroup, nil
}

// DeleteNotificationGroup deletes a notification group
func (c *Client) DeleteNotificationGroup(id string) (*NotificationGroup, error) {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/notification-groups/%s", id)).String(),
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

	var notificationGroup NotificationGroup
	if err := json.NewDecoder(resp.Body).Decode(&notificationGroup); err != nil {
		return nil, err
	}
	return &notificationGroup, nil
}
