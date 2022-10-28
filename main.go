package main

import (
	"github.com/gabe565/ascii-telnet-go/cmd"
	"os"
)

//go:generate go run ./internal/cmd/generate_frames

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
