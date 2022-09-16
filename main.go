package main

import "net"

//go:generate go run ./internal/cmd/generate_frames

func main() {
	listen, err := net.Listen("tcp", ":23")
	if err != nil {
		panic(err)
	}
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		go Serve(conn)
	}
}
