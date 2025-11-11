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

## Image information

Moved to: http://aerokube.com/images/latest/#_browser_image_information
