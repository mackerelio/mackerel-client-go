package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GraphAnnotation represents parameters to post graph annotation.
type GraphAnnotation struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	From        int64    `json:"from,omitempty"`
	To          int64    `json:"to,omitempty"`
	Service     string   `json:"service,omitempty"`
	Roles       []string `json:"roles,omitempty"`
}

// CreateGraphAnnotation creates graph annotation.
func (c *Client) CreateGraphAnnotation(annotation *GraphAnnotation) (*GraphAnnotation, error) {
	resp, err := c.PostJSON("/api/v0/graph-annotations", annotation)
	defer closeResponse(resp)

	if err != nil {
		return nil, err
	}

	var anno GraphAnnotation
	err = json.NewDecoder(resp.Body).Decode(&anno)
	if err != nil {
		return nil, err
	}
	return &anno, nil
}

// FindGraphAnnotations fetches graph annotation.
func (c *Client) FindGraphAnnotations(service string, from int64, to int64) ([]GraphAnnotation, error) {
	v := url.Values{}
	v.Add("service", service)
	v.Add("from", strconv.FormatInt(from, 10))
	v.Add("to", strconv.FormatInt(to, 10))

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", c.urlFor("/api/v0/graph-annotations").String(), v.Encode()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		GraphAnnotations []GraphAnnotation `json:"graphAnnotations"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.GraphAnnotations, err
}

// UpdateGraphAnnotation updates graph annotation.
func (c *Client) UpdateGraphAnnotation(annotationID string, annotation *GraphAnnotation) (*GraphAnnotation, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/api/v0/graph-annotations/%s", annotationID), annotation)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var updatedAnnotation = GraphAnnotation{}
	err = json.NewDecoder(resp.Body).Decode(&updatedAnnotation)
	if err != nil {
		return nil, err
	}

	return &updatedAnnotation, nil
}

// DeleteGraphAnnotation deletes graph annotation.
func (c *Client) DeleteGraphAnnotation(annotationID string) (*GraphAnnotation, error) {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/graph-annotations/%s", annotationID)).String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var annotation GraphAnnotation
	err = json.NewDecoder(resp.Body).Decode(&annotation)
	if err != nil {
		return nil, err
	}
	return &annotation, nil
}
