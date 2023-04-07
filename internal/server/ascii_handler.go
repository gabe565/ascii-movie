package server

import (
	"bytes"
	_ "embed"
	"github.com/ahmetb/go-cursor"
	"github.com/gabe565/ascii-telnet-go/generated_frames"
	flag "github.com/spf13/pflag"
	"io"
	"time"
)

var ClearExtraLinesFlag = "clear-extra-lines"

func Flags(flags *flag.FlagSet) {
	flags.Int(
		ClearExtraLinesFlag,
		0,
		"Clears extra lines between each frame. Should typically be ignored.",
	)
}

func New(flags *flag.FlagSet) (handler Handler, err error) {
	handler.ClearExtraLines, err = flags.GetInt(ClearExtraLinesFlag)
	if err != nil {
		return handler, err
	}
	return handler, nil
}

type Handler struct {
	ClearExtraLines int
}

func (s *Handler) ServeAscii(w io.Writer) error {
	var buf bytes.Buffer
	buf.Grow(generated_frames.Cap)

	for _, f := range generated_frames.List {
		buf.WriteString(f.Data)

		if _, err := io.Copy(w, &buf); err != nil {
			return err
		}

		time.Sleep(f.Sleep)

		buf.Reset()
		buf.WriteString(cursor.MoveUp(f.Height+s.ClearExtraLines) + cursor.ClearScreenDown())
	}
	return nil
}
