package server

import (
	"context"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type ApiServer struct {
	Server
}

func NewApi(flags *flag.FlagSet) ApiServer {
	return ApiServer{Server: NewServer(flags, ApiFlagPrefix)}
}

func (s *ApiServer) Listen(ctx context.Context) error {
	s.Log.WithField("address", s.Address).Info("Starting API server")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.Status)
	server := http.Server{Addr: s.Address, Handler: mux}
	go func() {
		<-ctx.Done()
		s.Log.Info("Stopping API server")
		defer s.Log.Info("Stopped API server")
		if err := server.Close(); err != nil {
			log.WithError(err).Error("Failed to close server")
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *ApiServer) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"healthy":true}`))
}
