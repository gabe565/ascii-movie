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

type TelnetServer struct {
	Server
}

func NewTelnet(flags *flag.FlagSet) TelnetServer {
	return TelnetServer{Server: NewServer(flags, TelnetFlagPrefix)}
}

func (s *TelnetServer) Listen(ctx context.Context, m *movie.Movie) error {
	s.Log.WithField("address", s.Address).Info("Starting Telnet server")

	listen, err := net.Listen("tcp", s.Address)
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
					s.Log.WithError(err).Error("Failed to accept connection")
					continue
				}
			}

			go s.ServeTelnet(conn, m)
		}
	}()

	<-ctx.Done()
	s.Log.Info("Stopping Telnet server")
	defer s.Log.Info("Stopped Telnet server")
	return listen.Close()
}

func (s *TelnetServer) ServeTelnet(conn net.Conn, m *movie.Movie) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	remoteIP := RemoteIp(conn.RemoteAddr().String())
	durationHook := log_hooks.NewDuration()
	sessionLog := s.Log.WithFields(log.Fields{
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
			if remoteIP == s.DefaultGateway || time.Since(durationHook.GetStart()) < s.LogExcludeFaster {
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
