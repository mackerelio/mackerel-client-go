package mackerel

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Host host information
type Host struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	DisplayName      string      `json:"displayName,omitempty"`
	CustomIdentifier string      `json:"customIdentifier,omitempty"`
	Size             string      `json:"size"`
	Status           string      `json:"status"`
	Memo             string      `json:"memo"`
	Roles            Roles       `json:"roles"`
	IsRetired        bool        `json:"isRetired"`
	CreatedAt        int32       `json:"createdAt"`
	Meta             HostMeta    `json:"meta"`
	Interfaces       []Interface `json:"interfaces"`
}

// Roles host role maps
type Roles map[string][]string

// HostMeta host meta informations
type HostMeta struct {
	AgentRevision string      `json:"agent-revision,omitempty"`
	AgentVersion  string      `json:"agent-version,omitempty"`
	AgentName     string      `json:"agent-name,omitempty"`
	BlockDevice   BlockDevice `json:"block_device,omitempty"`
	CPU           CPU         `json:"cpu,omitempty"`
	Filesystem    FileSystem  `json:"filesystem,omitempty"`
	Kernel        Kernel      `json:"kernel,omitempty"`
	Memory        Memory      `json:"memory,omitempty"`
	Cloud         *Cloud      `json:"cloud,omitempty"`
}

// BlockDevice blockdevice
type BlockDevice map[string]map[string]interface{}

// CPU cpu
type CPU []map[string]interface{}

// FileSystem filesystem
type FileSystem map[string]interface{}

// Kernel kernel
type Kernel map[string]string

// Memory memory
type Memory map[string]string

// Cloud cloud
type Cloud struct {
	Provider string      `json:"provider,omitempty"`
	MetaData interface{} `json:"metadata,omitempty"`
}

// Interface network interface
type Interface struct {
	Name          string   `json:"name,omitempty"`
	IPAddress     string   `json:"ipAddress,omitempty"`
	IPv4Addresses []string `json:"ipv4Addresses,omitempty"`
	IPv6Addresses []string `json:"ipv6Addresses,omitempty"`
	MacAddress    string   `json:"macAddress,omitempty"`
}

// FindHostsParam parameters for FindHosts
type FindHostsParam struct {
	Service          string
	Roles            []string
	Name             string
	Statuses         []string
	CustomIdentifier string
}

// CreateHostParam parameters for CreateHost
type CreateHostParam struct {
	Name             string        `json:"name"`
	DisplayName      string        `json:"displayName,omitempty"`
	Memo             string        `json:"memo,omitempty"`
	Meta             HostMeta      `json:"meta"`
	Interfaces       []Interface   `json:"interfaces"`
	RoleFullnames    []string      `json:"roleFullnames"`
	Checks           []CheckConfig `json:"checks"`
	CustomIdentifier string        `json:"customIdentifier,omitempty"`
}

// CheckConfig is check plugin name and memo
type CheckConfig struct {
	Name string `json:"name,omitempty"`
	Memo string `json:"memo,omitempty"`
}

// UpdateHostParam parameters for UpdateHost
type UpdateHostParam CreateHostParam

// MonitoredStatus monitored status
type MonitoredStatus struct {
	MonitorID string                `json:"monitorId"`
	Status    string                `json:"status"`
	Detail    MonitoredStatusDetail `json:"detail,omitempty"`
}

// MonitoredStatusDetail monitored status detail
type MonitoredStatusDetail struct {
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
	Memo    string `json:"memo,omitempty"`
}

// FindHostByCustomIdentifierParam parameters for FindHostByCustomIdentifier
type FindHostByCustomIdentifierParam struct {
	CaseInsensitive bool
}

const (
	// HostStatusWorking represents "working" status
	HostStatusWorking = "working"
	// HostStatusStandby represents "standby" status
	HostStatusStandby = "standby"
	// HostStatusMaintenance represents "maintenance" status
	HostStatusMaintenance = "maintenance"
	// HostStatusPoweroff represents "poeroff" status
	HostStatusPoweroff = "poweroff"
)

// GetRoleFullnames gets role-full-names
func (h *Host) GetRoleFullnames() []string {
	if len(h.Roles) < 1 {
		return nil
	}

	var fullnames []string
	for service, roles := range h.Roles {
		for _, role := range roles {
			fullname := strings.Join([]string{service, role}, ":")
			fullnames = append(fullnames, fullname)
		}
	}

	return fullnames
}

// DateFromCreatedAt returns time.Time
func (h *Host) DateFromCreatedAt() time.Time {
	return time.Unix(int64(h.CreatedAt), 0)
}

