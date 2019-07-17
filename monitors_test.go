package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func pfloat64(x float64) *float64 {
	return &x
}

func puint64(x uint64) *uint64 {
	return &x
}

func TestFindMonitors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/monitors" {
			t.Error("request URL should be /api/v0/monitors but: ", req.URL.Path)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"monitors": {
				{
					"id":            "2cSZzK3XfmG",
					"type":          "connectivity",
					"memo":          "connectivity monitor",
					"scopes":        []string{},
					"excludeScopes": []string{},
				},
				{
					"id":                              "2c5bLca8d",
					"type":                            "external",
					"name":                            "testMonitorExternal",
					"memo":                            "this monitor checks example.com.",
					"method":                          "GET",
					"url":                             "https://www.example.com/",
					"maxCheckAttempts":                3,
					"service":                         "someService",
					"notificationInterval":            60,
					"responseTimeCritical":            5000,
					"responseTimeWarning":             10000,
					"responseTimeDuration":            5,
					"certificationExpirationCritical": 15,
					"certificationExpirationWarning":  30,
					"containsString":                  "Foo Bar Baz",
					"skipCertificateVerification":     true,
					"headers": []map[string]interface{}{
						{"name": "Cache-Control", "value": "no-cache"},
					},
				},
				{
					"id":         "2DujfcR2kA9",
					"name":       "expression test",
					"memo":       "a monitor for expression",
					"type":       "expression",
					"expression": "avg(roleSlots('service:role','loadavg5'))",
					"operator":   ">",
					"warning":    20,
					"critical":   30,
				},
				{
					"id":                 "3CSsK3HKiHb",
					"type":               "anomalyDetection",
					"isMute":             false,
					"name":               "My first anomaly detection",
					"trainingPeriodFrom": 1561429260,
					"scopes": []string{
						"myService: myRole",
					},
					"maxCheckAttempts":   3,
					"warningSensitivity": "insensitive",
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	monitors, err := client.FindMonitors()

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	{
		m, ok := monitors[0].(*MonitorConnectivity)
		if !ok || m.Type != "connectivity" {
			t.Error("request sends json including type but: ", m)
		}
		if m.Memo != "connectivity monitor" {
			t.Error("request sends json including memo but: ", m)
		}
	}

	{
		m, ok := monitors[1].(*MonitorExternalHTTP)
		if !ok || m.Type != "external" {
			t.Error("request sends json including type but: ", m)
		}
		if m.Memo != "this monitor checks example.com." {
			t.Error("request sends json including memo but: ", m)
		}
		if m.Service != "someService" {
			t.Error("request sends json including service but: ", m)
		}
		if m.NotificationInterval != 60 {
			t.Error("request sends json including notificationInterval but: ", m)
		}

		if m.URL != "https://www.example.com/" {
			t.Error("request sends json including url but: ", m)
		}
		if m.MaxCheckAttempts != 3 {
			t.Error("request sends json including maxCheckAttempts but: ", m)
		}
		if *m.ResponseTimeCritical != 5000 {
			t.Error("request sends json including responseTimeCritical but: ", m)
		}

		if *m.ResponseTimeWarning != 10000 {
			t.Error("request sends json including responseTimeWarning but: ", m)
		}

		if *m.ResponseTimeDuration != 5 {
			t.Error("request sends json including responseTimeDuration but: ", m)
		}

		if *m.CertificationExpirationCritical != 15 {
			t.Error("request sends json including certificationExpirationCritical but: ", m)
		}

		if *m.CertificationExpirationWarning != 30 {
			t.Error("request sends json including certificationExpirationWarning but: ", m)
		}

		if m.ContainsString != "Foo Bar Baz" {
			t.Error("request sends json including containsString but: ", m)
		}

		if m.SkipCertificateVerification != true {
			t.Error("request sends json including skipCertificateVerification but: ", m)
		}

		if !reflect.DeepEqual(m.Headers, []HeaderField{{Name: "Cache-Control", Value: "no-cache"}}) {
			t.Error("request sends json including headers but: ", m)
		}
	}

	{
		m, ok := monitors[2].(*MonitorExpression)
		if !ok || m.Type != "expression" {
			t.Error("request sends json including expression but: ", monitors[2])
		}
		if m.Memo != "a monitor for expression" {
			t.Error("request sends json including memo but: ", m)
		}
	}
	{
		m, ok := monitors[3].(*MonitorAnomalyDetection)
		if !ok || m.Type != "anomalyDetection" {
			t.Error("request sends json including anomalyDetection but: ", monitors[3])
		}
		if m.TrainingPeriodFrom != 1561429260 {
			t.Error("request sends json including trainingPeriodFrom but: ", m)
		}
		if m.MaxCheckAttempts != 3 {
			t.Error("request sends json including maxCheckAttempts but: ", m)
		}
		if m.WarningSensitivity != "insensitive" {
			t.Error("request sends json including warningSensitivity but: ", m)
		}
	}
}

func TestGetMonitor(t *testing.T) {
	var (
		monitorID = "2cSZzK3XfmG"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		u := fmt.Sprintf("/api/v0/monitors/%s", monitorID)
		if req.URL.Path != u {
			t.Errorf("request URL should be %v but %v:", u, req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]map[string]interface{}{
			"monitor": {
				"id":            monitorID,
				"type":          "connectivity",
				"memo":          "connectivity monitor",
				"scopes":        []string{},
				"excludeScopes": []string{},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	monitor, err := client.GetMonitor(monitorID)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	m, ok := monitor.(*MonitorConnectivity)
	if !ok || m.Type != "connectivity" {
		t.Error("request sends json including type but: ", m)
	}
	if m.Memo != "connectivity monitor" {
		t.Error("request sends json including memo but: ", m)
	}
}

// ensure that it supports `"headers":[]` and headers must be nil by default.
func TestMonitorExternalHTTP_headers(t *testing.T) {
	tests := []struct {
		name string
		in   *MonitorExternalHTTP
		want string
	}{
		{
			name: "default",
			in:   &MonitorExternalHTTP{},
			want: `{"headers":null}`,
		},
		{
			name: "empty list",
			in:   &MonitorExternalHTTP{Headers: []HeaderField{}},
			want: `{"headers":[]}`,
		},
	}

	for _, tt := range tests {
		b, err := json.Marshal(tt.in)
		if err != nil {
			t.Error(err)
			continue
		}
		if got := string(b); got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

var testCases = []struct {
	title   string
	monitor Monitor
	json    string
}{
	{
		"connectivity",
		&MonitorConnectivity{
			ID:                   "2cSZzK3XfmA",
			Name:                 "",
			Type:                 "connectivity",
			IsMute:               false,
			NotificationInterval: 0,
			Scopes:               nil,
			ExcludeScopes:        nil,
		},
		`{
			"id": "2cSZzK3XfmA",
			"type": "connectivity"
		}`,
	},
	{
		"host metric monitor",
		&MonitorHostMetric{
			ID:                   "2cSZzK3XfmB",
			Name:                 "disk.aa-00.writes.delta",
			Type:                 "host",
			IsMute:               false,
			NotificationInterval: 0,
			Metric:               "disk.aa-00.writes.delta",
			Operator:             ">",
			Warning:              pfloat64(20000.000000),
			Critical:             pfloat64(400000.000000),
			Duration:             3,
			MaxCheckAttempts:     3,
			Scopes: []string{
				"Hatena-Blog",
			},
			ExcludeScopes: []string{
				"Hatena-Bookmark: db-master",
			},
		},
		`{
			"id": "2cSZzK3XfmB",
			"type": "host",
			"name": "disk.aa-00.writes.delta",
			"duration": 3,
			"metric": "disk.aa-00.writes.delta",
			"operator": ">",
			"warning": 20000,
			"critical": 400000,
			"maxCheckAttempts": 3,
			"scopes": [
			"Hatena-Blog"
			],
			"excludeScopes": [
			"Hatena-Bookmark: db-master"
			]
		}`,
	},
	{
		"host metric monitor without critical threshold",
		&MonitorHostMetric{
			ID:                   "2cSZzK3XfmF",
			Name:                 "Foo Bar",
			Type:                 "host",
			IsMute:               false,
			NotificationInterval: 0,
			Metric:               "custom.foo.bar",
			Operator:             ">",
			Warning:              pfloat64(200.0),
			Critical:             nil,
			Duration:             3,
			MaxCheckAttempts:     5,
			Scopes:               nil,
			ExcludeScopes:        nil,
		},
		`{
			"id": "2cSZzK3XfmF",
			"type": "host",
			"name": "Foo Bar",
			"duration": 3,
			"metric": "custom.foo.bar",
			"operator": ">",
			"warning": 200.0,
			"critical": null,
			"maxCheckAttempts": 5
		}`,
	},
	{
		"host metric monitor without warning threshold",
		&MonitorHostMetric{
			ID:                   "2cSZzK3XfmX",
			Name:                 "Foo Baz",
			Type:                 "host",
			IsMute:               false,
			NotificationInterval: 0,
			Metric:               "custom.foo.baz",
			Operator:             "<",
			Warning:              nil,
			Critical:             pfloat64(300),
			Duration:             7,
			MaxCheckAttempts:     2,
			Scopes:               nil,
			ExcludeScopes:        nil,
		},
		`{
			"id": "2cSZzK3XfmX",
			"type": "host",
			"name": "Foo Baz",
			"duration": 7,
			"metric": "custom.foo.baz",
			"operator": "<",
			"warning": null,
			"critical": 300.0,
			"maxCheckAttempts": 2
		}`,
	},
	{
		"service metric monitor",
		&MonitorServiceMetric{
			ID:                   "2cSZzK3XfmC",
			Name:                 "Hatena-Blog - access_num.4xx_count",
			Type:                 "service",
			IsMute:               false,
			NotificationInterval: 60,
			Service:              "Hatena-Blog",
			Metric:               "access_num.4xx_count",
			Operator:             ">",
			Warning:              pfloat64(50.000000),
			Critical:             pfloat64(100.000000),
			Duration:             1,
			MaxCheckAttempts:     5,
		},
		`{
			"id"  : "2cSZzK3XfmC",
			"type": "service",
			"name": "Hatena-Blog - access_num.4xx_count",
			"service": "Hatena-Blog",
			"duration": 1,
			"metric": "access_num.4xx_count",
			"operator": ">",
			"warning": 50.0,
			"critical": 100.0,
			"maxCheckAttempts": 5,
			"notificationInterval": 60
		}`,
	},
	{
		"service metric monitor without warning threshold",
		&MonitorServiceMetric{
			ID:                   "2cSZzK3XfmG",
			Name:                 "Hatena-Blog - access_num.5xx_count",
			Type:                 "service",
			IsMute:               false,
			NotificationInterval: 60,
			Service:              "Hatena-Blog",
			Metric:               "access_num.5xx_count",
			Operator:             ">",
			Warning:              nil,
			Critical:             pfloat64(0.0),
			Duration:             3,
			MaxCheckAttempts:     3,
		},
		`{
			"id"  : "2cSZzK3XfmG",
			"type": "service",
			"name": "Hatena-Blog - access_num.5xx_count",
			"service": "Hatena-Blog",
			"duration": 3,
			"metric": "access_num.5xx_count",
			"operator": ">",
			"critical": 0.0,
			"warning": null,
			"maxCheckAttempts": 3,
			"notificationInterval": 60
		}`,
	},
	{
		"service metric monitor with missing duration thresholds",
		&MonitorServiceMetric{
			ID:                      "2cSZzK3XfmG",
			Name:                    "Hatena-Blog - access_num.5xx_count",
			Type:                    "service",
			IsMute:                  false,
			NotificationInterval:    60,
			Service:                 "Hatena-Blog",
			Metric:                  "access_num.5xx_count",
			Operator:                ">",
			Warning:                 nil,
			MissingDurationWarning:  360,
			MissingDurationCritical: 720,
			Duration:                3,
			MaxCheckAttempts:        3,
		},
		`{
			"id"  : "2cSZzK3XfmG",
			"type": "service",
			"name": "Hatena-Blog - access_num.5xx_count",
			"service": "Hatena-Blog",
			"duration": 3,
			"metric": "access_num.5xx_count",
			"operator": ">",
			"critical": null,
			"warning": null,
			"maxCheckAttempts": 3,
			"missingDurationWarning": 360,
			"missingDurationCritical": 720,
			"notificationInterval": 60
		}`,
	},
	{
		"external monitor",
		&MonitorExternalHTTP{
			ID:                              "2cSZzK3XfmD",
			Name:                            "example.com",
			Type:                            "external",
			IsMute:                          false,
			NotificationInterval:            0,
			Method:                          "POST",
			URL:                             "https://example.com",
			MaxCheckAttempts:                7,
			Service:                         "Hatena-Blog",
			ResponseTimeCritical:            pfloat64(3000.0),
			ResponseTimeWarning:             pfloat64(2000.0),
			ResponseTimeDuration:            puint64(7),
			RequestBody:                     "Request Body",
			ContainsString:                  "",
			CertificationExpirationCritical: puint64(60),
			CertificationExpirationWarning:  puint64(90),
			SkipCertificateVerification:     false,
			Headers: []HeaderField{
				{
					Name:  "Cache-Control",
					Value: "no-cache",
				},
			},
		},
		`{
			"id"  : "2cSZzK3XfmD",
			"type": "external",
			"name": "example.com",
			"method": "POST",
			"url": "https://example.com",
			"service": "Hatena-Blog",
			"headers": [{"name":"Cache-Control", "value":"no-cache"}],
			"requestBody": "Request Body",
			"maxCheckAttempts": 7,
			"responseTimeCritical": 3000,
			"responseTimeWarning": 2000,
			"responseTimeDuration": 7,
			"certificationExpirationCritical": 60,
			"certificationExpirationWarning": 90
		}`,
	},
	{
		"external monitor without service",
		&MonitorExternalHTTP{
			ID:               "2cSZzK3XfmY",
			Name:             "POST example.com",
			Type:             "external",
			Method:           "POST",
			URL:              "https://example.com",
			MaxCheckAttempts: 5,
			RequestBody:      "Request Body",
			ContainsString:   "",
			Headers:          []HeaderField{},
		},
		`{
			"id"  : "2cSZzK3XfmY",
			"type": "external",
			"name": "POST example.com",
			"method": "POST",
			"url": "https://example.com",
			"headers": [],
			"requestBody": "Request Body",
			"maxCheckAttempts": 5
		}`,
	},
	{
		"external monitor with empty threshold",
		&MonitorExternalHTTP{
			ID:                              "2cSZzK3XfmH",
			Name:                            "example.com",
			Type:                            "external",
			IsMute:                          false,
			NotificationInterval:            0,
			Method:                          "GET",
			URL:                             "https://example.com",
			MaxCheckAttempts:                5,
			Service:                         "Hatena-Blog",
			ResponseTimeCritical:            nil,
			ResponseTimeWarning:             pfloat64(3000.0),
			ResponseTimeDuration:            puint64(7),
			RequestBody:                     "Request Body",
			ContainsString:                  "",
			CertificationExpirationCritical: puint64(30),
			CertificationExpirationWarning:  nil,
			SkipCertificateVerification:     false,
			Headers: []HeaderField{
				{
					Name:  "Cache-Control",
					Value: "no-cache",
				},
			},
		},
		`{
			"id"  : "2cSZzK3XfmH",
			"type": "external",
			"name": "example.com",
			"method": "GET",
			"url": "https://example.com",
			"service": "Hatena-Blog",
			"headers": [{"name":"Cache-Control", "value":"no-cache"}],
			"requestBody": "Request Body",
			"maxCheckAttempts": 5,
			"responseTimeWarning": 3000,
			"responseTimeDuration": 7,
			"certificationExpirationCritical": 30
		}`,
	},
	{
		"expression monitor",
		&MonitorExpression{
			ID:                   "2cSZzK3XfmE",
			Name:                 "role average",
			Type:                 "expression",
			IsMute:               false,
			NotificationInterval: 60,
			Expression:           "avg(roleSlots(\"server:role\",\"loadavg5\"))",
			Operator:             ">",
			Warning:              pfloat64(5.000000),
			Critical:             pfloat64(10.000000),
		},
		`{
			"id"  : "2cSZzK3XfmE",
			"type": "expression",
			"name": "role average",
			"expression": "avg(roleSlots(\"server:role\",\"loadavg5\"))",
			"operator": ">",
			"warning": 5.0,
			"critical": 10.0,
			"notificationInterval": 60
		}`,
	},
	{
		"expression monitor without thresholds",
		&MonitorExpression{
			ID:                   "2cSZzK3XfmE",
			Name:                 "role average",
			Type:                 "expression",
			IsMute:               false,
			NotificationInterval: 60,
			Expression:           "avg(roleSlots(\"server:role\",\"loadavg5\"))",
			Operator:             ">",
			Warning:              nil,
			Critical:             nil,
		},
		`{
			"id"  : "2cSZzK3XfmE",
			"type": "expression",
			"name": "role average",
			"expression": "avg(roleSlots(\"server:role\",\"loadavg5\"))",
			"operator": ">",
			"warning": null,
			"critical": null,
			"notificationInterval": 60
		}`,
	},
}

func TestDecodeEncodeMonitor(t *testing.T) {
	for _, testCase := range testCases {
		gotMonitor, err := decodeMonitorReader(strings.NewReader(testCase.json))
		if err != nil {
			t.Errorf("%s: err should be nil but: %v", testCase.title, err)
		}
		if !reflect.DeepEqual(gotMonitor, testCase.monitor) {
			t.Errorf("%s: fail to get correct data: diff: (-got +want)\n%v", testCase.title, pretty.Compare(gotMonitor, testCase.monitor))
		}

		b, err := json.MarshalIndent(testCase.monitor, "", "    ")
		if err != nil {
			t.Errorf("%s: err should be nil but: %v", testCase.title, err)
		}
		if gotJSON := string(b); !equalJSON(gotJSON, testCase.json) {
			t.Errorf("%s: got %v, want %v", testCase.title, gotJSON, testCase.json)
		}
	}
}

func equalJSON(x, y string) bool {
	var xval, yval interface{}
	json.Unmarshal([]byte(x), &xval)
	json.Unmarshal([]byte(y), &yval)
	return reflect.DeepEqual(xval, yval)
}
