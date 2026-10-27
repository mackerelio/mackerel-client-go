package mackerel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// APIError represents the error type from Mackerel API.
type APIError struct {
	StatusCode int
	Message    string
	// RetryAfter is the duration in the Retry-After header, or zero if the header is absent or invalid.
	RetryAfter time.Duration
}

func (err *APIError) Error() string {
	return fmt.Sprintf("API request failed: %s", err.Message)
}

func extractErrorMessage(r io.Reader) (errorMessage string, err error) {
	bs, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	var data struct{ Error struct{ Message string } }
	err = json.Unmarshal(bs, &data)
	if err != nil {
		return "", err
	}
	return data.Error.Message, nil
}

func parseRetryAfter(s string) time.Duration {
	if s == "" {
		return 0
	}
	if secs, err := strconv.ParseInt(s, 10, 32); err == nil {
		return max(time.Duration(secs)*time.Second, 0)
	}
	if t, err := http.ParseTime(s); err == nil {
		return max(time.Until(t), 0)
	}
	return 0
}
