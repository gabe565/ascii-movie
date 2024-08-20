//go:build !gzip

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/gabe565/ascii-movie/movies"
)

func main() {
	config.InitLog(os.Stderr, slog.LevelInfo, config.FormatAuto)

	if err := fs.WalkDir(movies.Movies, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		outPath := filepath.Join("movies", path+".gz")
		slog.Debug("Create output", "path", outPath)
		out, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create output: %w", err)
		}

		slog.Debug("Open input", "path", path)
		in, err := movies.Movies.Open(path)
		if err != nil {
			return fmt.Errorf("open input: %w", err)
		}

		slog.Debug("Copy input to gzip writer")
		gz := gzip.NewWriter(out)
		if _, err := io.Copy(gz, in); err != nil {
			return fmt.Errorf("copy input to gzip writer: %w", err)
		}

		if err := gz.Close(); err != nil {
			return fmt.Errorf("close gzip writer: %w", err)
		}

		slog.Debug("Close output")
		if err := out.Close(); err != nil {
			return fmt.Errorf("close output: %w", err)
		}

		return nil
	}); err != nil {
		slog.Error("Failed to gzip movies", "error", err)
		os.Exit(1)
	}
}
