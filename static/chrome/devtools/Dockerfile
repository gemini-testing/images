FROM golang:1.22-bullseye as go

COPY . /devtools

RUN \
    cd /devtools && \
    go test -race && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"

FROM scratch

COPY --from=go /devtools/devtools /usr/bin/
