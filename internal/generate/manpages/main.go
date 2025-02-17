package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gabe565.com/ascii-movie/cmd"
	"gabe565.com/ascii-movie/internal/config"
	"gabe565.com/utils/slogx"
	"github.com/spf13/cobra/doc"
	flag "github.com/spf13/pflag"
)

func main() {
	config.InitLog(os.Stderr, slogx.LevelInfo, slogx.FormatAuto)

	flags := flag.NewFlagSet("", flag.ContinueOnError)

	var version string
	flags.StringVar(&version, "version", "beta", "Version")

	var dateParam string
	flags.StringVar(&dateParam, "date", time.Now().Format(time.RFC3339), "Build date")

	if err := flags.Parse(os.Args); err != nil {
		slog.Error("Failed to parse arguments", "error", err)
		os.Exit(1)
	}

	if err := run(version, dateParam); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(version, dateParam string) error {
	if err := os.RemoveAll("manpages"); err != nil {
		return fmt.Errorf("failed to remove manpages dir: %w", err)
	}

	if err := os.MkdirAll("manpages", 0o755); err != nil {
		return fmt.Errorf("failed to create manpages dir: %w", err)
	}

	rootCmd := cmd.NewCommand()
	rootName := rootCmd.Name()

	date, err := time.Parse(time.RFC3339, dateParam)
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}

	header := doc.GenManHeader{
		Title:   strings.ToUpper(rootName),
		Section: "1",
		Date:    &date,
		Source:  rootName + " " + version,
		Manual:  "User Commands",
	}

	if err := doc.GenManTree(rootCmd, &header, "manpages"); err != nil {
		return fmt.Errorf("failed to generate manpages: %w", err)
	}

	if err := filepath.WalkDir("manpages", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
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
		return fmt.Errorf("failed to gzip manpages: %w", err)
	}

	return nil
}
