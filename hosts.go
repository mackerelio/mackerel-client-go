package mackerel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

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

type FindHostsParam struct {
	Service  string
	Roles    []string
	Name     string
	Statuses []string
}

type CreateHostParam struct {
	Name          string      `json:"name,omitempty"`
	Meta          HostMeta    `json:"meta,omitempty"`
	Interfaces    []Interface `json:"interfaces,omitempty"`
	RoleFullnames []string    `json:"roleFullnames,omitempty"`
}

type UpdateHostParam CreateHostParam

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

func (c *Client) FindHost(id string) (*Host, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/api/v0/hosts/%s", id)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Host *Host `json:"host"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Host, err
}

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

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", c.urlFor("/api/v0/hosts.json").String(), v.Encode()), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Hosts []*(Host) `json:"hosts"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Hosts, err
}

func (c *Client) CreateHost(param *CreateHostParam) (string, error) {
	requestJson, err := json.Marshal(param)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		c.urlFor("/api/v0/hosts").String(),
		bytes.NewReader(requestJson),
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data.Id, nil
}

func (c *Client) UpdateHost(hostId string, param *UpdateHostParam) (string, error) {
	requestJson, err := json.Marshal(param)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"PUT",
		c.urlFor(fmt.Sprintf("/api/v0/hosts/%s", hostId)).String(),
		bytes.NewReader(requestJson),
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data.Id, nil
}

func (c *Client) UpdateHostStatus(hostId string, status string) error {
	requestJson, err := json.Marshal(map[string]string{
		"status": status,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"PUT",
		c.urlFor(fmt.Sprintf("/api/v0/hosts/%s/status", hostId)).String(),
		bytes.NewReader(requestJson),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) RetireHost(id string) error {
	requestJson, _ := json.Marshal("{}")

	req, err := http.NewRequest(
		"POST",
		c.urlFor(fmt.Sprintf("/api/v0/hosts/%s/retire", id)).String(),
		bytes.NewReader(requestJson),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("status code is not 200")
	}

	return nil
}
