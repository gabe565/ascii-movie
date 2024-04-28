package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	LogLevelFlag    = "log-level"
	DefaultLogLevel = zerolog.InfoLevel

	LogFormatFlag    = "log-format"
	DefaultLogFormat = "auto"
)

func RegisterLogFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		LogLevelFlag,
		"l",
		DefaultLogLevel.String(),
		"log level (trace, debug, info, warn, error, fatal, panic)",
	)
	if err := cmd.RegisterFlagCompletionFunc(
		LogLevelFlag,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{
				zerolog.TraceLevel.String(),
				zerolog.DebugLevel.String(),
				zerolog.InfoLevel.String(),
				zerolog.WarnLevel.String(),
				zerolog.ErrorLevel.String(),
				zerolog.FatalLevel.String(),
				zerolog.PanicLevel.String(),
			}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}

	cmd.PersistentFlags().String(LogFormatFlag, DefaultLogFormat, "log formatter (auto, color, plain, json)")
	if err := cmd.RegisterFlagCompletionFunc(
		LogFormatFlag,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"auto", "color", "plain", "json"}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}
}

func logLevel(level string) zerolog.Level {
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil || parsedLevel == zerolog.NoLevel {
		if level == "warning" {
			parsedLevel = zerolog.WarnLevel
		} else {
			log.Warn().Str("value", level).Msg("Invalid log level. Defaulting to info.")
			parsedLevel = zerolog.InfoLevel
		}
	}
	return parsedLevel
}

func logFormat(out io.Writer, format string) io.Writer {
	switch format {
	case "json", "j":
		return out
	default:
		style := lipgloss.NewStyle().Bold(true)
		var useColor bool
		switch format {
		case "auto", "text":
			if w, ok := out.(*os.File); ok {
				useColor = isatty.IsTerminal(w.Fd())
			}
		case "color":
			useColor = true
		case "plain":
		default:
			log.Warn().Str("value", format).Msg("Invalid log formatter. Defaulting to auto.")
		}

		return zerolog.ConsoleWriter{
			Out:        out,
			NoColor:    !useColor,
			TimeFormat: time.DateTime,
			FormatMessage: func(i interface{}) string {
				msg := fmt.Sprintf("%-25s", i)
				if useColor {
					return style.Render(msg)
				}
				return msg
			},
		}
	}
}

func InitLog(cmd *cobra.Command) {
	level, err := cmd.Flags().GetString("log-level")
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(logLevel(level))

	format, err := cmd.Flags().GetString("log-format")
	if err != nil {
		panic(err)
	}
	log.Logger = log.Output(logFormat(cmd.ErrOrStderr(), format))
}
