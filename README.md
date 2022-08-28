# PacGen

[![Build and publish docker image](https://github.com/nnemirovsky/pacgen/actions/workflows/docker.yml/badge.svg)](https://github.com/nnemirovsky/pacgen/actions/workflows/docker.yml)
[![Test and lint code](https://github.com/nnemirovsky/pacgen/actions/workflows/golang.yml/badge.svg)](https://github.com/nnemirovsky/pacgen/actions/workflows/golang.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/nnemirovsky/pacgen/blob/main/LICENSE)

Simple PAC file generator with REST API for managing proxy profiles
and routing rules. The API reference is described in [swagger spec](api/swagger.yml).
Application uses sqlite3 as config storage.

## Usage

At this point the best way to try this tool is to use the prebuilt docker image.
Simplified example (if possible, use docker-compose instead):

```shell
$ docker pull ghcr.io/nnemirovsky/pacgen:latest
$ docker run -d \
    -v $(pwd)/data:/app/data \
    ghcr.io/nnemirovsky/pacgen:latest migrate up
$ docker run -d \
    -p 8080:8080 \
    -v $(pwd)/data:/app/data \
    ghcr.io/nnemirovsky/pacgen:latest
```

Next you can specify `http(s)://{host:port}/proxy.pac` as a PAC file address.
