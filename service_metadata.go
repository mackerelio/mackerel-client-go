package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// https://mackerel.io/ja/api-docs/entry/metadata

// ServiceMetaDataResp represents response for service metadata.
type ServiceMetaDataResp struct {
	ServiceMetaData ServiceMetaData
	LastModified    time.Time
}

// ServiceMetaData represents service metadata body.
type ServiceMetaData interface{}

// GetServiceMetaData find service metadata.
func (c *Client) GetServiceMetaData(serviceName, namespace string) (*ServiceMetaDataResp, error) {
	url := c.urlFor(fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace))
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	var data ServiceMetaDataResp
	if err := json.NewDecoder(resp.Body).Decode(&data.ServiceMetaData); err != nil {
		return nil, err
	}
	data.LastModified, err = http.ParseTime(resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetServiceMetaDataNameSpaces fetches namespaces of service metadata.
func (c *Client) GetServiceMetaDataNameSpaces(serviceName string) ([]string, error) {
	url := c.urlFor(fmt.Sprintf("/api/v0/services/%s/metadata", serviceName))
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	var data struct {
		MetaDatas []struct {
			NameSpace string `json:"namespace"`
		} `json:"metadata"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	namespaces := make([]string, 0, len(data.MetaDatas))
	for _, metadata := range data.MetaDatas {
		namespaces = append(namespaces, metadata.NameSpace)
	}
	return namespaces, nil
}

// PutServiceMetaData put service metadata.
func (c *Client) PutServiceMetaData(serviceName, namespace string, metadata ServiceMetaData) error {
	path := fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace)
	resp, err := c.PutJSON(path, metadata)
	defer closeResponse(resp)
	return err
}

// DeleteServiceMetaData delete service metadata.
func (c *Client) DeleteServiceMetaData(serviceName, namespace string) error {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace)).String(),
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	defer closeResponse(resp)
	return err
}
