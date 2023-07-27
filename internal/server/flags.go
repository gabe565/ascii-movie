package server

import (
	"time"

	flag "github.com/spf13/pflag"
)

const (
	SSHFlagPrefix    = "ssh"
	TelnetFlagPrefix = "telnet"
	ApiFlagPrefix    = "api"
	EnabledFlag      = "-enabled"
	AddressFlag      = "-address"

	SSHHostKeyPathFlag = SSHFlagPrefix + "-host-key"
	SSHHostKeyDataFlag = SSHFlagPrefix + "-host-key-data"

	ConcurrentStreamsFlag = "concurrent-streams"
	TimeoutFlag           = "timeout"
)

var (
	concurrentStreams uint
	timeout           time.Duration
)

func Flags(flags *flag.FlagSet) {
	flags.Bool(SSHFlagPrefix+EnabledFlag, true, "Enables SSH listener")
	flags.String(SSHFlagPrefix+AddressFlag, ":22", "SSH listen address")
	flags.StringSlice(SSHHostKeyPathFlag, []string{}, "SSH host key file path")
	flags.StringSlice(SSHHostKeyDataFlag, []string{}, "SSH host key PEM data")

	flags.Bool(TelnetFlagPrefix+EnabledFlag, true, "Enables Telnet listener")
	flags.String(TelnetFlagPrefix+AddressFlag, ":23", "Telnet listen address")

	flags.Bool(ApiFlagPrefix+EnabledFlag, true, "Enables API listener")
	flags.String(ApiFlagPrefix+AddressFlag, "127.0.0.1:1977", "API listen address")

	flags.UintVar(&concurrentStreams, ConcurrentStreamsFlag, 10, "Number of concurrent streams allowed from an IP address. Set to 0 to disable.")
	flags.DurationVar(&timeout, TimeoutFlag, time.Hour, "Maximum amount of time that a connection may stay active.")
}
