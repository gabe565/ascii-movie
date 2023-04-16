package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/gabe565/ascii-movie/internal/log_hooks"
	"github.com/gabe565/ascii-movie/internal/movie"
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

	sshOptions := []ssh.Option{
		wish.WithAddress(s.Address),
		wish.WithMiddleware(
			s.ServeSSH(m),
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

func (s *SSHServer) ServeSSH(m *movie.Movie) wish.Middleware {
	return func(handle ssh.Handler) ssh.Handler {
		return func(session ssh.Session) {
			remoteIP := RemoteIp(session.RemoteAddr().String())
			durationHook := log_hooks.NewDuration()
			sessionLog := s.Log.WithFields(log.Fields{
				"remote_ip": remoteIP,
				"user":      session.User(),
				"duration":  durationHook,
			})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go HandleInput(ctx, cancel, session, nil)

			level := log.InfoLevel
			var status StreamStatus

			if err := m.Stream(ctx, session); err == nil {
				status = StreamSuccess
			} else {
				if errors.Is(err, context.Canceled) {
					if remoteIP == s.DefaultGateway || time.Since(durationHook.GetStart()) < s.LogExcludeFaster {
						level = log.TraceLevel
					}
					status = StreamDisconnect
				} else {
					sessionLog = sessionLog.WithError(err)
					level = log.ErrorLevel
					status = StreamFailed
				}
			}

			sessionLog.Log(level, status)
			handle(session)
		}
	}
}
