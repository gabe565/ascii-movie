package server

import (
	"context"
	"errors"
	"github.com/gabe565/ascii-telnet-go/internal/log_hooks"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"net"
	"syscall"
)

func (s *Handler) Listen(ctx context.Context, addr string) error {
	log.WithField("address", addr).Info("listening for connections")

	listen, err := net.Listen("tcp", addr)
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
					log.WithError(err).Error("failed to accept connection")
					continue
				}
			}

			go s.Serve(conn)
		}
	}()

	<-ctx.Done()
	log.Info("closing server")
	return listen.Close()
}

func (s *Handler) Serve(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	log := logrus.WithField("duration", log_hooks.NewDuration())

	remoteIP, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		remoteIP = conn.RemoteAddr().String()
	}
	log = log.WithField("remote_ip", remoteIP)

	if err := s.ServeAscii(conn); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			log.Info("disconnected early")
		} else {
			log.WithError(err).Error("failed to serve")
		}
		return
	}

	log.Info("finished movie")
}
