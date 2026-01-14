package mackerel

import (
	"context"
	"fmt"
)

// Channel represents a Mackerel notification channel.
// ref. https://mackerel.io/api-docs/entry/channels
type Channel struct {
	// ID is excluded when used to call CreateChannel
	ID string `json:"id,omitempty"`

	Name string `json:"name"`
	Type string `json:"type"`

	SuspendedAt *int64 `json:"suspendedAt,omitempty"`

	// Exists when the type is "email"
	Emails  *[]string `json:"emails,omitempty"`
	UserIDs *[]string `json:"userIds,omitempty"`

	// Exists when the type is "slack"
	Mentions Mentions `json:"mentions,omitempty"`
	// In order to support both 'not setting this field' and 'setting the field as false',
	// this field needed to be *bool not bool.
	EnabledGraphImage *bool `json:"enabledGraphImage,omitempty"`

	// Exists when the type is "slack" or "webhook"
	URL string `json:"url,omitempty"`

	// Exists when the type is "email", "slack", or "webhook"
	Events *[]string `json:"events,omitempty"`
}

// Mentions represents the structure used for slack channel mentions
type Mentions struct {
	OK       string `json:"ok,omitempty"`
	Warning  string `json:"warning,omitempty"`
	Critical string `json:"critical,omitempty"`
}

// FindChannels finds channels.
func (c *Client) FindChannels() ([]*Channel, error) {
	return c.FindChannelsContext(context.Background())
}

// FindChannelsContext finds channels.
func (c *Client) FindChannelsContext(ctx context.Context) ([]*Channel, error) {
	data, err := requestGetContext[struct {
		Channels []*Channel `json:"channels"`
	}](ctx, c, "/api/v0/channels")
	if err != nil {
		return nil, err
	}
	return data.Channels, nil
}

// CreateChannel creates a channel.
func (c *Client) CreateChannel(param *Channel) (*Channel, error) {
	return c.CreateChannelContext(context.Background(), param)
}

// CreateChannelContext creates a channel.
func (c *Client) CreateChannelContext(ctx context.Context, param *Channel) (*Channel, error) {
	return requestPostContext[Channel](ctx, c, "/api/v0/channels", param)
}

// UpdateChannel updates a specific channel
func (c *Client) UpdateChannel(channelId string, param *Channel) (*Channel, error) {
	return c.UpdateChannelContext(context.Background(), channelId, param)
}

// UpdateChannelContext is like [UpdateChannel]
func (c *Client) UpdateChannelContext(ctx context.Context, channelID string, param *Channel) (*Channel, error) {
	data, err := requestPutWithContext[Channel](ctx, c, fmt.Sprintf("/api/v0/channels/%s", channelID), param)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteChannel deletes a channel.
func (c *Client) DeleteChannel(channelID string) (*Channel, error) {
	return c.DeleteChannelContext(context.Background(), channelID)
}

// DeleteChannelContext deletes a channel.
func (c *Client) DeleteChannelContext(ctx context.Context, channelID string) (*Channel, error) {
	path := fmt.Sprintf("/api/v0/channels/%s", channelID)
	return requestDeleteContext[Channel](ctx, c, path)
}
