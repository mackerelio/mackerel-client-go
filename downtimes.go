package mackerel

import (
	"encoding/json"
	"fmt"
	"time"
)

// Downtime information
type Downtime struct {
	ID                   string              `json:"id,omitempty"`
	Name                 string              `json:"name"`
	Memo                 string              `json:"memo,omitempty"`
	Start                int64               `json:"start"`
	Duration             int64               `json:"duration"`
	Recurrence           *DowntimeRecurrence `json:"recurrence,omitempty"`
	ServiceScopes        []string            `json:"serviceScopes,omitempty"`
	ServiceExcludeScopes []string            `json:"serviceExcludeScopes,omitempty"`
	RoleScopes           []string            `json:"roleScopes,omitempty"`
	RoleExcludeScopes    []string            `json:"roleExcludeScopes,omitempty"`
	MonitorScopes        []string            `json:"monitorScopes,omitempty"`
	MonitorExcludeScopes []string            `json:"monitorExcludeScopes,omitempty"`
}

// DowntimeRecurrence ...
type DowntimeRecurrence struct {
	Type     DowntimeRecurrenceType `json:"type"`
	Interval int64                  `json:"interval"`
	Weekdays []DowntimeWeekday      `json:"weekdays,omitempty"`
	Until    int64                  `json:"until,omitempty"`
}

// DowntimeRecurrenceType ...
type DowntimeRecurrenceType int

// DowntimeRecurrenceType ...
const (
	DowntimeRecurrenceTypeHourly DowntimeRecurrenceType = iota + 1
	DowntimeRecurrenceTypeDaily
	DowntimeRecurrenceTypeWeekly
	DowntimeRecurrenceTypeMonthly
	DowntimeRecurrenceTypeYearly
)

var stringToRecurrenceType = map[string]DowntimeRecurrenceType{
	"hourly":  DowntimeRecurrenceTypeHourly,
	"daily":   DowntimeRecurrenceTypeDaily,
	"weekly":  DowntimeRecurrenceTypeWeekly,
	"monthly": DowntimeRecurrenceTypeMonthly,
	"yearly":  DowntimeRecurrenceTypeYearly,
}

var recurrenceTypeToString = map[DowntimeRecurrenceType]string{
	DowntimeRecurrenceTypeHourly:  "hourly",
	DowntimeRecurrenceTypeDaily:   "daily",
	DowntimeRecurrenceTypeWeekly:  "weekly",
	DowntimeRecurrenceTypeMonthly: "monthly",
	DowntimeRecurrenceTypeYearly:  "yearly",
}

// UnmarshalJSON implements json.Unmarshaler
func (t *DowntimeRecurrenceType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if x, ok := stringToRecurrenceType[s]; ok {
		*t = x
		return nil
	}
	return fmt.Errorf("unknown downtime recurrence type: %q", s)
}

// MarshalJSON implements json.Marshaler
func (t DowntimeRecurrenceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(recurrenceTypeToString[t])
}

// String implements Stringer
func (t DowntimeRecurrenceType) String() string {
	return recurrenceTypeToString[t]
}

// DowntimeWeekday ...
type DowntimeWeekday time.Weekday

var stringToWeekday = map[string]DowntimeWeekday{
	"Sunday":    DowntimeWeekday(time.Sunday),
	"Monday":    DowntimeWeekday(time.Monday),
	"Tuesday":   DowntimeWeekday(time.Tuesday),
	"Wednesday": DowntimeWeekday(time.Wednesday),
	"Thursday":  DowntimeWeekday(time.Thursday),
	"Friday":    DowntimeWeekday(time.Friday),
	"Saturday":  DowntimeWeekday(time.Saturday),
}

var weekdayToString = map[DowntimeWeekday]string{
	DowntimeWeekday(time.Sunday):    "Sunday",
	DowntimeWeekday(time.Monday):    "Monday",
	DowntimeWeekday(time.Tuesday):   "Tuesday",
	DowntimeWeekday(time.Wednesday): "Wednesday",
	DowntimeWeekday(time.Thursday):  "Thursday",
	DowntimeWeekday(time.Friday):    "Friday",
	DowntimeWeekday(time.Saturday):  "Saturday",
}

// UnmarshalJSON implements json.Unmarshaler
func (w *DowntimeWeekday) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if x, ok := stringToWeekday[s]; ok {
		*w = x
		return nil
	}
	return fmt.Errorf("unknown downtime weekday: %q", s)
}

// MarshalJSON implements json.Marshaler
func (w DowntimeWeekday) MarshalJSON() ([]byte, error) {
	return json.Marshal(weekdayToString[w])
}

// String implements Stringer
func (w DowntimeWeekday) String() string {
	return weekdayToString[w]
}

// FindDowntimes finds downtimes.
func (c *Client) FindDowntimes() ([]*Downtime, error) {
	data, err := requestGet[struct {
		Downtimes []*Downtime `json:"downtimes"`
	}](c, "/api/v0/downtimes")
	if err != nil {
		return nil, err
	}
	return data.Downtimes, nil
}

// CreateDowntime creates a downtime.
func (c *Client) CreateDowntime(param *Downtime) (*Downtime, error) {
	return requestPost[Downtime](c, "/api/v0/downtimes", param)
}

// UpdateDowntime updates a downtime.
func (c *Client) UpdateDowntime(downtimeID string, param *Downtime) (*Downtime, error) {
	path := fmt.Sprintf("/api/v0/downtimes/%s", downtimeID)
	return requestPut[Downtime](c, path, param)
}

// DeleteDowntime deletes a downtime.
func (c *Client) DeleteDowntime(downtimeID string) (*Downtime, error) {
	path := fmt.Sprintf("/api/v0/downtimes/%s", downtimeID)
	return requestDelete[Downtime](c, path)
}
