package mackerel

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// https://mackerel.io/ja/api-docs/entry/metadata

// RoleMetaDataResp represents response for role metadata.
type RoleMetaDataResp struct {
	RoleMetaData RoleMetaData
	LastModified time.Time
}

// RoleMetaData represents role metadata body.
type RoleMetaData any

// GetRoleMetaData gets a role metadata.
func (c *Client) GetRoleMetaData(serviceName, roleName, namespace string) (*RoleMetaDataResp, error) {
	return c.GetRoleMetaDataContext(context.Background(), serviceName, roleName, namespace)
}

// GetRoleMetaDataContext gets a role metadata.
func (c *Client) GetRoleMetaDataContext(ctx context.Context, serviceName, roleName, namespace string) (*RoleMetaDataResp, error) {
	path := fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata/%s", serviceName, roleName, namespace)
	metadata, header, err := requestGetAndReturnHeaderContext[HostMetaData](ctx, c, path)
	if err != nil {
		return nil, err
	}
	lastModified, err := http.ParseTime(header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	return &RoleMetaDataResp{RoleMetaData: *metadata, LastModified: lastModified}, nil
}

// GetRoleMetaDataNameSpaces fetches namespaces of role metadata.
func (c *Client) GetRoleMetaDataNameSpaces(serviceName, roleName string) ([]string, error) {
	return c.GetRoleMetaDataNameSpacesContext(context.Background(), serviceName, roleName)
}

// GetRoleMetaDataNameSpacesContext fetches namespaces of role metadata.
func (c *Client) GetRoleMetaDataNameSpacesContext(ctx context.Context, serviceName, roleName string) ([]string, error) {
	data, err := requestGetContext[struct {
		MetaDatas []struct {
			NameSpace string `json:"namespace"`
		} `json:"metadata"`
	}](ctx, c, fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata", serviceName, roleName))
	if err != nil {
		return nil, err
	}
	namespaces := make([]string, len(data.MetaDatas))
	for i, metadata := range data.MetaDatas {
		namespaces[i] = metadata.NameSpace
	}
	return namespaces, nil
}

// PutRoleMetaData puts a role metadata.
func (c *Client) PutRoleMetaData(serviceName, roleName, namespace string, metadata RoleMetaData) error {
	return c.PutRoleMetaDataContext(context.Background(), serviceName, roleName, namespace, metadata)
}

// PutRoleMetaDataContext puts a role metadata.
func (c *Client) PutRoleMetaDataContext(ctx context.Context, serviceName, roleName, namespace string, metadata RoleMetaData) error {
	path := fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata/%s", serviceName, roleName, namespace)
	_, err := requestPutContext[any](ctx, c, path, metadata)
	return err
}

// DeleteRoleMetaData deletes a role metadata.
func (c *Client) DeleteRoleMetaData(serviceName, roleName, namespace string) error {
	path := fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata/%s", serviceName, roleName, namespace)
	_, err := requestDeleteContext[any](context.Background(), c, path)
	return err
}

// DeleteRoleMetaDataContext is like [DeleteRoleMetaData].
func (c *Client) DeleteRoleMetaDataContext(ctx context.Context, serviceName, roleName, namespace string) error {
	path := fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata/%s", serviceName, roleName, namespace)
	_, err := requestDeleteContext[any](ctx, c, path)
	return err
}
