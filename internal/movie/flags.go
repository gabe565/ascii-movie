package movie

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"gabe565.com/ascii-movie/internal/config"
	"gabe565.com/ascii-movie/movies"
	flag "github.com/spf13/pflag"
)

const SpeedFlag = "speed"

var ErrInvalidSpeed = errors.New("speed must be greater than 0")

func Flags(flags *flag.FlagSet) {
	flags.Float64(SpeedFlag, 1, "Playback speed multiplier. Must be greater than 0.")
}

func FromFlags(flags *flag.FlagSet, path string) (Movie, error) {
	var err error

	slog.Info("Loading movie...")
	start := time.Now()

	movie := NewMovie()

	var src io.ReadCloser
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
		slog.Debug("Using embedded movie", "name", embeddedPath)

		if strings.HasSuffix(embeddedPath, ".gz") {
			src, err = gzip.NewReader(src)
			if err != nil {
				return movie, err
			}
		}
	} else {
		if errors.Is(err, os.ErrNotExist) {
			// Fallback to loading file
			slog.Log(context.Background(), config.LevelTrace, "No embedded movie matches name. Searching filesystem.")
			f, err := os.Open(path)
			if err != nil {
				return movie, err
			}
			slog.Debug("Found movie file", "name", path)

			src = f
			defer func(f *os.File) {
				_ = f.Close()
			}(f)

			if strings.HasSuffix(path, ".gz") {
				src, err = gzip.NewReader(src)
				if err != nil {
					return movie, err
				}
			}
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

	if err := movie.LoadFile(path, src, speed); err != nil {
		return movie, err
	}

	if err := src.Close(); err != nil {
		return movie, err
	}

	slog.Info("Movie loaded",
		"duration", movie.Duration().Round(time.Second),
		"took", time.Since(start).Round(time.Microsecond),
	)

	return movie, nil
}
