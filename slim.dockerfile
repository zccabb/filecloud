# docker build -f slim.dockerfile -t filecloud/filecloud .
FROM node:16-bullseye AS frontend-builder

COPY frontend /work/filecloud
WORKDIR /work/filecloud

RUN npm --registry=https://registry.npmmirror.com install
RUN npm run build

FROM golang:1.18-bullseye AS backend-builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://goproxy.io

WORKDIR /go/src/github.com/filecloud/filecloud

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /work/filecloud/dist ./frontend/dist
RUN go build -o /go/bin/filecloud \
    -trimpath \
    -ldflags "-s -w" \
    .

FROM alpine:3.17
WORKDIR /

RUN apk --update add ca-certificates \
    mailcap \
    curl

HEALTHCHECK --start-period=2s --interval=5s --timeout=3s \
    CMD curl -f http://localhost/health || exit 1

VOLUME /srv
EXPOSE 80

COPY --from=backend-builder /go/src/github.com/filecloud/filecloud/docker_config.json /.filecloud.json
COPY --from=backend-builder /go/bin/filebcloud /

ENTRYPOINT [ "/filecloud" ]