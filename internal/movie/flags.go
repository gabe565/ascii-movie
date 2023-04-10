package movie

import (
	"errors"
	"fmt"
	"github.com/gabe565/ascii-movie/movies"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"io"
	"os"
	"strings"
)

var (
	ClearExtraLinesFlag = "clear-extra-lines"

	SpeedFlag       = "speed"
	ErrInvalidSpeed = errors.New("speed must be greater than 0")

	FrameHeightFlag = "frame-height"

	PadTopFlag            = "pad-top"
	PadBottomFlag         = "pad-bottom"
	PadLeftFlag           = "pad-left"
	ProgressPadBottomFlag = "progress-pad-bottom"
)

func Flags(flags *flag.FlagSet) {
	flags.Int(
		ClearExtraLinesFlag,
		0,
		"Clears extra lines between each frame. Should typically be ignored.",
	)
	if err := flags.MarkHidden(ClearExtraLinesFlag); err != nil {
		panic(err)
	}

	flags.Float64(
		SpeedFlag,
		1,
		"Playback speed multiplier. Must be greater than 0.",
	)

	flags.Int(FrameHeightFlag, 14, "Height of the movie frames")

	flags.Int(PadTopFlag, 3, "Padding above the movie")
	flags.Int(PadBottomFlag, 2, "Padding below the movie")
	flags.Int(PadLeftFlag, 6, "Padding left of the movie")
	flags.Int(ProgressPadBottomFlag, 3, "Padding below the progress bar")
}

func FromFlags(flags *flag.FlagSet, path string) (*Movie, error) {
	var err error

	log.Info("Loading movie...")

	var movie *Movie

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
	if src, err = movies.Movies.Open(embeddedPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Fallback to loading file
			f, err := os.Open(path)
			if err != nil {
				return movie, err
			}

			src = f
			defer func(f *os.File) {
				_ = f.Close()
			}(f)
		} else {
			return movie, err
		}
	}

	var pad Padding
	if pad.Top, err = flags.GetInt(PadTopFlag); err != nil {
		panic(err)
	}
	if pad.Bottom, err = flags.GetInt(PadBottomFlag); err != nil {
		panic(err)
	}
	if pad.Left, err = flags.GetInt(PadLeftFlag); err != nil {
		panic(err)
	}

	progressPad := pad
	if progressPad.Bottom, err = flags.GetInt(ProgressPadBottomFlag); err != nil {
		panic(err)
	}

	frameHeight, err := flags.GetInt(FrameHeightFlag)
	if err != nil {
		panic(err)
	}

	movie, err = NewFromFile(path, src, frameHeight, pad, progressPad)
	if err != nil {
		return movie, err
	}

	movie.ClearExtraLines, err = flags.GetInt(ClearExtraLinesFlag)
	if err != nil {
		return movie, err
	}

	movie.Speed, err = flags.GetFloat64(SpeedFlag)
	if err != nil {
		return movie, err
	}
	if movie.Speed <= 0 {
		return movie, fmt.Errorf("%w: %f", ErrInvalidSpeed, movie.Speed)
	}

	log.WithField("duration", movie.Duration()).Info("Movie loaded")

	return movie, nil
}