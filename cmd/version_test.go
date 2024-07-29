package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildVersion(t *testing.T) {
	got, _ := buildVersion("0.0.0-next")
	assert.Equal(t, "0.0.0-next", got)
}
