package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// User information
type User struct {
	ID         string `json:"id,omitempty"`
	ScreenName string `json:"screenName,omitempty"`
	Email      string `json:"email,omitempty"`
	Authority  string `json:"authority,omitempty"`

	IsInRegistrationProcess bool     `json:"isInRegistrationProcess,omitempty"`
	IsMFAEnabled            bool     `json:"isMFAEnabled,omitempty"`
	AuthenticationMethods   []string `json:"authenticationMethods,omitempty"`
	JoinedAt                int64    `json:"joinedAt,omitempty"`
}

// FindUsers find users.
func (c *Client) FindUsers() ([]*User, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/users").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Users []*User `json:"users"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Users, err
}

// DeleteUser delete users.
func (c *Client) DeleteUser(userID string) (*User, error) {
	req, err := http.NewRequest("DELETE", c.urlFor(fmt.Sprintf("/api/v0/users/%s", userID)).String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
