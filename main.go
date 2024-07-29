package main

import (
	"os"

	"github.com/gabe565/ascii-movie/cmd"
)

//go:generate go run ./internal/generate/gzip

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
