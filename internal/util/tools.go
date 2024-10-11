//go:build tools

package util

import (
	// Protect this entry in go.mod from being removed by go mod tidy.
	_ "github.com/dmarkham/enumer"
)
