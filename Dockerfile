FROM --platform=$BUILDPLATFORM golang:1.21.6 as build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETPLATFORM
# Set Golang build envs based on Docker platform string
RUN --mount=type=cache,target=/root/.cache \
    set -x \
    && case "$TARGETPLATFORM" in \
        'linux/amd64') export GOARCH=amd64 ;; \
        'linux/arm/v6') export GOARCH=arm GOARM=6 ;; \
        'linux/arm/v7') export GOARCH=arm GOARM=7 ;; \
        'linux/arm64' | 'linux/arm64/v8') export GOARCH=arm64 ;; \
        *) echo "Unsupported target: $TARGETPLATFORM" && exit 1 ;; \
    esac \
    && go generate \
    && CGO_ENABLED=0 go build -ldflags='-w -s' -trimpath -tags gzip


FROM alpine:3.19

RUN apk add --no-cache tzdata

ARG USERNAME=ascii-movie
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"

COPY --from=build /app/ascii-movie /bin
ENV TERM=xterm-256color
USER $UID
ENTRYPOINT ["ascii-movie"]
CMD ["serve"]
