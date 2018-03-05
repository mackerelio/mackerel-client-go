package mackerel

import (
	"encoding/json"
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
type HostMetaData map[string]interface{}

// GetHostMetaData find host metadata.
func (c *Client) GetHostMetaData(hostID, namespace string) (*HostMetaDataResp, error) {
	url := c.urlFor(fmt.Sprintf("/api/v0/hosts/%s/metadata/%s", hostID, namespace))
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	var data HostMetaDataResp
	if err := json.NewDecoder(resp.Body).Decode(&data.HostMetaData); err != nil {
		return nil, err
	}
	data.LastModified, err = http.ParseTime(resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	return &data, nil
}
