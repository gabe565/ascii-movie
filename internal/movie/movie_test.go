package movie

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMovie_Duration(t *testing.T) {
	type fields struct {
		Filename        string
		Cap             int
		Frames          []Frame
		Speed           float64
		ClearExtraLines int
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{"0", fields{Speed: 1, Frames: []Frame{{}}}, 0},
		{"2 1s frames", fields{Speed: 1, Frames: []Frame{{Duration: time.Second}, {Duration: time.Second}}}, 2 * time.Second},
		{"2 1s frames 2x multiplier", fields{Speed: 2, Frames: []Frame{{Duration: time.Second}, {Duration: time.Second}}}, time.Second},
		{"2 1s frames 0.5x multiplier", fields{Speed: 0.5, Frames: []Frame{{Duration: time.Second}, {Duration: time.Second}}}, 4 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Movie{
				Filename:        tt.fields.Filename,
				Cap:             tt.fields.Cap,
				Frames:          tt.fields.Frames,
				Speed:           tt.fields.Speed,
				ClearExtraLines: tt.fields.ClearExtraLines,
			}
			assert.Equal(t, tt.want, m.Duration())
		})
	}
}

func TestMovie_Stream(t *testing.T) {
	t.Run("stream output", func(t *testing.T) {
		t.Parallel()

		movie := Movie{Speed: 1, Frames: []Frame{
			{Duration: time.Millisecond, Height: 1, Data: "Test frame 1"},
			{Duration: time.Millisecond, Height: 1, Data: "Test frame 2"},
		}}

		var writeBuf bytes.Buffer
		if err := movie.Stream(context.Background(), &writeBuf); !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, "Test frame 1\x1B[1A\x1b[0JTest frame 2", writeBuf.String())
	})

	t.Run("cancel context", func(t *testing.T) {
		t.Parallel()

		movie := Movie{Speed: 1, Frames: []Frame{
			{Duration: 5 * time.Second, Height: 1, Data: "Test frame 1"},
			{Duration: 5 * time.Second, Height: 1, Data: "Test frame 2"},
		}}

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		var writeBuf bytes.Buffer
		start := time.Now()
		if err := movie.Stream(ctx, &writeBuf); !assert.ErrorIs(t, err, context.Canceled) {
			return
		}

		assert.Less(t, time.Since(start), time.Second)
	})

	t.Run("sleep time", func(t *testing.T) {
		t.Parallel()

		movie := Movie{Speed: 1, Frames: []Frame{
			{Duration: 100 * time.Millisecond, Height: 1, Data: "Test frame 1"},
			{Duration: 250 * time.Millisecond, Height: 1, Data: "Test frame 2"},
			{Duration: 50 * time.Millisecond, Height: 1, Data: "Test frame 3"},
			{Duration: 500 * time.Millisecond, Height: 1, Data: "Test frame 4"},
		}}

		r, w := io.Pipe()

		start := time.Now()
		go func() {
			if err := movie.Stream(context.Background(), w); !assert.NoError(t, err) {
				return
			}
		}()

		var elapsed time.Duration
		var prevDuration time.Duration
		data := make([]byte, 20)
	outer:
		for _, frame := range movie.Frames {
			if _, err := r.Read(data); !assert.NoError(t, err) {
				return
			}
			elapsed = time.Since(start)
			start = time.Now()
			assert.Equal(t, prevDuration.Truncate(10*time.Millisecond), elapsed.Truncate(10*time.Millisecond))
			assert.Contains(t, string(data), frame.Data)
			prevDuration = frame.Duration
			continue outer
		}
	})
}
