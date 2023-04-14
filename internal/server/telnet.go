package server

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/gabe565/ascii-movie/internal/log_hooks"
	"github.com/gabe565/ascii-movie/internal/movie"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type Telnet struct {
	Config
}

func NewTelnet(flags *flag.FlagSet) Telnet {
	return Telnet{Config: NewConfig(flags, TelnetFlagPrefix)}
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
	durationHook := log_hooks.NewDuration()
	sessionLog := t.Log.WithFields(log.Fields{
		"remote_ip": remoteIP,
		"duration":  durationHook,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go HandleTelnetInput(ctx, cancel, conn, conn)

	level := log.InfoLevel
	var status StreamStatus

	if err := m.Stream(ctx, conn); err == nil {
		status = StreamSuccess
	} else {
		if errors.Is(err, context.Canceled) {
			if remoteIP == t.DefaultGateway || time.Since(durationHook.GetStart()) < t.LogExcludeFaster {
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
}
