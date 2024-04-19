package server

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Stream struct {
	Server    string    `json:"server"`
	RemoteIP  string    `json:"remote_ip"`
	Connected time.Time `json:"connected"`
}

func NewInfo() Info {
	return Info{
		streams:    make(map[uint]Stream, 64),
		concurrent: make(map[string]uint, 64),

		activeConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "ascii_movie",
				Name:      "connections_active",
				Help:      "Count of active connections",
			},
			[]string{"server"},
		),
		totalConnections: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "ascii_movie",
				Name:      "connections_total",
				Help:      "Total connections",
			},
			[]string{"server"},
		),
		rateLimitedConnections: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "ascii_movie",
				Name:      "rate_limited_connections_total",
				Help:      "Total number of rate limited connections",
			},
			[]string{"server"},
		),
		connectionDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "ascii_movie",
			Name:      "connection_duration_seconds",
			Help:      "Connection duration in seconds",
			Buckets: []float64{
				100 * time.Millisecond.Seconds(),
				1 * time.Second.Seconds(),
				10 * time.Second.Seconds(),
				30 * time.Second.Seconds(),
				1 * time.Minute.Seconds(),
				3 * time.Minute.Seconds(),
				6 * time.Minute.Seconds(),
				9 * time.Minute.Seconds(),
				12 * time.Minute.Seconds(),
				15 * time.Minute.Seconds(),
				18 * time.Minute.Seconds(),
			},
		}, []string{"server"}),
	}
}

type Info struct {
	streams    map[uint]Stream
	totalCount atomic.Uint32
	concurrent map[string]uint
	nextID     uint
	mu         sync.Mutex

	activeConnections      *prometheus.GaugeVec
	totalConnections       *prometheus.CounterVec
	rateLimitedConnections *prometheus.CounterVec
	connectionDuration     *prometheus.HistogramVec
}

var ErrRateLimited = errors.New("rate limited")

func (s *Info) StreamConnect(server, remoteIP string) (uint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	prometheusLabels := prometheus.Labels{"server": server}

	if concurrentStreams != 0 && s.concurrent[remoteIP]+1 > concurrentStreams {
		s.rateLimitedConnections.With(prometheusLabels).Inc()
		return 0, ErrRateLimited
	}

	s.totalCount.Add(1)
	s.activeConnections.With(prometheusLabels).Inc()
	s.totalConnections.With(prometheusLabels).Inc()

	defer func() {
		s.nextID++
	}()
	s.streams[s.nextID] = Stream{
		Server:    server,
		RemoteIP:  remoteIP,
		Connected: time.Now(),
	}
	s.concurrent[remoteIP]++
	return s.nextID, nil
}

func (s *Info) StreamDisconnect(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, ok := s.streams[id]
	if !ok {
		return
	}

	prometheusLabels := prometheus.Labels{"server": stream.Server}

	s.connectionDuration.With(prometheusLabels).
		Observe(time.Since(stream.Connected).Seconds())

	s.concurrent[stream.RemoteIP]--
	if s.concurrent[stream.RemoteIP] == 0 {
		delete(s.concurrent, stream.RemoteIP)
	}
	delete(s.streams, id)

	s.activeConnections.With(prometheusLabels).Dec()
}

func (s *Info) NumActive() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.streams)
}

func (s *Info) GetStreams() []Stream {
	result := make([]Stream, 0, s.NumActive())
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, stream := range s.streams {
		result = append(result, stream)
	}
	return result
}

func ErrorText(err error) string {
	if errors.Is(err, ErrRateLimited) {
		return "409: Too many concurrent streams"
	}
	return err.Error()
}
