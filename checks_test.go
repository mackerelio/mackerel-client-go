package mackerel

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCheckReports_MarshalJSON(t *testing.T) {
	s := "1980-06-05T15:04:05+09:00"
	ti, _ := time.Parse(time.RFC3339, s)
	crs := &CheckReports{
		Reports: []*CheckReport{
			{
				Source:     NewCheckSourceHost("hogee"),
				Name:       "chchch",
				Status:     CheckStatusCritical,
				OccurredAt: Time(ti),
				Message:    "OKOK",
			},
		},
	}
	expect := `{"reports":[{"source":{"type":"host","hostId":"hogee"},"name":"chchch","status":"CRITICAL","message":"OKOK","occurredAt":329033045}]}`
	bs, _ := json.Marshal(crs)
	got := string(bs)

	if got != expect {
		t.Errorf("expect: %s, but: %s", expect, got)
	}
}
