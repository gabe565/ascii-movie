package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/gabe565/ascii-movie/internal/log_hooks"
	"github.com/gabe565/ascii-movie/internal/movie"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"io"
)

type SSH struct {
	Config
	HostKeyPEM  string
	HostKeyPath string
}

func NewSSH(flags *flag.FlagSet) SSH {
	var ssh SSH
	var err error

	if ssh.Enabled, err = flags.GetBool(SSHEnabledFlag); err != nil {
		panic(err)
	}

	if ssh.Address, err = flags.GetString(SSHAddressFlag); err != nil {
		panic(err)
	}

	if ssh.HostKeyPath, err = flags.GetString(SSHHostKeyPathFlag); err != nil {
		panic(err)
	}

	if ssh.HostKeyPEM, err = flags.GetString(SSHHostKeyFlag); err != nil {
		panic(err)
	}

	ssh.Log = log.WithField("server", "ssh")

	return ssh
}

func (s *SSH) Listen(ctx context.Context, m *movie.Movie) error {
	s.Log.WithField("address", s.Address).Info("Starting SSH server")

	sshOptions := []ssh.Option{
		wish.WithAddress(s.Address),
		wish.WithMiddleware(
			s.ServeSSH(m),
		),
	}

	switch {
	case s.HostKeyPEM != "":
		sshOptions = append(sshOptions, wish.WithHostKeyPEM([]byte(s.HostKeyPEM)))
	case s.HostKeyPath != "":
		sshOptions = append(sshOptions, wish.WithHostKeyPath(s.HostKeyPath))
	}

	server, err := wish.NewServer(sshOptions...)
	if err != nil {
		return err
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

func (s *SSH) ServeSSH(m *movie.Movie) wish.Middleware {
	return func(handle ssh.Handler) ssh.Handler {
		return func(session ssh.Session) {
			remoteIP := RemoteIp(session.RemoteAddr().String())
			sessionLog := s.Log.WithFields(log.Fields{
				"remote_ip": remoteIP,
				"duration":  log_hooks.NewDuration(),
			})

			go func() {
				// Exit on user input
				b := make([]byte, 1)
				if _, err := session.Read(b); err == nil {
					sessionLog.Info("Disconnected early")
					if err := session.Close(); err != nil {
						log.WithError(err).Warn("failed to close session on user input")
					}
				}
			}()

			err := m.Stream(session)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					sessionLog.WithError(err).Error("Failed to serve")
				}
			} else {
				sessionLog.Info("Finished movie")
			}

			handle(session)
		}
	}
}
