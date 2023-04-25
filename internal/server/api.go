package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var streamList = NewStreamList()

type ApiServer struct {
	Server
	TelnetEnabled bool
	SSHEnabled    bool
}

func NewApi(flags *flag.FlagSet) ApiServer {
	return ApiServer{Server: NewServer(flags, ApiFlagPrefix)}
}

func (s *ApiServer) Listen(ctx context.Context) error {
	s.Log.WithField("address", s.Address).Info("Starting API server")

	http.HandleFunc("/health", s.Health)
	http.HandleFunc("/streams", s.Streams)
	server := http.Server{Addr: s.Address}
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

type HealthResponse struct {
	Healthy bool `json:"healthy"`
	SSH     bool `json:"ssh"`
	Telnet  bool `json:"telnet"`
}

func (s *ApiServer) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Telnet: telnetListeners == 1,
		SSH:    sshListeners == 1,
	}
	response.Healthy = (!s.SSHEnabled || response.SSH) && (!s.TelnetEnabled || response.Telnet)

	buf, err := json.Marshal(response)
	if err != nil {
		s.Log.WithError(err).Error("Failed to marshal API response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Healthy {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	if _, err := w.Write([]byte(buf)); err != nil {
		s.Log.WithError(err).Error("Failed to write API response")
	}
}

type StreamsResponse struct {
	Count int `json:"count"`
}

func (s *ApiServer) Streams(w http.ResponseWriter, r *http.Request) {
	response := StreamsResponse{Count: streamList.Len()}

	buf, err := json.Marshal(response)
	if err != nil {
		s.Log.WithError(err).Error("Failed to marshal API response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(buf)); err != nil {
		s.Log.WithError(err).Error("Failed to write API response")
	}
}
