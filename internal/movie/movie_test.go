package movie

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMovie_Duration(t *testing.T) {
	type fields struct {
		Filename string
		Cap      int
		Frames   []Frame
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{"0", fields{Frames: []Frame{{}}}, 0},
		{"2 1s frames", fields{Frames: []Frame{{Duration: time.Second}, {Duration: time.Second}}}, 2 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Movie{
				Filename: tt.fields.Filename,
				Cap:      tt.fields.Cap,
				Frames:   tt.fields.Frames,
			}
			assert.Equal(t, tt.want, m.Duration())
		})
	}
}
