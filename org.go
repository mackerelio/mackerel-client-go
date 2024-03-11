package mackerel

// Org information
type Org struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

// GetOrg gets the org.
func (c *Client) GetOrg() (*Org, error) {
	return requestGet[Org](c, "/api/v0/org")
}
