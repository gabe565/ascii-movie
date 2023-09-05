package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

	return ssh
}

func (s *SSHServer) Listen(ctx context.Context, m *movie.Movie) error {
	s.Log.WithField("address", s.Address).Info("Starting SSH server")

	sshOptions := []ssh.Option{
		wish.WithAddress(s.Address),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(s.Handler(m), termenv.ANSI256),
			s.TrackStream,
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

		sshListeners += 1
		defer func() {
			sshListeners -= 1
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

func (s *SSHServer) Handler(m *movie.Movie) bubbletea.ProgramHandler {
	return func(session ssh.Session) *tea.Program {
		remoteIP := RemoteIp(session.RemoteAddr().String())
		logger := s.Log.WithFields(log.Fields{
			"remote_ip": remoteIP,
			"user":      session.User(),
		})

		player := movie.NewPlayer(m, logger)
		program := tea.NewProgram(
			player,
			tea.WithInput(session),
			tea.WithOutput(session),
		)

		if timeout != 0 {
			go func() {
				timer := time.NewTimer(timeout)
				defer timer.Stop()
				select {
				case <-timer.C:
					program.Send(movie.Quit())
				case <-session.Context().Done():

				}
			}()
		}

		return program
	}
}

func (s *SSHServer) TrackStream(handler ssh.Handler) ssh.Handler {
	return func(session ssh.Session) {
		remoteIP := RemoteIp(session.RemoteAddr().String())
		id, err := serverInfo.StreamConnect("ssh", remoteIP)
		if err != nil {
			logger := s.Log.WithFields(log.Fields{
				"remote_ip": remoteIP,
				"user":      session.User(),
			})
			logger.Error(err)
			_, _ = session.Write([]byte(ErrorText(err) + "\n"))
			return
		}
		defer serverInfo.StreamDisconnect(id)
		handler(session)
	}
}
