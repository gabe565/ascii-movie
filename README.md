# ascii-telnet-go

<img src="./assets/icon.svg" alt="ascii-telnet logo" width="92" align="right">

[![Build](https://github.com/gabe565/ascii-telnet-go/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/ascii-telnet-go/actions/workflows/build.yml)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/gabe565)](https://artifacthub.io/packages/helm/gabe565/ascii-telnet)

The original Star Wars telnet server is currently down, so why not implement it in Go? This server will open a TCP server on `0.0.0.0:23` which streams the original Star Wars ASCII movie over telnet.

See it in action by running `telnet gabecook.com` or `nc gabecook.com 23`.

<p align="center">
  <a href="https://asciinema.org/a/431278"><img src="https://asciinema.org/a/431278.svg"/></a>
</p>

## Running

The app supports building locally or in a Docker container at `ghcr.io/gabe565/ascii-telnet-go`.

### Local
```shell
# Generate the movie frames
go generate

# Build the app
go build -ldflags='-w -s' -o ascii-telnet

# Run the app in your terminal
./ascii-telnet play

# Or run it as a server
./ascii-telnet serve

# Now, run `telnet localhost` to watch the movie!
```

### Docker
An image is available at [`ghcr.io/gabe565/ascii-telnet-go`](ghcr.io/gabe565/ascii-telnet-go).
```shell
docker run --rm -it -p '23:23' ghcr.io/gabe565/ascii-telnet-go
```

### Kubernetes

A Helm chart is available for Kubernetes deployment.
For more information, go to
[Artifact Hub](https://artifacthub.io/packages/helm/gabe565/ascii-telnet) or
[gabe565/charts](https://github.com/gabe565/charts/tree/main/charts/ascii-telnet).
