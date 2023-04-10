package server

import (
	"context"
	"errors"
	"github.com/gabe565/ascii-movie/internal/log_hooks"
	"github.com/gabe565/ascii-movie/internal/movie"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"net"
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

	remoteIP := RemoteIp(conn.RemoteAddr().String())
	sessionLog := t.Log.WithFields(log.Fields{
		"remote_ip": remoteIP,
		"duration":  log_hooks.NewDuration(),
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Exit on user input
		b := make([]byte, 1)
		for {
			if _, err := conn.Read(b); err != nil {
				cancel()
				return
			}
		}
	}()

	if err := m.Stream(ctx, conn); err == nil {
		sessionLog.Info("Finished movie")
	} else {
		if errors.Is(err, context.Canceled) {
			sessionLog.Info("Disconnected early")
		}
		return
	}
}
