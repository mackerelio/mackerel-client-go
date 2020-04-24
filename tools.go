// +build tools

package mackerel

import (
	_ "github.com/mattn/goveralls"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/cover"
)
