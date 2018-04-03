package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Host host information
type Host struct {
	ID               string      `json:"id,omitempty"`
	Name             string      `json:"name,omitempty"`
	DisplayName      string      `json:"displayName,omitempty"`
	CustomIdentifier string      `json:"customIdentifier,omitempty"`
	Type             string      `json:"type,omitempty"`
	Status           string      `json:"status,omitempty"`
	Memo             string      `json:"memo,omitempty"`
	Roles            Roles       `json:"roles,omitempty"`
	RoleFullnames    []string    `json:"roleFullnames,omitempty"`
	IsRetired        bool        `json:"isRetired,omitempty"`
	CreatedAt        int32       `json:"createdAt,omitempty"`
	Meta             HostMeta    `json:"meta,omitempty"`
	Interfaces       []Interface `json:"interfaces,omitempty"`
}

// Roles host role maps
type Roles map[string][]string

// HostMeta host meta informations
type HostMeta struct {
	AgentRevision string      `json:"agent-revision,omitempty"`
	AgentVersion  string      `json:"agent-version,omitempty"`
	BlockDevice   BlockDevice `json:"block_device,omitempty"`
	CPU           CPU         `json:"cpu,omitempty"`
	Filesystem    FileSystem  `json:"filesystem,omitempty"`
	Kernel        Kernel      `json:"kernel,omitempty"`
	Memory        Memory      `json:"memory,omitempty"`
	Cloud         Cloud       `json:"cloud,omitempty"`
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
	Name             string      `json:"name,omitempty"`
	DisplayName      string      `json:"displayName,omitempty"`
	Meta             HostMeta    `json:"meta,omitempty"`
	Interfaces       []Interface `json:"interfaces,omitempty"`
	RoleFullnames    []string    `json:"roleFullnames,omitempty"`
	CustomIdentifier string      `json:"customIdentifier,omitempty"`
}

// UpdateHostParam parameters for UpdateHost
type UpdateHostParam CreateHostParam

// GetRoleFullnames getrolefullnames
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

// DateStringFromCreatedAt returns date string
func (h *Host) DateStringFromCreatedAt() string {
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	return h.DateFromCreatedAt().Format(layout)
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

// FindHost find the host
func (c *Client) FindHost(id string) (*Host, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/api/v0/hosts/%s", id)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Host *Host `json:"host"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Host, err
}

// FindHosts find hosts
func (c *Client) FindHosts(param *FindHostsParam) ([]*Host, error) {
	v := url.Values{}
	if param.Service != "" {
		v.Set("service", param.Service)
	}
	if len(param.Roles) >= 1 {
		for _, role := range param.Roles {
			v.Add("role", role)
		}
	}
	if param.Name != "" {
		v.Set("name", param.Name)
	}
	if len(param.Statuses) >= 1 {
		for _, status := range param.Statuses {
			v.Add("status", status)
		}
	}
	if param.CustomIdentifier != "" {
		v.Set("customIdentifier", param.CustomIdentifier)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", c.urlFor("/api/v0/hosts").String(), v.Encode()), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Hosts []*(Host) `json:"hosts"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Hosts, err
}

// CreateHost creating host
func (c *Client) CreateHost(param *CreateHostParam) (string, error) {
	resp, err := c.PostJSON("/api/v0/hosts", param)
	defer closeResponse(resp)
	if err != nil {
		return "", err
	}

	var data struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.ID, nil
}

// UpdateHost updates host
func (c *Client) UpdateHost(hostID string, param *UpdateHostParam) (string, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/api/v0/hosts/%s", hostID), param)
	defer closeResponse(resp)
	if err != nil {
		return "", err
	}

	var data struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

// UpdateHostStatus updates host status
func (c *Client) UpdateHostStatus(hostID string, status string) error {
	resp, err := c.PostJSON(fmt.Sprintf("/api/v0/hosts/%s/status", hostID), map[string]string{
		"status": status,
	})
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}

// UpdateHostRoleFullnames updates host roles
func (c *Client) UpdateHostRoleFullnames(hostID string, roleFullnames []string) error {
	resp, err := c.PutJSON(fmt.Sprintf("/api/v0/hosts/%s/role-fullnames", hostID), map[string][]string{
		"roleFullnames": roleFullnames,
	})
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}

// RetireHost retuire the host
func (c *Client) RetireHost(id string) error {
	resp, err := c.PostJSON(fmt.Sprintf("/api/v0/hosts/%s/retire", id), "{}")
	defer closeResponse(resp)
	return err
}
