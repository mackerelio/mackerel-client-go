package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// https://mackerel.io/ja/api-docs/entry/metadata

// RoleMetaDataResp represents response for role metadata.
type RoleMetaDataResp struct {
	RoleMetaData RoleMetaData
	LastModified time.Time
}

// RoleMetaData represents role metadata body.
type RoleMetaData interface{}

// GetRoleMetaData find role metadata.
func (c *Client) GetRoleMetaData(serviceName, roleName, namespace string) (*RoleMetaDataResp, error) {
	url := c.urlFor(fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata/%s", serviceName, roleName, namespace))
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	var data RoleMetaDataResp
	if err := json.NewDecoder(resp.Body).Decode(&data.RoleMetaData); err != nil {
		return nil, err
	}
	data.LastModified, err = http.ParseTime(resp.Header.Get("Last-Modified"))
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetRoleMetaDataNameSpaces fetches namespaces of role metadata.
func (c *Client) GetRoleMetaDataNameSpaces(serviceName, roleName string) ([]string, error) {
	url := c.urlFor(fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata", serviceName, roleName))
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

// PutRoleMetaData put role metadata.
func (c *Client) PutRoleMetaData(serviceName, roleName, namespace string, metadata RoleMetaData) error {
	path := fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata/%s", serviceName, roleName, namespace)
	resp, err := c.PutJSON(path, metadata)
	defer closeResponse(resp)
	return err
}

// DeleteRoleMetaData delete role metadata.
func (c *Client) DeleteRoleMetaData(serviceName, roleName, namespace string) error {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/services/%s/roles/%s/metadata/%s", serviceName, roleName, namespace)).String(),
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
