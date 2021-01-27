package mackerel

import (
	"context"

	v2 "github.com/mackerelio/mackerel-client-go/v2"
)

type Org = v2.Org

// GetOrg get the org
func (c *Client) GetOrg() (*Org, error) {
	return c.OrgService.Get(context.Background())
}
