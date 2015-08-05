package mackerel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
{
  "monitors": [
    {
      "id": "2cSZzK3XfmG",
      "type": "passive",
      "name": "connectivity",
      "duration": 5
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "host",
      "name": "disk.aa-00.writes.delta",
      "duration": 3,
      "metric": "disk.aa-00.writes.delta",
      "operator": ">",
      "warning": 20000.0,
      "critical": 400000.0,
      "scopes": [
        "SomeService"
      ],
      "excludeScopes": [
        "SomeService: db-slave-backup"
      ]
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "service",
      "name": "SomeService - custom.access_num.4xx_count",
      "service": "SomeService",
      "duration": 1,
      "metric": "custom.access_num.4xx_count",
      "operator": ">",
      "warning": 50.0,
      "critical": 100.0
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "external",
      "name": "example.com",
      "url": "http://www.example.com"
    }
  ]
}
*/

// Monitor information
type Monitor struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Type          string   `json:"type,omitempty"`
	Duration      uint64   `json:"duration,omitempty"`
	Scopes        []string `json:"scopes,omitempty"`
	ExcludeScopes []string `json:"excludeScopes,omitempty"`
}

// FindMonitors find monitors
func (c *Client) FindMonitors() ([]*Monitor, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s", c.urlFor("/api/v0/monitors").String()), nil)
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
		Monitors []*(Monitor) `json:"monitors"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Monitors, err
}
