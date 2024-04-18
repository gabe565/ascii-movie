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
				scanner := bufio.NewScanner(reader)
				scanner.Split(ScanIacSe)
				if scanner.Scan() {
					command := scanner.Bytes()

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
				}
				if scanner.Err() != nil {
					return scanner.Err()
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

func ScanIacSe(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{byte(Iac), byte(Se)}); i != 0 {
		// We have a full newline-terminated line.
		return i + 2, trimIac(data[:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), trimIac(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func trimIac(data []byte) []byte {
	if len(data) != 0 && data[len(data)-1] == byte(Iac) {
		return data[:len(data)-1]
	}
	return data
}
