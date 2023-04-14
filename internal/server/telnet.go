package server

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/gabe565/ascii-movie/internal/log_hooks"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/jackpal/gateway"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
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

	logExcludeGateway, err := flags.GetBool(LogExcludeGatewayFlag)
	if err != nil {
		panic(err)
	}
	if logExcludeGateway {
		if defaultGateway, err := gateway.DiscoverGateway(); err == nil {
			telnet.DefaultGateway = defaultGateway.String()
		} else {
			telnet.Log.Warn("Failed to discover default gateway")
		}
	}

	telnet.LogExcludeFaster, err = flags.GetDuration(LogExcludeFaster)
	if err != nil {
		panic(err)
	}

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
			status = StreamDisconnect
			if remoteIP == t.DefaultGateway || time.Since(durationHook.GetStart()) < t.LogExcludeFaster {
				level = log.TraceLevel
			}
		}
	}

	sessionLog.Log(level, status)
}
