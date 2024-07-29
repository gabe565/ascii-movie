package main

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/gabe565/ascii-movie/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := os.RemoveAll("completions"); err != nil {
		log.Fatal().Err(err).Msg("failed to remove completions dir")
	}

	if err := os.MkdirAll("completions", 0o777); err != nil {
		log.Fatal().Err(err).Msg("failed to create completions dir")
	}

	rootCmd := cmd.NewCommand()
	name := rootCmd.Name()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	for _, shell := range []string{"bash", "zsh", "fish"} {
		rootCmd.SetArgs([]string{"completion", shell})
		if err := rootCmd.Execute(); err != nil {
			log.Fatal().Err(err).Msg("failed to generate completion")
		}

		err := os.WriteFile(filepath.Join("completions", name+"."+shell), buf.Bytes(), 0o644)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to write completion")
		}
	}
}
