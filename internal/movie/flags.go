package movie

import (
	"errors"
	"fmt"
	flag "github.com/spf13/pflag"
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

	fileFlag, err := flags.GetString(FileFlag)
	if err != nil {
		return nil, err
	}

	var movie *Movie
	if fileFlag == "" {
		movie = Generated
	} else {
		movie, err = NewFromFile(fileFlag)
		if err != nil {
			return movie, err
		}
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

	return movie, nil
}
