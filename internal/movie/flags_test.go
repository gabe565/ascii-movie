package movie

import (
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestFromFlags(t *testing.T) {
	testMovie := func(t *testing.T, path string) {
		flags := flag.NewFlagSet(t.Name(), flag.PanicOnError)
		Flags(flags)

		_, err := FromFlags(flags, path)
		if !assert.NoError(t, err) {
			return
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

		flags := flag.NewFlagSet(t.Name(), flag.PanicOnError)
		Flags(flags)

		if err := flags.Set(SpeedFlag, "-1"); !assert.NoError(t, err) {
			return
		}

		if _, err := FromFlags(flags, ""); !assert.Error(t, err) {
			return
		}
	})
}
