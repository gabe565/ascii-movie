## ascii-movie serve

Serve an ASCII movie over Telnet and SSH.

```
ascii-movie serve [movie] [flags]
```

### Options

```
      --frame-height int           Height of the movie frames (default 14)
  -h, --help                       help for serve
      --pad-bottom int             Padding below the movie (default 2)
      --pad-left int               Padding left of the movie (default 6)
      --pad-top int                Padding above the movie (default 3)
      --progress-pad-bottom int    Padding below the progress bar (default 3)
      --speed float                Playback speed multiplier. Must be greater than 0. (default 1)
      --ssh-address string         SSH listen address (default ":22")
      --ssh-enabled                Enables SSH listener (default true)
      --ssh-host-key string        SSH host key PEM
      --ssh-host-key-path string   SSH host key file path
      --telnet-address string      Telnet listen address (default ":23")
      --telnet-enabled             Enables Telnet listener (default true)
```

### SEE ALSO

* [ascii-movie](ascii-movie.md)	 - Command line ASCII movie player.
