package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabe565/ascii-movie/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra/doc"
	flag "github.com/spf13/pflag"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flags := flag.NewFlagSet("", flag.ContinueOnError)

	var version string
	flags.StringVar(&version, "version", "beta", "Version")

	var dateParam string
	flags.StringVar(&dateParam, "date", time.Now().Format(time.RFC3339), "Build date")

	if err := flags.Parse(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to parse arguments")
	}

	if err := os.RemoveAll("manpages"); err != nil {
		log.Fatal().Err(err).Msg("failed to remove manpages dir")
	}

	if err := os.MkdirAll("manpages", 0o755); err != nil {
		log.Fatal().Err(err).Msg("failed to create manpages dir")
	}

	rootCmd := cmd.NewCommand("beta", "")
	rootName := rootCmd.Name()

	date, err := time.Parse(time.RFC3339, dateParam)
	if err != nil {
		log.Fatal().Err(err).Str("raw", dateParam).Msg("failed to parse date")
	}

	header := doc.GenManHeader{
		Title:   strings.ToUpper(rootName),
		Section: "1",
		Date:    &date,
		Source:  rootName + " " + version,
		Manual:  "User Commands",
	}

	if err := doc.GenManTree(rootCmd, &header, "manpages"); err != nil {
		log.Fatal().Err(err).Msg("failed to generate manpages")
	}

	if err := filepath.Walk("manpages", func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open input: %w", err)
		}

		out, err := os.Create(path + ".gz")
		if err != nil {
			return fmt.Errorf("create output: %w", err)
		}
		gz := gzip.NewWriter(out)

		if _, err := io.Copy(gz, in); err != nil {
			return fmt.Errorf("copy input to gzip writer: %w", err)
		}

		if err := in.Close(); err != nil {
			return fmt.Errorf("close input: %w", err)
		}
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("remove input: %w", err)
		}

		if err := gz.Close(); err != nil {
			return fmt.Errorf("close gzip: %w", err)
		}
		if err := out.Close(); err != nil {
			return fmt.Errorf("close output: %w", err)
		}

		return nil
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to gzip manpages")
	}
}
