package mackerel

import (
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

// GetServiceMetaData gets service metadata.
func (c *Client) GetServiceMetaData(serviceName, namespace string) (*ServiceMetaDataResp, error) {
	path := fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace)
	metadata, header, err := requestGetAndReturnHeader[HostMetaData](c, path)
	if err != nil {
		return nil, err
	}
	lastModified, err := http.ParseTime(header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	return &ServiceMetaDataResp{ServiceMetaData: *metadata, LastModified: lastModified}, nil
}

// GetServiceMetaDataNameSpaces fetches namespaces of service metadata.
func (c *Client) GetServiceMetaDataNameSpaces(serviceName string) ([]string, error) {
	data, err := requestGet[struct {
		MetaDatas []struct {
			NameSpace string `json:"namespace"`
		} `json:"metadata"`
	}](c, fmt.Sprintf("/api/v0/services/%s/metadata", serviceName))
	if err != nil {
		return nil, err
	}
	namespaces := make([]string, len(data.MetaDatas))
	for i, metadata := range data.MetaDatas {
		namespaces[i] = metadata.NameSpace
	}
	return namespaces, nil
}

// PutServiceMetaData puts a service metadata.
func (c *Client) PutServiceMetaData(serviceName, namespace string, metadata ServiceMetaData) error {
	path := fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace)
	_, err := requestPut[any](c, path, metadata)
	return err
}

// DeleteServiceMetaData deletes a service metadata.
func (c *Client) DeleteServiceMetaData(serviceName, namespace string) error {
	path := fmt.Sprintf("/api/v0/services/%s/metadata/%s", serviceName, namespace)
	_, err := requestDelete[any](c, path)
	return err
}
