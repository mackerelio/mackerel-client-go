package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AWSIntegration aws integration information
type AWSIntegration struct {
	ID           string                            `json:"id"`
	Name         string                            `json:"name"`
	Memo         string                            `json:"memo"`
	Key          string                            `json:"key,omitempty"`
	RoleArn      string                            `json:"roleArn,omitempty"`
	ExternalID   string                            `json:"externalId,omitempty"`
	Region       string                            `json:"region"`
	IncludedTags string                            `json:"includedTags"`
	ExcludedTags string                            `json:"excludedTags"`
	Services     map[string]*AWSIntegrationService `json:"services"`
}

// AWSIntegrationService integration settings for each AWS service
type AWSIntegrationService struct {
	Enable              bool     `json:"enable"`
	Role                *string  `json:"role"`
	ExcludedMetrics     []string `json:"excludedMetrics"`
	RetireAutomatically bool     `json:"retireAutomatically,omitempty"`
}

// CreateAWSIntegrationParam  parameters for CreateAWSIntegration
type CreateAWSIntegrationParam struct {
	Name         string                            `json:"name"`
	Memo         string                            `json:"memo"`
	Key          string                            `json:"key,omitempty"`
	SecretKey    string                            `json:"secretKey,omitempty"`
	RoleArn      string                            `json:"roleArn,omitempty"`
	ExternalID   string                            `json:"externalId,omitempty"`
	Region       string                            `json:"region"`
	IncludedTags string                            `json:"includedTags"`
	ExcludedTags string                            `json:"excludedTags"`
	Services     map[string]*AWSIntegrationService `json:"services"`
}

// UpdateAWSIntegrationParam parameters for UpdateAwsIntegration
type UpdateAWSIntegrationParam CreateAWSIntegrationParam

// ListAWSIntegrationExcludableMetrics List of excludeable metric names for aws integration
type ListAWSIntegrationExcludableMetrics map[string][]string

// FindAWSIntegrations finds AWS Integration Settings
func (c *Client) FindAWSIntegrations() ([]*AWSIntegration, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/aws-integrations").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		AWSIntegrations []*AWSIntegration `json:"aws_integrations"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.AWSIntegrations, err
}

// FindAWSIntegration finds AWS Integration Setting
func (c *Client) FindAWSIntegration(awsIntegrationID string) (*AWSIntegration, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/api/v0/aws-integrations/%s", awsIntegrationID)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var awsIntegration *AWSIntegration
	err = json.NewDecoder(resp.Body).Decode(&awsIntegration)
	if err != nil {
		return nil, err
	}
	return awsIntegration, err
}

// CreateAWSIntegration creates AWS Integration Setting
func (c *Client) CreateAWSIntegration(param *CreateAWSIntegrationParam) (*AWSIntegration, error) {
	resp, err := c.PostJSON("/api/v0/aws-integrations", param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var awsIntegration *AWSIntegration
	err = json.NewDecoder(resp.Body).Decode(&awsIntegration)
	if err != nil {
		return nil, err
	}
	return awsIntegration, err
}

// UpdateAWSIntegration updates AWS Integration Setting
func (c *Client) UpdateAWSIntegration(awsIntegrationID string, param *UpdateAWSIntegrationParam) (*AWSIntegration, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/api/v0/aws-integrations/%s", awsIntegrationID), param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var awsIntegration *AWSIntegration
	err = json.NewDecoder(resp.Body).Decode(&awsIntegration)
	if err != nil {
		return nil, err
	}
	return awsIntegration, err
}

// DeleteAWSIntegration deletes AWS Integration Setting
func (c *Client) DeleteAWSIntegration(awsIntegrationID string) (*AWSIntegration, error) {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/aws-integrations/%s", awsIntegrationID)).String(),
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

	var awsIntegration *AWSIntegration
	err = json.NewDecoder(resp.Body).Decode(&awsIntegration)
	if err != nil {
		return nil, err
	}
	return awsIntegration, err
}

// CreateAWSIntegrationExternalID creates AWS Integration External ID
func (c *Client) CreateAWSIntegrationExternalID() (string, error) {
	resp, err := c.PostJSON("/api/v0/aws-integrations-external-id", nil)
	defer closeResponse(resp)
	if err != nil {
		return "", err
	}

	var data struct {
		ExternalID string `json:"externalId"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.ExternalID, nil
}

// ListAWSIntegrationExcludableMetrics lists excludable metrics for AWS Integration
func (c *Client) ListAWSIntegrationExcludableMetrics() (*ListAWSIntegrationExcludableMetrics, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/aws-integrations-excludable-metrics").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var listAWSIntegrationExcludableMetrics *ListAWSIntegrationExcludableMetrics
	err = json.NewDecoder(resp.Body).Decode(&listAWSIntegrationExcludableMetrics)
	if err != nil {
		return nil, err
	}
	return listAWSIntegrationExcludableMetrics, err
}
