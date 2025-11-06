# dtail

A small CLI tool to tail logs from Docker containers by name or prefix.

This repository contains a lightweight command-line utility written in Go that watches and colorizes logs from Docker containers. It expects a `DockerService` type to provide container discovery and log streaming. A simple stub implementation can be used for development.

Features
- Tail the last N lines from one or more containers.
- Colorized per-container output for easy visual separation.
- Graceful shutdown on Ctrl+C / SIGTERM.

Requirements
- Go 1.20+ (or compatible)
- Docker daemon if you use a real Docker integration

Quick build

```bash
go build -o dtail
```

Quick run

```bash
# run with go
go run ./main.go example-container

# run built binary with 50 lines of tail
./dtail -n 50 example-container another-prefix
```

Flags
- `-n`, `--tail` — Number of lines to show from the end of the logs per container (default: 10)

How it matches containers
- You may pass one or more container names or prefixes as CLI arguments.
- The program looks up running containers and matches either exact names or prefixes.

License
- MIT

