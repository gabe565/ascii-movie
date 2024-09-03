package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/gabe565/ascii-movie/internal/player"
	"github.com/gabe565/ascii-movie/internal/util"
	"github.com/muesli/termenv"
	flag "github.com/spf13/pflag"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

//nolint:gochecknoglobals
var sshListeners uint8

type SSHServer struct {
	MovieServer
	HostKeyPath []string
	HostKeyPEM  []string
}

func NewSSH(flags *flag.FlagSet) SSHServer {
	ssh := SSHServer{MovieServer: NewMovieServer(flags, SSHFlagPrefix)}
	var err error

	if ssh.HostKeyPath, err = flags.GetStringSlice(SSHHostKeyPathFlag); err != nil {
		panic(err)
	}

	if ssh.HostKeyPEM, err = flags.GetStringSlice(SSHHostKeyDataFlag); err != nil {
		panic(err)
	}

	if len(ssh.HostKeyPath) == 0 && len(ssh.HostKeyPEM) == 0 {
		ssh.HostKeyPath = []string{"$HOME/.ssh/ascii_movie_ed25519", "$HOME/.ssh/ascii_movie_rsa"}
	}

	return ssh
}

func (s *SSHServer) Listen(ctx context.Context, m *movie.Movie) error {
	s.Log.Info("Starting SSH server", "address", s.Address)

	sshOptions := []ssh.Option{
		wish.WithAddress(s.Address),
		wish.WithIdleTimeout(idleTimeout),
		wish.WithMaxTimeout(maxTimeout),
		wish.WithMiddleware(
			bubbletea.Middleware(s.Handler(m)),
			s.TrackStream,
		),
	}

	for _, pem := range s.HostKeyPEM {
		sshOptions = append(sshOptions, wish.WithHostKeyPEM([]byte(pem)))
	}

	for _, path := range s.HostKeyPath {
		if strings.Contains(path, "$HOME") {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			path = strings.ReplaceAll(path, "$HOME", home)
		}
		sshOptions = append(sshOptions, wish.WithHostKeyPath(path))
	}

	server, err := wish.NewServer(sshOptions...)
	if err != nil {
		return err
	}

	for _, signer := range server.HostSigners {
		s.Log.Debug("Using host key",
			"type", signer.PublicKey().Type(),
			"fingerprint", gossh.FingerprintSHA256(signer.PublicKey()),
		)
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		sshListeners++
		defer func() {
			sshListeners--
		}()

		if err = server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		s.Log.Info("Stopping SSH server")
		defer s.Log.Info("Stopped SSH server")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		return server.Shutdown(shutdownCtx)
	})

	return group.Wait()
}

func (s *SSHServer) Handler(m *movie.Movie) bubbletea.Handler {
	return func(session ssh.Session) (tea.Model, []tea.ProgramOption) {
		remoteIP := RemoteIP(session.RemoteAddr())
		logger := s.Log.With(
			"remoteIP", remoteIP,
			"user", session.User(),
		)

		renderer := bubbletea.MakeRenderer(session)
		if renderer.ColorProfile() == termenv.Ascii {
			if pty, _, ok := session.Pty(); ok {
				renderer.SetColorProfile(util.Profile(pty.Term))
			}
		}

		p := player.NewPlayer(m, logger, renderer)
		go func() {
			<-session.Context().Done()
			p.Close()
		}()
		return p, []tea.ProgramOption{
			tea.WithFPS(30),
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		}
	}
}

func (s *SSHServer) TrackStream(handler ssh.Handler) ssh.Handler {
	return func(session ssh.Session) {
		remoteIP := RemoteIP(session.RemoteAddr())
		id, err := serverInfo.StreamConnect("ssh", remoteIP)
		if err != nil {
			s.Log.Error("Failed to begin stream",
				"remoteIP", remoteIP,
				"user", session.User(),
			)
			_, _ = session.Write([]byte(ErrorText(err) + "\n"))
			return
		}
		defer serverInfo.StreamDisconnect(id)
		handler(session)
	}
}
