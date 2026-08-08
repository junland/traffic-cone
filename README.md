# traffic-cone

## Overview

A lightweight sidecar process that bridges container runtimes with HAProxy to automate traffic management. By monitoring container lifecycles via the Docker and Podman sockets, it dynamically updates routing rules using the HAProxy DataPlane API—delivering Traefik-like automatic service discovery with the performance, reliability, and enterprise feature set of HAProxy.

## Key Features

- **Dynamic Service Discovery:** Listens to real-time events from Docker and Podman container engines to discover, register, and deregister microservices instantly.
- **HAProxy DataPlane API Integration:** Programmatically configures HAProxy backends, frontends, and ACLs on the fly without requiring process reloads or dropping active connections.
- **Multi-Runtime Support:** Seamless compatibility with both Docker and daemonless Podman environments out of the box.
- **Label-Driven Configuration:** Declarative routing rules defined directly via container labels (e.g., host matching, path prefixing, port mapping).
- **Zero-Downtime Reloads:** Updates traffic routing rules hot-in-memory to ensure zero dropped packets during container deployments or scaling events.
- **Sidecar Deployment Model:** Operates as a decoupled sidecar alongside HAProxy, maint

## Quick start

```bash
go mod tidy
go build -o bin/traffic-cone ./cmd/traffic-cone
```

## CLI usage

```text
traffic-cone [flags]

Behavior:
  Starts the program in the foreground.
  Use Ctrl+C to stop.
```

Common flags:

- `-pid-file` Path to PID file (default: `<temp-dir>/traffic-cone.pid`)
- `-docker-socket` Path to Docker socket (default: `/var/run/docker.sock`)
- `-haproxy-data-plane-api-address` HAProxy Data Plane API service address (default: `http://127.0.0.1:5555`)
- `-haproxy-data-plane-api-username` HAProxy Data Plane API username
- `-haproxy-data-plane-api-password-file` Path to file containing HAProxy Data Plane API password (or set `HAPROXY_DATA_PLANE_API_PASSWORD`)

## Build with version metadata

```bash
go build -ldflags "-X traffic-cone/internal/version.Value=v1.0.0" -o bin/traffic-cone ./cmd/traffic-cone
```
