package mackerel

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// APIError represents the error type from Mackerel API.
type APIError struct {
	StatusCode int
	Message    string
}

func (err *APIError) Error() string {
	return fmt.Sprintf("API request failed: %s", err.Message)
}

func extractErrorMessage(r io.Reader) (errorMessage string, err error) {
	bs, err := ioutil.ReadAll(r)
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
