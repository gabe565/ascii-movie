package telnet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

type WindowSize struct {
	Width, Height uint16
}

func Proxy(conn net.Conn, proxy io.Writer, termCh chan string, sizeCh chan WindowSize) error {
	reader := bufio.NewReaderSize(conn, 64)
	var wroteTelnetCommands bool
	var wroteTermType bool

	// Gets Telnet to send option negotiation commands if explicit port was given.
	// Also clears the line in case the client isn't Telnet
	// https://ibm.com/docs/zos/2.5.0?topic=problems-telnet-commands-options
	if _, err := WriteAndClear(conn, Iac, Do, Linemode); err != nil {
		return err
	}

outer:
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
				command, err := reader.ReadBytes(byte(Se))
				if err != nil {
					return err
				}
				if len(command) < 2 {
					continue
				}
				for i := 0; command[len(command)-2] != byte(Iac); i++ {
					commandExt, err := reader.ReadBytes(byte(Se))
					if err != nil {
						return err
					}
					command = append(command, commandExt...)
					if i == 2 {
						continue outer
					}
				}

				command = bytes.ReplaceAll(command, []byte{byte(Iac), byte(Iac)}, []byte{byte(Iac)})

				if len(command) != 0 {
					switch Operator(command[0]) {
					case TerminalType:
						if !wroteTermType && len(command) > 2 {
							log.Trace("Got terminal type")
							termCh <- string(command[2:])
							wroteTermType = true
						}
					case NegotiateAboutWindowSize:
						if len(command) >= 5 {
							log.Trace("Got window size")
							r := bytes.NewReader(command[1:])
							var size WindowSize
							if err := binary.Read(r, binary.BigEndian, &size); err != nil {
								return err
							}
							sizeCh <- size
						}
					}
				}
			case Will:
				if b, err = reader.ReadByte(); err != nil {
					return err
				}

				if b == byte(TerminalType) {
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
