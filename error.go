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

func extractErrorMessage(r io.Reader) (errorMessage string) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	var data struct{ Error struct{ Message string } }
	err = json.Unmarshal(bs, &data)
	if err == nil {
		errorMessage = data.Error.Message
	} else {
		var data struct{ Error string }
		json.Unmarshal(bs, &data)
		errorMessage = data.Error
	}
	return
}
