package main

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/ahmetb/go-cursor"
	"github.com/gabe565/ascii-telnet-go/generated_frames"
	"github.com/reiver/go-telnet"
	"syscall"
	"time"
)

type AsciiHandler struct{}

func (handler AsciiHandler) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	var buf bytes.Buffer
	for _, f := range generated_frames.List {
		buf.WriteString(f.Data)

		_, err := w.Write(buf.Bytes())
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				return
			}
			panic(err)
		}

		time.Sleep(f.Sleep)

		buf.Reset()
		buf.WriteString(cursor.MoveUp(f.Height) + cursor.ClearScreenDown())
	}
}
