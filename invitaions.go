package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Invitation struct {
	Email     string `json:"email,omitempty"`
	Authority string `json:"authority,omitempty"`
	ExpiresAt int64  `json:"expiresAt,omitempty"`
}

func (c *Client) FindInvitations() ([]*Invitation, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/api/v0/invitations")).String(), nil)

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
