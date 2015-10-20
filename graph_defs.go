package mackerel

// GraphDefsParam parameters for posting graph definitions
type GraphDefsParam struct {
	Name        string             `json:"name"`
	DisplayName string             `json:"displayName"`
	Unit        string             `json:"unit"`
	Metrics     []*GraphDefsMetric `json:"metrics"`
}

// GraphDefsMetric graph metric
type GraphDefsMetric struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IsStacked   bool   `json:"isStacked"`
}

// CreateGraphDefs create graph defs
func (c *Client) CreateGraphDefs(payloads []*GraphDefsParam) error {
	resp, err := c.PostJSON("/api/v0/graph-defs/create", payloads)
	defer closeResponse(resp)
	return err
}
