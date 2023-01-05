package mackerel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
)

func TestFindDowntimes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/downtimes" {
			t.Error("request URL should be /api/v0/downtimes but: ", req.URL.Path)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"downtimes": []interface{}{
				map[string]interface{}{
					"id":       "abcde0",
					"name":     "Maintenance #0",
					"memo":     "Memo #0",
					"start":    1563600000,
					"duration": 120,
				},
				map[string]interface{}{
					"id":       "abcde1",
					"name":     "Maintenance #1",
					"memo":     "Memo #1",
					"start":    1563700000,
					"duration": 60,
					"recurrence": map[string]interface{}{
						"interval": 3,
						"type":     "weekly",
						"weekdays": []string{
							"Monday",
							"Thursday",
							"Saturday",
						},
					},
					"serviceScopes": []string{
						"service1",
					},
					"serviceExcludeScopes": []string{
						"service2",
					},
					"roleScopes": []string{
						"service3: role1",
					},
					"roleExcludeScopes": []string{
						"service1: role1",
					},
					"monitorScopes": []string{
						"monitor0",
					},
					"monitorExcludeScopes": []string{
						"monitor1",
					},
				}},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	downtimes, err := client.FindDowntimes()

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	expected := []*Downtime{
		{
			ID:       "abcde0",
			Name:     "Maintenance #0",
			Memo:     "Memo #0",
			Start:    1563600000,
			Duration: 120,
		},
		{
			ID:       "abcde1",
			Name:     "Maintenance #1",
			Memo:     "Memo #1",
			Start:    1563700000,
			Duration: 60,
			Recurrence: &DowntimeRecurrence{
				Type:     DowntimeRecurrenceTypeWeekly,
				Interval: 3,
				Weekdays: []DowntimeWeekday{
					DowntimeWeekday(time.Monday),
					DowntimeWeekday(time.Thursday),
					DowntimeWeekday(time.Saturday),
				},
			},
			ServiceScopes: []string{
				"service1",
			},
			ServiceExcludeScopes: []string{
				"service2",
			},
			RoleScopes: []string{
				"service3: role1",
			},
			RoleExcludeScopes: []string{
				"service1: role1",
			},
			MonitorScopes: []string{
				"monitor0",
			},
			MonitorExcludeScopes: []string{
				"monitor1",
			},
		},
	}

	if !reflect.DeepEqual(downtimes, expected) {
		t.Errorf("fail to get correct data: diff: (-got +want)\n%v", pretty.Compare(downtimes, expected))
	}
}

