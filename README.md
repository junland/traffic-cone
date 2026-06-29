# traffic-cone

Generic Go CLI daemon scaffold where passing a daemon name starts it in foreground.

## Features

- Foreground daemon runtime driven by a daemon name argument
- Docker event watching via configured Docker socket
- Periodic heartbeat loop (replace with your real workload)
- Cross-platform process status checks (Windows + Unix)
- Simple file-based logging

## Quick start

```bash
go mod tidy
go build -o bin/traffic-cone ./cmd/traffic-cone
```

## CLI usage

```text
traffic-cone <daemon-name> [flags]

Behavior:
  Any daemon name starts that daemon in foreground.
  Use Ctrl+C to stop.
```

Common flags:

- `-pid-file` Path to PID file (default: `./<daemon-name>.pid`)
- `-log-file` Path to log file (default: `./<daemon-name>.log`)
- `-docker-socket` Path to Docker socket (default: `/var/run/docker.sock`)
- `-tick` Heartbeat interval (default: `5s`)

## Project structure

```text
cmd/traffic-cone/main.go       # Entrypoint
internal/app/app.go            # CLI parsing and foreground daemon launch
internal/daemon/manager.go     # Daemon lifecycle and runtime loop
internal/daemon/process_*.go   # Platform-specific process checks
internal/version/version.go    # Build-time version variable
```

## Build with version metadata

```bash
go build -ldflags "-X traffic-cone/internal/version.Value=v1.0.0" -o bin/traffic-cone ./cmd/traffic-cone
```
