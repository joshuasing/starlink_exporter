# Starlink Prometheus Exporter

[![Go Reference](https://pkg.go.dev/badge/github.com/joshuasing/starlink_exporter.svg)](https://pkg.go.dev/github.com/joshuasing/starlink_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/joshuasing/starlink_exporter)](https://goreportcard.com/report/github.com/joshuasing/starlink_exporter)
[![Go Build Status](https://github.com/joshuasing/starlink_exporter/actions/workflows/go.yml/badge.svg)](https://github.com/joshuasing/starlink_exporter/actions/workflows/go.yml)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/joshuasing/starlink_exporter/badge)](https://scorecard.dev/viewer/?uri=github.com/joshuasing/starlink_exporter)
[![Starlink Dishy Software Version](https://img.shields.io/badge/Starlink_Dishy_Version-2026.05.30.cr80794_(API_v42)-blue)](internal/spacex_api/README.md)
[![MIT License](https://img.shields.io/badge/license-MIT-2155cc)](LICENSE)

A simple Starlink exporter for Prometheus. *Not affiliated with Starlink or SpaceX.*

## Metrics

The following metrics are exposed by this exporter:

| Metric name                                                       | Description                                                                                 |
|-------------------------------------------------------------------|---------------------------------------------------------------------------------------------|
| `starlink_exporter_scrapes_total`                                 | Total number of Starlink dish scrapes                                                       |
| `starlink_exporter_scrape_duration_seconds`                       | Time taken to scrape metrics from the Starlink dish                                         |
| `starlink_dish_up`                                                | Whether scraping metrics from the Starlink dish was successful                              |
| `starlink_dish_info`                                              | Starlink dish software information                                                          |
| `starlink_dish_uptime_seconds`                                    | Starlink dish uptime in seconds                                                             |
| `starlink_dish_mobility_class`                                    | Starlink dish mobility class                                                                |
| `starlink_dish_snr_above_noise_floor`                             | Whether Starlink dish signal-to-noise ratio is above noise floor                            |
| `starlink_dish_snr_persistently_low`                              | Whether Starlink dish signal-to-noise ratio is persistently low                             |
| `starlink_dish_uplink_throughput_bps`                             | Starlink dish uplink throughput in bits/sec                                                 |
| `starlink_dish_downlink_throughput_bps`                           | Starlink dish downlink throughput in bit/sec                                                |
| `starlink_dish_downlink_throughput_bps_histogram`                 | Histogram of Starlink dish downlink throughput over last 15 minutes                         |
| `starlink_dish_uplink_throughput_bps_histogram`                   | Histogram of Starlink dish uplink throughput in bits/sec over last 15 minutes               |
| `starlink_dish_pop_ping_drop_ratio`                               | Starlink PoP ping drop ratio                                                                |
| `starlink_dish_pop_ping_latency_seconds`                          | Starlink PoP ping latency in seconds                                                        |
| `starlink_dish_pop_ping_latency_seconds_histogram`                | Histogram of Starlink dish PoP ping latency in seconds over last 15 minutes                 |
| `starlink_dish_software_update_state`                             | Starlink dish update state                                                                  |
| `starlink_dish_software_update_reboot_ready`                      | Whether the Starlink dish is ready to reboot to apply a software update                     |
| `starlink_dish_gps_valid`                                         | Whether the Starlink dish GPS is valid                                                      |
| `starlink_dish_gps_satellites`                                    | Number of GPS satellites visible to the Starlink dish                                       |
| `starlink_dish_tilt_angle_deg`                                    | Starlink dish tilt angle degrees                                                            |
| `starlink_dish_boresight_azimuth_deg`                             | Starlink dish boresight azimuth degrees                                                     |
| `starlink_dish_boresight_elevation_deg`                           | Starlink dish boresight elevation degrees                                                   |
| `starlink_dish_desired_boresight_azimuth_deg`                     | Starlink dish desired boresight azimuth degrees                                             |
| `starlink_dish_desired_boresight_elevation_deg`                   | Starlink dish desired boresight elevation degrees                                           |
| `starlink_dish_currently_obstructed`                              | Whether the Starlink dish is currently obstructed                                           |
| `starlink_dish_fraction_obstruction_ratio`                        | Fraction of Starlink dish that is obstructed                                                |
| `starlink_dish_last_24h_obstructed_seconds`                       | Number of seconds the Starlink dish was obstructed in the past 24 hours                     |
| `starlink_dish_power_input_watts_histogram`                       | Histogram of Starlink dish power input in watts over last 15 minutes                        |
| `starlink_dish_power_input_watts`                                 | Current power input for the Starlink dish                                                   |
| `starlink_dish_location_info`                                     | Dish location information                                                                   |
| `starlink_dish_location_latitude_deg`                             | Location latitude in degrees                                                                |
| `starlink_dish_location_longitude_deg`                            | Location longitude in degrees                                                               |
| `starlink_dish_location_altitude_meters`                          | Location altitude in meters above sea level                                                 |
| `starlink_dish_alert_unexpected_location`                         | Whether the Starlink dish is in an unexpected location                                      |
| `starlink_dish_alert_install_pending`                             | Whether a Starlink Dish software update is pending installation                             |
| `starlink_dish_alert_is_heating`                                  | Whether the Starlink dish is heating (snow melting)                                         |
| `starlink_dish_alert_is_power_save_idle`                          | Whether the Starlink dish is currently in power saving mode                                 |
| `starlink_dish_alert_signal_lower_than_predicted`                 | Whether the Starlink dish signal is lower than predicted                                    |
| `starlink_dish_alert_motors_stuck`                                | Whether the Starlink dish motors are stuck                                                  |
| `starlink_dish_alert_thermal_throttle`                            | Whether the Starlink dish is thermally throttled                                            |
| `starlink_dish_alert_thermal_shutdown`                            | Whether the Starlink dish has thermally shut down                                           |
| `starlink_dish_alert_mast_not_near_vertical`                      | Whether the Starlink dish mast is not near vertical                                         |
| `starlink_dish_alert_slow_ethernet_speeds`                        | Whether the Starlink dish ethernet link is negotiated below gigabit                         |
| `starlink_dish_alert_slow_ethernet_speeds_100`                    | Whether the Starlink dish ethernet link is negotiated at 100 Mbps or lower                  |
| `starlink_dish_alert_roaming`                                     | Whether the Starlink dish is roaming                                                        |
| `starlink_dish_alert_power_supply_thermal_throttle`               | Whether the Starlink dish power supply is thermally throttled                               |
| `starlink_dish_alert_dbf_telem_stale`                             | Whether the Starlink dish DBF telemetry is stale                                            |
| `starlink_dish_alert_low_motor_current`                           | Whether the Starlink dish has low motor current                                             |
| `starlink_dish_alert_obstruction_map_reset`                       | Whether the Starlink dish obstruction map was reset                                         |
| `starlink_dish_alert_dish_water_detected`                         | Whether water has been detected inside the Starlink dish                                    |
| `starlink_dish_alert_router_water_detected`                       | Whether water has been detected inside the Starlink router                                  |
| `starlink_dish_alert_upsu_router_port_slow`                       | Whether the Starlink dish UPSU router port is negotiated below the expected speed           |
| `starlink_dish_alert_no_ethernet_link`                            | Whether the Starlink dish has no ethernet link                                              |
| `starlink_dish_seconds_to_first_nonempty_slot`                    | Seconds until the next non-empty network schedule slot                                      |
| `starlink_dish_stow_requested`                                    | Whether a stow has been requested for the Starlink dish                                     |
| `starlink_dish_eth_speed_mbps`                                    | Negotiated speed of the Starlink dish ethernet link in Mbps                                 |
| `starlink_dish_class_of_service`                                  | Starlink dish class of service                                                              |
| `starlink_dish_reboot_reason`                                     | Reason for the most recent Starlink dish reboot                                             |
| `starlink_dish_disablement_code`                                  | Reason the Starlink dish is disabled, if any                                                |
| `starlink_dish_dl_bandwidth_restricted_reason`                    | Reason downlink bandwidth is currently restricted, if any                                   |
| `starlink_dish_ul_bandwidth_restricted_reason`                    | Reason uplink bandwidth is currently restricted, if any                                     |
| `starlink_dish_is_cell_disabled`                                  | Whether the Starlink dish cell is disabled                                                  |
| `starlink_dish_seconds_until_swupdate_reboot_possible`            | Seconds until the Starlink dish can reboot to apply a software update                       |
| `starlink_dish_high_power_test_mode`                              | Whether the Starlink dish is in high-power test mode                                        |
| `starlink_dish_is_moving_fast_persisted`                          | Whether the Starlink dish has persistently detected fast movement                           |
| `starlink_dish_mac_flag`                                          | Starlink dish MAC flag                                                                      |
| `starlink_dish_nat_flag`                                          | Starlink dish NAT flag                                                                      |
| `starlink_dish_account_shard`                                     | Starlink dish account shard                                                                 |
| `starlink_dish_connected_routers_count`                           | Number of routers currently connected to the Starlink dish                                  |
| `starlink_dish_has_actuators`                                     | Whether the Starlink dish has actuators                                                     |
| `starlink_dish_ned2dish_quaternion_w`                             | W component of the NED-to-dish orientation quaternion                                       |
| `starlink_dish_ned2dish_quaternion_x`                             | X component of the NED-to-dish orientation quaternion                                       |
| `starlink_dish_ned2dish_quaternion_y`                             | Y component of the NED-to-dish orientation quaternion                                       |
| `starlink_dish_ned2dish_quaternion_z`                             | Z component of the NED-to-dish orientation quaternion                                       |
| `starlink_dish_outage_info`                                       | Information about the current Starlink dish outage, if any                                  |
| `starlink_dish_outage_start_timestamp_seconds`                    | Unix timestamp at which the current outage started                                          |
| `starlink_dish_outage_duration_seconds`                           | Duration of the current outage in seconds                                                   |
| `starlink_dish_gps_no_sats_after_ttff`                            | Whether the Starlink dish has no GPS satellites after time-to-first-fix                     |
| `starlink_dish_gps_inhibit`                                       | Whether GPS is currently inhibited on the Starlink dish                                     |
| `starlink_dish_obstruction_valid_seconds`                         | Seconds the Starlink dish obstruction map has been collecting data                          |
| `starlink_dish_obstruction_patches_valid`                         | Number of valid patches in the Starlink dish obstruction map                                |
| `starlink_dish_avg_prolonged_obstruction_duration_seconds`        | Average duration of a prolonged obstruction in seconds                                      |
| `starlink_dish_avg_prolonged_obstruction_interval_seconds`        | Average interval between prolonged obstructions in seconds                                  |
| `starlink_dish_avg_prolonged_obstruction_valid`                   | Whether the average prolonged obstruction metrics are currently valid                       |
| `starlink_dish_ready_state_cady`                                  | Whether the cady subsystem is ready                                                         |
| `starlink_dish_ready_state_scp`                                   | Whether the SCP subsystem is ready                                                          |
| `starlink_dish_ready_state_l1l2`                                  | Whether the L1/L2 subsystem is ready                                                        |
| `starlink_dish_ready_state_xphy`                                  | Whether the XPHY subsystem is ready                                                         |
| `starlink_dish_ready_state_aap`                                   | Whether the AAP subsystem is ready                                                          |
| `starlink_dish_ready_state_rf`                                    | Whether the RF subsystem is ready                                                           |
| `starlink_dish_software_update_progress`                          | Progress of the in-flight Starlink dish software update (0-1)                               |
| `starlink_dish_software_update_requires_reboot`                   | Whether the pending Starlink dish software update requires a reboot                         |
| `starlink_dish_software_update_reboot_scheduled_utc_seconds`      | Unix timestamp at which the Starlink dish is scheduled to reboot to apply a software update |
| `starlink_dish_actuator_state`                                    | Current state of the Starlink dish actuators                                                |
| `starlink_dish_attitude_estimation_state`                         | Current state of the Starlink dish attitude estimator                                       |
| `starlink_dish_attitude_uncertainty_deg`                          | Starlink dish attitude uncertainty in degrees                                               |
| `starlink_dish_initialization_attitude_seconds`                   | Seconds spent on attitude initialization during dish boot                                   |
| `starlink_dish_initialization_burst_detected_seconds`             | Seconds until first burst was detected during dish boot                                     |
| `starlink_dish_initialization_ekf_converged_seconds`              | Seconds until the EKF converged during dish boot                                            |
| `starlink_dish_initialization_first_cplane_seconds`               | Seconds until first cplane message during dish boot                                         |
| `starlink_dish_initialization_first_pop_ping_seconds`             | Seconds until first PoP ping during dish boot                                               |
| `starlink_dish_initialization_gps_valid_seconds`                  | Seconds until GPS became valid during dish boot                                             |
| `starlink_dish_initialization_initial_network_entry_seconds`      | Seconds until initial network entry during dish boot                                        |
| `starlink_dish_initialization_network_schedule_seconds`           | Seconds until first network schedule was received during dish boot                          |
| `starlink_dish_initialization_rf_ready_seconds`                   | Seconds until the RF subsystem was ready during dish boot                                   |
| `starlink_dish_initialization_stable_connection_seconds`          | Seconds until a stable connection was established during dish boot                          |
| `starlink_dish_plc_receiving`                                     | Whether the dish is receiving Power-Line Communication data                                 |
| `starlink_dish_plc_average_time_to_empty_seconds`                 | Average time until the battery is empty in seconds                                          |
| `starlink_dish_plc_average_time_to_full_seconds`                  | Average time until the battery is full in seconds                                           |
| `starlink_dish_plc_battery_health`                                | Battery health as reported by PLC                                                           |
| `starlink_dish_plc_permanent_failure`                             | Whether the PLC battery has reported a permanent failure                                    |
| `starlink_dish_plc_safety_mode_active`                            | Whether the PLC battery is in safety mode                                                   |
| `starlink_dish_plc_state_of_charge_percent`                       | PLC battery state of charge in percent                                                      |
| `starlink_dish_plc_thermal_throttle_level`                        | PLC thermal throttle level                                                                  |
| `starlink_dish_plc_revision`                                      | PLC protocol revision                                                                       |
| `starlink_dish_upsu_dish_power_watts`                             | Dish power draw as reported by the UPSU in watts                                            |
| `starlink_dish_upsu_router_power_watts`                           | Router power draw as reported by the UPSU in watts                                          |
| `starlink_dish_upsu_uptime_seconds`                               | UPSU uptime in seconds                                                                      |
| `starlink_dish_upsu_board_rev`                                    | UPSU board revision                                                                         |
| `starlink_dish_aps_dish_power_watts`                              | Dish power draw as reported by the APS in watts                                             |
| `starlink_dish_aps_uptime_seconds`                                | APS uptime in seconds                                                                       |
| `starlink_dish_aps_board_rev`                                     | APS board revision                                                                          |
| `starlink_dish_pop_ping_drop_ratio_histogram`                     | Histogram of Starlink dish PoP ping drop ratio over the last 15 minutes                     |
| `starlink_dish_history_outages_count`                             | Number of outages retained in the Starlink dish history buffer                              |
| `starlink_dish_location_uncertainty_meters`                       | 1-sigma uncertainty of the reported location in meters                                      |
| `starlink_dish_horizontal_speed_mps`                              | Horizontal speed of the Starlink dish in meters per second                                  |
| `starlink_dish_vertical_speed_mps`                                | Vertical speed of the Starlink dish in meters per second                                    |
| `starlink_wifi_up`                                                | Whether scraping metrics from the Starlink WiFi router was successful                       |
| `starlink_wifi_info`                                              | Starlink WiFi router software information                                                   |
| `starlink_wifi_uptime_seconds`                                    | Starlink WiFi router uptime in seconds                                                      |
| `starlink_wifi_hops_from_controller`                              | Number of mesh hops between this router and the controller                                  |
| `starlink_wifi_no_wan_link`                                       | Whether the Starlink WiFi router has no WAN link                                            |
| `starlink_wifi_is_aviation`                                       | Whether the Starlink WiFi router is in aviation mode                                        |
| `starlink_wifi_is_aviation_conformed`                             | Whether the Starlink WiFi router is aviation conformed                                      |
| `starlink_wifi_using_individualized_calibration`                  | Whether the Starlink WiFi router is using individualized calibration                        |
| `starlink_wifi_calibration_partitions_state`                      | Starlink WiFi router calibration partitions state                                           |
| `starlink_wifi_dish_disablement_code`                             | Disablement code reported by the dish as seen by the WiFi router                            |
| `starlink_wifi_seconds_since_last_public_ipv4_change`             | Seconds since the public IPv4 address last changed                                          |
| `starlink_wifi_clients_count`                                     | Number of clients currently connected to the Starlink WiFi router                           |
| `starlink_wifi_dhcp_servers_count`                                | Number of DHCP servers reported by the Starlink WiFi router                                 |
| `starlink_wifi_ping_drop_ratio`                                   | Router WAN ping drop ratio                                                                  |
| `starlink_wifi_ping_drop_ratio_5m`                                | Router WAN ping drop ratio over the last 5 minutes                                          |
| `starlink_wifi_ping_latency_seconds`                              | Router WAN ping latency in seconds                                                          |
| `starlink_wifi_dish_ping_drop_ratio`                              | Router-to-dish ping drop ratio                                                              |
| `starlink_wifi_dish_ping_drop_ratio_5m`                           | Router-to-dish ping drop ratio over the last 5 minutes                                      |
| `starlink_wifi_dish_ping_latency_seconds`                         | Router-to-dish ping latency in seconds                                                      |
| `starlink_wifi_pop_ping_drop_ratio`                               | Router-to-PoP ping drop ratio                                                               |
| `starlink_wifi_pop_ping_drop_ratio_5m`                            | Router-to-PoP ping drop ratio over the last 5 minutes                                       |
| `starlink_wifi_pop_ping_latency_seconds`                          | Router-to-PoP ping latency in seconds                                                       |
| `starlink_wifi_pop_ipv6_ping_drop_ratio`                          | Router-to-PoP IPv6 ping drop ratio                                                          |
| `starlink_wifi_pop_ipv6_ping_drop_ratio_5m`                       | Router-to-PoP IPv6 ping drop ratio over the last 5 minutes                                  |
| `starlink_wifi_pop_ipv6_ping_latency_seconds`                     | Router-to-PoP IPv6 ping latency in seconds                                                  |
| `starlink_wifi_alert_thermal_throttle`                            | Whether the Starlink WiFi router is thermally throttled                                     |
| `starlink_wifi_alert_install_pending`                             | Whether a Starlink WiFi router software update is pending installation                      |
| `starlink_wifi_alert_freshly_fused`                               | Whether the Starlink WiFi router was recently fused                                         |
| `starlink_wifi_alert_lan_eth_slow_link_10`                        | Whether a Starlink WiFi router LAN ethernet link is negotiated at 10 Mbps                   |
| `starlink_wifi_alert_lan_eth_slow_link_100`                       | Whether a Starlink WiFi router LAN ethernet link is negotiated at 100 Mbps                  |
| `starlink_wifi_alert_high_cable_ping_drop_rate`                   | Whether the Starlink WiFi router is reporting a high cable ping drop rate                   |
| `starlink_wifi_alert_wan_eth_poor_connection`                     | Whether the Starlink WiFi router WAN ethernet connection is poor                            |
| `starlink_wifi_alert_mesh_topology_changing_often`                | Whether the Starlink WiFi mesh topology is changing often                                   |
| `starlink_wifi_alert_mesh_unreliable_backhaul`                    | Whether the Starlink WiFi mesh backhaul is unreliable                                       |
| `starlink_wifi_alert_radius_missing_process`                      | Whether the Starlink WiFi router RADIUS process is missing                                  |
| `starlink_wifi_alert_eth_switch_error`                            | Whether the Starlink WiFi router ethernet switch has reported an error                      |
| `starlink_wifi_alert_poe_on_dish_unreachable`                     | Whether the dish is unreachable while PoE is on                                             |
| `starlink_wifi_alert_poe_fuse_blown`                              | Whether the Starlink WiFi router PoE fuse is blown                                          |
| `starlink_wifi_alert_poe_router_overcurrent`                      | Whether the Starlink WiFi router PoE has detected an overcurrent                            |
| `starlink_wifi_alert_poe_off_current_nominal`                     | Whether the Starlink WiFi router is drawing nominal current while PoE is off                |
| `starlink_wifi_alert_poe_vin_overvoltage`                         | Whether the Starlink WiFi router PoE input voltage is over the threshold                    |
| `starlink_wifi_alert_poe_vin_undervoltage`                        | Whether the Starlink WiFi router PoE input voltage is under the threshold                   |
| `starlink_wifi_alert_sandbox_disabled`                            | Whether the Starlink WiFi router client sandbox is disabled                                 |
| `starlink_wifi_alert_only_overflight_blocked`                     | Whether only overflight is blocked on the Starlink WiFi router                              |
| `starlink_wifi_alert_offline_networks_disabled`                   | Whether offline networks are disabled on the Starlink WiFi router                           |
| `starlink_wifi_alert_wired_mesh_not_using_wan_iface`              | Whether wired mesh is not using the WAN interface                                           |
| `starlink_wifi_software_update_state`                             | Starlink WiFi router software update state                                                  |
| `starlink_wifi_software_update_download_progress`                 | Progress of the in-flight Starlink WiFi router software update download (0-1)               |
| `starlink_wifi_software_update_seconds_since_get_target_versions` | Seconds since the WiFi router last fetched its target software versions                     |
| `starlink_wifi_software_update_info`                              | Starlink WiFi router software update version information                                    |
| `starlink_wifi_poe_state`                                         | Starlink WiFi router PoE state                                                              |
| `starlink_wifi_poe_power_watts`                                   | Starlink WiFi router PoE power draw in watts                                                |
| `starlink_wifi_poe_faults_fast_overcurrent`                       | Number of fast overcurrent PoE faults reported by the WiFi router                           |
| `starlink_wifi_poe_faults_slow_overcurrent`                       | Number of slow overcurrent PoE faults reported by the WiFi router                           |
| `starlink_wifi_poe_faults_overvoltage`                            | Number of overvoltage PoE faults reported by the WiFi router                                |
| `starlink_wifi_poe_faults_undervoltage`                           | Number of undervoltage PoE faults reported by the WiFi router                               |
| `starlink_wifi_poe_vsns_vin_volts`                                | Starlink WiFi router PoE input voltage in volts                                             |
| `starlink_wifi_setup_requirement_state`                           | Starlink WiFi router setup requirement state                                                |
| `starlink_wifi_setup_requirement_pause_countdown_seconds`         | Seconds remaining before the WiFi router setup requirement pause expires                    |

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
#   -router string
#         WiFi router address (set empty to disable WiFi metrics) (default "192.168.1.1:9000")
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

- Go v1.25 or newer (https://go.dev/dl/)

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
