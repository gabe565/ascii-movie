package main

import (
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabe565/ascii-movie/cmd"
	"github.com/spf13/cobra/doc"
	flag "github.com/spf13/pflag"
)

func main() {
	flags := flag.NewFlagSet("", flag.ContinueOnError)

	var version string
	flags.StringVar(&version, "version", "beta", "Version")

	var dateParam string
	flags.StringVar(&dateParam, "date", time.Now().Format(time.RFC3339), "Build date")

	if err := flags.Parse(os.Args); err != nil {
		panic(err)
	}

	if err := os.RemoveAll("manpages"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("manpages", 0o755); err != nil {
		panic(err)
	}

	rootCmd := cmd.NewCommand("beta", "")
	rootName := rootCmd.Name()

	date, err := time.Parse(time.RFC3339, dateParam)
	if err != nil {
		panic(err)
	}

	header := doc.GenManHeader{
		Title:   strings.ToUpper(rootName),
		Section: "1",
		Date:    &date,
		Source:  rootName + " " + version,
		Manual:  "User Commands",
	}

	if err := doc.GenManTree(rootCmd, &header, "manpages"); err != nil {
		panic(err)
	}

	if err := filepath.Walk("manpages", func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}

		out, err := os.Create(path + ".gz")
		if err != nil {
			return err
		}
		gz := gzip.NewWriter(out)

		if _, err := io.Copy(gz, in); err != nil {
			return err
		}

		if err := in.Close(); err != nil {
			return err
		}
		if err := os.Remove(path); err != nil {
			return err
		}

		if err := gz.Close(); err != nil {
			return err
		}
		if err := out.Close(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(err)
	}
}
