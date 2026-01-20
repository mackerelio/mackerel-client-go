package mackerel

import "context"

// Org information
type Org struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

// GetOrg gets the org.
func (c *Client) GetOrg() (*Org, error) {
	return c.GetOrgContext(context.Background())
}

// GetOrgContext gets the org.
func (c *Client) GetOrgContext(ctx context.Context) (*Org, error) {
	return requestGetContext[Org](ctx, c, "/api/v0/org")
}
