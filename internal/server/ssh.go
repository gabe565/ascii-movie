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
	"github.com/jackpal/gateway"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

type SSH struct {
	Config
	HostKeyPath []string
	HostKeyPEM  []string
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

	if ssh.HostKeyPath, err = flags.GetStringSlice(SSHHostKeyPathFlag); err != nil {
		panic(err)
	}

	if ssh.HostKeyPEM, err = flags.GetStringSlice(SSHHostKeyDataFlag); err != nil {
		panic(err)
	}

	ssh.Log = log.WithField("server", "ssh")

	logExcludeGateway, err := flags.GetBool(LogExcludeGatewayFlag)
	if err != nil {
		panic(err)
	}
	if logExcludeGateway {
		if defaultGateway, err := gateway.DiscoverGateway(); err == nil {
			ssh.DefaultGateway = defaultGateway.String()
		} else {
			ssh.Log.Warn("Failed to discover default gateway")
		}
	}

	ssh.LogExcludeFaster, err = flags.GetDuration(LogExcludeFaster)
	if err != nil {
		panic(err)
	}

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

func (s *SSH) ServeSSH(m *movie.Movie) wish.Middleware {
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

			if err := m.Stream(ctx, session); err == nil {
				sessionLog.Info("Finished movie")
			} else {
				if errors.Is(err, context.Canceled) {
					switch {
					case remoteIP == s.DefaultGateway,
						time.Since(durationHook.GetStart()) < s.LogExcludeFaster:
						sessionLog.Trace("Disconnected early")
					default:
						sessionLog.Info("Disconnected early")
					}
				}
				return
			}

			handle(session)
		}
	}
}
