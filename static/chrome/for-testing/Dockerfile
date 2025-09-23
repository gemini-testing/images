FROM browsers/base:7.4.2

ARG VERSION
ARG PACKAGE=chrome-for-testing

LABEL browser=$PACKAGE:$VERSION

RUN \
    apt-get update && \
    apt-get -y --no-install-recommends install gconf-service \
        libasound2 \
        libatk1.0-0 \
        libc6 \
        libcairo2 \
        libcups2 \
        libdbus-1-3 \
        libexpat1 \
        libfontconfig1 \
        libfreetype6 \
        libgcc1 \
        libgconf-2-4 \
        libgdk-pixbuf2.0-0 \
        libglib2.0-0 \
        libgtk2.0-0 \
        libgtk-3-0 \
        libnspr4 \
        libnss3 \
        libpango1.0-0 \
        libstdc++6 \
        libx11-6 \
        libx11-xcb1 \
        libxcb1 \
        libxcomposite1 \
        libxcursor1 \
        libxdamage1 \
        libxext6 \
        libxfixes3 \
        libxi6 \
        libxrandr2 \
        libxrender1 \
        libxss1 \
        libxtst6 \
        libgbm1 \
        libu2f-udev \
        libvulkan1 \
        ca-certificates \
        fonts-liberation \
        libappindicator3-1 \
        libnss3 \
        lsb-base \
        xdg-utils \
        libcurl4 \
        iproute2 \
        unzip && \
        rm -Rf /tmp/* && \
        rm -Rf /var/lib/apt/lists/*

RUN \
    CHROME_MAJOR_VERSION=$(echo $VERSION | sed -E "s/([0-9]+)\..*/\1/") \
    && CHROME_VERSION="$(wget -qO- https://googlechromelabs.github.io/chrome-for-testing/LATEST_RELEASE_${CHROME_MAJOR_VERSION} | sed 's/\r$//')" \
    && wget -qO chrome-linux64.zip "https://storage.googleapis.com/chrome-for-testing-public/${CHROME_VERSION}/linux64/chrome-linux64.zip" \
    && mkdir -p /opt/google/chrome \
    && unzip chrome-linux64.zip -d /tmp \
    && cp -a /tmp/chrome-linux64/. /opt/google/chrome/ \
    && rm -rf /var/lib/apt/lists/* /var/cache/apt/* /tmp/* \
    && ln -s /opt/google/chrome/chrome /usr/bin/google-chrome
