package mackerel

import (
	"encoding/json"
	"testing"
)

func TestCheckReports_MarshalJSON(t *testing.T) {
	crs := &CheckReports{
		Reports: []*CheckReport{
			{
				Source:     NewCheckSourceHost("hogee"),
				Name:       "chchch",
				Status:     CheckStatusCritical,
				OccurredAt: 100,
				Message:    "OKOK",
			},
		},
	}
	expect := `{"reports":[{"source":{"type":"host","hostId":"hogee"},"name":"chchch","status":"CRITICAL","message":"OKOK","occurredAt":100}]}`
	bs, _ := json.Marshal(crs)
	got := string(bs)

	if got != expect {
		t.Errorf("expect: %s, but: %s", expect, got)
	}
}
