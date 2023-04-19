package movie

import (
	"strings"
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestFromFlags(t *testing.T) {
	testMovie := func(t *testing.T, path string) {
		flags := flag.NewFlagSet("FromFlags", flag.ContinueOnError)
		Flags(flags)
		movie, err := FromFlags(flags, path)
		if !assert.NoError(t, err) {
			return
		}

		speed, err := flags.GetFloat64(SpeedFlag)
		if !assert.NoError(t, err) {
			return
		}
		assert.EqualValues(t, speed, movie.Speed)

		padTop, err := flags.GetInt(PadTopFlag)
		if !assert.NoError(t, err) {
			return
		}

		padLeft, err := flags.GetInt(PadLeftFlag)
		if !assert.NoError(t, err) {
			return
		}

		padBottom, err := flags.GetInt(PadBottomFlag)
		if !assert.NoError(t, err) {
			return
		}

		progressPadBottom, err := flags.GetInt(ProgressPadBottomFlag)
		if !assert.NoError(t, err) {
			return
		}

		for _, frame := range movie.Frames {
			var numPadTop, numPadBottom, numProgressPadBottom int
			for i, line := range strings.Split(frame.Data, "\n") {
				if i < padTop {
					// Line is above frame data. Count the top padding.
					numPadTop += 1
				} else if i < frame.Height-padBottom-progressPadBottom-1 {
					// Line is within frame data. Count the left padding.
					assert.True(t, strings.HasPrefix(line, strings.Repeat(" ", padLeft)), "Incorrect left padding in frame contents")
				} else if i < frame.Height-progressPadBottom-1 {
					// Line is before frame data. Count the bottom padding.
					numPadBottom += 1
				} else if i < frame.Height-progressPadBottom {
					// Line is progress bar. Count the left padding.
					assert.True(t, strings.HasPrefix(line, strings.Repeat(" ", padLeft)), "Incorrect left padding in progress bar")
				} else if i <= frame.Height {
					// Line is below progress bar. Count the progress bottom padding.
					numProgressPadBottom += 1
				} else {
					// This should never be hit
					assert.LessOrEqual(t, i, frame.Height, "Frame data had more lines than the height indicated")
				}
			}
			assert.Equal(t, padTop, numPadTop, "Incorrect top padding")
			assert.Equal(t, padBottom, numPadBottom, "Incorrect bottom padding")
			assert.Equal(t, progressPadBottom, numProgressPadBottom-1, "Incorrect progress bottom padding")
		}
	}

	t.Run("default embedded", func(t *testing.T) {
		t.Parallel()
		testMovie(t, "")
	})

	t.Run("short_intro embedded", func(t *testing.T) {
		t.Parallel()
		testMovie(t, "short_intro")
	})

	t.Run("short_intro file", func(t *testing.T) {
		t.Parallel()
		testMovie(t, "../../movies/short_intro.txt")
	})

	t.Run("invalid speed", func(t *testing.T) {
		t.Parallel()

		flags := flag.NewFlagSet("FromFlags", flag.ContinueOnError)
		Flags(flags)

		if err := flags.Set(SpeedFlag, "-1"); !assert.NoError(t, err) {
			return
		}

		if _, err := FromFlags(flags, ""); !assert.Error(t, err) {
			return
		}
	})
}
