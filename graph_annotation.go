package mackerel

import (
	"fmt"
	"net/url"
	"strconv"
)

// GraphAnnotation represents parameters to post a graph annotation.
type GraphAnnotation struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	From        int64    `json:"from,omitempty"`
	To          int64    `json:"to,omitempty"`
	Service     string   `json:"service,omitempty"`
	Roles       []string `json:"roles,omitempty"`
}

// FindGraphAnnotations fetches graph annotations.
func (c *Client) FindGraphAnnotations(service string, from int64, to int64) ([]*GraphAnnotation, error) {
	params := url.Values{}
	params.Add("service", service)
	params.Add("from", strconv.FormatInt(from, 10))
	params.Add("to", strconv.FormatInt(to, 10))

	data, err := requestGetWithParams[struct {
		GraphAnnotations []*GraphAnnotation `json:"graphAnnotations"`
	}](c, "/api/v0/graph-annotations", params)
	if err != nil {
		return nil, err
	}
	return data.GraphAnnotations, nil
}

// CreateGraphAnnotation creates a graph annotation.
func (c *Client) CreateGraphAnnotation(annotation *GraphAnnotation) (*GraphAnnotation, error) {
	return requestPost[GraphAnnotation](c, "/api/v0/graph-annotations", annotation)
}

// UpdateGraphAnnotation updates a graph annotation.
func (c *Client) UpdateGraphAnnotation(annotationID string, annotation *GraphAnnotation) (*GraphAnnotation, error) {
	path := fmt.Sprintf("/api/v0/graph-annotations/%s", annotationID)
	return requestPut[GraphAnnotation](c, path, annotation)
}

// DeleteGraphAnnotation deletes a graph annotation.
func (c *Client) DeleteGraphAnnotation(annotationID string) (*GraphAnnotation, error) {
	path := fmt.Sprintf("/api/v0/graph-annotations/%s", annotationID)
	return requestDelete[GraphAnnotation](c, path)
}
