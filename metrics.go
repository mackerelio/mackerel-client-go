package mackerel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type MetricValue struct {
	Name  string      `json:"name,omitempty"`
	Time  int64       `json:"time,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type HostMetricValue struct {
	HostID string      `json:"hostID,omitempty"`
	Name   string      `json:"name,omitempty"`
	Time   int64       `json:"time,omitempty"`
	Value  interface{} `json:"value,omitempty"`
}

type LatestMetricValues map[string]map[string]*MetricValue

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
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API result failed: %s", resp.Status)
	}

	return nil
}

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
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API result failed: %s", resp.Status)
	}

	return nil
}

func (c *Client) FetchLatestMetricValues(hostIDs []string, metricNames []string) (LatestMetricValues, error) {
	v := url.Values{}
	for _, hostID := range hostIDs {
		v.Add("hostID", hostID)
	}
	for _, metricName := range metricNames {
		v.Add("name", metricName)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", c.urlFor("/api/v0/tsdb/latest").String(), v.Encode()), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is not 200")
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
