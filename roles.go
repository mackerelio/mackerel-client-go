package mackerel

import (
	"context"
	"fmt"
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
	return c.FindRolesContext(context.Background(), serviceName)
}

// FindRolesContext finds roles.
func (c *Client) FindRolesContext(ctx context.Context, serviceName string) ([]*Role, error) {
	data, err := requestGetContext[struct {
		Roles []*Role `json:"roles"`
	}](ctx, c, fmt.Sprintf("/api/v0/services/%s/roles", serviceName))
	if err != nil {
		return nil, err
	}
	return data.Roles, nil
}

// CreateRole creates a role.
func (c *Client) CreateRole(serviceName string, param *CreateRoleParam) (*Role, error) {
	return c.CreateRoleContext(context.Background(), serviceName, param)
}

// CreateRoleContext creates a role.
func (c *Client) CreateRoleContext(ctx context.Context, serviceName string, param *CreateRoleParam) (*Role, error) {
	path := fmt.Sprintf("/api/v0/services/%s/roles", serviceName)
	return requestPostContext[Role](ctx, c, path, param)
}

// DeleteRole deletes a role.
func (c *Client) DeleteRole(serviceName, roleName string) (*Role, error) {
	path := fmt.Sprintf("/api/v0/services/%s/roles/%s", serviceName, roleName)
	return requestDeleteContext[Role](context.Background(), c, path)
}

// DeleteRoleContext is like [DeleteRole].
func (c *Client) DeleteRoleContext(ctx context.Context, serviceName, roleName string) (*Role, error) {
	path := fmt.Sprintf("/api/v0/services/%s/roles/%s", serviceName, roleName)
	return requestDeleteContext[Role](ctx, c, path)
}
