package mackerel

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// MetricValue metric value
type MetricValue struct {
	Name  string `json:"name,omitempty"`
	Time  int64  `json:"time,omitempty"`
	Value any    `json:"value,omitempty"`
}

// HostMetricValue host metric value
type HostMetricValue struct {
	HostID string `json:"hostId,omitempty"`
	*MetricValue
}

// LatestMetricValues latest metric value
type LatestMetricValues map[string]map[string]*MetricValue

// PostHostMetricValues post host metrics
func (c *Client) PostHostMetricValues(metricValues []*HostMetricValue) error {
	return c.PostHostMetricValuesContext(context.Background(), metricValues)
}

// PostHostMetricValuesContext post host metrics
func (c *Client) PostHostMetricValuesContext(ctx context.Context, metricValues []*HostMetricValue) error {
	_, err := requestPostContext[any](ctx, c, "/api/v0/tsdb", metricValues)
	return err
}

// PostHostMetricValuesByHostID post host metrics
func (c *Client) PostHostMetricValuesByHostID(hostID string, metricValues []*MetricValue) error {
	return c.PostHostMetricValuesByHostIDContext(context.Background(), hostID, metricValues)
}

// PostHostMetricValuesByHostIDContext post host metrics
func (c *Client) PostHostMetricValuesByHostIDContext(ctx context.Context, hostID string, metricValues []*MetricValue) error {
	var hostMetricValues []*HostMetricValue
	for _, metricValue := range metricValues {
		hostMetricValues = append(hostMetricValues, &HostMetricValue{
			HostID:      hostID,
			MetricValue: metricValue,
		})
	}
	return c.PostHostMetricValuesContext(ctx, hostMetricValues)
}

// PostServiceMetricValues posts service metrics.
func (c *Client) PostServiceMetricValues(serviceName string, metricValues []*MetricValue) error {
	return c.PostServiceMetricValuesContext(context.Background(), serviceName, metricValues)
}

// PostServiceMetricValuesContext posts service metrics.
func (c *Client) PostServiceMetricValuesContext(ctx context.Context, serviceName string, metricValues []*MetricValue) error {
	path := fmt.Sprintf("/api/v0/services/%s/tsdb", serviceName)
	_, err := requestPostContext[any](ctx, c, path, metricValues)
	return err
}

// FetchLatestMetricValues fetches latest metrics.
func (c *Client) FetchLatestMetricValues(hostIDs []string, metricNames []string) (LatestMetricValues, error) {
	return c.FetchLatestMetricValuesContext(context.Background(), hostIDs, metricNames)
}

// FetchLatestMetricValuesContext fetches latest metrics.
func (c *Client) FetchLatestMetricValuesContext(ctx context.Context, hostIDs []string, metricNames []string) (LatestMetricValues, error) {
	params := url.Values{}
	for _, hostID := range hostIDs {
		params.Add("hostId", hostID)
	}
	for _, metricName := range metricNames {
		params.Add("name", metricName)
	}

	data, err := requestGetWithParamsContext[struct {
		LatestMetricValues LatestMetricValues `json:"tsdbLatest"`
	}](ctx, c, "/api/v0/tsdb/latest", params)
	if err != nil {
		return nil, err
	}
	return data.LatestMetricValues, nil
}

// FetchHostMetricValues fetches the metric values for a host.
func (c *Client) FetchHostMetricValues(hostID string, metricName string, from int64, to int64) ([]MetricValue, error) {
	return c.fetchMetricValues(context.Background(), hostID, "", metricName, from, to)
}

// FetchHostMetricValuesContext fetches the metric values for a host.
func (c *Client) FetchHostMetricValuesContext(ctx context.Context, hostID string, metricName string, from int64, to int64) ([]MetricValue, error) {
	return c.fetchMetricValues(ctx, hostID, "", metricName, from, to)
}

// FetchServiceMetricValues fetches the metric values for a service.
func (c *Client) FetchServiceMetricValues(serviceName string, metricName string, from int64, to int64) ([]MetricValue, error) {
	return c.fetchMetricValues(context.Background(), "", serviceName, metricName, from, to)
}

// FetchServiceMetricValuesContext fetches the metric values for a service.
func (c *Client) FetchServiceMetricValuesContext(ctx context.Context, serviceName string, metricName string, from int64, to int64) ([]MetricValue, error) {
	return c.fetchMetricValues(ctx, "", serviceName, metricName, from, to)
}

func (c *Client) fetchMetricValues(ctx context.Context, hostID string, serviceName string, metricName string, from int64, to int64) ([]MetricValue, error) {
	params := url.Values{}
	params.Add("name", metricName)
	params.Add("from", strconv.FormatInt(from, 10))
	params.Add("to", strconv.FormatInt(to, 10))

	path := ""
	if hostID != "" {
		path = fmt.Sprintf("/api/v0/hosts/%s/metrics", hostID)
	} else if serviceName != "" {
		path = fmt.Sprintf("/api/v0/services/%s/metrics", serviceName)
	} else {
		return nil, errors.New("specify either host or service")
	}

	data, err := requestGetWithParamsContext[struct {
		MetricValues []MetricValue `json:"metrics"`
	}](ctx, c, path, params)
	if err != nil {
		return nil, err
	}
	return data.MetricValues, nil
}
