package mackerel

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// https://mackerel.io/ja/api-docs/entry/metadata

// ServiceMetaDataResp represents response for service metadata.
type ServiceMetaDataResp struct {
	ServiceMetaData ServiceMetaData
	LastModified    time.Time
}

// ServiceMetaData represents service metadata body.
type ServiceMetaData interface{}

// GetServiceMetaData gets service metadata.
func (c *Client) GetServiceMetaData(serviceName, namespace string) (*ServiceMetaDataResp, error) {
	return c.GetServiceMetaDataContext(context.Background(), serviceName, namespace)
}

// GetServiceMetaDataContext gets service metadata.
func (c *Client) GetServiceMetaDataContext(ctx context.Context, serviceName, namespace string) (*ServiceMetaDataResp, error) {
	path := fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace)
	metadata, header, err := requestGetAndReturnHeaderContext[HostMetaData](ctx, c, path)
	if err != nil {
		return nil, err
	}
	lastModified, err := http.ParseTime(header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	return &ServiceMetaDataResp{ServiceMetaData: *metadata, LastModified: lastModified}, nil
}

// GetServiceMetaDataNameSpaces fetches namespaces of service metadata.
func (c *Client) GetServiceMetaDataNameSpaces(serviceName string) ([]string, error) {
	return c.GetServiceMetaDataNameSpacesContext(context.Background(), serviceName)
}

// GetServiceMetaDataNameSpacesContext fetches namespaces of service metadata.
func (c *Client) GetServiceMetaDataNameSpacesContext(ctx context.Context, serviceName string) ([]string, error) {
	data, err := requestGetContext[struct {
		MetaDatas []struct {
			NameSpace string `json:"namespace"`
		} `json:"metadata"`
	}](ctx, c, fmt.Sprintf("/api/v0/services/%s/metadata", serviceName))
	if err != nil {
		return nil, err
	}
	namespaces := make([]string, len(data.MetaDatas))
	for i, metadata := range data.MetaDatas {
		namespaces[i] = metadata.NameSpace
	}
	return namespaces, nil
}

// PutServiceMetaData puts a service metadata.
func (c *Client) PutServiceMetaData(serviceName, namespace string, metadata ServiceMetaData) error {
	return c.PutServiceMetaDataContext(context.Background(), serviceName, namespace, metadata)
}

// v puts a service metadata.
func (c *Client) PutServiceMetaDataContext(ctx context.Context, serviceName, namespace string, metadata ServiceMetaData) error {
	path := fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace)
	_, err := requestPutWithContext[any](ctx, c, path, metadata)
	return err
}

// DeleteServiceMetaData deletes a service metadata.
func (c *Client) DeleteServiceMetaData(serviceName, namespace string) error {
	return c.DeleteServiceMetaDataContext(context.Background(), serviceName, namespace)
}

// DeleteServiceMetaDataContext is like [DeleteServiceMetaData].
func (c *Client) DeleteServiceMetaDataContext(ctx context.Context, serviceName, namespace string) error {
	path := fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace)
	_, err := requestDeleteContext[any](ctx, c, path)
	return err
}
