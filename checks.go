package mackerel

const (
	CheckStatusOK       = "OK"
	CheckStatusWarning  = "WARNING"
	CheckStatusCritical = "CRITICAL"
	CheckStatusUnknown  = "UNKNOWN"
)

type CheckReport struct {
	Source               CheckSource `json:"source"`
	Name                 string      `json:"name"`
	Status               string      `json:"status"`
	Message              string      `json:"message"`
	OccurredAt           int64       `json:"occurredAt"`
	NotificationInterval uint        `json:"notificationInterval,omitempty"`
	MaxCheckAttempts     uint        `json:"maxCheckAttempts,omitempty"`
}

type CheckSource interface {
	CheckType() string

	isCheck()
}

const checkTypeHost = "host"

// Ensure each check type conforms to the CheckSource interface.
var _ CheckSource = (*checkSourceHost)(nil)

// Ensure only checkSource types defined in this package can be assigned to the
// CheckSource interface.
func (m *checkSourceHost) isCheck() {}

type checkSourceHost struct {
	Type   string `json:"type"`
	HostID string `json:"hostId"`
}

func (cs *checkSourceHost) CheckType() string {
	return checkTypeHost
}

func NewCheckSourceHost(hostID string) CheckSource {
	return &checkSourceHost{
		Type:   checkTypeHost,
		HostID: hostID,
	}
}

type CheckReports struct {
	Reports []*CheckReport `json:"reports"`
}

func (c *Client) ReportCheckMonitors(crs *CheckReports) error {
	resp, err := c.PostJSON("/api/v0/monitoring/checks/report", crs)
	defer closeResponse(resp)
	return err
}
