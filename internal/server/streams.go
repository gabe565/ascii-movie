package server

import (
	"sync"
	"time"
)

type Stream struct {
	RemoteIp  string    `json:"remote_ip"`
	Connected time.Time `json:"connected"`
}

func NewStreamList() StreamList {
	return StreamList{
		streams: make(map[string]Stream),
	}
}

type StreamList struct {
	streams map[string]Stream
	mu      sync.Mutex
}

func (s *StreamList) Connect(remoteIp string) (ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.streams[remoteIp]; ok {
		return false
	}

	s.streams[remoteIp] = Stream{
		RemoteIp:  remoteIp,
		Connected: time.Now(),
	}
	return true
}

func (s *StreamList) Disconnect(remoteIp string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.streams, remoteIp)
}

func (s *StreamList) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.streams)
}

func (s *StreamList) Streams() []Stream {
	result := make([]Stream, 0, s.Len())
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, stream := range s.streams {
		result = append(result, stream)
	}
	return result
}
