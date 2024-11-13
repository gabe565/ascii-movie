package main

import (
	"log/slog"
	"os"

	"gabe565.com/ascii-movie/cmd"
	"gabe565.com/utils/cobrax"
)

//go:generate go run ./internal/generate/gzip

var version = "beta"

func main() {
	root := cmd.NewCommand(cobrax.WithVersion(version))
	if err := root.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
