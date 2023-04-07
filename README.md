# ascii-movie

<img src="./assets/icon.svg" alt="ascii-movie logo" width="92" align="right">

[![Build](https://github.com/gabe565/ascii-movie/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/ascii-movie/actions/workflows/build.yml)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/gabe565)](https://artifacthub.io/packages/helm/gabe565/ascii-telnet)

The original Star Wars telnet server is currently down, so why not implement it in Go? This server will open a TCP server on `0.0.0.0:23` which streams the original Star Wars ASCII movie over telnet.

See it in action by running `telnet gabecook.com` or `nc gabecook.com 23`.

<div align="center">
  <video src="https://user-images.githubusercontent.com/7717888/230577875-ef2e19bb-a804-40a1-9990-84a4ccff29df.mp4"></video>
</div>

<details>
  <summary>Also available on asciinema</summary>

  <p align="center">
    <a href="https://asciinema.org/a/431278"><img src="https://asciinema.org/a/431278.svg"/></a>
  </p>
</details>

## Running

The app supports building locally or in a Docker container at `ghcr.io/gabe565/ascii-movie`.

See generated [docs](./docs/ascii-movie.md) for command line usage.

### Local
```shell
# Generate the movie frames
go generate

# Build the app
go build -ldflags='-w -s'

# Run the app in your terminal
./ascii-movie play

# Or run it as a server
./ascii-movie serve

# Now, run `telnet localhost` to watch the movie!
```

### Docker
An image is available at [`ghcr.io/gabe565/ascii-movie`](ghcr.io/gabe565/ascii-movie).

The following command would run a Telnet server on port `23` and an SSH server on port `2222`.
```shell
docker run --rm -it -p 23:23 -p 2222:22 ghcr.io/gabe565/ascii-movie
```

### Kubernetes

A Helm chart is available for Kubernetes deployment.
For more information, go to
[Artifact Hub](https://artifacthub.io/packages/helm/gabe565/ascii-telnet) or
[gabe565/charts](https://github.com/gabe565/charts/tree/main/charts/ascii-telnet).
