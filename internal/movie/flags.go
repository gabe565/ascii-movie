package movie

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gabe565/ascii-movie/config"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"io"
	"os"
)

var (
	FileFlag = "file"

	ClearExtraLinesFlag = "clear-extra-lines"

	SpeedFlag = "speed"

	ErrInvalidSpeed = errors.New("speed must be greater than 0")
)

func Flags(flags *flag.FlagSet) {
	flags.String(
		FileFlag,
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
		src = bytes.NewReader(config.DefaultMovie)
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

	movie, err = NewFromFile(fileFlag, src)
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
