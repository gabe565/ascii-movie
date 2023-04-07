package movie

import (
	"bytes"
	"github.com/ahmetb/go-cursor"
	"io"
	"time"
)

type Movie struct {
	Filename string
	Cap      int
	Frames   []Frame

	Speed float64

	ClearExtraLines int
}

func (m Movie) Duration() time.Duration {
	var totalDuration time.Duration
	for _, f := range m.Frames {
		totalDuration += f.CalcDuration(m.Speed)
	}
	return totalDuration
}

func (m *Movie) Stream(w io.Writer) error {
	var buf bytes.Buffer
	buf.Grow(m.Cap)

	for _, f := range m.Frames {
		buf.WriteString(f.Data)

		if _, err := io.Copy(w, &buf); err != nil {
			return err
		}

		time.Sleep(f.CalcDuration(m.Speed))

		buf.Reset()
		buf.WriteString(cursor.MoveUp(f.Height+m.ClearExtraLines) + cursor.ClearScreenDown())
	}
	return nil
}
