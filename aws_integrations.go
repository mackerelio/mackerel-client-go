package mackerel

import (
	"encoding/json"
	"fmt"
)

// AWSIntegration AWS integration information
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
	IncludedMetrics     []string `json:"includedMetrics"`
	ExcludedMetrics     []string `json:"excludedMetrics"`
	RetireAutomatically bool     `json:"retireAutomatically,omitempty"`
}

type awsIntegrationService = AWSIntegrationService

type awsIntegrationServiceWithIncludedMetrics struct {
	Enable              bool     `json:"enable"`
	Role                *string  `json:"role"`
	IncludedMetrics     []string `json:"includedMetrics"`
	RetireAutomatically bool     `json:"retireAutomatically,omitempty"`
}

type awsIntegrationServiceWithExcludedMetrics struct {
	Enable              bool     `json:"enable"`
	Role                *string  `json:"role"`
	ExcludedMetrics     []string `json:"excludedMetrics"`
	RetireAutomatically bool     `json:"retireAutomatically,omitempty"`
}

// MarshalJSON implements json.Marshaler
func (a *AWSIntegrationService) MarshalJSON() ([]byte, error) {
	// AWS integration create/update APIs only accept either includedMetrics or excludedMetrics
	if a.ExcludedMetrics != nil && a.IncludedMetrics == nil {
		return json.Marshal(awsIntegrationServiceWithExcludedMetrics{
			Enable:              a.Enable,
			Role:                a.Role,
			ExcludedMetrics:     a.ExcludedMetrics,
			RetireAutomatically: a.RetireAutomatically,
		})
	}
	if a.ExcludedMetrics == nil && a.IncludedMetrics != nil {
		return json.Marshal(awsIntegrationServiceWithIncludedMetrics{
			Enable:              a.Enable,
			Role:                a.Role,
			IncludedMetrics:     a.IncludedMetrics,
			RetireAutomatically: a.RetireAutomatically,
		})
	}
	return json.Marshal(awsIntegrationService(*a))
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

// ListAWSIntegrationExcludableMetrics List of excludeable metric names for AWS integration
type ListAWSIntegrationExcludableMetrics map[string][]string

// FindAWSIntegrations finds AWS integration settings.
func (c *Client) FindAWSIntegrations() ([]*AWSIntegration, error) {
	data, err := requestGet[struct {
		AWSIntegrations []*AWSIntegration `json:"aws_integrations"`
	}](c, "/api/v0/aws-integrations")
	if err != nil {
		return nil, err
	}
	return data.AWSIntegrations, nil
}

// CreateAWSIntegration creates an AWS integration setting.
func (c *Client) CreateAWSIntegration(param *CreateAWSIntegrationParam) (*AWSIntegration, error) {
	return requestPost[AWSIntegration](c, "/api/v0/aws-integrations", param)
}

// FindAWSIntegration finds an AWS integration setting.
func (c *Client) FindAWSIntegration(awsIntegrationID string) (*AWSIntegration, error) {
	path := fmt.Sprintf("/api/v0/aws-integrations/%s", awsIntegrationID)
	return requestGet[AWSIntegration](c, path)
}

// UpdateAWSIntegration updates an AWS integration setting.
func (c *Client) UpdateAWSIntegration(awsIntegrationID string, param *UpdateAWSIntegrationParam) (*AWSIntegration, error) {
	path := fmt.Sprintf("/api/v0/aws-integrations/%s", awsIntegrationID)
	return requestPut[AWSIntegration](c, path, param)
}

// DeleteAWSIntegration deletes an AWS integration setting.
func (c *Client) DeleteAWSIntegration(awsIntegrationID string) (*AWSIntegration, error) {
	path := fmt.Sprintf("/api/v0/aws-integrations/%s", awsIntegrationID)
	return requestDelete[AWSIntegration](c, path)
}

// CreateAWSIntegrationExternalID creates an AWS integration External ID.
func (c *Client) CreateAWSIntegrationExternalID() (string, error) {
	data, err := requestPost[struct {
		ExternalID string `json:"externalId"`
	}](c, "/api/v0/aws-integrations-external-id", nil)
	if err != nil {
		return "", err
	}
	return data.ExternalID, nil
}

// ListAWSIntegrationExcludableMetrics lists excludable metrics for AWS integration.
func (c *Client) ListAWSIntegrationExcludableMetrics() (*ListAWSIntegrationExcludableMetrics, error) {
	return requestGet[ListAWSIntegrationExcludableMetrics](c, "/api/v0/aws-integrations-excludable-metrics")
}
