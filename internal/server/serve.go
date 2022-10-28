package server

import (
	"errors"
	"github.com/gabe565/ascii-telnet-go/internal/log_hooks"
	"github.com/sirupsen/logrus"
	"net"
	"syscall"
)

func Serve(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	log := logrus.WithField("duration", log_hooks.NewDuration())

	remoteIP, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		remoteIP = conn.RemoteAddr().String()
	}
	log = log.WithField("remote_ip", remoteIP)

	if err := ServeAscii(conn); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			log.Info("disconnected early")
		} else {
			log.WithError(err).Error("failed to serve")
		}
		return
	}

	log.Info("finished movie")
}
