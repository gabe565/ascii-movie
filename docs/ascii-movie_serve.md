## ascii-movie serve

Serve an ASCII movie over Telnet and SSH.

```
ascii-movie serve [movie] [flags]
```

### Options

```
  -h, --help                        help for serve
      --pad-bottom int              Padding below the movie (default 2)
      --pad-left int                Padding left of the movie (default 6)
      --pad-top int                 Padding above the movie (default 3)
      --progress-pad-bottom int     Padding below the progress bar (default 2)
      --speed float                 Playback speed multiplier. Must be greater than 0. (default 1)
      --ssh-address string          SSH listen address (default ":22")
      --ssh-enabled                 Enables SSH listener (default true)
      --ssh-host-key strings        SSH host key PEM
      --ssh-host-key-path strings   SSH host key file path
      --telnet-address string       Telnet listen address (default ":23")
      --telnet-enabled              Enables Telnet listener (default true)
```

### Options inherited from parent commands

```
      --log-format string   log formatter (text, json) (default "text")
  -l, --log-level string    log level (trace, debug, info, warning, error, fatal, panic) (default "info")
```

### SEE ALSO

* [ascii-movie](ascii-movie.md)	 - Command line ASCII movie player.

