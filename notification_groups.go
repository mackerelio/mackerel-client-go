package mackerel

import (
	"context"
	"fmt"
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

// FindNotificationGroups finds notification groups.
func (c *Client) FindNotificationGroups() ([]*NotificationGroup, error) {
	return c.FindNotificationGroupsContext(context.Background())
}

// FindNotificationGroupsContext finds notification groups.
func (c *Client) FindNotificationGroupsContext(ctx context.Context) ([]*NotificationGroup, error) {
	data, err := requestGetContext[struct {
		NotificationGroups []*NotificationGroup `json:"notificationGroups"`
	}](ctx, c, "/api/v0/notification-groups")
	if err != nil {
		return nil, err
	}
	return data.NotificationGroups, nil
}

// CreateNotificationGroup creates a notification group.
func (c *Client) CreateNotificationGroup(param *NotificationGroup) (*NotificationGroup, error) {
	return c.CreateNotificationGroupContext(context.Background(), param)
}

// CreateNotificationGroupContext creates a notification group.
func (c *Client) CreateNotificationGroupContext(ctx context.Context, param *NotificationGroup) (*NotificationGroup, error) {
	return requestPostContext[NotificationGroup](ctx, c, "/api/v0/notification-groups", param)
}

// UpdateNotificationGroup updates a notification group.
func (c *Client) UpdateNotificationGroup(id string, param *NotificationGroup) (*NotificationGroup, error) {
	return c.UpdateNotificationGroupContext(context.Background(), id, param)
}

// UpdateNotificationGroupContext updates a notification group.
func (c *Client) UpdateNotificationGroupContext(ctx context.Context, id string, param *NotificationGroup) (*NotificationGroup, error) {
	path := fmt.Sprintf("/api/v0/notification-groups/%s", id)
	return requestPutWithContext[NotificationGroup](ctx, c, path, param)
}

// DeleteNotificationGroup deletes a notification group.
func (c *Client) DeleteNotificationGroup(id string) (*NotificationGroup, error) {
	return c.DeleteNotificationGroupContext(context.Background(), id)
}

// DeleteNotificationGroupContext deletes a notification group.
func (c *Client) DeleteNotificationGroupContext(ctx context.Context, id string) (*NotificationGroup, error) {
	path := fmt.Sprintf("/api/v0/notification-groups/%s", id)
	return requestDeleteContext[NotificationGroup](ctx, c, path)
}
