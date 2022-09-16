package main

import (
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"net"
)

//go:generate go run ./internal/cmd/generate_frames

func main() {
	var addr string
	flag.StringVarP(&addr, "address", "a", ":23", "Listen address")

	flag.Parse()

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	log.WithField("address", addr).Info("listening for connections")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.WithError(err).Error("failed to accept connection")
			continue
		}

		go Serve(conn)
	}
}
