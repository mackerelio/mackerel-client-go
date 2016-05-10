package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// MetricValue metric value
type MetricValue struct {
	Name  string      `json:"name,omitempty"`
	Time  int64       `json:"time,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// HostMetricValue host metric value
type HostMetricValue struct {
	HostID string `json:"hostId,omitempty"`
	*MetricValue
}

// LatestMetricValues latest metric value
type LatestMetricValues map[string]map[string]*MetricValue

// PostHostMetricValues post host metrics
func (c *Client) PostHostMetricValues(metricValues [](*HostMetricValue)) error {
	resp, err := c.PostJSON("/api/v0/tsdb", metricValues)
	defer closeResponse(resp)
	return err
}

// PostHostMetricValuesByHostID post host metrics
func (c *Client) PostHostMetricValuesByHostID(hostID string, metricValues [](*MetricValue)) error {
	var hostMetricValues []*HostMetricValue
	for _, metricValue := range metricValues {
		hostMetricValues = append(hostMetricValues, &HostMetricValue{
			HostID:      hostID,
			MetricValue: metricValue,
		})
	}
	return c.PostHostMetricValues(hostMetricValues)
}

// PostServiceMetricValues post service metrics
func (c *Client) PostServiceMetricValues(serviceName string, metricValues [](*MetricValue)) error {
	resp, err := c.PostJSON(fmt.Sprintf("/api/v0/services/%s/tsdb", serviceName), metricValues)
	defer closeResponse(resp)
	return err
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
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		LatestMetricValues *LatestMetricValues `json:"tsdbLatest"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return *(data.LatestMetricValues), err
}

// FetchLatestHostMetricValues fetch latest metrics
func (c *Client) FetchLatestHostMetricValues(hostID string, metricName string, from int64, to int64) ([]MetricValue, error) {
	v := url.Values{}
	v.Add("name", metricName)
	v.Add("from", strconv.FormatInt(from, 10))
	v.Add("to", strconv.FormatInt(to, 10))

	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s?%s",
			c.urlFor(fmt.Sprintf("/api/v0/hosts/%s/metrics", hostID)).String(),
			v.Encode()),
		nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		MetricValues *[]MetricValue `json:"metrics"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return *(data.MetricValues), err
}
