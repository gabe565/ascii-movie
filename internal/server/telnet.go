package server

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"gabe565.com/ascii-movie/internal/movie"
	"gabe565.com/ascii-movie/internal/player"
	"gabe565.com/ascii-movie/internal/server/idleconn"
	"gabe565.com/ascii-movie/internal/server/telnet"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	flag "github.com/spf13/pflag"
)

//nolint:gochecknoglobals
var telnetListeners uint8

type TelnetServer struct {
	MovieServer
}

func NewTelnet(flags *flag.FlagSet) TelnetServer {
	return TelnetServer{MovieServer: NewMovieServer(flags, TelnetFlagPrefix)}
}

func (s *TelnetServer) Listen(ctx context.Context, m *movie.Movie) error {
	s.Log.Info("Starting telnet server", "address", s.Address)

	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	var serveGroup sync.WaitGroup
	serveCtx, serveCancel := context.WithCancel(context.Background())
	defer serveCancel()

	go func() {
		telnetListeners++
		defer func() {
			telnetListeners--
		}()

		for {
			conn, err := listen.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					s.Log.Error("Failed to accept connection", "error", err)
					continue
				}
			}

			serveGroup.Add(1)
			go func() {
				defer serveGroup.Done()
				ctx, cancel := context.WithCancel(serveCtx)
				conn = idleconn.New(conn, idleTimeout, maxTimeout, cancel)
				s.Handler(ctx, conn, m)
			}()
		}
	}()

	<-ctx.Done()
	s.Log.Info("Stopping Telnet server")
	defer s.Log.Info("Stopped Telnet server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	go func() {
		serveCancel()
		serveGroup.Wait()
		shutdownCancel()
	}()
	<-shutdownCtx.Done()

	return listen.Close()
}

func (s *TelnetServer) Handler(ctx context.Context, conn net.Conn, m *movie.Movie) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	remoteIP := RemoteIP(conn.RemoteAddr())
	logger := s.Log.With("remoteIP", remoteIP)

	id, err := serverInfo.StreamConnect("telnet", remoteIP)
	if err != nil {
		logger.Error("Failed to begin stream", "error", err)
		_, _ = conn.Write([]byte(ErrorText(err) + "\n"))
		return
	}
	defer serverInfo.StreamDisconnect(id)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	in, profile, sizeCh, errCh := telnet.Proxy(conn)
	defer func() {
		_ = in.Close()
	}()

	gotProfile := profile != -1
	if !gotProfile {
		profile = termenv.ANSI256
	}

	p := player.NewPlayer(m, logger, telnet.MakeRenderer(conn, profile))
	defer p.Close()

	opts := []tea.ProgramOption{
		tea.WithInput(in),
		tea.WithOutput(conn),
		tea.WithFPS(30),
	}
	if gotProfile {
		opts = append(opts, tea.WithAltScreen(), tea.WithMouseCellMotion())
	}
	program := tea.NewProgram(p, opts...)

	go func() {
		for {
			select {
			case <-ctx.Done():
				program.Quit()
				return
			case <-errCh:
				cancel()
			case info := <-sizeCh:
				program.Send(tea.WindowSizeMsg{
					Width:  int(info.Width),
					Height: int(info.Height),
				})
			}
		}
	}()

	if _, err := program.Run(); err != nil && !errors.Is(err, tea.ErrProgramKilled) {
		logger.Error("Program failed", "error", err)
	}

	// p.Kill() will force kill the program if it's still running,
	// and restore the terminal to its original state in case of a
	// tui crash
	program.Kill()
}
