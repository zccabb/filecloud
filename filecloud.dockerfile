FROM golang:1.18-bullseye AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://goproxy.io

WORKDIR /go/src/github.com/filecloud/filecloud

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/filecloud \
    -trimpath \
    -ldflags "-s -w" \
    .

FROM alpine:3.17
WORKDIR /

HEALTHCHECK --start-period=2s --interval=5s --timeout=3s \
    CMD curl -f http://localhost/health || exit 1

VOLUME /srv
EXPOSE 80

COPY --from=builder /go/src/github.com/filecloud/filecloud/docker_config.json /.filecloud.json
COPY --from=builder /go/bin/filecloud /

ENTRYPOINT [ "/filecloud" ]