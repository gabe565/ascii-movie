## ascii-movie serve

Serve an ASCII movie over Telnet and SSH.

```
ascii-movie serve [movie] [flags]
```

### Options

```
      --body-pad ints                 Body padding (default [3,6,2,6])
  -h, --help                          help for serve
      --log-exclude-faster duration   Makes early disconnect logs faster than the value be trace level. Useful for excluding health checks from logs.
      --progress-pad ints             Progress bar padding (default [2,0,1,0])
      --speed float                   Playback speed multiplier. Must be greater than 0. (default 1)
      --ssh-address string            SSH listen address (default ":22")
      --ssh-enabled                   Enables SSH listener (default true)
      --ssh-host-key strings          SSH host key file path
      --ssh-host-key-data strings     SSH host key PEM data
      --telnet-address string         Telnet listen address (default ":23")
      --telnet-enabled                Enables Telnet listener (default true)
```

### Options inherited from parent commands

```
      --log-format string   log formatter (text, json) (default "text")
  -l, --log-level string    log level (trace, debug, info, warning, error, fatal, panic) (default "info")
```

### SEE ALSO

* [ascii-movie](ascii-movie.md)	 - Command line ASCII movie player.

