# traffic-cone

A daemon that listens for Docker events and prints the container names to stdout.

## Features

- Docker event watching via configured Docker socket

## Quick start

```bash
go mod tidy
go build -o bin/traffic-cone ./cmd/traffic-cone
```

## CLI usage

```text
traffic-cone <daemon-name> [flags]

Behavior:
  Starts the program in the foreground.
  Use Ctrl+C to stop.
```

Common flags:

- `-pid-file` Path to PID file (default: `./<daemon-name>.pid`)
- `-log-file` Path to log file (default: `./<daemon-name>.log`)
- `-docker-socket` Path to Docker socket (default: `/var/run/docker.sock`)

## Build with version metadata

```bash
go build -ldflags "-X traffic-cone/internal/version.Value=v1.0.0" -o bin/traffic-cone ./cmd/traffic-cone
```
