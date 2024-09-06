# ASCII Movie

<img src="./assets/icon.svg" alt="ascii-movie logo" width="92" align="right">

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/gabe565/ascii-movie)](https://github.com/gabe565/ascii-movie/releases)
[![Build](https://github.com/gabe565/ascii-movie/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/ascii-movie/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabe565/ascii-movie)](https://goreportcard.com/report/github.com/gabe565/ascii-movie)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gabe565_ascii-movie&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gabe565_ascii-movie)

Stream the original Star Wars ASCII movie to command-line clients via SSH or Telnet.

Inspired by [asciimation](https://asciimation.co.nz) and the iconic [towel.blinkenlights.nl](https://web.archive.org/web/20021205144143/http://www.blinkenlights.nl/thereg/), this Go rewrite introduces an interactive UI with both keyboard and mouse support.

## Try It

Run one of these commands in a terminal to see it in action:
- **SSH:** `ssh starwarstel.net`
- **Telnet:** `telnet starwarstel.net`
- **Docker:** `docker run --rm -it ghcr.io/gabe565/ascii-movie play`

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

The app can play a movie directly on your terminal with the [`play`](docs/ascii-movie_play.md) subcommand, or it can host SSH and Telnet servers with the [`serve`](docs/ascii-movie_serve.md) subcommand.

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
