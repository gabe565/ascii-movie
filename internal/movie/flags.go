package movie

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gabe565/ascii-movie/movies"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	SpeedFlag       = "speed"
	ErrInvalidSpeed = errors.New("speed must be greater than 0")

	PadFlag         = "body-pad"
	ProgressPadFlag = "progress-pad"
)

func Flags(flags *flag.FlagSet) {
	flags.Float64(
		SpeedFlag,
		1,
		"Playback speed multiplier. Must be greater than 0.",
	)

	flags.IntSlice(PadFlag, []int{3, 6, 2, 6}, "Body padding")
	flags.IntSlice(ProgressPadFlag, []int{2, 0, 1, 0}, "Progress bar padding")
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
	if !strings.HasSuffix(embeddedPath, ".txt") {
		embeddedPath += ".txt"
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

	if err := movie.LoadFile(path, src, speed); err != nil {
		return movie, err
	}

	bodyPad, err := flags.GetIntSlice(PadFlag)
	if err != nil {
		panic(err)
	}
	movie.BodyStyle = movie.BodyStyle.Padding(bodyPad...)

	progressPad, err := flags.GetIntSlice(ProgressPadFlag)
	if err != nil {
		panic(err)
	}
	movie.ProgressStyle = movie.ProgressStyle.Padding(progressPad...)

	log.WithField("duration", movie.Duration()).Info("Movie loaded")

	return movie, nil
}
