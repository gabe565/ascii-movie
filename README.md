# ascii-movie

<img src="./assets/icon.svg" alt="ascii-movie logo" width="92" align="right">

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/gabe565/ascii-movie)](https://github.com/gabe565/ascii-movie/releases)
[![Build](https://github.com/gabe565/ascii-movie/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/ascii-movie/actions/workflows/build.yml)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/gabe565)](https://artifacthub.io/packages/helm/gabe565/ascii-movie)

The original Star Wars telnet server is currently down, so why not implement it in Go? This server will open a TCP server on `0.0.0.0:23` and an SSH server on `0.0.0.0:22` which stream the original Star Wars ASCII movie to command line clients.

Run one of these commands in a terminal to see it in action:
- **SSH:** `ssh movie.gabe565.com`
- **Telnet:** `telnet movie.gabe565.com`
- **Docker:** `docker run --rm -it ghcr.io/gabe565/ascii-movie play`

## Demo
<div align="center">
  <video src="https://user-images.githubusercontent.com/7717888/233742309-4eeace5e-9a7c-41c6-9fc5-21ea19728f77.mp4"></video>
</div>

<details>
  <summary>Also available on asciinema</summary>

  <p align="center">
    <a href="https://asciinema.org/a/431278"><img src="https://asciinema.org/a/431278.svg"/></a>
  </p>
</details>

## Installation

See [Installation](https://github.com/gabe565/ascii-movie/wiki/Installation).

## Usage

The app can play a movie directly on your terminal with the [`play`](docs/ascii-movie_play.md) subcommand, or it can host an SSH and Telnet stream server with the [`serve`](docs/ascii-movie_serve.md) subcommand.

See generated [docs](./docs/ascii-movie.md) for command line usage information.

### Docker (Suggested)
An image is available at [`ghcr.io/gabe565/ascii-movie`](https://ghcr.io/gabe565/ascii-movie).

#### Watch Locally
The following command will run a container that plays the movie directly in your terminal.

```shell
sudo docker run --rm -it ghcr.io/gabe565/ascii-movie play
```

#### Serve Movie over Telnet and SSH
The following command will run a Telnet server on port `23` and an SSH server on port `2222`.
```shell
sudo docker run --port=22:22 --port=23:23 ghcr.io/gabe565/ascii-movie serve
```

### Other

See [Usage](https://github.com/gabe565/ascii-movie/wiki/Usage).
