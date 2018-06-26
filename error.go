package mackerel

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

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
		err = json.Unmarshal(bs, &data)
		errorMessage = data.Error
	}
	return
}
