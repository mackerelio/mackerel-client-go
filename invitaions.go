package mackerel

import (
	"encoding/json"
	"net/http"
)

// Invitation information
type Invitation struct {
	Email     string `json:"email,omitempty"`
	Authority string `json:"authority,omitempty"`
	ExpiresAt int64  `json:"expiresAt,omitempty"`
}

// FindInvitations find invitations.
func (c *Client) FindInvitations() ([]*Invitation, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/invitations").String(), nil)

	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Invitations []*Invitation `json:"invitations"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Invitations, err
}
