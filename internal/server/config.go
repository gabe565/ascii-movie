package server

import (
	"log/slog"

	"gabe565.com/ascii-movie/internal/movie"
	"gabe565.com/utils/must"
	flag "github.com/spf13/pflag"
)

type Server struct {
	Enabled bool
	Address string
	Log     *slog.Logger
}

type MovieServer struct {
	Server
	Movie *movie.Movie
}

func NewServer(flags *flag.FlagSet, prefix string) Server {
	return Server{
		Enabled: must.Must2(flags.GetBool(prefix + EnabledFlag)),
		Address: must.Must2(flags.GetString(prefix + AddressFlag)),
		Log:     slog.With("server", prefix),
	}
}

func NewMovieServer(flags *flag.FlagSet, prefix string) MovieServer {
	return MovieServer{Server: NewServer(flags, prefix)}
}
