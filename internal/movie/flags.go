package movie

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gabe565/ascii-movie/movies"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"io"
	"os"
)

var (
	FileFlag = "file"

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
	flags.StringP(
		FileFlag,
		"f",
		"",
		"Movie file path. If left blank, Star Wars will be played.",
	)

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

func FromFlags(flags *flag.FlagSet) (*Movie, error) {
	var err error

	log.Info("Loading movie...")

	fileFlag, err := flags.GetString(FileFlag)
	if err != nil {
		return nil, err
	}

	var movie *Movie

	var src io.Reader
	if fileFlag == "" {
		src = bytes.NewReader(movies.Default)
	} else {
		f, err := os.Open(fileFlag)
		if err != nil {
			return movie, err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		src = f
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

	movie, err = NewFromFile(fileFlag, src, frameHeight, pad, progressPad)
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
