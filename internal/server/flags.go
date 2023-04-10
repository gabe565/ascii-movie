package server

import (
	"github.com/gabe565/ascii-movie/internal/movie"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type Config struct {
	Enabled bool
	Address string
	Log     *log.Entry
	Movie   *movie.Movie
}

var (
	SSHEnabledFlag     = "ssh-enabled"
	SSHAddressFlag     = "ssh-address"
	SSHHostKeyFlag     = "ssh-host-key"
	SSHHostKeyPathFlag = "ssh-host-key-path"

	TelnetEnabledFlag = "telnet-enabled"
	TelnetAddressFlag = "telnet-address"
)

func Flags(flags *flag.FlagSet) {
	flags.Bool(SSHEnabledFlag, true, "Enables SSH listener")
	flags.String(SSHAddressFlag, ":22", "SSH listen address")
	flags.String(SSHHostKeyFlag, "", "SSH host key PEM")
	flags.String(SSHHostKeyPathFlag, "", "SSH host key file path")

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
