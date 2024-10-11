package config

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

const (
	LogLevelFlag  = "log-level"
	LogFormatFlag = "log-format"
	LevelTrace    = slog.Level(-5)
)

//go:generate go run github.com/dmarkham/enumer -type LogFormat -trimprefix Format -transform lower -text

type LogFormat uint8

const (
	FormatAuto LogFormat = iota
	FormatColor
	FormatPlain
	FormatJSON
)

func RegisterLogFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(LogLevelFlag, "l", strings.ToLower(slog.LevelInfo.String()), "log level (one of debug, info, warn, error)")
	cmd.PersistentFlags().String(LogFormatFlag, FormatAuto.String(), "log formatter (one of "+strings.Join(LogFormatStrings(), ", ")+")")

	if err := errors.Join(
		cmd.RegisterFlagCompletionFunc(LogLevelFlag,
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return []string{
					strings.ToLower(slog.LevelDebug.String()),
					strings.ToLower(slog.LevelInfo.String()),
					strings.ToLower(slog.LevelWarn.String()),
					strings.ToLower(slog.LevelError.String()),
				}, cobra.ShellCompDirectiveNoFileComp
			},
		),
		cmd.RegisterFlagCompletionFunc(LogFormatFlag,
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return LogFormatStrings(), cobra.ShellCompDirectiveNoFileComp
			},
		),
	); err != nil {
		panic(err)
	}
}

func InitLogCmd(cmd *cobra.Command) {
	levelStr, err := cmd.Flags().GetString("log-level")
	if err != nil {
		panic(err)
	}
	var level slog.Level
	if v, err := strconv.Atoi(levelStr); err == nil {
		level = slog.Level(v)
	} else if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		defer func() {
			slog.Warn("Invalid log level. Defaulting to info.", "value", levelStr)
		}()
		level = slog.LevelInfo
	}

	formatStr, err := cmd.Flags().GetString("log-format")
	if err != nil {
		panic(err)
	}
	var format LogFormat
	if err := format.UnmarshalText([]byte(formatStr)); err != nil {
		defer func() {
			slog.Warn("Invalid log format. Defaulting to auto.", "value", formatStr)
		}()
		format = FormatAuto
	}

	InitLog(cmd.ErrOrStderr(), level, format)
}

func InitLog(w io.Writer, level slog.Level, format LogFormat) {
	switch format {
	case FormatJSON:
		slog.SetDefault(slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: level,
		})))
	default:
		var color bool
		switch format {
		case FormatAuto:
			if f, ok := w.(*os.File); ok {
				color = isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
			}
		case FormatColor:
			color = true
		}

		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      level,
				TimeFormat: time.DateTime,
				NoColor:    !color,
			}),
		))
	}
}
