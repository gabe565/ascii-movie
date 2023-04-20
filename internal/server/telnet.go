package server

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type TelnetServer struct {
	Server
}

func NewTelnet(flags *flag.FlagSet) TelnetServer {
	return TelnetServer{Server: NewServer(flags, TelnetFlagPrefix)}
}

func (s *TelnetServer) Listen(ctx context.Context, m *movie.Movie) error {
	s.Log.WithField("address", s.Address).Info("Starting Telnet server")

	lipgloss.SetColorProfile(termenv.ANSI256)

	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					s.Log.WithError(err).Error("Failed to accept connection")
					continue
				}
			}

			go s.Handler(conn, m)
		}
	}()

	<-ctx.Done()
	s.Log.Info("Stopping Telnet server")
	defer s.Log.Info("Stopped Telnet server")
	return listen.Close()
}

func (s *TelnetServer) Handler(conn net.Conn, m *movie.Movie) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	remoteIP := RemoteIp(conn.RemoteAddr().String())
	logger := s.Log.WithField("remote_ip", remoteIP)

	inR, inW := io.Pipe()
	outR, outW := io.Pipe()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	player := movie.NewPlayer(m, logger)
	player.LogExcludeFaster = s.LogExcludeFaster
	program := tea.NewProgram(
		player,
		tea.WithInput(inR),
		tea.WithOutput(outW),
		tea.WithContext(ctx),
	)

	go func() {
		// Proxy output to client
		_, _ = io.Copy(conn, outR)
		program.Send(movie.Quit())
	}()

	go func() {
		// Proxy input to program
		_ = proxyTelnetInput(ctx, conn, inW)
		program.Send(movie.Quit())
	}()

	if _, err := program.Run(); err != nil && !errors.Is(err, tea.ErrProgramKilled) {
		logger.WithError(err).Error("Stream failed")
	}
}

func proxyTelnetInput(ctx context.Context, conn io.ReadWriter, proxy io.Writer) error {
	buf := make([]byte, 64)
	var proxyBuf bytes.Buffer
	var skip uint8
	var subNegotiation bool
	var wroteTelnetCommands bool

	// Gets Telnet to send option negotiation commands if explicit port was given.
	// IAC DO LINEMODE (Then clear the line in case client isn't Telnet)
	// https://ibm.com/docs/zos/2.5.0?topic=problems-telnet-commands-options
	if _, err := conn.Write([]byte("\xFF\xFD\x22\r\x1B[K")); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			n, err := conn.Read(buf)
			if err != nil {
				return err
			}

			for _, c := range buf[:n] {
				switch c {
				case 0xFF:
					// IAC DO LINEMODE IAC WILL Echo IAC WILL Suppress Go Ahead
					// https://ibm.com/docs/zos/2.5.0?topic=problems-telnet-commands-options
					if conn != nil && !wroteTelnetCommands {
						log.Trace("Writing Telnet commands")
						if _, err := conn.Write([]byte("\xFF\xFD\x22\xFF\xFB\x01\xFF\xFB\x03")); err != nil {
							log.WithError(err).Error("Failed to write Telnet commands")
						}
						wroteTelnetCommands = true
					}
					skip = 3
				case 0xFA:
					subNegotiation = true
				case 0xF0:
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
}