// IPAddresses returns ipaddresses
func (h *Host) IPAddresses() map[string]string {
	if len(h.Interfaces) < 1 {
		return nil
	}

	ipAddresses := make(map[string]string, 0)
	for _, iface := range h.Interfaces {
		ipAddresses[iface.Name] = iface.IPAddress
	}
	return ipAddresses
}

// FindHost finds the host.
func (c *Client) FindHost(hostID string) (*Host, error) {
	data, err := requestGet[struct {
		Host *Host `json:"host"`
	}](c, fmt.Sprintf("/api/v0/hosts/%s", hostID))
	if err != nil {
		return nil, err
	}
	return data.Host, nil
}

// FindHosts finds hosts.
func (c *Client) FindHosts(param *FindHostsParam) ([]*Host, error) {
	params := url.Values{}
	if param.Service != "" {
		params.Set("service", param.Service)
	}
	for _, role := range param.Roles {
		params.Add("role", role)
	}
	if param.Name != "" {
		params.Set("name", param.Name)
	}
	for _, status := range param.Statuses {
		params.Add("status", status)
	}
	if param.CustomIdentifier != "" {
		params.Set("customIdentifier", param.CustomIdentifier)
	}

	data, err := requestGetWithParams[struct {
		Hosts []*Host `json:"hosts"`
	}](c, "/api/v0/hosts", params)
	if err != nil {
		return nil, err
	}
	return data.Hosts, nil
}

// FindHostByCustomIdentifier finds a host by the custom identifier.
func (c *Client) FindHostByCustomIdentifier(customIdentifier string, param *FindHostByCustomIdentifierParam) (*Host, error) {
	path := "/api/v0/hosts-by-custom-identifier/" + url.PathEscape(customIdentifier)
	params := url.Values{}
	if param.CaseInsensitive {
		params.Set("caseInsensitive", "true")
	}
	data, err := requestGetWithParams[struct {
		Host *Host `json:"host"`
	}](c, path, params)
	if err != nil {
		return nil, err
	}
	return data.Host, nil
}

// CreateHost creates a host.
func (c *Client) CreateHost(param *CreateHostParam) (string, error) {
	data, err := requestPost[struct {
		ID string `json:"id"`
	}](c, "/api/v0/hosts", param)
	if err != nil {
		return "", err
	}
	return data.ID, nil
}

// UpdateHost updates a host.
func (c *Client) UpdateHost(hostID string, param *UpdateHostParam) (string, error) {
	path := fmt.Sprintf("/api/v0/hosts/%s", hostID)
	data, err := requestPut[struct {
		ID string `json:"id"`
	}](c, path, param)
	if err != nil {
		return "", err
	}
	return data.ID, nil
}

// UpdateHostStatus updates a host status.
func (c *Client) UpdateHostStatus(hostID string, status string) error {
	path := fmt.Sprintf("/api/v0/hosts/%s/status", hostID)
	_, err := requestPost[any](c, path, map[string]string{"status": status})
	return err
}

// UpdateHostRoleFullnames updates host roles.
func (c *Client) UpdateHostRoleFullnames(hostID string, roleFullnames []string) error {
	path := fmt.Sprintf("/api/v0/hosts/%s/role-fullnames", hostID)
	_, err := requestPut[any](c, path, map[string][]string{"roleFullnames": roleFullnames})
	return err
}

// RetireHost retires the host.
func (c *Client) RetireHost(hostID string) error {
	path := fmt.Sprintf("/api/v0/hosts/%s/retire", hostID)
	_, err := requestPost[any](c, path, nil)
	return err
}

// BulkRetireHosts retires the hosts.
func (c *Client) BulkRetireHosts(ids []string) error {
	_, err := requestPost[any](c, "/api/v0/hosts/bulk-retire", map[string][]string{"ids": ids})
	return err
}

// ListHostMetricNames lists metric names of a host.
func (c *Client) ListHostMetricNames(hostID string) ([]string, error) {
	data, err := requestGet[struct {
		Names []string `json:"names"`
	}](c, fmt.Sprintf("/api/v0/hosts/%s/metric-names", hostID))
	if err != nil {
		return nil, err
	}
	return data.Names, nil
}

// ListMonitoredStatues lists monitored statues of a host.
func (c *Client) ListMonitoredStatues(hostID string) ([]MonitoredStatus, error) {
	data, err := requestGet[struct {
		MonitoredStatuses []MonitoredStatus `json:"monitoredStatuses"`
	}](c, fmt.Sprintf("/api/v0/hosts/%s/monitored-statuses", hostID))
	if err != nil {
		return nil, err
	}
	return data.MonitoredStatuses, nil
}
