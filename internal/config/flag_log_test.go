package config

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
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
		defer cleanupFunc(log.GetLevel(), log.StandardLogger().Formatter)

		InitLog(cmd.PersistentFlags())
		assert.Equal(t, log.GetLevel(), DefaultLogLevel)
		assert.Equal(t, log.StandardLogger().Formatter, &log.TextFormatter{})
	})

	t.Run("warn level/json formatter", func(t *testing.T) {
		defer cleanupFunc(log.GetLevel(), log.StandardLogger().Formatter)

		if err := cmd.PersistentFlags().Set(LogLevelFlag, log.WarnLevel.String()); !assert.NoError(t, err) {
			return
		}
		if err := cmd.PersistentFlags().Set(LogFormatFlag, "json"); !assert.NoError(t, err) {
			return
		}
		InitLog(cmd.PersistentFlags())
		assert.Equal(t, log.GetLevel(), log.WarnLevel)
		assert.Equal(t, log.StandardLogger().Formatter, &log.JSONFormatter{})
	})

	t.Run("invalid level/invalid formatter", func(t *testing.T) {
		defer cleanupFunc(log.GetLevel(), log.StandardLogger().Formatter)

		formatter := log.StandardLogger().Formatter

		if err := cmd.PersistentFlags().Set(LogLevelFlag, "invalid"); !assert.NoError(t, err) {
			return
		}
		if err := cmd.PersistentFlags().Set(LogFormatFlag, "invalid"); !assert.NoError(t, err) {
			return
		}
		InitLog(cmd.PersistentFlags())
		assert.Equal(t, log.GetLevel(), log.InfoLevel)
		assert.Equal(t, log.StandardLogger().Formatter, formatter)
	})
}
