package telnet

import (
	"bufio"
	"errors"
	"io"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func Proxy(conn net.Conn, proxy io.Writer) error {
	reader := bufio.NewReaderSize(conn, 64)
	var wroteTelnetCommands bool

	// Gets Telnet to send option negotiation commands if explicit port was given.
	// Also clears the line in case the client isn't Telnet
	// https://ibm.com/docs/zos/2.5.0?topic=problems-telnet-commands-options
	if _, err := WriteAndClear(conn, Iac, Do, Linemode); err != nil {
		return err
	}

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}

		switch Operator(b) {
		case BinaryTransmission:
		case Iac:
			// https://ibm.com/docs/zos/2.5.0?topic=problems-telnet-commands-options
			if conn != nil && !wroteTelnetCommands {
				log.Trace("Writing Telnet commands")
				if _, err := Write(conn,
					Iac, Will, Echo,
					Iac, Will, SuppressGoAhead,
				); err != nil {
					log.WithError(err).Error("Failed to write Telnet commands")
				}

				wroteTelnetCommands = true
			}

			if b, err = reader.ReadByte(); err != nil {
				return err
			}

			switch Operator(b) {
			case Subnegotiation:
				if err := conn.SetReadDeadline(time.Now().Add(250 * time.Millisecond)); err != nil {
					return err
				}
				_, err := reader.ReadBytes(byte(Se))
				if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
					return err
				}
				if err := conn.SetReadDeadline(time.Time{}); err != nil {
					return err
				}
			case Will, Wont, Do, Dont:
				if _, err := reader.Discard(1); err != nil {
					return err
				}
			}
		default:
			if err := reader.UnreadByte(); err != nil {
				return err
			}
			if _, err := io.CopyN(proxy, reader, int64(reader.Buffered())); err != nil {
				return err
			}
		}
	}
}
