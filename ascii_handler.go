package main

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/ahmetb/go-cursor"
	"github.com/gabe565/ascii-telnet-go/generated_frames"
	"net"
	"syscall"
	"time"
)

func Serve(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	var buf bytes.Buffer
	buf.Grow(generated_frames.Cap)

	for _, f := range generated_frames.List {
		buf.WriteString(f.Data)

		_, err := conn.Write(buf.Bytes())
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
