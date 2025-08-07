## ascii-movie serve

Serve an ASCII movie over Telnet and SSH.

```
ascii-movie serve [movie] [flags]
```

### Options

```
      --api-address string          API listen address (default "127.0.0.1:1977")
      --api-enabled                 Enables API listener (default true)
      --concurrent-streams uint     Number of concurrent streams allowed from an IP address. Set to 0 to disable. (default 10)
  -h, --help                        help for serve
      --idle-timeout duration       Idle connection timeout. (default 15m0s)
      --max-timeout duration        Absolute connection timeout. (default 2h0m0s)
      --no-controls                 Disable all UI controls, resulting in an experience more faithful to the original.
      --speed float                 Playback speed multiplier. Must be greater than 0. (default 1)
      --ssh-address string          SSH listen address (default ":22")
      --ssh-enabled                 Enables SSH listener (default true)
      --ssh-host-key strings        SSH host key file path
      --ssh-host-key-data strings   SSH host key PEM data
      --telnet-address string       Telnet listen address (default ":23")
      --telnet-enabled              Enables Telnet listener (default true)
```

### Options inherited from parent commands

```
      --log-format string   Log format (one of auto, color, plain, json) (default "auto")
  -l, --log-level string    Log level (one of trace, debug, info, warn, error) (default "info")
```

### SEE ALSO

* [ascii-movie](ascii-movie.md)	 - Command line ASCII movie player.

