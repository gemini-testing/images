# Browser Images
[![Build Status](https://github.com/aerokube/images/workflows/build/badge.svg)](https://github.com/aerokube/images/actions?query=workflow%3Abuild)
[![Release](https://img.shields.io/github/release/aerokube/images.svg)](https://github.com/aerokube/images/releases/latest)

**UNMAINTAINED**. Consider https://aerokube.com/moon/latest as alternative.

This repository contains [Docker](http://docker.com/) build files to be used for [Selenoid](http://github.com/aerokube/selenoid) and [Moon](http://github.com/aerokube/moon) projects. You can find prebuilt images [here](https://hub.docker.com/u/selenoid/).

## Download Statistics

### Firefox: [![Firefox Docker Pulls](https://img.shields.io/docker/pulls/selenoid/firefox.svg)](https://hub.docker.com/r/selenoid/firefox)

### Chrome: [![Chrome Docker Pulls](https://img.shields.io/docker/pulls/selenoid/chrome.svg)](https://hub.docker.com/r/selenoid/chrome)

### Opera: [![Opera Docker Pulls](https://img.shields.io/docker/pulls/selenoid/opera.svg)](https://hub.docker.com/r/selenoid/opera)

### Android: [![Android Docker Pulls](https://img.shields.io/docker/pulls/selenoid/android.svg)](https://hub.docker.com/r/selenoid/android)

## Building Images

Moved to: http://aerokube.com/images/latest/#_building_images

### How to build chrome for testing

To build a chrome image for testing, you must specify `--source-dir` option. For example:

```bash
./images chrome -b 138.0.7204.183 -d 138.0.7204.183 -t selenoid/chrome-ft:138.0 --source-dir chrome/for-testing
```

To get the latest version number, visit: https://googlechromelabs.github.io/chrome-for-testing/LATEST_RELEASE_$MAJOR_VERSION (replace `$MAJOR_VERSION` with the desired major version, e.g., `138`).

### How to build Firefox with BiDi support

1. Download the Firefox deb package from https://ftp.mozilla.org/pub/firefox/releases/

2. Build the image with the `--with-bidi-proxy` option:
    ```bash
    ./images firefox -b /path/to/firefox_145.0+build1-0ubuntu0.18.04.1.deb -t selenoid/firefox:145.0-bidi-proxy --with-bidi-proxy
    ```

The `--with-bidi-proxy` option builds an image with a reverse proxy needed for BiDi to work. Selenoid expects both HTTP and BiDi WebSocket endpoints to be available at `<hostname>:4444/session/<session-id>` and chromedriver does exactly that. Geckodriver, however, can only serve them on different ports: one for HTTP and one for BiDi WebSocket. The reverse proxy makes this behavior similar to chromedriver's.

## Image information

Moved to: http://aerokube.com/images/latest/#_browser_image_information
