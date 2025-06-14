FROM browsers/base:7.4.2

ARG VERSION
ARG PACKAGE=google-chrome-stable
ARG INSTALL_DIR=chrome

LABEL browser=$PACKAGE:$VERSION

# Install Chrome
# Ref https://github.com/SeleniumHQ/docker-selenium/blob/trunk/NodeChrome/Dockerfile

RUN \
    CHROME_MAJOR_VERSION=$(echo $VERSION | sed -E "s/([0-9]+)\..*/\1/") \
    && CHROME_VERSION="${PACKAGE}_$(wget -qO- https://googlechromelabs.github.io/chrome-for-testing/LATEST_RELEASE_${CHROME_MAJOR_VERSION} | sed 's/\r$//')-1" \
    && apt-get update \
    && wget -qO google-chrome.deb "https://dl.google.com/linux/chrome/deb/pool/main/g/${PACKAGE}/${CHROME_VERSION}_amd64.deb" \
    && apt-get -qqy --no-install-recommends install --allow-downgrades ./google-chrome.deb \
    && rm -rf google-chrome.deb \
    && chown root:root /opt/google/chrome/chrome-sandbox \
    && chmod 4755 /opt/google/chrome/chrome-sandbox \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /var/cache/apt/* /tmp/*
