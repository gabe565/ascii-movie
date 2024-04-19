package cmd

import (
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_loadFlagEnvs(t *testing.T) {
	flags := flag.NewFlagSet("loadFlagEnvs", flag.ContinueOnError)
	flags.Bool("worked", false, "Test flag")
	t.Setenv("ASCII_MOVIE_WORKED", "true")

	require.NoError(t, loadFlagEnvs(flags))

	worked, err := flags.GetBool("worked")
	require.NoError(t, err)

	assert.True(t, worked)
}
