package mackerel

import (
	"context"
	"net/http"
)

// GraphDefsParam parameters for posting graph definitions
type GraphDefsParam struct {
	Name        string             `json:"name"`
	DisplayName string             `json:"displayName,omitempty"`
	Unit        string             `json:"unit,omitempty"`
	Metrics     []*GraphDefsMetric `json:"metrics"`
}

// GraphDefsMetric graph metric
type GraphDefsMetric struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	IsStacked   bool   `json:"isStacked"`
}

// CreateGraphDefs creates graph definitions.
func (c *Client) CreateGraphDefs(graphDefs []*GraphDefsParam) error {
	return c.CreateGraphDefsContext(context.Background(), graphDefs)
}

// CreateGraphDefsContext creates graph definitions.
func (c *Client) CreateGraphDefsContext(ctx context.Context, graphDefs []*GraphDefsParam) error {
	_, err := requestPostContext[any](ctx, c, "/api/v0/graph-defs/create", graphDefs)
	return err
}

// DeleteGraphDef deletes a graph definition.
func (c *Client) DeleteGraphDef(name string) error {
	return c.DeleteGraphDefContext(context.Background(), name)
}

// DeleteGraphDefContext deletes a graph definition.
func (c *Client) DeleteGraphDefContext(ctx context.Context, name string) error {
	_, err := requestJSON[any](ctx, c, http.MethodDelete, "/api/v0/graph-defs", map[string]string{"name": name})
	return err
}
