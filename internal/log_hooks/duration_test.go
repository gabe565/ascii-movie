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

func TestDuration_GetStart(t *testing.T) {
	start := time.Now()

	type fields struct {
		start time.Time
		trunc time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{"simple", fields{start: start}, start},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Duration{
				start: tt.fields.start,
				trunc: tt.fields.trunc,
			}
			assert.Equalf(t, tt.want, d.GetStart(), "GetStart()")
		})
	}
}
