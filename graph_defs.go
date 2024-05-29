package mackerel

import "net/http"

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
	_, err := requestPost[any](c, "/api/v0/graph-defs/create", graphDefs)
	return err
}

// DeleteGraphDef deletes a graph definition.
func (c *Client) DeleteGraphDef(name string) error {
	_, err := requestJSON[any](c, http.MethodDelete, "/api/v0/graph-defs/delete", map[string]string{"name": name})
	return err
}
