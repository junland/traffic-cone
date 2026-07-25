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
- `-docker-socket` Path to Docker socket (default: `/var/run/docker.sock`)
- `-haproxy-data-plane-api-address` HAProxy Data Plane API service address (default: `http://127.0.0.1:5555`)
- `-haproxy-data-plane-api-username` HAProxy Data Plane API username
- `-haproxy-data-plane-api-password` HAProxy Data Plane API password

## Build with version metadata

```bash
go build -ldflags "-X traffic-cone/internal/version.Value=v1.0.0" -o bin/traffic-cone ./cmd/traffic-cone
```
