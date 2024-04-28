//go:build !gzip

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gabe565/ascii-movie/movies"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := fs.WalkDir(movies.Movies, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		outPath := filepath.Join("movies", path+".gz")
		log.Debug().Str("path", outPath).Msg("Create output")
		out, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create output: %w", err)
		}

		log.Debug().Str("path", path).Msg("Open input")
		in, err := movies.Movies.Open(path)
		if err != nil {
			return fmt.Errorf("open input: %w", err)
		}

		log.Debug().Msg("Copy input to gzip writer")
		gz := gzip.NewWriter(out)
		if _, err := io.Copy(gz, in); err != nil {
			return fmt.Errorf("copy input to gzip writer: %w", err)
		}

		if err := gz.Close(); err != nil {
			return fmt.Errorf("close gzip writer: %w", err)
		}

		log.Debug().Msg("Close output")
		if err := out.Close(); err != nil {
			return fmt.Errorf("close output: %w", err)
		}

		return nil
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to gzip movies")
	}
}
