package mackerel

import (
	"net/http"
	"testing"
	"time"
)

func TestParseRetryAfter(t *testing.T) {
	tests := []struct {
		value string
		want  time.Duration
	}{
		{"", 0},
		{"1", time.Second},
		{"120", 2 * time.Minute},
		{"-1", 0},
		{"soon", 0},
		{"9999999999", 0},
		{time.Now().Add(-time.Hour).UTC().Format(http.TimeFormat), 0},
	}
	for _, tt := range tests {
		if got := parseRetryAfter(tt.value); got != tt.want {
			t.Errorf("parseRetryAfter(%q) should be %s but: %s", tt.value, tt.want, got)
		}
	}

	future := time.Now().Add(time.Hour).UTC().Format(http.TimeFormat)
	if got := parseRetryAfter(future); got <= 58*time.Minute || got > time.Hour {
		t.Errorf("parseRetryAfter(%q) should be about 1h but: %s", future, got)
	}
}
