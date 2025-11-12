package mackerel

import (
	"net/url"
	"strconv"
)

// HTTPServerStats represents HTTP server statistics
type HTTPServerStats struct {
	Method              string  `json:"method"`
	Route               string  `json:"route"`
	TotalMillis         float64 `json:"totalMillis"`
	AverageMillis       float64 `json:"averageMillis"`
	ApproxP95Millis     float64 `json:"approxP95Millis"`
	ErrorRatePercentage float64 `json:"errorRatePercentage"`
	RequestCount        int64   `json:"requestCount"`
}

// HTTPServerStatsPageConnection represents a paginated response of HTTP server statistics
type HTTPServerStatsPageConnection struct {
	Results     []*HTTPServerStats `json:"results"`
	HasNextPage bool               `json:"hasNextPage"`
}

// ListHTTPServerStatsParam represents parameters for listing HTTP server statistics
type ListHTTPServerStatsParam struct {
	ServiceName      string
	From             int64
	To               int64
	ServiceNamespace *string
	Environment      *string
	Version          *string
	OrderColumn      *string
	OrderDirection   *string
	Method           *string
	Route            *string
	Page             *int
	PerPage          *int
}

// ListHTTPServerStats retrieves HTTP server statistics
func (c *Client) ListHTTPServerStats(param *ListHTTPServerStatsParam) (*HTTPServerStatsPageConnection, error) {
	params := url.Values{}
	params.Set("serviceName", param.ServiceName)
	params.Set("from", strconv.FormatInt(param.From, 10))
	params.Set("to", strconv.FormatInt(param.To, 10))

	if param.ServiceNamespace != nil {
		params.Set("serviceNamespace", *param.ServiceNamespace)
	}
	if param.Environment != nil {
		params.Set("environment", *param.Environment)
	}
	if param.Version != nil {
		params.Set("version", *param.Version)
	}
	if param.OrderColumn != nil {
		params.Set("orderColumn", *param.OrderColumn)
	}
	if param.OrderDirection != nil {
		params.Set("orderDirection", *param.OrderDirection)
	}
	if param.Method != nil {
		params.Set("method", *param.Method)
	}
	if param.Route != nil {
		params.Set("route", *param.Route)
	}
	if param.Page != nil {
		params.Set("page", strconv.Itoa(*param.Page))
	}
	if param.PerPage != nil {
		params.Set("perPage", strconv.Itoa(*param.PerPage))
	}

	return requestGetWithParams[HTTPServerStatsPageConnection](c, "/api/v0/apm/http-server-stats", params)
}
