package movie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListEmbedded(t *testing.T) {
	got, err := ListEmbedded()
	if !assert.NoError(t, err) {
		return
	}

	assert.NotEmpty(t, got)
}
