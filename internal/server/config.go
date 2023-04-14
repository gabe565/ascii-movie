package server

import (
	"time"

	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/jackpal/gateway"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type Config struct {
	Enabled          bool
	Address          string
	Log              *log.Entry
	Movie            *movie.Movie
	DefaultGateway   string
	LogExcludeFaster time.Duration
}

func NewConfig(flags *flag.FlagSet, prefix string) Config {
	var config Config
	var err error

	config.Log = log.WithField("server", prefix)

	if config.Enabled, err = flags.GetBool(prefix + EnabledFlag); err != nil {
		panic(err)
	}

	if config.Address, err = flags.GetString(prefix + AddressFlag); err != nil {
		panic(err)
	}

	logExcludeGateway, err := flags.GetBool(LogExcludeGatewayFlag)
	if err != nil {
		panic(err)
	}
	if logExcludeGateway {
		if defaultGateway, err := gateway.DiscoverGateway(); err == nil {
			config.DefaultGateway = defaultGateway.String()
		} else {
			config.Log.Warn("Failed to discover default gateway")
		}
	}

	config.LogExcludeFaster, err = flags.GetDuration(LogExcludeFaster)
	if err != nil {
		panic(err)
	}

	return config
}
