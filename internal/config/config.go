package config

import (
	"time"

	"gabe565.com/utils/slogx"
)

type Config struct {
	Speed     float64
	LogLevel  slogx.Level
	LogFormat slogx.Format

	Server Server
}

type Server struct {
	ConcurrentStreams uint
	IdleTimeout       time.Duration
	MaxTimeout        time.Duration

	SSH    SSH
	Telnet Telnet
	API    API
}

type Listener struct {
	Enabled bool
	Address string
}

type SSH struct {
	Listener
	HostKeyPath []string
	HostKeyPEM  []string
}

type Telnet struct {
	Listener
}

type API struct {
	Listener
}

func New() *Config {
	return &Config{
		Speed:     1,
		LogLevel:  slogx.LevelInfo,
		LogFormat: slogx.FormatAuto,

		Server: Server{
			ConcurrentStreams: 10,
			IdleTimeout:       15 * time.Minute,
			MaxTimeout:        2 * time.Hour,

			SSH: SSH{
				Listener: Listener{
					Enabled: true,
					Address: ":22",
				},
			},

			Telnet: Telnet{
				Listener: Listener{
					Enabled: true,
					Address: ":23",
				},
			},

			API: API{
				Listener: Listener{
					Enabled: true,
					Address: "127.0.0.1:1977",
				},
			},
		},
	}
}
