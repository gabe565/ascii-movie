FROM alpine:3.20

RUN apk add --no-cache tzdata

ARG USERNAME=ascii-movie
ARG UID=1000
ARG GID=$UID
RUN addgroup -g "$GID" "$USERNAME" \
    && adduser -S -u "$UID" -G "$USERNAME" "$USERNAME"

COPY ascii-movie /bin
ENV TERM=xterm-256color
ENV ASCII_MOVIE_SSH_HOST_KEY=/data/id_ed25519,/data/id_rsa
VOLUME /data
USER $UID
ENTRYPOINT ["ascii-movie"]
CMD ["serve"]
