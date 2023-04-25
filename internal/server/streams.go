package server

import (
	"sync"
	"time"
)

type Stream struct {
	Server    string    `json:"server"`
	RemoteIp  string    `json:"remote_ip"`
	Connected time.Time `json:"connected"`
}

func NewStreamList() StreamList {
	return StreamList{
		streams:    make(map[uint]Stream),
		concurrent: make(map[string]uint),
	}
}

type StreamList struct {
	streams    map[uint]Stream
	concurrent map[string]uint
	nextId     uint
	mu         sync.Mutex
}

func (s *StreamList) Connect(server, remoteIp string) (id, concurrent uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

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

func (s *StreamList) Disconnect(id uint) {
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
