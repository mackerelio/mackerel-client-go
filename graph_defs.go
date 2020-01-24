package mackerel

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

// CreateGraphDefs create graph defs
func (c *Client) CreateGraphDefs(payloads []*GraphDefsParam) error {
	resp, err := c.PostJSON("/api/v0/graph-defs/create", payloads)
	defer closeResponse(resp)
	return err
}
