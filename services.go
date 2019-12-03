package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Service represents Mackerel "service".
type Service struct {
	Name  string   `json:"name"`
	Memo  string   `json:"memo"`
	Roles []string `json:"roles"`
}

// CreateServiceParam parameters for CreateService
type CreateServiceParam struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
}

// FindServices finds services.
func (c *Client) FindServices() ([]*Service, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/services").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Services []*Service `json:"services"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Services, err
}

// CreateService creates service
func (c *Client) CreateService(param *CreateServiceParam) (*Service, error) {
	resp, err := c.PostJSON("/api/v0/services", param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	service := &Service{}
	err = json.NewDecoder(resp.Body).Decode(service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// DeleteService deletes service
func (c *Client) DeleteService(serviceName string) (*Service, error) {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/services/%s", serviceName)).String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	service := &Service{}
	err = json.NewDecoder(resp.Body).Decode(service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// ListServiceMetricNames lists metric names of a service
func (c *Client) ListServiceMetricNames(serviceName string) ([]string, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/api/v0/services/%s/metric-names", serviceName)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Names []string `json:"names"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Names, err
}
