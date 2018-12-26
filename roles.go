package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Role represents Mackerel "role".
type Role struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

// CreateRoleParam parameters for CreateRole
type CreateRoleParam Role

// FindRoles finds roles.
func (c *Client) FindRoles(serviceName string) ([]*Role, error) {
	uri := fmt.Sprintf("/api/v0/services/%s/roles", serviceName)
	req, err := http.NewRequest("GET", c.urlFor(uri).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Roles []*Role `json:"roles"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Roles, err
}

// CreateRole creates role.
func (c *Client) CreateRole(serviceName string, param *CreateRoleParam) (*Role, error) {
	uri := fmt.Sprintf("/api/v0/services/%s/roles", serviceName)
	resp, err := c.PostJSON(uri, param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	role := &Role{}
	err = json.NewDecoder(resp.Body).Decode(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// DeleteRole deletes role.
func (c *Client) DeleteRole(serviceName, roleName string) (*Role, error) {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/services/%s/roles/%s", serviceName, roleName)).String(),
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

	role := &Role{}
	err = json.NewDecoder(resp.Body).Decode(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}
