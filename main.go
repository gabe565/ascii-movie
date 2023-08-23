package main

import (
	"os"

	"github.com/gabe565/ascii-movie/cmd"
)

//go:generate go run ./internal/generate/gzip

var (
	version = "beta"
	commit  = ""
)

func main() {
	if err := cmd.NewCommand(version, commit).Execute(); err != nil {
		os.Exit(1)
	}
}
