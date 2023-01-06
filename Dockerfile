FROM --platform=$BUILDPLATFORM golang:1.19-alpine as go-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go generate

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
    && go build -ldflags='-w -s' -o ascii-telnet


FROM alpine
LABEL org.opencontainers.image.source="https://github.com/gabe565/ascii-telnet-go"

RUN apk add --no-cache tzdata

COPY --from=go-builder /app/ascii-telnet /usr/local/bin

ARG USERNAME=ascii-telnet
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"
USER $UID

CMD ["ascii-telnet", "serve"]
