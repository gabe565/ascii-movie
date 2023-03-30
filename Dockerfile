FROM --platform=$BUILDPLATFORM golang:1.20 as build
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
    && CGO_ENABLED=0 go build -ldflags='-w -s' -o ascii-telnet


FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=build /app/ascii-telnet /
CMD ["/ascii-telnet", "serve"]
