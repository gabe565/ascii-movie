#syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.6.1 AS xx

FROM --platform=$BUILDPLATFORM golang:1.25.0-alpine AS build
WORKDIR /app

COPY --from=xx / /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go generate -x ./...

ARG TARGETPLATFORM
RUN --mount=type=cache,target=/root/.cache \
  CGO_ENABLED=0 xx-go build -ldflags='-w -s' -trimpath -tags gzip


FROM alpine:3.22.1

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
