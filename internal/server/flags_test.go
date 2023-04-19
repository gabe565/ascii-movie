package server

import (
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestFlags(t *testing.T) {
	t.Run("doesn't panic", func(t *testing.T) {
		flags := flag.NewFlagSet(t.Name(), flag.PanicOnError)
		Flags(flags)
		if err := flags.Parse([]string{}); !assert.NoError(t, err) {
			return
		}
	})
}
