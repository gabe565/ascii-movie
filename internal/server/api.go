package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

//nolint:gochecknoglobals
var serverInfo = NewInfo()

type APIServer struct {
	Server
	TelnetEnabled bool
	SSHEnabled    bool
}

func NewAPI(flags *flag.FlagSet) APIServer {
	return APIServer{Server: NewServer(flags, APIFlagPrefix)}
}

func (s *APIServer) Listen(ctx context.Context) error {
	s.Log.Info().Str("address", s.Address).Msg("Starting API server")

	http.HandleFunc("/health", s.Health)
	http.HandleFunc("/streams", s.Streams)
	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:        s.Address,
		ReadTimeout: 5 * time.Second,
	}

	var group errgroup.Group

	group.Go(func() error {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		s.Log.Info().Msg("Stopping API server")
		defer s.Log.Info().Msg("Stopped API server")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		return server.Shutdown(shutdownCtx)
	})

	return group.Wait()
}

type HealthResponse struct {
	Healthy bool `json:"healthy"`
	SSH     bool `json:"ssh"`
	Telnet  bool `json:"telnet"`
}

func (s *APIServer) Health(w http.ResponseWriter, _ *http.Request) {
	response := HealthResponse{
		Telnet: telnetListeners == 1,
		SSH:    sshListeners == 1,
	}
	response.Healthy = (!s.SSHEnabled || response.SSH) && (!s.TelnetEnabled || response.Telnet)

	buf, err := json.Marshal(response)
	if err != nil {
		s.Log.Err(err).Msg("Failed to marshal API response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Healthy {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	if _, err := w.Write(buf); err != nil {
		s.Log.Err(err).Msg("Failed to write API response")
	}
}

type StreamsResponse struct {
	Active  *int      `json:"active,omitempty"`
	Total   *uint32   `json:"total,omitempty"`
	Streams *[]Stream `json:"streams,omitempty"`
}

func (s *APIServer) Streams(w http.ResponseWriter, r *http.Request) {
	var response StreamsResponse

	fields := r.URL.Query().Get("fields")

	if fields == "" || strings.Contains(fields, "total") {
		total := serverInfo.totalCount.Load()
		response.Total = &total
	}

	if fields == "" || strings.Contains(fields, "active") {
		count := serverInfo.NumActive()
		response.Active = &count
	}

	if fields == "" || strings.Contains(fields, "streams") {
		streams := serverInfo.GetStreams()
		response.Streams = &streams
	}

	buf, err := json.Marshal(response)
	if err != nil {
		s.Log.Err(err).Msg("Failed to marshal API response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(buf); err != nil {
		s.Log.Err(err).Msg("Failed to write API response")
	}
}
