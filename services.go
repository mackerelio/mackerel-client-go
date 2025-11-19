package mackerel

import (
	"context"
	"fmt"
)

// Service represents Mackerel "service".
type Service struct {
	Name  string   `json:"name"`
	Memo  string   `json:"memo"`
	Roles []string `json:"roles"`
}

// CreateServiceParam parameters for CreateService
type CreateServiceParam struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

// FindServices finds services.
func (c *Client) FindServices() ([]*Service, error) {
	data, err := requestGet[struct {
		Services []*Service `json:"services"`
	}](c, "/api/v0/services")
	if err != nil {
		return nil, err
	}
	return data.Services, nil
}

// CreateService creates a service.
func (c *Client) CreateService(param *CreateServiceParam) (*Service, error) {
	return requestPost[Service](c, "/api/v0/services", param)
}

// DeleteService deletes a service.
func (c *Client) DeleteService(serviceName string) (*Service, error) {
	path := fmt.Sprintf("/api/v0/services/%s", serviceName)
	return requestDeleteContext[Service](context.Background(), c, path)
}

// DeleteServiceContext is like [DeleteService].
func (c *Client) DeleteServiceContext(ctx context.Context, serviceName string) (*Service, error) {
	path := fmt.Sprintf("/api/v0/services/%s", serviceName)
	return requestDeleteContext[Service](ctx, c, path)
}

// ListServiceMetricNames lists metric names of a service.
func (c *Client) ListServiceMetricNames(serviceName string) ([]string, error) {
	data, err := requestGet[struct {
		Names []string `json:"names"`
	}](c, fmt.Sprintf("/api/v0/services/%s/metric-names", serviceName))
	if err != nil {
		return nil, err
	}
	return data.Names, nil
}

// DeleteServiceGraphDef deletes a service metrics graph definition.
func (c *Client) DeleteServiceGraphDef(serviceName string, graphName string) error {
	path := fmt.Sprintf("/api/v0/services/%s/graph-defs/%s", serviceName, graphName)
	_, err := requestDeleteContext[any](context.Background(), c, path)
	return err
}

// DeleteServiceGraphDefContext is like [DeleteServiceGraphDef].
func (c *Client) DeleteServiceGraphDefContext(ctx context.Context, serviceName string, graphName string) error {
	path := fmt.Sprintf("/api/v0/services/%s/graph-defs/%s", serviceName, graphName)
	_, err := requestDeleteContext[any](ctx, c, path)
	return err
}
