package mackerel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// MetricValue metric value
type MetricValue struct {
	Name  string      `json:"name,omitempty"`
	Time  int64       `json:"time,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// HostMetricValue host metric value
type HostMetricValue struct {
	HostID string      `json:"hostId,omitempty"`
	Name   string      `json:"name,omitempty"`
	Time   int64       `json:"time,omitempty"`
	Value  interface{} `json:"value,omitempty"`
}

// LatestMetricValues latest metric value
type LatestMetricValues map[string]map[string]*MetricValue

// PostHostMetricValues post host metrics
func (c *Client) PostHostMetricValues(metricValues [](*HostMetricValue)) error {
	requestJSON, err := json.Marshal(metricValues)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		c.urlFor("/api/v0/tsdb").String(),
		bytes.NewReader(requestJSON),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	defer c.CloseReponse(resp)
	if err != nil {
		return err
	}

	return nil
}

// PostHostMetricValuesByHostID post host metrics
func (c *Client) PostHostMetricValuesByHostID(hostID string, metricValues [](*MetricValue)) error {
	var hostMetricValues []*HostMetricValue
	for _, metricValue := range metricValues {
		hostMetricValues = append(hostMetricValues, &HostMetricValue{
			HostID: hostID,
			Name:   metricValue.Name,
			Value:  metricValue.Value,
			Time:   metricValue.Time,
		})
	}
	return c.PostHostMetricValues(hostMetricValues)
}

// PostServiceMetricValues post service metrics
func (c *Client) PostServiceMetricValues(serviceName string, metricValues [](*MetricValue)) error {
	requestJSON, err := json.Marshal(metricValues)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		c.urlFor(fmt.Sprintf("/api/v0/services/%s/tsdb", serviceName)).String(),
		bytes.NewReader(requestJSON),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	defer c.CloseReponse(resp)
	if err != nil {
		return err
	}

	return nil
}

// FetchLatestMetricValues fetch latest metrics
func (c *Client) FetchLatestMetricValues(hostIDs []string, metricNames []string) (LatestMetricValues, error) {
	v := url.Values{}
	for _, hostID := range hostIDs {
		v.Add("hostId", hostID)
	}
	for _, metricName := range metricNames {
		v.Add("name", metricName)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", c.urlFor("/api/v0/tsdb/latest").String(), v.Encode()), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer c.CloseReponse(resp)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		LatestMetricValues *LatestMetricValues `json:"tsdbLatest"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return *(data.LatestMetricValues), err
}
