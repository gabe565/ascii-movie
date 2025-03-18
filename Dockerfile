#syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.24.1-alpine AS build
WORKDIR /app

COPY --from=tonistiigi/xx:1.6.1 / /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETPLATFORM
# Set Golang build envs based on Docker platform string
RUN --mount=type=cache,target=/root/.cache <<EOT
  set -x
  xx-go generate -x ./...
  CGO_ENABLED=0 xx-go build -ldflags='-w -s' -trimpath -tags gzip
EOT


FROM alpine:3.21.3

RUN apk add --no-cache tzdata

ARG USERNAME=ascii-movie
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"

COPY --from=build /app/ascii-movie /bin
ENV TERM=xterm-256color
ENV ASCII_MOVIE_SSH_HOST_KEY=/data/id_ed25519,/data/id_rsa
VOLUME /data
USER $UID
ENTRYPOINT ["ascii-movie"]
CMD ["serve"]
