package main

import (
	"os"

	"github.com/gabe565/ascii-movie/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra/doc"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	output := "./docs"

	if err := os.RemoveAll(output); err != nil {
		log.Fatal().Err(err).Msg("failed to remove existing dir")
	}

	if err := os.MkdirAll(output, 0o755); err != nil {
		log.Fatal().Err(err).Msg("failed to create directory")
	}

	rootCmd := cmd.NewCommand("latest", "")
	if err := doc.GenMarkdownTree(rootCmd, output); err != nil {
		log.Fatal().Err(err).Msg("failed to generate docs")
	}
}
