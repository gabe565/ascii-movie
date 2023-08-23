//go:build !gzip

package main

import (
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gabe565/ascii-movie/movies"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := fs.WalkDir(movies.Movies, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		outPath := filepath.Join("movies", path+".gz")
		log.WithField("path", outPath).Debug("Create output")
		out, err := os.Create(outPath)
		if err != nil {
			return err
		}

		log.WithField("path", path).Debug("Open input")
		in, err := movies.Movies.Open(path)
		if err != nil {
			return err
		}

		log.Debug("Copy input to gzip writer")
		gz := gzip.NewWriter(out)
		if _, err := io.Copy(gz, in); err != nil {
			return err
		}

		if err := gz.Close(); err != nil {
			return err
		}

		log.Debug("Close output")
		return out.Close()
	}); err != nil {
		log.Fatal(err)
	}
}
