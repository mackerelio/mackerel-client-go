package mackerel

// GraphAnnotation parameters to post graph annotation
type GraphAnnotation struct {
	Service     string   `json:"service,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	From        int64    `json:"from,omitempty"`
	To          int64    `json:"to,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
}

// CreateGraphAnnotation create graph annotation
func (c *Client) CreateGraphAnnotation(payloads *GraphAnnotation) error {
	resp, err := c.PostJSON("/api/v0/graph-annotations", payloads)
	defer closeResponse(resp)
	return err
}
