package telnet

import (
	"bytes"
	"io"

	log "github.com/sirupsen/logrus"
)

func Proxy(conn io.ReadWriter, proxy io.Writer) error {
	buf := make([]byte, 64)
	var proxyBuf bytes.Buffer
	var skip uint8
	var subNegotiation bool
	var wroteTelnetCommands bool

	// Gets Telnet to send option negotiation commands if explicit port was given.
	// Also clears the line in case the client isn't Telnet
	// https://ibm.com/docs/zos/2.5.0?topic=problems-telnet-commands-options
	if _, err := conn.Write([]byte{Iac, Do, Linemode}); err != nil {
		return err
	}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}

		for _, c := range buf[:n] {
			switch c {
			case Iac:
				// https://ibm.com/docs/zos/2.5.0?topic=problems-telnet-commands-options
				if conn != nil && !wroteTelnetCommands {
					log.Trace("Writing Telnet commands")
					if _, err := conn.Write([]byte{
						Iac, Do, Linemode,
						Iac, Will, Echo,
						Iac, Will, SuppressGoAhead,
					}); err != nil {
						log.WithError(err).Error("Failed to write Telnet commands")
					}

					wroteTelnetCommands = true
				}
				skip = 3
			case Subnegotiation:
				subNegotiation = true
			case Se:
				if skip == 2 {
					skip = 0
					subNegotiation = false
				}
			default:
				if skip == 0 && !subNegotiation {
					proxyBuf.WriteByte(c)
				}
			}

			if skip != 0 {
				skip -= 1
			}
		}

		if proxyBuf.Len() != 0 {
			if _, err := proxyBuf.WriteTo(proxy); err != nil {
				return err
			}
			proxyBuf.Reset()
		}
	}
}
