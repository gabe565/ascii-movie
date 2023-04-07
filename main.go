package main

import (
	"github.com/gabe565/ascii-movie/cmd"
	"os"
)

//go:generate go run -tags generate ./internal/cmd/generate_movie

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
