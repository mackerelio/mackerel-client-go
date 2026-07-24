package mackerel

import (
	"context"
	"fmt"
	"net/url"
)

// CheckStatus represents check monitoring status
type CheckStatus string

// CheckStatuses
const (
	CheckStatusOK       CheckStatus = "OK"
	CheckStatusWarning  CheckStatus = "WARNING"
	CheckStatusCritical CheckStatus = "CRITICAL"
	CheckStatusUnknown  CheckStatus = "UNKNOWN"
)

// CheckReport represents a report of check monitoring
type CheckReport struct {
	Source               CheckSource `json:"source"`
	Name                 string      `json:"name"`
	Status               CheckStatus `json:"status"`
	Message              string      `json:"message"`
	OccurredAt           int64       `json:"occurredAt"`
	NotificationInterval uint        `json:"notificationInterval,omitempty"`
	MaxCheckAttempts     uint        `json:"maxCheckAttempts,omitempty"`
}

// CheckSource represents interface to which each check source type must confirm to
type CheckSource interface {
	CheckType() string

	isCheckSource()
}

const checkTypeHost = "host"

// Ensure each check type conforms to the CheckSource interface.
var _ CheckSource = (*checkSourceHost)(nil)

// Ensure only checkSource types defined in this package can be assigned to the
// CheckSource interface.
func (cs *checkSourceHost) isCheckSource() {}

type checkSourceHost struct {
	Type   string `json:"type"`
	HostID string `json:"hostId"`
}

// CheckType is for satisfying CheckSource interface
func (cs *checkSourceHost) CheckType() string {
	return checkTypeHost
}

// NewCheckSourceHost returns new CheckSource which check type is "host"
func NewCheckSourceHost(hostID string) CheckSource {
	return &checkSourceHost{
		Type:   checkTypeHost,
		HostID: hostID,
	}
}

// CheckReports represents check reports for API
type CheckReports struct {
	Reports []*CheckReport `json:"reports"`
}

// CheckMonitor represents a check monitor
type CheckMonitor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// FindCheckMonitorsParam is the parameters for FindCheckMonitors
type FindCheckMonitorsParam struct {
	NextID *string
	Limit  *int
}

func (p FindCheckMonitorsParam) toValues() url.Values {
	values := url.Values{}
	if p.NextID != nil {
		values.Set("nextId", *p.NextID)
	}
	if p.Limit != nil {
		values.Set("limit", fmt.Sprintf("%d", *p.Limit))
	}
	return values
}

// FindCheckMonitorsResp is for FindCheckMonitors
type FindCheckMonitorsResp struct {
	Checks []*CheckMonitor `json:"checks"`
	NextID string          `json:"nextId,omitempty"`
}

// FindCheckMonitors finds check monitors.
func (c *Client) FindCheckMonitors(params *FindCheckMonitorsParam) (*FindCheckMonitorsResp, error) {
	return c.FindCheckMonitorsContext(context.Background(), params)
}

// FindCheckMonitorsContext finds check monitors.
func (c *Client) FindCheckMonitorsContext(ctx context.Context, params *FindCheckMonitorsParam) (*FindCheckMonitorsResp, error) {
	if params == nil {
		return requestGetContext[FindCheckMonitorsResp](ctx, c, "/api/v0/monitoring/checks")
	}
	return requestGetWithParamsContext[FindCheckMonitorsResp](ctx, c, "/api/v0/monitoring/checks", params.toValues())
}

// PostCheckReports reports check monitoring results.
func (c *Client) PostCheckReports(checkReports *CheckReports) error {
	return c.PostCheckReportsContext(context.Background(), checkReports)
}

// PostCheckReportsContext reports check monitoring results.
func (c *Client) PostCheckReportsContext(ctx context.Context, checkReports *CheckReports) error {
	_, err := requestPostContext[any](ctx, c, "/api/v0/monitoring/checks/report", checkReports)
	return err
}
