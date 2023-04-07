package server

import (
	"context"
	"errors"
	"github.com/gabe565/ascii-movie/internal/log_hooks"
	log "github.com/sirupsen/logrus"
	"net"
	"syscall"
)

var telnetLog = log.WithField("server", "telnet")

func (s *Handler) ListenTelnet(ctx context.Context) error {
	telnetLog.WithField("address", s.TelnetConfig.Address).Info("Starting Telnet server")

	listen, err := net.Listen("tcp", s.TelnetConfig.Address)
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
					telnetLog.WithError(err).Error("Failed to accept connection")
					continue
				}
			}

			go s.Serve(conn)
		}
	}()

	<-ctx.Done()
	telnetLog.Info("Stopping Telnet server")
	defer telnetLog.Info("Stopped Telnet server")
	return listen.Close()
}

func (s *Handler) Serve(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	sessionLog := telnetLog.WithFields(log.Fields{
		"duration": log_hooks.NewDuration(),
	})

	remoteIP := RemoteIp(conn.RemoteAddr().String())
	sessionLog = sessionLog.WithField("remote_ip", remoteIP)

	go func() {
		// Exit on user input
		b := make([]byte, 1)
		_, _ = conn.Read(b)
		if err := conn.Close(); err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			} else {
				log.WithError(err).Warn("failed to close session on user input")
			}
		}
		sessionLog.Info("Disconnected early")
	}()

	if err := s.ServeAscii(conn); err != nil {
		if !errors.Is(err, net.ErrClosed) && !errors.Is(err, syscall.EPIPE) {
			sessionLog.WithError(err).Error("Failed to serve")
		}
		return
	}

	sessionLog.Info("Finished movie")
}
