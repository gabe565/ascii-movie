package server

import (
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

	// Deprecated --address flag
	flags.StringP("address", "a", ":23", "Telnet listen address")
	if err := flags.MarkDeprecated(
		"address",
		"please use --telnet-address instead.",
	); err != nil {
		panic(err)
	}
}
