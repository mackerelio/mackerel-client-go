package mackerel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFindAWSIntegrations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/aws-integrations" {
			t.Error("request URL should be /api/v0/aws-integrations but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"aws_integrations": {
				{
					"id":           "9rxGOHfVF8F",
					"name":         "my-aws-integrations-1",
					"memo":         "my-aws-integrations-1",
					"key":          "ADIUHIBCGY6VXDMAUAA5E",
					"roleArn":      "",
					"externalId":   "",
					"region":       "ap-northeast-1",
					"includedTags": "Name:web-server,Environment:staging,Product:web",
					"excludedTags": "Name:test-server,Environment:staging,Product:test",
					"services": map[string]map[string]interface{}{
						"EC2": {
							"enable":              true,
							"role":                "web-group",
							"excludedMetrics":     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
							"retireAutomatically": true,
						},
						"ALB": {
							"enable":          true,
							"role":            "web-group",
							"excludedMetrics": []string{"alb.request.count"},
						},
						"RDS": {
							"enable":          true,
							"role":            "db-group",
							"excludedMetrics": []string{"rds.cpu.used", "rds.aurora.row_lock_time.row_lock"},
						},
					},
				},
				{
					"id":           "9rxGOHfb12F",
					"name":         "my-aws-integrations-2",
					"memo":         "",
					"key":          "",
					"roleArn":      "arn:aws:iam::111111111111:role/MackerelIntegrationRole",
					"externalId":   "PyrtkY42H8poFvRBU42dNL12BIPd9dF9QaCe1pgoXK2",
					"region":       "eu-central-1",
					"includedTags": "Name:web-server,Environment:staging,Product:web",
					"excludedTags": "Name:test-server,Environment:staging,Product:test",
					"services": map[string]map[string]interface{}{
						"EC2": {
							"enable":              false,
							"role":                (*string)(nil),
							"excludedMetrics":     []string{""},
							"retireAutomatically": false,
						},
						"ALB": {
							"enable":          false,
							"role":            (*string)(nil),
							"excludedMetrics": []string{""},
						},
						"RDS": {
							"enable":          false,
							"role":            (*string)(nil),
							"excludedMetrics": []string{""},
						},
					},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	awsIntegrations, err := client.ListAWSIntegrations()
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if awsIntegrations[0].ID != "9rxGOHfVF8F" {
		t.Error("aws integrations id should be empty but: ", awsIntegrations[0].ID)
	}

	if reflect.DeepEqual(awsIntegrations[0].Services["EC2"], &AWSIntegrationService{
		Enable:              true,
		Role:                toPointer("web-group"),
		ExcludedMetrics:     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
		RetireAutomatically: true,
	}) != true {
		t.Errorf("Wrong data for aws integrations services: %v", awsIntegrations[0].Services["EC2"])
	}
}

func TestFindAWSIntegration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		url := fmt.Sprintf("/api/v0/aws-integrations/%s", "9rxGOHfVF8F")
		if req.URL.Path != url {
			t.Error("request URL should be /api/v0/aws-integrations/<ID> but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":           "9rxGOHfVF8F",
			"name":         "my-aws-integrations-1",
			"memo":         "my-aws-integrations-1",
			"key":          "ADIUHIBCGY6VXDMAUAA5E",
			"roleArn":      "",
			"externalId":   "",
			"region":       "ap-northeast-1",
			"includedTags": "Name:web-server,Environment:staging,Product:web",
			"excludedTags": "Name:test-server,Environment:staging,Product:test",
			"services": map[string]map[string]interface{}{
				"EC2": {
					"enable":              true,
					"role":                "web-group",
					"excludedMetrics":     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
					"retireAutomatically": true,
				},
				"ALB": {
					"enable":          true,
					"role":            "web-group",
					"excludedMetrics": []string{"alb.request.count"},
				},
				"RDS": {
					"enable":          true,
					"role":            "db-group",
					"excludedMetrics": []string{"rds.cpu.used", "rds.aurora.row_lock_time.row_lock"},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	awsIntegration, err := client.FindAWSIntegration("9rxGOHfVF8F")
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if awsIntegration.ID != "9rxGOHfVF8F" {
		t.Error("aws integrations id should be empty but: ", awsIntegration.ID)
	}

	if reflect.DeepEqual(awsIntegration.Services["EC2"], &AWSIntegrationService{
		Enable:              true,
		Role:                toPointer("web-group"),
		ExcludedMetrics:     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
		RetireAutomatically: true,
	}) != true {
		t.Errorf("Wrong data for aws integration services: %v", awsIntegration.Services["EC2"])
	}
}

func TestCreateAWSIntegration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/aws-integrations" {
			t.Error("request URL should be /api/v0/aws-integrations but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var awsIntegration *AWSIntegration
		err := json.Unmarshal(body, &awsIntegration)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		if awsIntegration.Name != "my-aws-integrations-1" {
			t.Error("request sends json including name but: ", awsIntegration.Name)
		}

		if reflect.DeepEqual(awsIntegration.Services["EC2"], &AWSIntegrationService{
			Enable:              true,
			Role:                toPointer("web-group"),
			ExcludedMetrics:     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
			RetireAutomatically: true,
		}) != true {
			t.Errorf("Wrong data for aws integration services: %v", awsIntegration.Services["EC2"])
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":           "9rxGOHfVF8F",
			"name":         "my-aws-integrations-1",
			"memo":         "my-aws-integrations-1",
			"key":          "",
			"secretKey":    "",
			"roleArn":      "arn:aws:iam::111111111111:role/MackerelIntegrationRole",
			"externalId":   "PyrtkY42H8poFvRBU42dNL12BIPd9dF9QaCe1pgoXK2",
			"region":       "ap-northeast-1",
			"includedTags": "Name:web-server,Environment:staging,Product:web",
			"excludedTags": "Name:test-server,Environment:staging,Product:test",
			"services": map[string]map[string]interface{}{
				"EC2": {
					"enable":              true,
					"role":                "web-group",
					"excludedMetrics":     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
					"retireAutomatically": true,
				},
				"ALB": {
					"enable":          true,
					"role":            "web-group",
					"excludedMetrics": []string{"alb.request.count"},
				},
				"RDS": {
					"enable":          true,
					"role":            "db-group",
					"excludedMetrics": []string{"rds.cpu.used", "rds.aurora.row_lock_time.row_lock"},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	awsIntegrations, err := client.CreateAWSIntegration(&CreateAWSIntegrationParam{
		Name:         "my-aws-integrations-1",
		Memo:         "my-aws-integrations-1",
		Key:          "",
		SecretKey:    "",
		RoleArn:      "arn:aws:iam::111111111111:role/MackerelIntegrationRole",
		ExternalID:   "PyrtkY42H8poFvRBU42dNL12BIPd9dF9QaCe1pgoXK2",
		Region:       "ap-northeast-1",
		IncludedTags: "Name:web-server,Environment:staging,Product:web",
		ExcludedTags: "Name:test-server,Environment:staging,Product:test",
		Services: map[string]*AWSIntegrationService{
			"EC2": {
				Enable:              true,
				Role:                toPointer("web-group"),
				ExcludedMetrics:     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
				RetireAutomatically: true,
			},
			"ALB": {
				Enable:          true,
				Role:            toPointer("web-group"),
				ExcludedMetrics: []string{"alb.request.count"},
			},
			"RDS": {
				Enable:          true,
				Role:            toPointer("db-group"),
				ExcludedMetrics: []string{"rds.cpu.used", "rds.aurora.row_lock_time.row_lock"},
			},
		},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if awsIntegrations.ID != "9rxGOHfVF8F" {
		t.Error("aws integrations id should be empty but ", awsIntegrations.ID)
	}

	if reflect.DeepEqual(awsIntegrations.Services["EC2"], &AWSIntegrationService{
		Enable:              true,
		Role:                toPointer("web-group"),
		ExcludedMetrics:     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
		RetireAutomatically: true,
	}) != true {
		t.Errorf("Wrong data for aws integration services: %v", awsIntegrations.Services["EC2"])
	}
}

func TestUpdateAWSIntegration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/aws-integrations/%s", "9rxGOHfVF8F") {
			t.Error("request URL should be /api/v0/aws-integrations/<ID> but: ", req.URL.Path)
		}

		if req.Method != "PUT" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)

		var awsIntegration *AWSIntegration
		err := json.Unmarshal(body, &awsIntegration)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		if awsIntegration.Name != "my-aws-integrations-1" {
			t.Error("request sends json including name but: ", awsIntegration.Name)
		}

		if reflect.DeepEqual(awsIntegration.Services["EC2"], &AWSIntegrationService{
			Enable:              true,
			Role:                toPointer("web-group"),
			ExcludedMetrics:     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
			RetireAutomatically: true,
		}) != true {
			t.Errorf("Wrong data for aws integration services: %v", awsIntegration.Services["EC2"])
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":           "9rxGOHfVF8F",
			"name":         "my-aws-integrations-1",
			"memo":         "my-aws-integrations-1",
			"key":          "",
			"secretKey":    "",
			"roleArn":      "arn:aws:iam::111111111111:role/MackerelIntegrationRole",
			"externalId":   "PyrtkY42H8poFvRBU42dNL12BIPd9dF9QaCe1pgoXK2",
			"region":       "ap-northeast-1",
			"includedTags": "Name:web-server,Environment:staging,Product:web",
			"excludedTags": "Name:test-server,Environment:staging,Product:test",
			"services": map[string]map[string]interface{}{
				"EC2": {
					"enable":              true,
					"role":                "web-group",
					"excludedMetrics":     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
					"retireAutomatically": true,
				},
				"ALB": {
					"enable":          true,
					"role":            "web-group",
					"excludedMetrics": []string{"alb.request.count"},
				},
				"RDS": {
					"enable":          true,
					"role":            "db-group",
					"excludedMetrics": []string{"rds.cpu.used", "rds.aurora.row_lock_time.row_lock"},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	awsIntegrations, err := client.UpdateAWSIntegration("9rxGOHfVF8F", &UpdateAWSIntegrationParam{
		Name:         "my-aws-integrations-1",
		Memo:         "my-aws-integrations-1",
		Key:          "",
		SecretKey:    "",
		RoleArn:      "arn:aws:iam::111111111111:role/MackerelIntegrationRole",
		ExternalID:   "PyrtkY42H8poFvRBU42dNL12BIPd9dF9QaCe1pgoXK2",
		Region:       "ap-northeast-1",
		IncludedTags: "Name:web-server,Environment:staging,Product:web",
		ExcludedTags: "Name:test-server,Environment:staging,Product:test",
		Services: map[string]*AWSIntegrationService{
			"EC2": {
				Enable:              true,
				Role:                toPointer("web-group"),
				ExcludedMetrics:     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
				RetireAutomatically: true,
			},
			"ALB": {
				Enable:          true,
				Role:            toPointer("web-group"),
				ExcludedMetrics: []string{"alb.request.count"},
			},
			"RDS": {
				Enable:          true,
				Role:            toPointer("db-group"),
				ExcludedMetrics: []string{"rds.cpu.used", "rds.aurora.row_lock_time.row_lock"},
			},
		},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if awsIntegrations.ID != "9rxGOHfVF8F" {
		t.Error("aws integrations id should be empty but: ", awsIntegrations.ID)
	}

	if reflect.DeepEqual(awsIntegrations.Services["EC2"], &AWSIntegrationService{
		Enable:              true,
		Role:                toPointer("web-group"),
		ExcludedMetrics:     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
		RetireAutomatically: true,
	}) != true {
		t.Errorf("Wrong data for aws integration services: %v", awsIntegrations.Services["EC2"])
	}
}

func TestDeleteAWSIntegration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/aws-integrations/%s", "9rxGOHfVF8F") {
			t.Error("request URL should be /api/v0/aws-integrations/<ID> but: ", req.URL.Path)
		}

		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":           "9rxGOHfVF8F",
			"name":         "my-aws-integrations-1",
			"memo":         "my-aws-integrations-1",
			"key":          "",
			"roleArn":      "arn:aws:iam::111111111111:role/MackerelIntegrationRole",
			"externalId":   "PyrtkY42H8poFvRBU42dNL12BIPd9dF9QaCe1pgoXK2",
			"region":       "ap-northeast-1",
			"includedTags": "Name:web-server,Environment:staging,Product:web",
			"excludedTags": "Name:test-server,Environment:staging,Product:test",
			"services": map[string]map[string]interface{}{
				"EC2": {
					"enable":              true,
					"role":                "web-group",
					"excludedMetrics":     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
					"retireAutomatically": true,
				},
				"ALB": {
					"enable":          true,
					"role":            "web-group",
					"excludedMetrics": []string{"alb.request.count"},
				},
				"RDS": {
					"enable":          true,
					"role":            "db-group",
					"excludedMetrics": []string{"rds.cpu.used", "rds.aurora.row_lock_time.row_lock"},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	awsIntegrations, err := client.DeleteAWSIntegration("9rxGOHfVF8F")
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if awsIntegrations.ID != "9rxGOHfVF8F" {
		t.Error("aws integrations id should be empty but: ", awsIntegrations.ID)
	}

	if reflect.DeepEqual(awsIntegrations.Services["EC2"], &AWSIntegrationService{
		Enable:              true,
		Role:                toPointer("web-group"),
		ExcludedMetrics:     []string{"ec2.cpu.used", "ec2.network.in", "ec2.network.out"},
		RetireAutomatically: true,
	}) != true {
		t.Errorf("Wrong data for aws integration services: %v", awsIntegrations.Services["EC2"])
	}
}

func TestCreateAWSIntegrationExternalID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/aws-integrations-external-id" {
			t.Error("request URL should be /api/v0/aws-integrations-external-id but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]string{
			"externalId": "PyrtkY42H8poFvRBU42dNL12BIPd9dF9QaCe1pgoXK2",
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	awsIntegrationExternalID, err := client.CreateAWSIntegrationExternalID()
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if awsIntegrationExternalID != "PyrtkY42H8poFvRBU42dNL12BIPd9dF9QaCe1pgoXK2" {
		t.Error("aws integration external id should be empty but: ", awsIntegrationExternalID)
	}
}

func TestListAWSIntegrationExcludableMetrics(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/aws-integrations-excludable-metrics" {
			t.Error("request URL should be /api/v0/aws-integrations-excludable-metrics but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]string{
			"EC2": {
				"ec2.cpu.used",
				"ec2.cpu_credit.used",
				"ec2.cpu_credit.balance",
			},
			"ELB": {
				"elb.count.request_count",
				"elb.host_count.healthy",
				"elb.host_count.unhealthy",
			},
			"ALB": {
				"alb.request.count",
				"alb.bytes.processed",
				"alb.httpcode_count.target_2xx",
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	listAWSIntegrationExcludableMetrics, err := client.ListAWSIntegrationExcludableMetrics()
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if reflect.DeepEqual((*listAWSIntegrationExcludableMetrics)["EC2"], []string{"ec2.cpu.used", "ec2.cpu_credit.used", "ec2.cpu_credit.balance"}) != true {
		t.Errorf("Wrong data for list of excludeable metric names for aws integration: %v", (*listAWSIntegrationExcludableMetrics)["EC2"])
	}
}

func toPointer(s string) *string {
	return &s
}
