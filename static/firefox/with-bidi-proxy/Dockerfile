ARG VERSION
FROM golang:1.22-bullseye as go

WORKDIR /proxy
COPY bidi-proxy /proxy

RUN \
    go test ./... && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /firefox-bidi-proxy

FROM selenoid/dev_firefox:$VERSION

COPY geckodriver /usr/bin/
COPY selenoid /usr/bin/
COPY --from=go /firefox-bidi-proxy /usr/bin/
COPY --chown=selenium:root browsers.json /home/selenium/
COPY entrypoint.sh /

RUN \
    chmod +x /usr/bin/geckodriver /usr/bin/selenoid /usr/bin/firefox-bidi-proxy /entrypoint.sh

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
