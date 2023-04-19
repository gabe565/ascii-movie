package server

import (
	flag "github.com/spf13/pflag"
)

const (
	LogExcludeGatewayFlag = "log-exclude-gateway"
	LogExcludeFaster      = "log-exclude-faster"

	SSHFlagPrefix    = "ssh"
	TelnetFlagPrefix = "telnet"
	EnabledFlag      = "-enabled"
	AddressFlag      = "-address"

	SSHHostKeyPathFlag = SSHFlagPrefix + "-host-key"
	SSHHostKeyDataFlag = SSHFlagPrefix + "-host-key-data"
)

func Flags(flags *flag.FlagSet) {
	flags.Bool(LogExcludeGatewayFlag, false, "Makes default gateway early disconnect logs be trace level. Useful for excluding health checks from logs.")
	if err := flags.MarkDeprecated(
		LogExcludeGatewayFlag,
		"please use --log-exclude-faster instead.",
	); err != nil {
		panic(err)
	}

	flags.Duration(LogExcludeFaster, 0, "Makes early disconnect logs faster than the value be trace level. Useful for excluding health checks from logs.")

	flags.Bool(SSHFlagPrefix+EnabledFlag, true, "Enables SSH listener")
	flags.String(SSHFlagPrefix+AddressFlag, ":22", "SSH listen address")
	flags.StringSlice(SSHHostKeyPathFlag, []string{}, "SSH host key file path")
	flags.StringSlice(SSHHostKeyDataFlag, []string{}, "SSH host key PEM data")

	flags.Bool(TelnetFlagPrefix+EnabledFlag, true, "Enables Telnet listener")
	flags.String(TelnetFlagPrefix+AddressFlag, ":23", "Telnet listen address")

	// Deprecated --address flag
	flags.StringP("address", "a", ":23", "Telnet listen address")
	if err := flags.MarkDeprecated(
		"address",
		"please use --telnet-address instead.",
	); err != nil {
		panic(err)
	}
}
