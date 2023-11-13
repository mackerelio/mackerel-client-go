package mackerel

import (
	"fmt"
	"net/http"
	"time"
)

// https://mackerel.io/ja/api-docs/entry/metadata

// HostMetaDataResp represents response for host metadata.
type HostMetaDataResp struct {
	HostMetaData HostMetaData
	LastModified time.Time
}

// HostMetaData represents host metadata body.
type HostMetaData interface{}

// GetHostMetaData gets a host metadata.
func (c *Client) GetHostMetaData(hostID, namespace string) (*HostMetaDataResp, error) {
	path := fmt.Sprintf("/api/v0/hosts/%s/metadata/%s", hostID, namespace)
	metadata, header, err := requestGetAndReturnHeader[HostMetaData](c, path)
	if err != nil {
		return nil, err
	}
	lastModified, err := http.ParseTime(header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	return &HostMetaDataResp{HostMetaData: *metadata, LastModified: lastModified}, nil
}

// GetHostMetaDataNameSpaces fetches namespaces of host metadata.
func (c *Client) GetHostMetaDataNameSpaces(hostID string) ([]string, error) {
	data, err := requestGet[struct {
		MetaDatas []struct {
			NameSpace string `json:"namespace"`
		} `json:"metadata"`
	}](c, fmt.Sprintf("/api/v0/hosts/%s/metadata", hostID))
	if err != nil {
		return nil, err
	}
	namespaces := make([]string, len(data.MetaDatas))
	for i, metadata := range data.MetaDatas {
		namespaces[i] = metadata.NameSpace
	}
	return namespaces, nil
}

// PutHostMetaData puts a host metadata.
func (c *Client) PutHostMetaData(hostID, namespace string, metadata HostMetaData) error {
	path := fmt.Sprintf("/api/v0/hosts/%s/metadata/%s", hostID, namespace)
	_, err := requestPut[any](c, path, metadata)
	return err
}

// DeleteHostMetaData deletes a host metadata.
func (c *Client) DeleteHostMetaData(hostID, namespace string) error {
	path := fmt.Sprintf("/api/v0/hosts/%s/metadata/%s", hostID, namespace)
	_, err := requestDelete[any](c, path)
	return err
}
