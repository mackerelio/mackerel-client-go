package mackerel

import (
	"context"
	"fmt"
	"net/http"
)

// GraphDefsParam parameters for posting graph definitions
type GraphDefsParam struct {
	Name        string             `json:"name"`
	DisplayName string             `json:"displayName,omitempty"`
	Unit        string             `json:"unit,omitempty"`
	Metrics     []*GraphDefsMetric `json:"metrics"`
	// If this field is nil, the server treats it as true.
	// false can be specified only if a metric name in Metrics contains a wildcard.
	IsAutoDeleteEnabled *bool `json:"isAutoDeleteEnabled,omitempty"`
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

type CreateGraphDefsResponse struct {
	InvalidGraphNames []string `json:"invalidGraphNames"`
}

func (e *CreateGraphDefsResponse) Error() string {
	return fmt.Sprintf("failed create graph defs. GraphNames: %q", e.InvalidGraphNames)
}

type createGraphDefsResp struct {
	Success           bool     `json:"success"`
	InvalidGraphNames []string `json:"invalidGraphNames"`
}

// CreateGraphDefsContext creates graph definitions.
func (c *Client) CreateGraphDefsContext(ctx context.Context, graphDefs []*GraphDefsParam) error {
	data, err := requestPostContext[createGraphDefsResp](ctx, c, "/api/v0/graph-defs/create", graphDefs)
	if err != nil {
		return err
	}
	if data.Success {
		return nil
	}
	return &CreateGraphDefsResponse{InvalidGraphNames: data.InvalidGraphNames}
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
