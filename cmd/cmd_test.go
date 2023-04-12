package cmd

import (
	"os"
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func Test_loadFlagEnvs(t *testing.T) {
	defer func() {
		_ = os.Unsetenv("ASCII_MOVIE_WORKED")
	}()

	flags := flag.NewFlagSet("loadFlagEnvs", flag.ContinueOnError)
	flags.Bool("worked", false, "Test flag")
	if err := os.Setenv("ASCII_MOVIE_WORKED", "true"); !assert.NoError(t, err) {
		return
	}

	if err := loadFlagEnvs(flags); !assert.NoError(t, err) {
		return
	}

	worked, err := flags.GetBool("worked")
	if !assert.NoError(t, err) {
		return
	}

	assert.True(t, worked)
}
