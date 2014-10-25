package mackerel

type Host struct {
	Id            string      `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	Type          string      `json:"type,omitempty"`
	Status        string      `json:"status,omitempty"`
	Memo          string      `json:"memo,omitempty"`
	Roles         Roles       `json:"roles,omitempty"`
	RoleFullnames []string    `json:"roleFullnames,omitempty"`
	IsRetired     bool        `json:"isRetired,omitempty"`
	CreatedAt     int32       `json:"createdAt,omitempty"`
	Meta          HostMeta    `json:"meta,omitempty"`
	Interfaces    []Interface `json:"interfaces,omitempty"`
}

type Roles map[string][]string

type HostMeta struct {
	AgentRevision string      `json:"agent-revision,omitempty"`
	AgentVersion  string      `json:"agent-version,omitempty"`
	BlockDevice   BlockDevice `json:"block_device,omitempty"`
	Cpu           CPU         `json:"cpu,omitempty"`
	Filesystem    FileSystem  `json:"filesystem,omitempty"`
	Kernel        Kernel      `json:"kernel,omitempty"`
	Memory        Memory      `json:"memory,omitempty"`
}

type BlockDevice map[string]map[string]interface{}
type CPU []map[string]interface{}
type FileSystem map[string]interface{}
type Kernel map[string]string
type Memory map[string]string

type Interface struct {
	Name       string `json:"name,omitempty"`
	IPAddress  string `json:"ipAddress,omitempty"`
	MacAddress string `json:"macAddress,omitempty"`
}
