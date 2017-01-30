package mackerel

// GraphAnnotation represents parameters to post graph annotation.
type GraphAnnotation struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	From        int64    `json:"from,omitempty"`
	To          int64    `json:"to,omitempty"`
	Service     string   `json:"service,omitempty"`
	Roles       []string `json:"roles,omitempty"`
}

// CreateGraphAnnotation creates graph annotation.
func (c *Client) CreateGraphAnnotation(annotation *GraphAnnotation) error {
	resp, err := c.PostJSON("/api/v0/graph-annotations", annotation)
	defer closeResponse(resp)
	return err
}
