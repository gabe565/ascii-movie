package config

import (
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const (
	LogLevelFlag  = "log-level"
	LogFormatFlag = "log-format"
)

func RegisterLogFlags(flags *flag.FlagSet) {
	flags.StringP(
		LogLevelFlag,
		"l",
		log.InfoLevel.String(),
		"log level (trace, debug, info, warning, error, fatal, panic)",
	)

	flags.String(LogFormatFlag, "text", "log formatter (text, json)")
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
