package movie

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gabe565/ascii-movie/movies"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	SpeedFlag       = "speed"
	ErrInvalidSpeed = errors.New("speed must be greater than 0")
)

func Flags(flags *flag.FlagSet) {
	flags.Float64(
		SpeedFlag,
		1,
		"Playback speed multiplier. Must be greater than 0.",
	)
}

func FromFlags(flags *flag.FlagSet, path string) (Movie, error) {
	var err error

	log.Info("Loading movie...")

	movie := NewMovie()

	var src io.Reader
	if path == "" {
		// Use default embedded movie
		path = movies.Default
	}
	// Load embedded movie
	embeddedPath := path
	if !strings.HasSuffix(embeddedPath, FileSuffix) {
		embeddedPath += FileSuffix
	}
	if src, err = movies.Movies.Open(embeddedPath); err == nil {
		log.WithField("name", embeddedPath).Debug("Using embedded movie")
	} else {
		if errors.Is(err, os.ErrNotExist) {
			// Fallback to loading file
			log.WithField("name", embeddedPath).Trace("No embedded movie matches name. Searching filesystem.")
			f, err := os.Open(path)
			if err != nil {
				return movie, err
			}
			log.WithField("name", path).Debug("Found movie file")

			src = f
			defer func(f *os.File) {
				_ = f.Close()
			}(f)
		} else {
			return movie, err
		}
	}

	speed, err := flags.GetFloat64(SpeedFlag)
	if err != nil {
		return movie, err
	}
	if speed <= 0 {
		return movie, fmt.Errorf("%w: %f", ErrInvalidSpeed, speed)
	}

	var r io.Reader = src
	if strings.HasSuffix(path, ".gz") {
		r, err = gzip.NewReader(src)
		if err != nil {
			return movie, err
		}
	}

	if err := movie.LoadFile(path, r, speed); err != nil {
		return movie, err
	}

	log.WithField("duration", movie.Duration().Round(time.Second)).Info("Movie loaded")

	return movie, nil
}
