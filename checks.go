package mackerel

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

// PostCheckReports reports check monitoring results
func (c *Client) PostCheckReports(crs *CheckReports) error {
	resp, err := c.PostJSON("/api/v0/monitoring/checks/report", crs)
	defer closeResponse(resp)
	return err
}
