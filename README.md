# Browser Images

This repository contains [Docker](http://docker.com/) build files to be used for testing, with [Selenoid](http://github.com/aerokube/selenoid) or on their own. You can find prebuilt images [here](https://github.com/orgs/gemini-testing/packages).

## Building Images

See also: http://aerokube.com/images/latest/#_building_images

### How to build chrome for testing

To build a chrome image for testing, you must specify `--source-dir` option. For example:

```bash
./images chrome -b 138.0.7204.183 -d 138.0.7204.183 -t selenoid/chrome-ft:138.0 --source-dir chrome/for-testing
```

To get the latest version number, visit: https://googlechromelabs.github.io/chrome-for-testing/LATEST_RELEASE_$MAJOR_VERSION (replace `$MAJOR_VERSION` with the desired major version, e.g., `138`).

### How to build cross-platform Chromium image

Note: this can be done automatically by the GitHub Actions workflow in this repository.

Use `--debian` together with a Chromium major and the native Docker architecture:

```bash
./images chromium --debian --architecture arm64 -b 126 -t ghcr.io/gemini-testing/browsers/chromium:126.0-arm64
```

Supported architectures are `amd64` and `arm64`. Without `--debian`, image won't be cross-platform and will be based on Ubuntu.

After building images for two platforms, you can merge them like this:
```bash
docker buildx imagetools create \
    --tag ghcr.io/gemini-testing/browsers/chromium:126.0 \
    ghcr.io/gemini-testing/browsers/chromium:126.0-amd64 \
    ghcr.io/gemini-testing/browsers/chromium:126.0-arm64
```

### How to build Firefox

1. Download the Firefox deb package from https://ftp.mozilla.org/pub/firefox/releases/

2. Build the image:
    ```bash
    # install pkger tool
    go install github.com/markbates/pkger/cmd/pkger@latest
    export PATH="$(go env GOPATH)/bin:$PATH"

    go generate ./...

    go build -o images .

    ./images firefox -b /path/to/firefox_145.0+build1-0ubuntu0.18.04.1.deb -t selenoid/firefox:145.0-bidi-proxy
    ```

For Firefox `94+` the build automatically uses an image with a reverse proxy needed for BiDi to work. Selenoid expects both HTTP and BiDi WebSocket endpoints to be available at `<hostname>:4444/session/<session-id>` and chromedriver does exactly that. Geckodriver, however, can only serve them on different ports: one for HTTP and one for BiDi WebSocket. The reverse proxy makes this behavior similar to chromedriver's.

## Image information

Moved to: http://aerokube.com/images/latest/#_browser_image_information
