package main

import (
	"github.com/reiver/go-telnet"
)

//go:generate go run ./internal/cmd/generate_frames

func main() {
	handler := AsciiHandler{}
	err := telnet.ListenAndServe(":23", handler)
	if err != nil {
		panic(err)
	}
}
