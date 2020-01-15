package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Channel represents a Mackerel notification channel.
// ref. https://mackerel.io/api-docs/entry/channels
type Channel struct {
	// ID is excluded when used to call CreateChannel
	ID string `json:"id,omitempty"`

	Name string `json:"name"`
	Type string `json:"type"`

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

// FindChannels requests the channels API and returns a list of Channel
func (c *Client) FindChannels() ([]*Channel, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/channels").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Channels []*Channel `json:"channels"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Channels, err
}

// CreateChannel requests the channels API with the given params to create a channel and returns the created channel.
func (c *Client) CreateChannel(param *Channel) (*Channel, error) {
	resp, err := c.PostJSON("/api/v0/channels", param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	channel := &Channel{}
	err = json.NewDecoder(resp.Body).Decode(channel)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

// DeleteChannel requests the channels API with the given id to delete the specified channel, and returns the deleted channel.
func (c *Client) DeleteChannel(id string) (*Channel, error) {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/channels/%s", id)).String(),
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

	channel := &Channel{}
	err = json.NewDecoder(resp.Body).Decode(channel)
	if err != nil {
		return nil, err
	}
	return channel, nil
}
