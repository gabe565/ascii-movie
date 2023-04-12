package movie

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFrame_CalcDuration(t *testing.T) {
	type fields struct {
		Duration time.Duration
		Height   int
		Data     string
	}
	type args struct {
		multiplier float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Duration
	}{
		{"0", fields{}, args{multiplier: 1}, 0},
		{"1 second", fields{Duration: time.Second}, args{multiplier: 1}, time.Second},
		{"1 second 2x multiplier", fields{Duration: time.Second}, args{multiplier: 2}, time.Second / 2},
		{"1 second 0.5x multiplier", fields{Duration: time.Second}, args{multiplier: 0.5}, 2 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Frame{
				Duration: tt.fields.Duration,
				Height:   tt.fields.Height,
				Data:     tt.fields.Data,
			}
			assert.Equal(t, tt.want, f.CalcDuration(tt.args.multiplier))
		})
	}
}
