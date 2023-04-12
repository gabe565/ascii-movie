package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	LogLevelFlag    = "log-level"
	DefaultLogLevel = log.InfoLevel

	LogFormatFlag    = "log-format"
	DefaultLogFormat = "text"
)

func RegisterLogFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		LogLevelFlag,
		"l",
		DefaultLogLevel.String(),
		"log level (trace, debug, info, warning, error, fatal, panic)",
	)
	if err := cmd.RegisterFlagCompletionFunc(
		LogLevelFlag,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				log.TraceLevel.String(),
				log.DebugLevel.String(),
				log.InfoLevel.String(),
				log.WarnLevel.String(),
				log.ErrorLevel.String(),
				log.FatalLevel.String(),
				log.PanicLevel.String(),
			}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().String(LogFormatFlag, DefaultLogFormat, "log formatter (text, json)")
	if err := cmd.RegisterFlagCompletionFunc(
		LogFormatFlag,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"text", "json"}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}
}

func InitLog(flags *flag.FlagSet) {
	logLevel, err := flags.GetString(LogLevelFlag)
	if err != nil {
		panic(err)
	}

	if parsedLevel, err := log.ParseLevel(logLevel); err == nil {
		log.SetLevel(parsedLevel)
	} else {
		parsedLevel = log.InfoLevel

		log.WithField(LogLevelFlag, logLevel).Warn("invalid log level. defaulting to info.")
		if err = flags.Set(LogLevelFlag, "info"); err != nil {
			panic(err)
		}
		log.SetLevel(parsedLevel)
	}

	logFormat, err := flags.GetString(LogFormatFlag)
	if err != nil {
		panic(err)
	}

	switch logFormat {
	case "text", "txt", "t":
		log.SetFormatter(&log.TextFormatter{})
	case "json", "j":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.WithField(LogFormatFlag, logFormat).Warn("invalid log formatter. defaulting to text.")
		if err = flags.Set(LogFormatFlag, "text"); err != nil {
			panic(err)
		}
	}
}
