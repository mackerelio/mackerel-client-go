package mackerel

import (
	"context"
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
	return c.FindGraphAnnotationsContext(context.Background(), service, from, to)
}

// FindGraphAnnotationsContext fetches graph annotations.
func (c *Client) FindGraphAnnotationsContext(ctx context.Context, service string, from int64, to int64) ([]*GraphAnnotation, error) {
	params := url.Values{}
	params.Add("service", service)
	params.Add("from", strconv.FormatInt(from, 10))
	params.Add("to", strconv.FormatInt(to, 10))

	data, err := requestGetWithParamsContext[struct {
		GraphAnnotations []*GraphAnnotation `json:"graphAnnotations"`
	}](ctx, c, "/api/v0/graph-annotations", params)
	if err != nil {
		return nil, err
	}
	return data.GraphAnnotations, nil
}

// CreateGraphAnnotation creates a graph annotation.
func (c *Client) CreateGraphAnnotation(annotation *GraphAnnotation) (*GraphAnnotation, error) {
	return c.CreateGraphAnnotationContext(context.Background(), annotation)
}

// CreateGraphAnnotationContext creates a graph annotation.
func (c *Client) CreateGraphAnnotationContext(ctx context.Context, annotation *GraphAnnotation) (*GraphAnnotation, error) {
	return requestPostContext[GraphAnnotation](ctx, c, "/api/v0/graph-annotations", annotation)
}

// UpdateGraphAnnotation updates a graph annotation.
func (c *Client) UpdateGraphAnnotation(annotationID string, annotation *GraphAnnotation) (*GraphAnnotation, error) {
	return c.UpdateGraphAnnotationContext(context.Background(), annotationID, annotation)
}

// UpdateGraphAnnotationContext updates a graph annotation.
func (c *Client) UpdateGraphAnnotationContext(ctx context.Context, annotationID string, annotation *GraphAnnotation) (*GraphAnnotation, error) {
	path := fmt.Sprintf("/api/v0/graph-annotations/%s", annotationID)
	return requestPutWithContext[GraphAnnotation](ctx, c, path, annotation)
}

// DeleteGraphAnnotation deletes a graph annotation.
func (c *Client) DeleteGraphAnnotation(annotationID string) (*GraphAnnotation, error) {
	return c.DeleteGraphAnnotationContext(context.Background(), annotationID)
}

// DeleteGraphAnnotationContext deletes a graph annotation.
func (c *Client) DeleteGraphAnnotationContext(ctx context.Context, annotationID string) (*GraphAnnotation, error) {
	path := fmt.Sprintf("/api/v0/graph-annotations/%s", annotationID)
	return requestDeleteContext[GraphAnnotation](ctx, c, path)
}
