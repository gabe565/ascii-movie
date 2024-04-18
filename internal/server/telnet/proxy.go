package telnet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type TermInfo struct {
	Term string
	WindowSize
}

type WindowSize struct {
	Width, Height uint16
}

func Proxy(conn net.Conn, proxy io.Writer, termCh chan TermInfo) error {
	reader := bufio.NewReaderSize(conn, 64)
	var info TermInfo
	var wroteTelnetCommands bool
	var wroteTermType bool

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
				log.Trace("Configuring Telnet")
				if _, err := Write(conn,
					Iac, Will, Echo,
					Iac, Will, SuppressGoAhead,
					Iac, Do, TerminalType,
					Iac, Do, NegotiateAboutWindowSize,
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
				command, err := reader.ReadBytes(byte(Se))
				if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
					return err
				}
				if err := conn.SetReadDeadline(time.Time{}); err != nil {
					return err
				}

				if len(command) != 0 {
					switch Operator(command[0]) {
					case TerminalType:
						if len(command) > 5 && !wroteTermType {
							wroteTermType = true
							info.Term = string(command[2 : len(command)-2])
							log.Trace("Got terminal type")
							termCh <- info
						}
					case NegotiateAboutWindowSize:
						if len(command) > 5 {
							log.Trace("Got window size")
							r := bytes.NewReader(command[1 : len(command)-2])
							if err := binary.Read(r, binary.BigEndian, &info.WindowSize); err != nil {
								return err
							}
							termCh <- info
						}
					}
				}
			case Will:
				if b, err = reader.ReadByte(); err != nil {
					return err
				}

				switch Operator(b) {
				case TerminalType:
					log.Trace("Requesting terminal type")
					if _, err := Write(conn, Iac, Subnegotiation, TerminalType, 1, Iac, Se); err != nil {
						return err
					}
				}
			case Wont, Do, Dont:
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
