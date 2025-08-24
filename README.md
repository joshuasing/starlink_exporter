# Starlink Prometheus Exporter

[![Go Reference](https://pkg.go.dev/badge/github.com/joshuasing/starlink_exporter.svg)](https://pkg.go.dev/github.com/joshuasing/starlink_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/joshuasing/starlink_exporter)](https://goreportcard.com/report/github.com/joshuasing/starlink_exporter)
[![Go Build Status](https://github.com/joshuasing/starlink_exporter/actions/workflows/go.yml/badge.svg)](https://github.com/joshuasing/starlink_exporter/actions/workflows/go.yml)
[![Starlink Dishy Software Version](https://img.shields.io/badge/Starlink_Dishy_Version-2025.08.11.cr61518_(API_v38)-blue)](internal/spacex/README.md)
[![MIT License](https://img.shields.io/badge/license-MIT-2155cc)](LICENSE)

A simple Starlink exporter for Prometheus. *Not affiliated with Starlink or SpaceX.*

## Metrics

The following metrics are exposed by this exporter:

| Metric name                                        | Description                                                                   |
|----------------------------------------------------|-------------------------------------------------------------------------------|
| `starlink_exporter_scrapes_total`                  | Total number of Starlink dish scrapes                                         |
| `starlink_exporter_scrape_duration_seconds`        | Time taken to scrape metrics from the Starlink dish                           |
| `starlink_dish_up`                                 | Whether scraping metrics from the Starlink dish was successful                |
| `starlink_dish_info`                               | Starlink dish software information                                            |
| `starlink_dish_uptime_seconds`                     | Starlink dish uptime in seconds                                               |
| `starlink_dish_snr_above_noise_floor`              | Whether Starlink dish signal-to-noise ratio is above noise floor              |
| `starlink_dish_snr_persistently_low`               | Whether Starlink dish signal-to-noise ratio is persistently low               |
| `starlink_dish_uplink_throughput_bps`              | Starlink dish uplink throughput in bits/sec                                   |
| `starlink_dish_downlink_throughput_bps`            | Starlink dish downlink throughput in bit/sec                                  |
| `starlink_dish_downlink_throughput_bps_histogram`  | Histogram of Starlink dish downlink throughput over last 15 minutes           |
| `starlink_dish_uplink_throughput_bps_histogram`    | Histogram of Starlink dish uplink throughput in bits/sec over last 15 minutes |
| `starlink_dish_pop_ping_drop_ratio`                | Starlink PoP ping drop ratio                                                  |
| `starlink_dish_pop_ping_latency_seconds`           | Starlink PoP ping latency in seconds                                          |
| `starlink_dish_pop_ping_latency_seconds_histogram` | Histogram of Starlink dish PoP ping latency in seconds over last 15 minutes   |
| `starlink_dish_software_update_reboot_ready`       | Whether the Starlink dish is ready to reboot to apply a software update       |
| `starlink_dish_gps_valid`                          | Whether the Starlink dish GPS is valid                                        |
| `starlink_dish_gps_satellites`                     | Number of GPS satellites visible to the Starlink dish                         |
| `starlink_dish_tilt_angle_deg`                     | Starlink dish tilt angle degrees                                              |
| `starlink_dish_boresight_azimuth_deg`              | Starlink dish boresight azimuth degrees                                       |
| `starlink_dish_boresight_elevation_deg`            | Starlink dish boresight elevation degrees                                     |
| `starlink_dish_desired_boresight_azimuth_deg`      | Starlink dish desired boresight azimuth degrees                               |
| `starlink_dish_desired_boresight_elevation_deg`    | Starlink dish desired boresight elevation degrees                             |
| `starlink_dish_currently_obstructed`               | Whether the Starlink dish is currently obstructed                             |
| `starlink_dish_fraction_obstruction_ratio`         | Fraction of Starlink dish that is obstructed                                  |
| `starlink_dish_last_24h_obstructed_seconds`        | Number of seconds the Starlink dish was obstructed in the past 24 hours       |
| `starlink_dish_power_input_watts_histogram`        | Histogram of Starlink dish power input in watts over last 15 minutes          |
| `starlink_dish_power_input_watts`                  | Current power input for the Starlink dish                                     |
| `starlink_dish_location_info`                      | Dish location information                                                     |
| `starlink_dish_location_latitude_deg`              | Location latitude in degrees                                                  |
| `starlink_dish_location_longitude_deg`             | Location longitude in degrees                                                 |
| `starlink_dish_location_altitude_meters`           | Location altitude in meters above sea level                                   |
| `starlink_dish_alert_unexpected_location`          | Whether the Starlink dish is in an unexpected location                        |
| `starlink_dish_alert_install_pending`              | Whether a Starlink Dish software update is pending installation               |
| `starlink_dish_alert_is_heating`                   | Whether the Starlink dish is heating (snow melting)                           |
| `starlink_dish_alert_is_power_save_idle`           | Whether the Starlink dish is currently in power saving mode                   |
| `starlink_dish_alert_signal_lower_than_predicted`  | Whether the Starlink dish signal is lower than predicted                      |

## Installation

### Binaries

Pre-built binaries for Linux, macOS, Windows and OpenBSD are available
from [GitHub Releases](https://github.com/joshuasing/starlink_exporter/releases).

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

- Go v1.23 or newer (https://go.dev/dl/)

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
