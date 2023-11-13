package mackerel

import "fmt"

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

// FindNotificationGroups finds notification groups.
func (c *Client) FindNotificationGroups() ([]*NotificationGroup, error) {
	data, err := requestGet[struct {
		NotificationGroups []*NotificationGroup `json:"notificationGroups"`
	}](c, "/api/v0/notification-groups")
	if err != nil {
		return nil, err
	}
	return data.NotificationGroups, nil
}

// CreateNotificationGroup creates a notification group.
func (c *Client) CreateNotificationGroup(param *NotificationGroup) (*NotificationGroup, error) {
	return requestPost[NotificationGroup](c, "/api/v0/notification-groups", param)
}

// UpdateNotificationGroup updates a notification group.
func (c *Client) UpdateNotificationGroup(id string, param *NotificationGroup) (*NotificationGroup, error) {
	path := fmt.Sprintf("/api/v0/notification-groups/%s", id)
	return requestPut[NotificationGroup](c, path, param)
}

// DeleteNotificationGroup deletes a notification group.
func (c *Client) DeleteNotificationGroup(id string) (*NotificationGroup, error) {
	path := fmt.Sprintf("/api/v0/notification-groups/%s", id)
	return requestDelete[NotificationGroup](c, path)
}
