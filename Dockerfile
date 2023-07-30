FROM i386/golang:1.20.6-alpine3.18 AS cli_builder

ENV CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /go/src
RUN apk update
RUN apk add --no-cache git
RUN --mount=type=bind,target=/go/src,source=. \
    version=$(git describe --tags || git rev-parse HEAD) && \
    go_module=$(grep "module" go.mod | head -1 | cut -d ' ' -f 2) && \
    go get -u && \
    go build \
        -ldflags "-s -w -X '${go_module}/shared.Version=${version}'" \
        -o /go/bin/steamcmd-cli


FROM i386/debian:bookworm-20230725-slim AS steamcmd_builder

RUN apt-get update
RUN apt-get install -y --no-install-recommends bash-static
RUN apt-get install -y --no-install-recommends wget
RUN apt-get install -y --no-install-recommends libtinfo6
RUN apt-get install -y --no-install-recommends libncurses6
RUN apt-get install -y --no-install-recommends ca-certificates

WORKDIR /steamcmd
RUN wget http://media.steampowered.com/installer/steamcmd_linux.tar.gz
RUN tar xf steamcmd_linux.tar.gz
RUN rm -rf steamcmd_linux.tar.gz
RUN ./steamcmd.sh +quit

WORKDIR /steamcmd-libs
RUN cp /lib/i386-linux-gnu/libdl.so.2 ./libdl.so.2
RUN cp /lib/i386-linux-gnu/libgcc_s.so.1 ./libgcc_s.so.1
RUN cp /lib/i386-linux-gnu/libtinfo.so.6.* ./libtinfo.so.5
RUN cp /lib/i386-linux-gnu/libncurses.so.6.* ./libncurses.so.5
RUN cp /usr/lib/i386-linux-gnu/librt.so.1 ./librt.so.1


FROM i386/alpine:3.18.2 AS sdl2_builder

RUN apk update
RUN apk add --no-cache sdl2
WORKDIR /steamcmd-libs

RUN cp /usr/lib/libSDL2* .


FROM i386/busybox:1.36.1-glibc

COPY --from=steamcmd_builder /etc/ssl /etc/ssl
COPY --from=steamcmd_builder /usr/share/zoneinfo/Etc/UTC /etc/localtime

COPY --from=steamcmd_builder /usr/bin/bash-static /bin/bash
COPY --from=steamcmd_builder /steamcmd /var/lib/steamcmd

COPY --from=steamcmd_builder /steamcmd-libs/* /lib
COPY --from=sdl2_builder /steamcmd-libs/* /lib

COPY --from=cli_builder /go/bin/steamcmd-cli /bin/steamcmd-cli

ENV LC_ALL=C

ENV STEAMCMD_HOME=/var/lib/steamcmd

ENV STEAMCMD_SH=${STEAMCMD_HOME}/steamcmd.sh \
    STEAMCMD_SERVER_HOME=${STEAMCMD_HOME}/server

RUN mkdir -p ~/.steam && \
    ln -sf ${STEAMCMD_HOME}/linux32 ~/.steam/sdk32

ENTRYPOINT [ "/bin/steamcmd-cli" ]
