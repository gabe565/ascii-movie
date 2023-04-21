package server

import (
	"testing"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	testGetConfig := func(t *testing.T, prefix string) {
		flags := flag.NewFlagSet(t.Name(), flag.PanicOnError)
		Flags(flags)

		if err := flags.Set(prefix+EnabledFlag, "true"); !assert.NoError(t, err) {
			return
		}

		if err := flags.Set(prefix+AddressFlag, "127.0.0.1:1977"); !assert.NoError(t, err) {
			return
		}

		if err := flags.Set(LogExcludeFaster, "1s"); !assert.NoError(t, err) {
			return
		}

		server := NewMovieServer(flags, prefix)
		assert.Equal(t, true, server.Enabled)
		assert.Equal(t, "127.0.0.1:1977", server.Address)
		assert.Equal(t, time.Second, server.LogExcludeFaster)
		assert.Equal(t, prefix, server.Log.Data["server"])
	}

	t.Run("SSH gets config from flags", func(t *testing.T) {
		testGetConfig(t, SSHFlagPrefix)
	})

	t.Run("Telnet gets config from flags", func(t *testing.T) {
		testGetConfig(t, TelnetFlagPrefix)
	})
}
