# Starlink Prometheus Exporter

[![Go Reference](https://pkg.go.dev/badge/github.com/joshuasing/starlink_exporter.svg)](https://pkg.go.dev/github.com/joshuasing/starlink_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/joshuasing/starlink_exporter)](https://goreportcard.com/report/github.com/joshuasing/starlink_exporter)
[![Go Build Status](https://github.com/joshuasing/starlink_exporter/actions/workflows/go.yml/badge.svg)](https://github.com/joshuasing/starlink_exporter/actions/workflows/go.yml)
[![MIT License](https://img.shields.io/badge/license-MIT-2155cc)](LICENSE)

A simple Starlink exporter for Prometheus. *Not affiliated with Starlink or SpaceX.*

## Installation

### Binaries

Pre-built binaries are available from [GitHub Releases](https://github.com/joshuasing/starlink_exporter/releases).

You can also use `go install` to build and install a binary from source:
```shell
go install github.com/joshuasing/starlink_exporter@latest
````

**Flags**

```shell
starlink_exporter --help
# Usage of starlink_exporter:
#   -dish string
#         Dish address (default "192.168.100.1:9200")
#   -listen string
#         Listen address (default ":9451")
```

**Example**

```shell
starlink_exporter
# 2024/11/05 12:03:48 INFO Starting Starlink exporter
# 2024/11/05 12:03:48 INFO Connecting to Starlink Dishy address=192.168.100.1:9200
# 2024/11/05 12:03:48 INFO HTTP server listening address=:9451
```

### Docker

Docker images are published to both [GitHub Container Registry (ghcr.io)](https://ghcr.io/joshuasing/starlink_exporter)
and [Docker Hub](https://hub.docker.com/r/joshuasing/starlink_exporter).

```shell
docker run -p 9451:9451 ghcr.io/joshuasing/starlink_exporter:latest
# Status: Downloaded newer image for ghcr.io/joshuasing/starlink_exporter:latest
# 2024/11/05 12:03:48 INFO Starting Starlink exporter
# 2024/11/05 12:03:48 INFO Connecting to Starlink Dishy address=192.168.100.1:9200
# 2024/11/05 12:03:48 INFO HTTP server listening address=:9451
```

### Prometheus

To use the Starlink Prometheus Exporter, you need to configure Prometheus to scrape from the exporter:

```yaml
scrape_configs:
  - job_name: "starlink"
    scrape_interval: 3s # This can be whatever you would like.
    static_configs:
      - targets: [ "localhost:9451" ]
```

*Change `scrape_interval` and the address to match your setup.*

## Contributing

All contributions are welcome! If you have found something you think could be improved, or have discovered additional
metrics you would like included, please feel free to participate by creating an issue or pull request!

### Building

Steps to build starlink_exporter.

**Prerequisites**

- Go v1.22 or newer (https://go.dev/dl/)

**Build**

- Make: `make` (`make deps lint-deps` if you are missing dependencies)
- Standalone: `go build ./cmd/starlink_exporter/`

### Contact

This project is maintained by Joshua Sing. You see a list of ways to contact me on my
website: https://joshuasing.dev/#contact

#### Security vulnerabilities

I take the security of my projects very seriously. As such, I strongly encourage responsible disclosure of security
vulnerabilities.

If you have discovered a security vulnerability in starlink_exporter, please report it in accordance with the
project [Security Policy](SECURITY.md#reporting-a-vulnerability). **Never use GitHub issues to report a security
vulnerability.**

### License

starlink_exporter is distributed under the terms of the MIT License.<br/>
For more information, please refer to the [LICENSE](LICENSE) file.

### Disclaimer

This project is an independent, open-source Prometheus exporter and is not officially associated with, endorsed by, or
in any way affiliated with SpaceX, Starlink, or any of their subsidiaries or affiliates. This project's purpose is to
provide a tool for easily monitoring your Starlink Dishy with Prometheus, and is not authorised or supported by SpaceX
or Starlink in any way.

*SpaceX, Starlink, and any related logos or trademarks are the property of Space Exploration Technologies Corp.*
