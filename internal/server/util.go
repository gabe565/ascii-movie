package server

import (
	"context"
	log "github.com/sirupsen/logrus"
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

func HandleInput(ctx context.Context, cancel context.CancelFunc, in io.Reader, out io.Writer) {
	b := make([]byte, 1)
	var skip int8
	var wroteTelnetCommands bool
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
				// IAC WILL Suppress Go Ahead IAC WON'T X Display Location
				// https://ibm.com/docs/zos/2.5.0?topic=problems-telnet-commands-options
				if out != nil && !wroteTelnetCommands {
					if _, err := out.Write([]byte{0xFF, 0xFB, 0x3, 0xFF, 0xFC, 0x23}); err != nil {
						log.WithError(err).Error("Failed to write Telnet commands")
					}
					wroteTelnetCommands = true
				}
				skip = 2
			}
		}
	}
}
