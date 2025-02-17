package main

import (
	"log/slog"
	"os"

	"gabe565.com/ascii-movie/cmd"
	"gabe565.com/ascii-movie/internal/config"
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/slogx"
)

//go:generate go run ./internal/generate/gzip

var version = "beta"

func main() {
	config.InitLog(os.Stderr, slogx.LevelInfo, slogx.FormatAuto)
	root := cmd.NewCommand(cobrax.WithVersion(version))
	if err := root.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