func TestCreateDowntime(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/downtimes" {
			t.Error("request URL should be /api/v0/downtimes but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		body, _ := io.ReadAll(req.Body)

		var d Downtime
		err := json.Unmarshal(body, &d)
		if err != nil {
			t.Fatal("request body should be decoded as downtime", string(body))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":       "abcde",
			"name":     "My maintenance",
			"memo":     "This is a memo",
			"start":    1563700000,
			"duration": 30,
			"recurrence": map[string]interface{}{
				"type":     "daily",
				"interval": 2,
				"until":    1573700000,
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	downtime, err := client.CreateDowntime(&Downtime{
		Name:     "My maintenance",
		Memo:     "This is a memo",
		Start:    1563700000,
		Duration: 30,
		Recurrence: &DowntimeRecurrence{
			Type:     DowntimeRecurrenceTypeDaily,
			Interval: 2,
			Until:    1573700000,
		},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	expected := &Downtime{
		ID:       "abcde",
		Name:     "My maintenance",
		Memo:     "This is a memo",
		Start:    1563700000,
		Duration: 30,
		Recurrence: &DowntimeRecurrence{
			Type:     DowntimeRecurrenceTypeDaily,
			Interval: 2,
			Until:    1573700000,
		},
	}
	if !reflect.DeepEqual(downtime, expected) {
		t.Errorf("fail to get correct data: diff: (-got +want)\n%v", pretty.Compare(downtime, expected))
	}
}

func TestUpdateDowntime(t *testing.T) {
	downtimeID := "abcde"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/downtimes/%s", downtimeID) {
			t.Error("request URL should be /api/v0/downtimes/<ID> but: ", req.URL.Path)
		}

		if req.Method != "PUT" {
			t.Error("request method should be PUT but: ", req.Method)
		}

		body, _ := io.ReadAll(req.Body)

		var d Downtime
		err := json.Unmarshal(body, &d)
		if err != nil {
			t.Fatal("request body should be decoded as json", string(body))
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":       "abcde",
			"name":     "Updated maintenance",
			"memo":     "This is a memo",
			"start":    1563700000,
			"duration": 30,
			"serviceScopes": []string{
				"service1",
				"service2",
			},
			"roleExcludeScopes": []string{
				"service1: role1",
				"service2: role2",
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	downtime, err := client.UpdateDowntime(downtimeID, &Downtime{
		Name:     "Updated maintenance",
		Memo:     "This is a memo",
		Start:    1563700000,
		Duration: 30,
		ServiceScopes: []string{
			"service1",
			"service2",
		},
		RoleExcludeScopes: []string{
			"service1: role1",
			"service2: role2",
		},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	expected := &Downtime{
		ID:       "abcde",
		Name:     "Updated maintenance",
		Memo:     "This is a memo",
		Start:    1563700000,
		Duration: 30,
		ServiceScopes: []string{
			"service1",
			"service2",
		},
		RoleExcludeScopes: []string{
			"service1: role1",
			"service2: role2",
		},
	}
	if !reflect.DeepEqual(downtime, expected) {
		t.Errorf("fail to get correct data: diff: (-got +want)\n%v", pretty.Compare(downtime, expected))
	}
}

func TestDeleteDowntime(t *testing.T) {
	downtimeID := "abcde"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/downtimes/%s", downtimeID) {
			t.Error("request URL should be /api/v0/downtimes/<ID> but: ", req.URL.Path)
		}

		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":       "abcde",
			"name":     "My maintenance",
			"memo":     "This is a memo",
			"start":    1563700000,
			"duration": 60,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	downtime, err := client.DeleteDowntime(downtimeID)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	expected := &Downtime{
		ID:       "abcde",
		Name:     "My maintenance",
		Memo:     "This is a memo",
		Start:    1563700000,
		Duration: 60,
	}
	if !reflect.DeepEqual(downtime, expected) {
		t.Errorf("fail to get correct data: diff: (-got +want)\n%v", pretty.Compare(downtime, expected))
	}
}

var downtimeRecurrenceTestCases = []struct {
	title      string
	recurrence *DowntimeRecurrence
	json       string
}{
	{
		"hourly",
		&DowntimeRecurrence{
			Type:     DowntimeRecurrenceTypeHourly,
			Interval: 1,
		},
		`{
			"type": "hourly",
			"interval": 1
		}`,
	},
	{
		"daily",
		&DowntimeRecurrence{
			Type:     DowntimeRecurrenceTypeDaily,
			Interval: 3,
			Until:    1573730000,
		},
		`{
			"type": "daily",
			"interval": 3,
			"until": 1573730000
		}`,
	},
	{
		"weekly",
		&DowntimeRecurrence{
			Type:     DowntimeRecurrenceTypeWeekly,
			Interval: 2,
			Weekdays: []DowntimeWeekday{
				DowntimeWeekday(time.Sunday),
				DowntimeWeekday(time.Monday),
				DowntimeWeekday(time.Tuesday),
				DowntimeWeekday(time.Wednesday),
				DowntimeWeekday(time.Thursday),
				DowntimeWeekday(time.Friday),
				DowntimeWeekday(time.Saturday),
			},
		},
		`{
			"type": "weekly",
			"interval": 2,
			"weekdays": [
				"Sunday",
				"Monday",
				"Tuesday",
				"Wednesday",
				"Thursday",
				"Friday",
				"Saturday"
			]
		}`,
	},
	{
		"monthly",
		&DowntimeRecurrence{
			Type:     DowntimeRecurrenceTypeMonthly,
			Interval: 2,
		},
		`{
			"type": "monthly",
			"interval": 2
		}`,
	},
	{
		"yearly",
		&DowntimeRecurrence{
			Type:     DowntimeRecurrenceTypeYearly,
			Interval: 1,
		},
		`{
			"type": "yearly",
			"interval": 1
		}`,
	},
}

func TestDecodeEncodeDowntimeRecurrence(t *testing.T) {
	for _, testCase := range downtimeRecurrenceTestCases {
		r := new(DowntimeRecurrence)
		err := json.Unmarshal([]byte(testCase.json), r)
		if err != nil {
			t.Errorf("%s: err should be nil but: %v", testCase.title, err)
		}
		if !reflect.DeepEqual(r, testCase.recurrence) {
			t.Errorf("%s: fail to get correct data: diff: (-got +want)\n%v", testCase.title, pretty.Compare(r, testCase.recurrence))
		}

		b, err := json.MarshalIndent(&testCase.recurrence, "", "    ")
		if err != nil {
			t.Errorf("%s: err should be nil but: %v", testCase.title, err)
		}
		if gotJSON := string(b); !equalJSON(gotJSON, testCase.json) {
			t.Errorf("%s: got %v, want %v", testCase.title, gotJSON, testCase.json)
		}
	}
}
