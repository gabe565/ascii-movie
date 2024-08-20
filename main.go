package main

import (
	"os"

	"github.com/gabe565/ascii-movie/cmd"
	"github.com/gabe565/ascii-movie/cmd/util"
)

//go:generate go run ./internal/generate/gzip

var version = "beta"

func main() {
	root := cmd.NewCommand(util.WithVersion(version))
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
