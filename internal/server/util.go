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
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := in.Read(b); err != nil || b[0] == CtrlC || b[0] == CtrlD {
				cancel()
				return
			}
		}
	}
}

func HandleTelnetInput(ctx context.Context, cancel context.CancelFunc, in io.Reader, out io.Writer) {
	b := make([]byte, 1)
	var skip int8
	var subNegotiation bool
	var wroteTelnetCommands bool
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := in.Read(b); err != nil || (skip == 0 && !subNegotiation && (b[0] == CtrlC || b[0] == CtrlD)) {
				cancel()
				return
			}

			if skip > 0 {
				skip -= 1
			}
			switch b[0] {
			case 0xFF:
				// IAC DO LINEMODE IAC WILL Suppress Go Ahead
				// https://ibm.com/docs/zos/2.5.0?topic=problems-telnet-commands-options
				if out != nil && !wroteTelnetCommands {
					log.Trace("Writing Telnet commands")
					if _, err := out.Write([]byte{0xFF, 0xFD, 0x22, 0xFF, 0xFB, 0x3}); err != nil {
						log.WithError(err).Error("Failed to write Telnet commands")
					}
					wroteTelnetCommands = true
				}
				skip = 2
			case 0xFA:
				subNegotiation = true
			case 0xF0:
				if skip == 1 {
					skip = 0
					subNegotiation = false
				}
			}
		}
	}
}
