package config

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitLog(t *testing.T) {
	cleanupFunc := func(level log.Level, formatter log.Formatter) {
		// Set back to default
		log.SetLevel(level)
		log.SetFormatter(formatter)
	}

	cmd := &cobra.Command{}
	RegisterLogFlags(cmd)

	t.Run("defaults", func(t *testing.T) {
		t.Cleanup(func() {
			cleanupFunc(log.GetLevel(), log.StandardLogger().Formatter)
		})

		InitLog(cmd.PersistentFlags())
		assert.Equal(t, DefaultLogLevel, log.GetLevel())
		assert.Equal(t, &log.TextFormatter{}, log.StandardLogger().Formatter)
	})

	t.Run("warn level/json formatter", func(t *testing.T) {
		t.Cleanup(func() {
			cleanupFunc(log.GetLevel(), log.StandardLogger().Formatter)
		})

		require.NoError(t, cmd.PersistentFlags().Set(LogLevelFlag, log.WarnLevel.String()))
		require.NoError(t, cmd.PersistentFlags().Set(LogFormatFlag, "json"))
		InitLog(cmd.PersistentFlags())
		assert.Equal(t, log.WarnLevel, log.GetLevel())
		assert.Equal(t, &log.JSONFormatter{}, log.StandardLogger().Formatter)
	})

	t.Run("invalid level/invalid formatter", func(t *testing.T) {
		t.Cleanup(func() {
			cleanupFunc(log.GetLevel(), log.StandardLogger().Formatter)
		})

		formatter := log.StandardLogger().Formatter

		require.NoError(t, cmd.PersistentFlags().Set(LogLevelFlag, "invalid"))
		require.NoError(t, cmd.PersistentFlags().Set(LogFormatFlag, "invalid"))
		InitLog(cmd.PersistentFlags())
		assert.Equal(t, log.InfoLevel, log.GetLevel())
		assert.Equal(t, formatter, log.StandardLogger().Formatter)
	})
}
