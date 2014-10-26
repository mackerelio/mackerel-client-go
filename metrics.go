package mackerel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type MetricValue struct {
	Name  string      `json:"name,omitempty"`
	Time  float64     `json:"time,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type HostMetricValue struct {
	HostId string      `json:"hostId",omitempty`
	Name   string      `json:"name,omitempty"`
	Time   float64     `json:"time,omitempty"`
	Value  interface{} `json:"value,omitempty"`
}

type ServiceMetricValue MetricValue

type LatestMetricValues map[string]map[string]*MetricValue

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

func (c *Client) FetchLatestMetricValues(hostIds []string, metricNames []string) (LatestMetricValues, error) {
	v := url.Values{}
	for _, hostId := range hostIds {
		v.Add("hostId", hostId)
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
		return nil, errors.New("status code is not 200")
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
