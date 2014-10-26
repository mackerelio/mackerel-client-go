package mackerel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HostMetricValue struct {
	HostId string      `json:"hostId",omitempty`
	Name   string      `json:"name,omitempty"`
	Time   float64     `json:"time,omitempty"`
	Value  interface{} `json:"value,omitempty"`
}

type ServiceMetricValue struct {
	Name  string      `json:"name,omitempty"`
	Time  float64     `json:"time,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

func (c *Client) PostHostMetricValues(metricValues [](*HostMetricValue)) error {
	requestJson, err := json.Marshal(metricValues)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		c.urlFor("/api/v0/tsdb").String(),
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

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("API result failed: %s", resp.Status))
	}

	return nil
}

func (c *Client) PostServiceMetricValues(serviceName string, metricValues [](*ServiceMetricValue)) error {
	requestJson, err := json.Marshal(metricValues)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		c.urlFor(fmt.Sprintf("/api/v0/services/%s/tsdb", serviceName)).String(),
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

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("API result failed: %s", resp.Status))
	}

	return nil
}
