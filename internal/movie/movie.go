package movie

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/ahmetb/go-cursor"
)

type Movie struct {
	Filename string
	Cap      int
	Frames   []Frame

	Speed float64
}

func (m Movie) Duration() time.Duration {
	var totalDuration time.Duration
	for _, f := range m.Frames {
		totalDuration += f.CalcDuration(m.Speed)
	}
	return totalDuration
}

func (m *Movie) Stream(ctx context.Context, w io.Writer) error {
	var buf bytes.Buffer
	buf.Grow(m.Cap)

	timer := time.NewTimer(0)
	for _, f := range m.Frames {
		buf.WriteString(f.Data)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if _, err := io.Copy(w, &buf); err != nil {
				return err
			}
			timer.Reset(f.CalcDuration(m.Speed))
		}

		buf.Reset()
		buf.WriteString(cursor.MoveUp(f.Height) + cursor.ClearScreenDown())
	}
	return nil
}
