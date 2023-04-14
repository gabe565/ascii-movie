package server

import (
	"github.com/gabe565/ascii-movie/internal/movie"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type Config struct {
	Enabled        bool
	Address        string
	Log            *log.Entry
	Movie          *movie.Movie
	DefaultGateway string
}

const (
	LogExcludeGatewayFlag = "log-exclude-gateway"

	SSHEnabledFlag     = "ssh-enabled"
	SSHAddressFlag     = "ssh-address"
	SSHHostKeyPathFlag = "ssh-host-key"
	SSHHostKeyDataFlag = "ssh-host-key-data"

	TelnetEnabledFlag = "telnet-enabled"
	TelnetAddressFlag = "telnet-address"
)

func Flags(flags *flag.FlagSet) {
	flags.Bool(LogExcludeGatewayFlag, false, "Makes default gateway early disconnect logs be trace level. Useful for excluding health checks from logs.")

	flags.Bool(SSHEnabledFlag, true, "Enables SSH listener")
	flags.String(SSHAddressFlag, ":22", "SSH listen address")
	flags.StringSlice(SSHHostKeyPathFlag, []string{}, "SSH host key file path")
	flags.StringSlice(SSHHostKeyDataFlag, []string{}, "SSH host key PEM data")

	flags.Bool(TelnetEnabledFlag, true, "Enables Telnet listener")
	flags.String(TelnetAddressFlag, ":23", "Telnet listen address")

	// Deprecated --address flag
	flags.StringP("address", "a", ":23", "Telnet listen address")
	if err := flags.MarkDeprecated(
		"address",
		"please use --telnet-address instead.",
	); err != nil {
		panic(err)
	}
}
