package server

import (
	"context"
	"errors"
	"github.com/gabe565/ascii-movie/internal/log_hooks"
	"github.com/gabe565/ascii-movie/internal/movie"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"net"
	"syscall"
)

type Telnet Config

func NewTelnet(flags *flag.FlagSet) Telnet {
	var telnet Telnet
	var err error

	telnet.Enabled, err = flags.GetBool(TelnetEnabledFlag)
	if err != nil {
		panic(err)
	}

	telnet.Address, err = flags.GetString(TelnetAddressFlag)
	if err != nil {
		panic(err)
	}

	telnet.Log = log.WithField("server", "telnet")

	return telnet
}

func (t *Telnet) Listen(ctx context.Context, m *movie.Movie) error {
	t.Log.WithField("address", t.Address).Info("Starting Telnet server")

	listen, err := net.Listen("tcp", t.Address)
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
					t.Log.WithError(err).Error("Failed to accept connection")
					continue
				}
			}

			go t.ServeTelnet(conn, m)
		}
	}()

	<-ctx.Done()
	t.Log.Info("Stopping Telnet server")
	defer t.Log.Info("Stopped Telnet server")
	return listen.Close()
}

func (t *Telnet) ServeTelnet(conn net.Conn, m *movie.Movie) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	sessionLog := t.Log.WithFields(log.Fields{
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

	if err := m.Stream(conn); err != nil {
		if !errors.Is(err, net.ErrClosed) && !errors.Is(err, syscall.EPIPE) {
			sessionLog.WithError(err).Error("Failed to serve")
		}
		return
	}

	sessionLog.Info("Finished movie")
}
