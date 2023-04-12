package log_hooks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration_String(t *testing.T) {
	duration := NewDuration()
	duration.trunc = 100 * time.Millisecond
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, (100 * time.Millisecond).String(), duration.String())
}
