package mackerel

import (
	"encoding/json"
	"fmt"
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
	HostID string `json:"hostID,omitempty"`
	*MetricValue
}

// LatestMetricValues latest metric value
type LatestMetricValues map[string]map[string]*MetricValue

// PostHostMetricValues post host metrics
func (c *Client) PostHostMetricValues(metricValues [](*HostMetricValue)) error {
	resp, err := c.PostJSON("/api/v0/tsdb", metricValues)
	defer closeResp(resp)
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
	defer closeResp(resp)
	return err
}

// FetchLatestMetricValues fetch latest metrics
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
	defer closeResp(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is not 200")
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
