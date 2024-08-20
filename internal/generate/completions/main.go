package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gabe565/ascii-movie/cmd"
	"github.com/gabe565/ascii-movie/internal/config"
)

func main() {
	config.InitLog(os.Stderr, slog.LevelInfo, config.FormatAuto)

	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	if err := os.RemoveAll("completions"); err != nil {
		return fmt.Errorf("failed to remove completions dir: %w", err)
	}

	if err := os.MkdirAll("completions", 0o777); err != nil {
		return fmt.Errorf("failed to create completions dir: %w", err)
	}

	rootCmd := cmd.NewCommand()
	name := rootCmd.Name()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	for _, shell := range []string{"bash", "zsh", "fish"} {
		rootCmd.SetArgs([]string{"completion", shell})
		if err := rootCmd.Execute(); err != nil {
			return fmt.Errorf("failed to generate completion: %w", err)
		}

		err := os.WriteFile(filepath.Join("completions", name+"."+shell), buf.Bytes(), 0o644)
		if err != nil {
			return fmt.Errorf("failed to write completion: %w", err)
		}
	}

	return nil
}
