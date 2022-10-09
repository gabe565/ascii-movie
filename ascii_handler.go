package main

import (
	"bytes"
	_ "embed"
	"github.com/ahmetb/go-cursor"
	"github.com/gabe565/ascii-telnet-go/generated_frames"
	"io"
	"time"
)

func ServeAscii(w io.Writer) error {
	var buf bytes.Buffer
	buf.Grow(generated_frames.Cap)

	for _, f := range generated_frames.List {
		buf.WriteString(f.Data)

		if _, err := io.Copy(w, &buf); err != nil {
			return err
		}

		time.Sleep(f.Sleep)

		buf.Reset()
		buf.WriteString(cursor.MoveUp(f.Height) + cursor.ClearScreenDown())
	}
	return nil
}
