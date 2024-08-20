package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gabe565/ascii-movie/cmd"
	"github.com/gabe565/ascii-movie/cmd/util"
	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/spf13/cobra/doc"
)

const output = "./docs"

func main() {
	config.InitLog(os.Stderr, slog.LevelInfo, config.FormatAuto)

	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	if err := os.RemoveAll(output); err != nil {
		return fmt.Errorf("failed to remove existing dir: %w", err)
	}

	if err := os.MkdirAll(output, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	rootCmd := cmd.NewCommand(util.WithVersion("beta"))
	if err := doc.GenMarkdownTree(rootCmd, output); err != nil {
		return fmt.Errorf("failed to generate docs: %w", err)
	}

	return nil
}
