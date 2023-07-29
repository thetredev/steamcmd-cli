FROM golang:1.20.6-alpine3.18 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /go/src

# TODO for releases: $(git describe --tags)

RUN apk add --no-cache git
RUN --mount=type=bind,target=/go/src,source=. \
    go get -u && \
    go build \
        -ldflags "-s -w -X 'github.com/thetredev/steamcmd-cli/shared.Version=$(git rev-parse HEAD)'" \
        -o /go/bin/steamcmd-cli


FROM scratch
COPY --from=builder /go/bin/steamcmd-cli /steamcmd-cli

ENTRYPOINT [ "/steamcmd-cli" ]
