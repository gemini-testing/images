# syntax=docker/dockerfile:1

ARG DEBIAN_RELEASE=trixie

FROM golang:1.22-bookworm AS tools

COPY base/xseld /src/xseld
COPY base/fileserver /src/fileserver

RUN \
    cd /src/xseld && CGO_ENABLED=0 go build -ldflags="-s -w" -o /usr/local/bin/xseld && \
    cd /src/fileserver && CGO_ENABLED=0 go build -ldflags="-s -w" -o /usr/local/bin/fileserver

FROM debian:bookworm-slim

ARG DEBIAN_RELEASE
ARG DEBIAN_SUITE
ARG SNAPSHOT
ARG VERSION

LABEL browser="chromium:${VERSION}"

RUN \
    rm -f /etc/apt/sources.list.d/* && \
    echo "deb [check-valid-until=no] http://snapshot.debian.org/archive/debian/${SNAPSHOT}/ ${DEBIAN_RELEASE} main" > /etc/apt/sources.list && \
    if [ "$DEBIAN_SUITE" != "$DEBIAN_RELEASE" ]; then \
        echo "deb [check-valid-until=no] http://snapshot.debian.org/archive/debian/${SNAPSHOT}/ ${DEBIAN_SUITE} main" >> /etc/apt/sources.list; \
    fi && \
    apt-get -o APT::Update::Error-Mode=any -o Acquire::Check-Valid-Until=false -o Acquire::Retries=5 update && \
    DEBIAN_FRONTEND=noninteractive apt-get -y --no-install-recommends install \
        ca-certificates \
        chromium="${VERSION}" \
        chromium-driver="${VERSION}" \
        feh \
        fluxbox \
        fonts-dejavu-core \
        fonts-liberation \
        fonts-noto-color-emoji \
        iproute2 \
        libnss3-tools \
        libnss-wrapper \
        openssl \
        procps \
        wmctrl \
        x11-utils \
        x11vnc \
        xsel \
        xvfb && \
    chromium --version && \
    chromedriver --version && \
    adduser --system --home /home/selenium --uid 4096 --ingroup root --disabled-password --shell /bin/bash selenium && \
    mkdir -p /home/selenium/Downloads /home/selenium/.fluxbox && \
    chgrp -R 0 /home/selenium && \
    chmod -R g=u /home/selenium && \
    ln -sf /bin/true /usr/bin/xdg-open && \
    apt-get clean && \
    rm -rf /tmp/* /var/lib/apt/lists/*

RUN \
    sed -i 's|http://snapshot.debian.org|https://snapshot.debian.org|' /etc/apt/sources.list && \
    apt-get -o APT::Update::Error-Mode=any -o Acquire::Check-Valid-Until=false -o Acquire::Retries=5 update && \
    DEBIAN_FRONTEND=noninteractive apt-get -y --no-install-recommends install xauth && \
    echo 'CHROMIUM_FLAGS="$CHROMIUM_FLAGS --no-sandbox --disable-gpu"' > /etc/chromium.d/docker && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /tmp/.X11-unix && \
    chown selenium:root /tmp/.X11-unix && \
    chmod 755 /tmp/.X11-unix

COPY base/fluxbox /usr/share/fluxbox/styles/
COPY --chown=selenium:root base/fluxbox /home/selenium/.fluxbox/
COPY base/aerokube.png /usr/share/images/fluxbox/
COPY --from=tools /usr/local/bin/fileserver /usr/local/bin/xseld /usr/bin/
