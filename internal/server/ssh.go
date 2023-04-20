package server

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

type SSHServer struct {
	Server
	HostKeyPath []string
	HostKeyPEM  []string
}

func NewSSH(flags *flag.FlagSet) SSHServer {
	ssh := SSHServer{Server: NewServer(flags, SSHFlagPrefix)}
	var err error

	if ssh.HostKeyPath, err = flags.GetStringSlice(SSHHostKeyPathFlag); err != nil {
		panic(err)
	}

	if ssh.HostKeyPEM, err = flags.GetStringSlice(SSHHostKeyDataFlag); err != nil {
		panic(err)
	}

	return ssh
}

func (s *SSHServer) Listen(ctx context.Context, m *movie.Movie) error {
	s.Log.WithField("address", s.Address).Info("Starting SSH server")

	lipgloss.SetColorProfile(termenv.ANSI256)

	sshOptions := []ssh.Option{
		wish.WithAddress(s.Address),
		wish.WithMiddleware(
			bubbletea.Middleware(s.Handler(m)),
		),
	}

	for _, pem := range s.HostKeyPEM {
		sshOptions = append(sshOptions, wish.WithHostKeyPEM([]byte(pem)))
	}

	for _, path := range s.HostKeyPath {
		sshOptions = append(sshOptions, wish.WithHostKeyPath(path))
	}

	server, err := wish.NewServer(sshOptions...)
	if err != nil {
		return err
	}

	for _, signer := range server.HostSigners {
		s.Log.WithFields(log.Fields{
			"type":        signer.PublicKey().Type(),
			"fingerprint": gossh.FingerprintSHA256(signer.PublicKey()),
		}).Debug("Using host key")
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		if ctx.Err() != nil {
			return nil
		}
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		s.Log.Info("Stopping SSH server")
		defer s.Log.Info("Stopped SSH server")
		return server.Close()
	})

	return group.Wait()
}

func (s *SSHServer) Handler(m *movie.Movie) bubbletea.Handler {
	return func(session ssh.Session) (tea.Model, []tea.ProgramOption) {
		remoteIP := RemoteIp(session.RemoteAddr().String())
		logger := s.Log.WithFields(log.Fields{
			"remote_ip": remoteIP,
			"user":      session.User(),
		})
		player := movie.NewPlayer(m, logger)
		player.LogExcludeFaster = s.LogExcludeFaster
		return player, []tea.ProgramOption{}
	}
}
