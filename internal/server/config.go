package server

import (
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
)

type Server struct {
	Enabled bool
	Address string
	Log     zerolog.Logger
}

type MovieServer struct {
	Server
	Movie *movie.Movie
}

func NewServer(flags *flag.FlagSet, prefix string) Server {
	var config Server
	var err error

	config.Log = log.With().Str("server", prefix).Logger()

	if config.Enabled, err = flags.GetBool(prefix + EnabledFlag); err != nil {
		panic(err)
	}

	if config.Address, err = flags.GetString(prefix + AddressFlag); err != nil {
		panic(err)
	}

	return config
}

func NewMovieServer(flags *flag.FlagSet, prefix string) MovieServer {
	var config MovieServer

	config.Server = NewServer(flags, prefix)

	return config
}
