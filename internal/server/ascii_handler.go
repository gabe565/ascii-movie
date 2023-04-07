package server

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/ahmetb/go-cursor"
	"github.com/gabe565/ascii-telnet-go/generated_frames"
	flag "github.com/spf13/pflag"
	"io"
	"time"
)

var (
	ClearExtraLinesFlag = "clear-extra-lines"

	SpeedFlag = "speed"

	ErrInvalidSpeed = errors.New("speed must be greater than 0")
)

func Flags(flags *flag.FlagSet) {
	flags.Int(
		ClearExtraLinesFlag,
		0,
		"Clears extra lines between each frame. Should typically be ignored.",
	)
	flags.Float64(
		SpeedFlag,
		1,
		"Playback speed multiplier. Must be greater than 0.",
	)
}

func New(flags *flag.FlagSet) (handler Handler, err error) {
	handler.ClearExtraLines, err = flags.GetInt(ClearExtraLinesFlag)
	if err != nil {
		return handler, err
	}

	handler.Speed, err = flags.GetFloat64(SpeedFlag)
	if err != nil {
		return handler, err
	}
	if handler.Speed <= 0 {
		return handler, fmt.Errorf("%w: %f", ErrInvalidSpeed, handler.Speed)
	}

	return handler, nil
}

type Handler struct {
	ClearExtraLines int

	Speed float64
}

func (s *Handler) ServeAscii(w io.Writer) error {
	var buf bytes.Buffer
	buf.Grow(generated_frames.Cap)

	for _, f := range generated_frames.List {
		buf.WriteString(f.Data)

		if _, err := io.Copy(w, &buf); err != nil {
			return err
		}

		time.Sleep(f.CalcDuration(s.Speed))

		buf.Reset()
		buf.WriteString(cursor.MoveUp(f.Height+s.ClearExtraLines) + cursor.ClearScreenDown())
	}
	return nil
}

func (s *Handler) MovieDuration() time.Duration {
	var totalDuration time.Duration
	for _, f := range generated_frames.List {
		totalDuration += f.CalcDuration(s.Speed)
	}
	return totalDuration
}
