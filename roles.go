package mackerel

import "fmt"

// Role represents Mackerel "role".
type Role struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

// CreateRoleParam parameters for CreateRole
type CreateRoleParam Role

// FindRoles finds roles.
func (c *Client) FindRoles(serviceName string) ([]*Role, error) {
	data, err := requestGet[struct {
		Roles []*Role `json:"roles"`
	}](c, fmt.Sprintf("/api/v0/services/%s/roles", serviceName))
	if err != nil {
		return nil, err
	}
	return data.Roles, nil
}

// CreateRole creates a role.
func (c *Client) CreateRole(serviceName string, param *CreateRoleParam) (*Role, error) {
	path := fmt.Sprintf("/api/v0/services/%s/roles", serviceName)
	return requestPost[Role](c, path, param)
}

// DeleteRole deletes a role.
func (c *Client) DeleteRole(serviceName, roleName string) (*Role, error) {
	path := fmt.Sprintf("/api/v0/services/%s/roles/%s", serviceName, roleName)
	return requestDelete[Role](c, path)
}
