package main

import (
	log "github.com/sirupsen/logrus"
	"net"
)

//go:generate go run ./internal/cmd/generate_frames

func main() {
	listen, err := net.Listen("tcp", ":23")
	if err != nil {
		log.Fatal(err)
	}
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.WithError(err).Error("failed to accept connection")
			continue
		}

		go Serve(conn)
	}
}
