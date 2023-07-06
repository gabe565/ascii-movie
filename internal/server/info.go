package server

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Stream struct {
	Server    string    `json:"server"`
	RemoteIp  string    `json:"remote_ip"`
	Connected time.Time `json:"connected"`
}

func NewInfo() Info {
	return Info{
		streams:    make(map[uint]Stream, 64),
		concurrent: make(map[string]uint, 64),

		activeConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "ascii_movie_active_connections",
				Help: "Count of active connections",
			},
		),
		totalConnections: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "ascii_movie_total_connections",
				Help: "Total connections",
			},
		),
	}
}

type Info struct {
	streams    map[uint]Stream
	totalCount atomic.Uint32
	concurrent map[string]uint
	nextId     uint
	mu         sync.Mutex

	activeConnections prometheus.Gauge
	totalConnections  prometheus.Counter
}

func (s *Info) StreamConnect(server, remoteIp string) (id, concurrent uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.totalCount.Add(1)
	s.activeConnections.Inc()
	s.totalConnections.Inc()

	defer func() {
		s.nextId += 1
	}()
	s.streams[s.nextId] = Stream{
		Server:    server,
		RemoteIp:  remoteIp,
		Connected: time.Now(),
	}
	s.concurrent[remoteIp] += 1
	return s.nextId, s.concurrent[remoteIp]
}

func (s *Info) StreamDisconnect(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, ok := s.streams[id]
	if !ok {
		return
	}

	s.concurrent[stream.RemoteIp] -= 1
	if s.concurrent[stream.RemoteIp] == 0 {
		delete(s.concurrent, stream.RemoteIp)
	}
	delete(s.streams, id)

	s.activeConnections.Dec()
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
