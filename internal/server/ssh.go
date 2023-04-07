package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/gabe565/ascii-movie/internal/log_hooks"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io"
)

var sshLog = log.WithField("server", "ssh")

func (s *Handler) ListenSSH(ctx context.Context) error {
	sshLog.WithField("address", s.SSHConfig.Address).Info("Starting SSH server")

	server, err := wish.NewServer(
		wish.WithAddress(s.SSHConfig.Address),
		wish.WithMiddleware(
			s.HandleSSH,
		),
	)
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
		sshLog.Info("Stopping SSH server")
		defer sshLog.Info("Stopped SSH server")
		return server.Close()
	})

	return group.Wait()
}

func (s *Handler) HandleSSH(handle ssh.Handler) ssh.Handler {
	return func(session ssh.Session) {
		remoteIP := RemoteIp(session.RemoteAddr().String())

		sessionLog := sshLog.WithFields(log.Fields{
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

		err := s.ServeAscii(session)
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
