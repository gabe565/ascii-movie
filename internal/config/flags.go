package config

import (
	"strings"
	"time"

	"gabe565.com/utils/must"
	"gabe565.com/utils/slogx"
	"github.com/spf13/cobra"
)

const (
	FlagSpeed      = "speed"
	FlagNoControls = "no-controls"
	FlagLogLevel   = "log-level"
	FlagLogFormat  = "log-format"

	FlagPrefixSSH    = "ssh"
	FlagPrefixTelnet = "telnet"
	FlagPrefixAPI    = "api"
	FlagEnabled      = "-enabled"
	FlagAddress      = "-address"

	FlagSSHHostKeyPath = FlagPrefixSSH + "-host-key"
	FlagSSHHostKeyData = FlagPrefixSSH + "-host-key-data"

	FlagConcurrentStreams = "concurrent-streams"
	FlagTimeout           = "timeout"
	FlagIdleTimeout       = "idle-timeout"
	FlagMaxTimeout        = "max-timeout"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	fs := cmd.PersistentFlags()
	fs.VarP(&c.LogLevel, FlagLogLevel, "l", "Log level (one of "+strings.Join(slogx.LevelStrings(), ", ")+")")
	fs.Var(&c.LogFormat, FlagLogFormat, "Log format (one of "+strings.Join(slogx.FormatStrings(), ", ")+")")
}

func (c *Config) RegisterPlayFlags(cmd *cobra.Command) {
	fs := cmd.Flags()
	fs.Float64Var(&c.Speed, FlagSpeed, c.Speed, "Playback speed multiplier. Must be greater than 0.")
	fs.BoolVar(&c.NoControls, FlagNoControls, c.NoControls,
		"Disable all UI controls, resulting in an experience more faithful to the original.",
	)
}

func (s *Server) RegisterFlags(cmd *cobra.Command) {
	fs := cmd.Flags()
	fs.BoolVar(&s.SSH.Enabled, FlagPrefixSSH+FlagEnabled, true, "Enables SSH listener")
	fs.StringVar(&s.SSH.Address, FlagPrefixSSH+FlagAddress, s.SSH.Address, "SSH listen address")
	fs.StringSliceVar(&s.SSH.HostKeyPath, FlagSSHHostKeyPath, s.SSH.HostKeyPath, "SSH host key file path")
	fs.StringSliceVar(&s.SSH.HostKeyPEM, FlagSSHHostKeyData, s.SSH.HostKeyPEM, "SSH host key PEM data")

	fs.BoolVar(&s.Telnet.Enabled, FlagPrefixTelnet+FlagEnabled, s.Telnet.Enabled, "Enables Telnet listener")
	fs.StringVar(&s.Telnet.Address, FlagPrefixTelnet+FlagAddress, s.Telnet.Address, "Telnet listen address")

	fs.BoolVar(&s.API.Enabled, FlagPrefixAPI+FlagEnabled, s.API.Enabled, "Enables API listener")
	fs.StringVar(&s.API.Address, FlagPrefixAPI+FlagAddress, s.API.Address, "API listen address")

	fs.UintVar(&s.ConcurrentStreams, FlagConcurrentStreams, s.ConcurrentStreams,
		"Number of concurrent streams allowed from an IP address. Set to 0 to disable.",
	)
	fs.DurationVar(&s.IdleTimeout, FlagIdleTimeout, s.IdleTimeout, "Idle connection timeout.")
	fs.DurationVar(&s.MaxTimeout, FlagMaxTimeout, s.MaxTimeout, "Absolute connection timeout.")

	fs.Duration(FlagTimeout, time.Hour, "Maximum amount of time that a connection may stay active.")
	must.Must(fs.MarkDeprecated(FlagTimeout, "please use --idle-timeout and --max-timeout instead."))
}
