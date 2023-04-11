package server

import (
	"context"
	"io"
	"net"
)

func RemoteIp(remoteIpPort string) string {
	remoteIP, _, err := net.SplitHostPort(remoteIpPort)
	if err != nil {
		remoteIP = remoteIpPort
	}
	return remoteIP
}

const (
	CtrlC byte = 0x3
	CtrlD byte = 0x4
)

func ListenForExit(ctx context.Context, cancel context.CancelFunc, in io.Reader) {
	b := make([]byte, 1)
	var skip int8
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := in.Read(b); err != nil || (skip == 0 && b[0] == CtrlC || b[0] == CtrlD) {
				cancel()
				return
			}

			if skip > 0 {
				skip -= 1
			}
			if b[0] == 0xFF {
				skip = 2
			}
		}
	}
}
