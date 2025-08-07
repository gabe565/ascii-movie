package server

import (
	"log/slog"

	"gabe565.com/ascii-movie/internal/config"
	"gabe565.com/ascii-movie/internal/movie"
)

type Server struct {
	conf  *config.Config
	Info  *Info
	Log   *slog.Logger
	Movie *movie.Movie
}

func NewServer(conf *config.Config, server string, info *Info) Server {
	return Server{
		conf: conf,
		Info: info,
		Log:  slog.With("server", server),
	}
}
