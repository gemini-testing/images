# Building a layer with `/usr/bin/devtools` that works with Rosetta

## Why?

Context:
- QEMU has issues with upx-packed binaries: https://github.com/upx/upx/issues/683
- Browser images built in this repo need to work on Mac, which uses emulation

To make existing images work on Mac, a separate layer can be built with devtools binary not packed, replacing the original binary before launch.

## Building & publishing

To build, run:

```shell
docker build --platform linux/amd64 .
```

To publish, first follow the [instructions](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#pushing-container-images) to log in to GitHub containers register and then run:
```shell
docker push ghcr.io/gemini-testing/devtools-proxy:latest
```

## Usage

This image is supposed to be used as patch on top of existing ones, for example, like this:

```dockerfile
FROM ghcr.io/gemini-testing/devtools-proxy:latest AS devtools

FROM ghcr.io/gemini-testing/browsers/chrome-ft:152.0

COPY --from=devtools /usr/bin/devtools /usr/bin/devtools
```

This Dockerfile overrides problematic upx-packed devtools proxy binary with the working one stored in our image.
