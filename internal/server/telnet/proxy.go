package telnet

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"time"

	"gabe565.com/ascii-movie/internal/config"
	"gabe565.com/ascii-movie/internal/util"
	"github.com/muesli/termenv"
)

type WindowSize struct {
	Width, Height uint16
}

func (w WindowSize) String() string {
	return fmt.Sprintf("%dx%d", w.Width, w.Height)
}

func Proxy(conn net.Conn) (io.ReadCloser, termenv.Profile, <-chan WindowSize, <-chan error) {
	pr, pw := io.Pipe()
	termCh := make(chan string, 1)
	sizeCh := make(chan WindowSize, 1)
	errCh := make(chan error)

	go func() {
		defer func() {
			_ = pw.Close()
			close(termCh)
			close(sizeCh)
			close(errCh)
		}()

		errCh <- proxy(conn, pw, termCh, sizeCh)
	}()

	profile := termenv.Profile(-1)
	select {
	case term := <-termCh:
		profile = util.Profile(term)
	case <-time.After(time.Second):
	}

	return pr, profile, sizeCh, errCh
}

//nolint:gocyclo
func proxy(conn net.Conn, proxy io.Writer, termCh chan<- string, sizeCh chan<- WindowSize) error {
	reader := bufio.NewReaderSize(conn, 64)
	var wroteTelnetCommands bool
	var wroteTermType bool
	var willTerminalType bool
	var willNegotiateAboutWindowSize bool

	// Gets Telnet to send option negotiation commands if explicit port was given.
	// Also clears the line in case the client isn't Telnet
	// https://ibm.com/docs/zos/latest?topic=problems-telnet-commands-options
	if _, err := WriteAndClear(conn, Iac, Do, LineMode); err != nil {
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
			// https://ibm.com/docs/zos/latest?topic=problems-telnet-commands-options
			if !wroteTelnetCommands {
				wroteTelnetCommands = true
				slog.Log(context.Background(), config.LevelTrace, "Configuring Telnet")
				if _, err := Write(conn,
					Iac, Will, Echo,
					Iac, Will, SuppressGoAhead,
					Iac, Do, TerminalType,
					Iac, Do, NegotiateAboutWindowSize,
				); err != nil {
					return err
				}
			}

			if b, err = reader.ReadByte(); err != nil {
				return err
			}

			switch Operator(b) {
			case SubNegotiation:
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
						if !wroteTermType && len(command) > 4 {
							wroteTermType = true
							term := strings.ToLower(string(command[2 : len(command)-2]))
							slog.Log(context.Background(), config.LevelTrace, "Got terminal type", "type", term)
							termCh <- term
						}
					case NegotiateAboutWindowSize:
						if len(command) >= 5 {
							r := bytes.NewReader(command[1:])
							var size WindowSize
							if err := binary.Read(r, binary.BigEndian, &size); err != nil {
								return err
							}
							slog.Log(context.Background(), config.LevelTrace, "Got window size", "size", size)
							if size.Width != 0 && size.Height != 0 {
								sizeCh <- size
							}
						}
					}
				}
			case Will:
				if b, err = reader.ReadByte(); err != nil {
					return err
				}

				switch Operator(b) {
				case TerminalType:
					if !willTerminalType {
						willTerminalType = true
						slog.Log(context.Background(), config.LevelTrace, "Requesting terminal type")
						if _, err := Write(conn, Iac, SubNegotiation, TerminalType, 1, Iac, Se); err != nil {
							return err
						}
					}
				case NegotiateAboutWindowSize:
					if !willNegotiateAboutWindowSize {
						willNegotiateAboutWindowSize = true
						if _, err := Write(conn, Iac, Do, NegotiateAboutWindowSize); err != nil {
							return err
						}
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
