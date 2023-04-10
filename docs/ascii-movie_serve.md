## ascii-movie serve

Serve an ASCII movie over Telnet and SSH.

```
ascii-movie serve [flags]
```

### Options

```
  -f, --file string                Movie file path. If left blank, Star Wars will be played.
  -h, --help                       help for serve
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

