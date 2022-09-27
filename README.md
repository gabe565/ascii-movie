# ascii-telnet-go

[![Build](https://github.com/gabe565/ascii-telnet-go/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/ascii-telnet-go/actions/workflows/build.yml)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/gabe565)](https://artifacthub.io/packages/helm/gabe565/ascii-telnet)

The original Star Wars telnet server is currently down, so why not implement it in Go? This server will open a TCP server on `0.0.0.0:23` which streams the original Star Wars ASCII movie over telnet.

<p align="center">
  <a href="https://asciinema.org/a/431278"><img src="https://asciinema.org/a/431278.svg"/></a>
</p>

## Running

The app supports building locally or in a Docker container at `ghcr.io/gabe565/ascii-telnet-go`.

### Local
```shell
$ # To build and run in one step
$ go run . 
INFO[0000] listening for connections                     address=":23"
$ # You can now run `telnet localhost` to see the movie.
$
$ # To get a release binary:
$ go build -ldflags='-w -s'
$ # The binary will be available at ./ascii-telnet-go.
```

### Docker
```shell
$ # An image is available at `ghcr.io/gabe565/ascii-telnet-go`
$ docker run --rm -it -p '23:23' ghcr.io/gabe565/ascii-telnet-go
```

### Kubernetes

A Helm chart is available for Kubernetes deployment.
For more information, go to
[Artifact Hub](https://artifacthub.io/packages/helm/gabe565/ascii-telnet) or
[gabe565/charts](https://github.com/gabe565/charts/tree/main/charts/ascii-telnet).
