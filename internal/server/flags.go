package server

import (
	"errors"
	flag "github.com/spf13/pflag"
)

var (
	ClearExtraLinesFlag = "clear-extra-lines"

	SpeedFlag = "speed"

	ErrInvalidSpeed = errors.New("speed must be greater than 0")
)

func PlayFlags(flags *flag.FlagSet) {
	flags.Int(
		ClearExtraLinesFlag,
		0,
		"Clears extra lines between each movie. Should typically be ignored.",
	)
	if err := flags.MarkHidden(ClearExtraLinesFlag); err != nil {
		panic(err)
	}

	flags.Float64(
		SpeedFlag,
		1,
		"Playback speed multiplier. Must be greater than 0.",
	)
}

type ServerConfig struct {
	Enabled bool
	Address string
}

var (
	SSHEnabledFlag = "ssh-enabled"
	SSHAddressFlag = "ssh-address"

	TelnetEnabledFlag = "telnet-enabled"
	TelnetAddressFlag = "telnet-address"
)

func ServeFlags(flags *flag.FlagSet) {
	flags.Bool(SSHEnabledFlag, true, "Enables SSH listener")
	flags.String(SSHAddressFlag, ":22", "SSH listen address")

	flags.Bool(TelnetEnabledFlag, true, "Enables Telnet listener")
	flags.String(TelnetAddressFlag, ":23", "Telnet listen address")
}
