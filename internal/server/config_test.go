package server

import (
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	testGetConfig := func(t *testing.T, prefix string) {
		flags := flag.NewFlagSet(t.Name(), flag.PanicOnError)
		Flags(flags)

		require.NoError(t, flags.Set(prefix+EnabledFlag, "true"))
		require.NoError(t, flags.Set(prefix+AddressFlag, "127.0.0.1:1977"))

		server := NewMovieServer(flags, prefix)
		assert.True(t, server.Enabled)
		assert.Equal(t, "127.0.0.1:1977", server.Address)
	}

	t.Run("SSH gets config from flags", func(t *testing.T) {
		testGetConfig(t, SSHFlagPrefix)
	})

	t.Run("Telnet gets config from flags", func(t *testing.T) {
		testGetConfig(t, TelnetFlagPrefix)
	})
}
