// Copyright (c) 2024-2025 Joshua Sing <joshua@joshuasing.dev>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package exporter exports Prometheus metrics from a Starlink dishy.
package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace         = "starlink"
	dishSubsystem     = "dish"
	wifiSubsystem     = "wifi"
	exporterSubsystem = "exporter"
)

var (
	// Exporter
	exporterScrapesTotal = &Desc{
		Namespace: namespace,
		Subsystem: exporterSubsystem,
		Name:      "scrapes_total",
		Help:      "Total number of Starlink dish scrapes",
	}
	exporterScrapeDurationSeconds = &Desc{
		Namespace: namespace,
		Subsystem: exporterSubsystem,
		Name:      "scrape_duration_seconds",
		Help:      "Time taken to scrape metrics from the Starlink dish",
	}

	// Informational
	dishUp = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "up",
		Help:      "Whether scraping metrics from the Starlink dish was successful",
	}
	dishInfo = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "info",
		Help:      "Starlink dish software information",
		Labels: []string{
			"device_id",
			"hardware_version",
			"board_rev",
			"software_version",
			"manufactured_version",
			"generation_number",
			"country_code",
			"utc_offset",
			"boot_count",
			"mobility_class",
		},
	}
	dishUptimeSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "uptime_seconds",
		Help:      "Starlink dish uptime in seconds",
	}
	dishMobilityClass = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "mobility_class",
		Help:      "Starlink dish mobility class",
	}

	// Signal-to-noise ratio
	dishSnrAboveNoiseFloor = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "snr_above_noise_floor",
		Help:      "Whether Starlink dish signal-to-noise ratio is above noise floor",
	}
	dishSnrPersistentlyLow = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "snr_persistently_low",
		Help:      "Whether Starlink dish signal-to-noise ratio is persistently low",
	}

	// Throughput
	dishUplinkThroughputBps = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "uplink_throughput_bps",
		Help:      "Starlink dish uplink throughput in bits/sec",
	}
	dishDownlinkThroughputBps = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "downlink_throughput_bps",
		Help:      "Starlink dish downlink throughput in bit/sec",
	}
	dishDownlinkThroughputHistogram = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "downlink_throughput_bps_histogram",
		Help:      "Histogram of Starlink dish downlink throughput over last 15 minutes",
	}
	dishUplinkThroughputHistogram = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "uplink_throughput_bps_histogram",
		Help:      "Histogram of Starlink dish uplink throughput in bits/sec over last 15 minutes",
	}

	// PoP ping
	dishPopPingDropRatio = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "pop_ping_drop_ratio",
		Help:      "Starlink PoP ping drop ratio",
	}
	dishPopPingLatencySeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "pop_ping_latency_seconds",
		Help:      "Starlink PoP ping latency in seconds",
	}
	dishPopPingLatencyHistogram = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "pop_ping_latency_seconds_histogram",
		Help:      "Histogram of Starlink dish PoP ping latency in seconds over last 15 minutes",
	}

	// Power In
	dishPowerInputHistogram = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "power_input_watts_histogram",
		Help:      "Histogram of Starlink dish power input in watts over last 15 minutes",
	}
	dishPowerInput = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "power_input_watts",
		Help:      "Current power input for the Starlink dish",
	}

	// Software update
	dishSoftwareUpdateState = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "software_update_state",
		Help:      "Starlink dish update state",
	}
	dishSoftwareUpdateRebootReady = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "software_update_reboot_ready",
		Help:      "Whether the Starlink dish is ready to reboot to apply a software update",
	}

	// GPS
	dishGPSValid = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "gps_valid",
		Help:      "Whether the Starlink dish GPS is valid",
	}
	dishGPSSatellites = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "gps_satellites",
		Help:      "Number of GPS satellites visible to the Starlink dish",
	}

	// Tilt Angle
	dishTiltAngleDeg = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "tilt_angle_deg",
		Help:      "Starlink dish tilt angle degrees",
	}

	// Boresight
	dishBoresightAzimuthDeg = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "boresight_azimuth_deg",
		Help:      "Starlink dish boresight azimuth degrees",
	}
	dishBoresightElevationDeg = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "boresight_elevation_deg",
		Help:      "Starlink dish boresight elevation degrees",
	}
	dishDesiredBoresightAzimuthDeg = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "desired_boresight_azimuth_deg",
		Help:      "Starlink dish desired boresight azimuth degrees",
	}
	dishDesiredBoresightElevationDeg = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "desired_boresight_elevation_deg",
		Help:      "Starlink dish desired boresight elevation degrees",
	}

	// Obstruction
	dishCurrentlyObstructed = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "currently_obstructed",
		Help:      "Whether the Starlink dish is currently obstructed",
	}
	dishFractionObstructionRatio = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "fraction_obstruction_ratio",
		Help:      "Fraction of Starlink dish that is obstructed",
	}
	dishLast24HoursObstructedSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "last_24h_obstructed_seconds",
		Help:      "Number of seconds the Starlink dish was obstructed in the past 24 hours",
	}

	// Location
	dishLocationInfo = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "location_info",
		Help:      "Dish location information",
		Labels:    []string{"location_source", "lat", "lon", "alt"},
	}
	dishLocationLatitude = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "location_latitude_deg",
		Help:      "Location latitude in degrees",
	}
	dishLocationLongitude = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "location_longitude_deg",
		Help:      "Location longitude in degrees",
	}
	dishLocationAltitude = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "location_altitude_meters",
		Help:      "Location altitude in meters above sea level",
	}

	// Alerts
	dishAlertUnexpectedLocation = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_unexpected_location",
		Help:      "Whether the Starlink dish is in an unexpected location",
	}
	dishAlertInstallPending = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_install_pending",
		Help:      "Whether a Starlink Dish software update is pending installation",
	}
	dishAlertIsHeating = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_is_heating",
		Help:      "Whether the Starlink dish is heating (snow melting)",
	}
	dishAlertIsPowerSaveIdle = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_is_power_save_idle",
		Help:      "Whether the Starlink dish is currently in power saving mode",
	}
	dishAlertSignalLowerThanPredicted = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_signal_lower_than_predicted",
		Help:      "Whether the Starlink dish signal is lower than predicted",
	}
	dishAlertMotorsStuck = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_motors_stuck",
		Help:      "Whether the Starlink dish motors are stuck",
	}
	dishAlertThermalThrottle = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_thermal_throttle",
		Help:      "Whether the Starlink dish is thermally throttled",
	}
	dishAlertThermalShutdown = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_thermal_shutdown",
		Help:      "Whether the Starlink dish has thermally shut down",
	}
	dishAlertMastNotNearVertical = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_mast_not_near_vertical",
		Help:      "Whether the Starlink dish mast is not near vertical",
	}
	dishAlertSlowEthernetSpeeds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_slow_ethernet_speeds",
		Help:      "Whether the Starlink dish ethernet link is negotiated below gigabit",
	}
	dishAlertSlowEthernetSpeeds100 = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_slow_ethernet_speeds_100",
		Help:      "Whether the Starlink dish ethernet link is negotiated at 100 Mbps or lower",
	}
	dishAlertRoaming = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_roaming",
		Help:      "Whether the Starlink dish is roaming",
	}
	dishAlertPowerSupplyThermalThrottle = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_power_supply_thermal_throttle",
		Help:      "Whether the Starlink dish power supply is thermally throttled",
	}
	dishAlertDbfTelemStale = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_dbf_telem_stale",
		Help:      "Whether the Starlink dish DBF telemetry is stale",
	}
	dishAlertLowMotorCurrent = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_low_motor_current",
		Help:      "Whether the Starlink dish has low motor current",
	}
	dishAlertObstructionMapReset = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_obstruction_map_reset",
		Help:      "Whether the Starlink dish obstruction map was reset",
	}
	dishAlertDishWaterDetected = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_dish_water_detected",
		Help:      "Whether water has been detected inside the Starlink dish",
	}
	dishAlertRouterWaterDetected = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_router_water_detected",
		Help:      "Whether water has been detected inside the Starlink router",
	}
	dishAlertUpsuRouterPortSlow = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_upsu_router_port_slow",
		Help:      "Whether the Starlink dish UPSU router port is negotiated below the expected speed",
	}
	dishAlertNoEthernetLink = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "alert_no_ethernet_link",
		Help:      "Whether the Starlink dish has no ethernet link",
	}

	// Status (top-level scalars)
	dishSecondsToFirstNonemptySlot = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "seconds_to_first_nonempty_slot",
		Help:      "Seconds until the next non-empty network schedule slot",
	}
	dishStowRequested = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "stow_requested",
		Help:      "Whether a stow has been requested for the Starlink dish",
	}
	dishEthSpeedMbps = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "eth_speed_mbps",
		Help:      "Negotiated speed of the Starlink dish ethernet link in Mbps",
	}
	dishClassOfService = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "class_of_service",
		Help:      "Starlink dish class of service",
	}
	dishRebootReason = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "reboot_reason",
		Help:      "Reason for the most recent Starlink dish reboot",
	}
	dishDisablementCode = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "disablement_code",
		Help:      "Reason the Starlink dish is disabled, if any",
	}
	dishDlBandwidthRestrictedReason = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "dl_bandwidth_restricted_reason",
		Help:      "Reason downlink bandwidth is currently restricted, if any",
	}
	dishUlBandwidthRestrictedReason = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ul_bandwidth_restricted_reason",
		Help:      "Reason uplink bandwidth is currently restricted, if any",
	}
	dishIsCellDisabled = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "is_cell_disabled",
		Help:      "Whether the Starlink dish cell is disabled",
	}
	dishSecondsUntilSwupdateRebootPossible = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "seconds_until_swupdate_reboot_possible",
		Help:      "Seconds until the Starlink dish can reboot to apply a software update",
	}
	dishHighPowerTestMode = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "high_power_test_mode",
		Help:      "Whether the Starlink dish is in high-power test mode",
	}
	dishIsMovingFastPersisted = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "is_moving_fast_persisted",
		Help:      "Whether the Starlink dish has persistently detected fast movement",
	}
	dishMacFlag = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "mac_flag",
		Help:      "Starlink dish MAC flag",
	}
	dishNatFlag = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "nat_flag",
		Help:      "Starlink dish NAT flag",
	}
	dishAccountShard = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "account_shard",
		Help:      "Starlink dish account shard",
	}
	dishConnectedRoutersCount = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "connected_routers_count",
		Help:      "Number of routers currently connected to the Starlink dish",
	}
	dishHasActuators = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "has_actuators",
		Help:      "Whether the Starlink dish has actuators",
	}
	dishNed2DishQuaternionW = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ned2dish_quaternion_w",
		Help:      "W component of the NED-to-dish orientation quaternion",
	}
	dishNed2DishQuaternionX = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ned2dish_quaternion_x",
		Help:      "X component of the NED-to-dish orientation quaternion",
	}
	dishNed2DishQuaternionY = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ned2dish_quaternion_y",
		Help:      "Y component of the NED-to-dish orientation quaternion",
	}
	dishNed2DishQuaternionZ = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ned2dish_quaternion_z",
		Help:      "Z component of the NED-to-dish orientation quaternion",
	}

	// Current outage
	dishOutageInfo = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "outage_info",
		Help:      "Information about the current Starlink dish outage, if any",
		Labels:    []string{"cause", "did_switch"},
	}
	dishOutageStartTimestampSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "outage_start_timestamp_seconds",
		Help:      "Unix timestamp at which the current outage started",
	}
	dishOutageDurationSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "outage_duration_seconds",
		Help:      "Duration of the current outage in seconds",
	}

	// GPS
	dishGPSNoSatsAfterTtff = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "gps_no_sats_after_ttff",
		Help:      "Whether the Starlink dish has no GPS satellites after time-to-first-fix",
	}
	dishGPSInhibit = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "gps_inhibit",
		Help:      "Whether GPS is currently inhibited on the Starlink dish",
	}

	// Obstruction (additional)
	dishObstructionValidSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "obstruction_valid_seconds",
		Help:      "Seconds the Starlink dish obstruction map has been collecting data",
	}
	dishObstructionPatchesValid = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "obstruction_patches_valid",
		Help:      "Number of valid patches in the Starlink dish obstruction map",
	}
	dishAvgProlongedObstructionDurationSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "avg_prolonged_obstruction_duration_seconds",
		Help:      "Average duration of a prolonged obstruction in seconds",
	}
	dishAvgProlongedObstructionIntervalSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "avg_prolonged_obstruction_interval_seconds",
		Help:      "Average interval between prolonged obstructions in seconds",
	}
	dishAvgProlongedObstructionValid = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "avg_prolonged_obstruction_valid",
		Help:      "Whether the average prolonged obstruction metrics are currently valid",
	}

	// Ready states
	dishReadyStateCady = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ready_state_cady",
		Help:      "Whether the cady subsystem is ready",
	}
	dishReadyStateScp = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ready_state_scp",
		Help:      "Whether the SCP subsystem is ready",
	}
	dishReadyStateL1L2 = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ready_state_l1l2",
		Help:      "Whether the L1/L2 subsystem is ready",
	}
	dishReadyStateXphy = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ready_state_xphy",
		Help:      "Whether the XPHY subsystem is ready",
	}
	dishReadyStateAap = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ready_state_aap",
		Help:      "Whether the AAP subsystem is ready",
	}
	dishReadyStateRf = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "ready_state_rf",
		Help:      "Whether the RF subsystem is ready",
	}

	// Software update (additional from SoftwareUpdateStats)
	dishSoftwareUpdateProgress = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "software_update_progress",
		Help:      "Progress of the in-flight Starlink dish software update (0-1)",
	}
	dishSoftwareUpdateRequiresReboot = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "software_update_requires_reboot",
		Help:      "Whether the pending Starlink dish software update requires a reboot",
	}
	dishSoftwareUpdateRebootScheduledUtcSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "software_update_reboot_scheduled_utc_seconds",
		Help:      "Unix timestamp at which the Starlink dish is scheduled to reboot to apply a software update",
	}

	// Alignment (additional)
	dishActuatorState = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "actuator_state",
		Help:      "Current state of the Starlink dish actuators",
	}
	dishAttitudeEstimationState = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "attitude_estimation_state",
		Help:      "Current state of the Starlink dish attitude estimator",
	}
	dishAttitudeUncertaintyDeg = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "attitude_uncertainty_deg",
		Help:      "Starlink dish attitude uncertainty in degrees",
	}

	// Initialization durations
	dishInitAttitudeSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "initialization_attitude_seconds",
		Help:      "Seconds spent on attitude initialization during dish boot",
	}
	dishInitBurstDetectedSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "initialization_burst_detected_seconds",
		Help:      "Seconds until first burst was detected during dish boot",
	}
	dishInitEkfConvergedSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "initialization_ekf_converged_seconds",
		Help:      "Seconds until the EKF converged during dish boot",
	}
	dishInitFirstCplaneSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "initialization_first_cplane_seconds",
		Help:      "Seconds until first cplane message during dish boot",
	}
	dishInitFirstPopPingSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "initialization_first_pop_ping_seconds",
		Help:      "Seconds until first PoP ping during dish boot",
	}
	dishInitGpsValidSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "initialization_gps_valid_seconds",
		Help:      "Seconds until GPS became valid during dish boot",
	}
	dishInitInitialNetworkEntrySeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "initialization_initial_network_entry_seconds",
		Help:      "Seconds until initial network entry during dish boot",
	}
	dishInitNetworkScheduleSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "initialization_network_schedule_seconds",
		Help:      "Seconds until first network schedule was received during dish boot",
	}
	dishInitRfReadySeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "initialization_rf_ready_seconds",
		Help:      "Seconds until the RF subsystem was ready during dish boot",
	}
	dishInitStableConnectionSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "initialization_stable_connection_seconds",
		Help:      "Seconds until a stable connection was established during dish boot",
	}

	// PLC (Mini battery / Power-Line Communication)
	dishPlcReceiving = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "plc_receiving",
		Help:      "Whether the dish is receiving Power-Line Communication data",
	}
	dishPlcAverageTimeToEmptySeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "plc_average_time_to_empty_seconds",
		Help:      "Average time until the battery is empty in seconds",
	}
	dishPlcAverageTimeToFullSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "plc_average_time_to_full_seconds",
		Help:      "Average time until the battery is full in seconds",
	}
	dishPlcBatteryHealth = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "plc_battery_health",
		Help:      "Battery health as reported by PLC",
	}
	dishPlcPermanentFailure = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "plc_permanent_failure",
		Help:      "Whether the PLC battery has reported a permanent failure",
	}
	dishPlcSafetyModeActive = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "plc_safety_mode_active",
		Help:      "Whether the PLC battery is in safety mode",
	}
	dishPlcStateOfChargePercent = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "plc_state_of_charge_percent",
		Help:      "PLC battery state of charge in percent",
	}
	dishPlcThermalThrottleLevel = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "plc_thermal_throttle_level",
		Help:      "PLC thermal throttle level",
	}
	dishPlcRevision = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "plc_revision",
		Help:      "PLC protocol revision",
	}

	// UPSU (Universal Power Supply Unit, Gen 3)
	dishUpsuDishPowerWatts = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "upsu_dish_power_watts",
		Help:      "Dish power draw as reported by the UPSU in watts",
	}
	dishUpsuRouterPowerWatts = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "upsu_router_power_watts",
		Help:      "Router power draw as reported by the UPSU in watts",
	}
	dishUpsuUptimeSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "upsu_uptime_seconds",
		Help:      "UPSU uptime in seconds",
	}
	dishUpsuBoardRev = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "upsu_board_rev",
		Help:      "UPSU board revision",
	}

	// APS (Auxiliary Power Supply)
	dishApsDishPowerWatts = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "aps_dish_power_watts",
		Help:      "Dish power draw as reported by the APS in watts",
	}
	dishApsUptimeSeconds = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "aps_uptime_seconds",
		Help:      "APS uptime in seconds",
	}
	dishApsBoardRev = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "aps_board_rev",
		Help:      "APS board revision",
	}

	// History (additional)
	dishPopPingDropRatioHistogram = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "pop_ping_drop_ratio_histogram",
		Help:      "Histogram of Starlink dish PoP ping drop ratio over the last 15 minutes",
	}
	dishHistoryOutagesCount = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "history_outages_count",
		Help:      "Number of outages retained in the Starlink dish history buffer",
	}

	// Location (additional)
	dishLocationUncertaintyMeters = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "location_uncertainty_meters",
		Help:      "1-sigma uncertainty of the reported location in meters",
	}
	dishHorizontalSpeedMps = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "horizontal_speed_mps",
		Help:      "Horizontal speed of the Starlink dish in meters per second",
	}
	dishVerticalSpeedMps = &Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "vertical_speed_mps",
		Help:      "Vertical speed of the Starlink dish in meters per second",
	}

	// WiFi router

	// Informational
	wifiUp = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "up",
		Help:      "Whether scraping metrics from the Starlink WiFi router was successful",
	}
	wifiInfo = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "info",
		Help:      "Starlink WiFi router software information",
		Labels: []string{
			"device_id",
			"hardware_version",
			"board_rev",
			"software_version",
			"manufactured_version",
			"generation_number",
			"country_code",
			"utc_offset",
			"boot_count",
		},
	}
	wifiUptimeSeconds = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "uptime_seconds",
		Help:      "Starlink WiFi router uptime in seconds",
	}
	wifiHopsFromController = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "hops_from_controller",
		Help:      "Number of mesh hops between this router and the controller",
	}
	wifiNoWanLink = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "no_wan_link",
		Help:      "Whether the Starlink WiFi router has no WAN link",
	}
	wifiIsAviation = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "is_aviation",
		Help:      "Whether the Starlink WiFi router is in aviation mode",
	}
	wifiIsAviationConformed = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "is_aviation_conformed",
		Help:      "Whether the Starlink WiFi router is aviation conformed",
	}
	wifiUsingIndividualizedCalibration = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "using_individualized_calibration",
		Help:      "Whether the Starlink WiFi router is using individualized calibration",
	}
	wifiCalibrationPartitionsState = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "calibration_partitions_state",
		Help:      "Starlink WiFi router calibration partitions state",
	}
	wifiDishDisablementCode = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "dish_disablement_code",
		Help:      "Disablement code reported by the dish as seen by the WiFi router",
	}
	wifiSecsSinceLastPublicIpv4Change = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "seconds_since_last_public_ipv4_change",
		Help:      "Seconds since the public IPv4 address last changed",
	}
	wifiClientsCount = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "clients_count",
		Help:      "Number of clients currently connected to the Starlink WiFi router",
	}
	wifiDhcpServersCount = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "dhcp_servers_count",
		Help:      "Number of DHCP servers reported by the Starlink WiFi router",
	}

	// Ping
	wifiPingDropRatio = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "ping_drop_ratio",
		Help:      "Router WAN ping drop ratio",
	}
	wifiPingDropRatio5m = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "ping_drop_ratio_5m",
		Help:      "Router WAN ping drop ratio over the last 5 minutes",
	}
	wifiPingLatencySeconds = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "ping_latency_seconds",
		Help:      "Router WAN ping latency in seconds",
	}
	wifiDishPingDropRatio = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "dish_ping_drop_ratio",
		Help:      "Router-to-dish ping drop ratio",
	}
	wifiDishPingDropRatio5m = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "dish_ping_drop_ratio_5m",
		Help:      "Router-to-dish ping drop ratio over the last 5 minutes",
	}
	wifiDishPingLatencySeconds = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "dish_ping_latency_seconds",
		Help:      "Router-to-dish ping latency in seconds",
	}
	wifiPopPingDropRatio = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "pop_ping_drop_ratio",
		Help:      "Router-to-PoP ping drop ratio",
	}
	wifiPopPingDropRatio5m = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "pop_ping_drop_ratio_5m",
		Help:      "Router-to-PoP ping drop ratio over the last 5 minutes",
	}
	wifiPopPingLatencySeconds = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "pop_ping_latency_seconds",
		Help:      "Router-to-PoP ping latency in seconds",
	}
	wifiPopIpv6PingDropRatio = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "pop_ipv6_ping_drop_ratio",
		Help:      "Router-to-PoP IPv6 ping drop ratio",
	}
	wifiPopIpv6PingDropRatio5m = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "pop_ipv6_ping_drop_ratio_5m",
		Help:      "Router-to-PoP IPv6 ping drop ratio over the last 5 minutes",
	}
	wifiPopIpv6PingLatencySeconds = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "pop_ipv6_ping_latency_seconds",
		Help:      "Router-to-PoP IPv6 ping latency in seconds",
	}

	// Alerts
	wifiAlertThermalThrottle = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_thermal_throttle",
		Help:      "Whether the Starlink WiFi router is thermally throttled",
	}
	wifiAlertInstallPending = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_install_pending",
		Help:      "Whether a Starlink WiFi router software update is pending installation",
	}
	wifiAlertFreshlyFused = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_freshly_fused",
		Help:      "Whether the Starlink WiFi router was recently fused",
	}
	wifiAlertLanEthSlowLink10 = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_lan_eth_slow_link_10",
		Help:      "Whether a Starlink WiFi router LAN ethernet link is negotiated at 10 Mbps",
	}
	wifiAlertLanEthSlowLink100 = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_lan_eth_slow_link_100",
		Help:      "Whether a Starlink WiFi router LAN ethernet link is negotiated at 100 Mbps",
	}
	wifiAlertHighCablePingDropRate = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_high_cable_ping_drop_rate",
		Help:      "Whether the Starlink WiFi router is reporting a high cable ping drop rate",
	}
	wifiAlertWanEthPoorConnection = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_wan_eth_poor_connection",
		Help:      "Whether the Starlink WiFi router WAN ethernet connection is poor",
	}
	wifiAlertMeshTopologyChangingOften = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_mesh_topology_changing_often",
		Help:      "Whether the Starlink WiFi mesh topology is changing often",
	}
	wifiAlertMeshUnreliableBackhaul = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_mesh_unreliable_backhaul",
		Help:      "Whether the Starlink WiFi mesh backhaul is unreliable",
	}
	wifiAlertRadiusMissingProcess = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_radius_missing_process",
		Help:      "Whether the Starlink WiFi router RADIUS process is missing",
	}
	wifiAlertEthSwitchError = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_eth_switch_error",
		Help:      "Whether the Starlink WiFi router ethernet switch has reported an error",
	}
	wifiAlertPoeOnDishUnreachable = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_poe_on_dish_unreachable",
		Help:      "Whether the dish is unreachable while PoE is on",
	}
	wifiAlertPoeFuseBlown = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_poe_fuse_blown",
		Help:      "Whether the Starlink WiFi router PoE fuse is blown",
	}
	wifiAlertPoeRouterOvercurrent = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_poe_router_overcurrent",
		Help:      "Whether the Starlink WiFi router PoE has detected an overcurrent",
	}
	wifiAlertPoeOffCurrentNominal = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_poe_off_current_nominal",
		Help:      "Whether the Starlink WiFi router is drawing nominal current while PoE is off",
	}
	wifiAlertPoeVinOvervoltage = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_poe_vin_overvoltage",
		Help:      "Whether the Starlink WiFi router PoE input voltage is over the threshold",
	}
	wifiAlertPoeVinUndervoltage = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_poe_vin_undervoltage",
		Help:      "Whether the Starlink WiFi router PoE input voltage is under the threshold",
	}
	wifiAlertSandboxDisabled = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_sandbox_disabled",
		Help:      "Whether the Starlink WiFi router client sandbox is disabled",
	}
	wifiAlertOnlyOverflightBlocked = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_only_overflight_blocked",
		Help:      "Whether only overflight is blocked on the Starlink WiFi router",
	}
	wifiAlertOfflineNetworksDisabled = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_offline_networks_disabled",
		Help:      "Whether offline networks are disabled on the Starlink WiFi router",
	}
	wifiAlertWiredMeshNotUsingWanIface = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "alert_wired_mesh_not_using_wan_iface",
		Help:      "Whether wired mesh is not using the WAN interface",
	}

	// Software update
	wifiSoftwareUpdateState = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "software_update_state",
		Help:      "Starlink WiFi router software update state",
	}
	wifiSoftwareUpdateDownloadProgress = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "software_update_download_progress",
		Help:      "Progress of the in-flight Starlink WiFi router software update download (0-1)",
	}
	wifiSoftwareUpdateSecondsSinceGetTargetVersions = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "software_update_seconds_since_get_target_versions",
		Help:      "Seconds since the WiFi router last fetched its target software versions",
	}
	wifiSoftwareUpdateInfo = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "software_update_info",
		Help:      "Starlink WiFi router software update version information",
		Labels:    []string{"running_version", "version_in_progress"},
	}

	// PoE
	wifiPoeState = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "poe_state",
		Help:      "Starlink WiFi router PoE state",
	}
	wifiPoePowerWatts = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "poe_power_watts",
		Help:      "Starlink WiFi router PoE power draw in watts",
	}
	wifiPoeFaultsFastOvercurrent = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "poe_faults_fast_overcurrent",
		Help:      "Number of fast overcurrent PoE faults reported by the WiFi router",
	}
	wifiPoeFaultsSlowOvercurrent = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "poe_faults_slow_overcurrent",
		Help:      "Number of slow overcurrent PoE faults reported by the WiFi router",
	}
	wifiPoeFaultsOvervoltage = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "poe_faults_overvoltage",
		Help:      "Number of overvoltage PoE faults reported by the WiFi router",
	}
	wifiPoeFaultsUndervoltage = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "poe_faults_undervoltage",
		Help:      "Number of undervoltage PoE faults reported by the WiFi router",
	}
	wifiPoeVsnsVinVolts = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "poe_vsns_vin_volts",
		Help:      "Starlink WiFi router PoE input voltage in volts",
	}

	// Setup requirement
	wifiSetupRequirementState = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "setup_requirement_state",
		Help:      "Starlink WiFi router setup requirement state",
	}
	wifiSetupRequirementPauseCountdownSeconds = &Desc{
		Namespace: namespace,
		Subsystem: wifiSubsystem,
		Name:      "setup_requirement_pause_countdown_seconds",
		Help:      "Seconds remaining before the WiFi router setup requirement pause expires",
	}
)

// Descs contains all Prometheus metrics descriptors for the exporter.
var Descs = []*Desc{
	exporterScrapesTotal,
	exporterScrapeDurationSeconds,
	dishUp,
	dishInfo,
	dishUptimeSeconds,
	dishMobilityClass,
	dishSnrAboveNoiseFloor,
	dishSnrPersistentlyLow,
	dishUplinkThroughputBps,
	dishDownlinkThroughputBps,
	dishDownlinkThroughputHistogram,
	dishUplinkThroughputHistogram,
	dishPopPingDropRatio,
	dishPopPingLatencySeconds,
	dishPopPingLatencyHistogram,
	dishSoftwareUpdateState,
	dishSoftwareUpdateRebootReady,
	dishGPSValid,
	dishGPSSatellites,
	dishTiltAngleDeg,
	dishBoresightAzimuthDeg,
	dishBoresightElevationDeg,
	dishDesiredBoresightAzimuthDeg,
	dishDesiredBoresightElevationDeg,
	dishCurrentlyObstructed,
	dishFractionObstructionRatio,
	dishLast24HoursObstructedSeconds,
	dishPowerInputHistogram,
	dishPowerInput,
	dishLocationInfo,
	dishLocationLatitude,
	dishLocationLongitude,
	dishLocationAltitude,
	dishAlertUnexpectedLocation,
	dishAlertInstallPending,
	dishAlertIsHeating,
	dishAlertIsPowerSaveIdle,
	dishAlertSignalLowerThanPredicted,
	dishAlertMotorsStuck,
	dishAlertThermalThrottle,
	dishAlertThermalShutdown,
	dishAlertMastNotNearVertical,
	dishAlertSlowEthernetSpeeds,
	dishAlertSlowEthernetSpeeds100,
	dishAlertRoaming,
	dishAlertPowerSupplyThermalThrottle,
	dishAlertDbfTelemStale,
	dishAlertLowMotorCurrent,
	dishAlertObstructionMapReset,
	dishAlertDishWaterDetected,
	dishAlertRouterWaterDetected,
	dishAlertUpsuRouterPortSlow,
	dishAlertNoEthernetLink,
	dishSecondsToFirstNonemptySlot,
	dishStowRequested,
	dishEthSpeedMbps,
	dishClassOfService,
	dishRebootReason,
	dishDisablementCode,
	dishDlBandwidthRestrictedReason,
	dishUlBandwidthRestrictedReason,
	dishIsCellDisabled,
	dishSecondsUntilSwupdateRebootPossible,
	dishHighPowerTestMode,
	dishIsMovingFastPersisted,
	dishMacFlag,
	dishNatFlag,
	dishAccountShard,
	dishConnectedRoutersCount,
	dishHasActuators,
	dishNed2DishQuaternionW,
	dishNed2DishQuaternionX,
	dishNed2DishQuaternionY,
	dishNed2DishQuaternionZ,
	dishOutageInfo,
	dishOutageStartTimestampSeconds,
	dishOutageDurationSeconds,
	dishGPSNoSatsAfterTtff,
	dishGPSInhibit,
	dishObstructionValidSeconds,
	dishObstructionPatchesValid,
	dishAvgProlongedObstructionDurationSeconds,
	dishAvgProlongedObstructionIntervalSeconds,
	dishAvgProlongedObstructionValid,
	dishReadyStateCady,
	dishReadyStateScp,
	dishReadyStateL1L2,
	dishReadyStateXphy,
	dishReadyStateAap,
	dishReadyStateRf,
	dishSoftwareUpdateProgress,
	dishSoftwareUpdateRequiresReboot,
	dishSoftwareUpdateRebootScheduledUtcSeconds,
	dishActuatorState,
	dishAttitudeEstimationState,
	dishAttitudeUncertaintyDeg,
	dishInitAttitudeSeconds,
	dishInitBurstDetectedSeconds,
	dishInitEkfConvergedSeconds,
	dishInitFirstCplaneSeconds,
	dishInitFirstPopPingSeconds,
	dishInitGpsValidSeconds,
	dishInitInitialNetworkEntrySeconds,
	dishInitNetworkScheduleSeconds,
	dishInitRfReadySeconds,
	dishInitStableConnectionSeconds,
	dishPlcReceiving,
	dishPlcAverageTimeToEmptySeconds,
	dishPlcAverageTimeToFullSeconds,
	dishPlcBatteryHealth,
	dishPlcPermanentFailure,
	dishPlcSafetyModeActive,
	dishPlcStateOfChargePercent,
	dishPlcThermalThrottleLevel,
	dishPlcRevision,
	dishUpsuDishPowerWatts,
	dishUpsuRouterPowerWatts,
	dishUpsuUptimeSeconds,
	dishUpsuBoardRev,
	dishApsDishPowerWatts,
	dishApsUptimeSeconds,
	dishApsBoardRev,
	dishPopPingDropRatioHistogram,
	dishHistoryOutagesCount,
	dishLocationUncertaintyMeters,
	dishHorizontalSpeedMps,
	dishVerticalSpeedMps,
	wifiUp,
	wifiInfo,
	wifiUptimeSeconds,
	wifiHopsFromController,
	wifiNoWanLink,
	wifiIsAviation,
	wifiIsAviationConformed,
	wifiUsingIndividualizedCalibration,
	wifiCalibrationPartitionsState,
	wifiDishDisablementCode,
	wifiSecsSinceLastPublicIpv4Change,
	wifiClientsCount,
	wifiDhcpServersCount,
	wifiPingDropRatio,
	wifiPingDropRatio5m,
	wifiPingLatencySeconds,
	wifiDishPingDropRatio,
	wifiDishPingDropRatio5m,
	wifiDishPingLatencySeconds,
	wifiPopPingDropRatio,
	wifiPopPingDropRatio5m,
	wifiPopPingLatencySeconds,
	wifiPopIpv6PingDropRatio,
	wifiPopIpv6PingDropRatio5m,
	wifiPopIpv6PingLatencySeconds,
	wifiAlertThermalThrottle,
	wifiAlertInstallPending,
	wifiAlertFreshlyFused,
	wifiAlertLanEthSlowLink10,
	wifiAlertLanEthSlowLink100,
	wifiAlertHighCablePingDropRate,
	wifiAlertWanEthPoorConnection,
	wifiAlertMeshTopologyChangingOften,
	wifiAlertMeshUnreliableBackhaul,
	wifiAlertRadiusMissingProcess,
	wifiAlertEthSwitchError,
	wifiAlertPoeOnDishUnreachable,
	wifiAlertPoeFuseBlown,
	wifiAlertPoeRouterOvercurrent,
	wifiAlertPoeOffCurrentNominal,
	wifiAlertPoeVinOvervoltage,
	wifiAlertPoeVinUndervoltage,
	wifiAlertSandboxDisabled,
	wifiAlertOnlyOverflightBlocked,
	wifiAlertOfflineNetworksDisabled,
	wifiAlertWiredMeshNotUsingWanIface,
	wifiSoftwareUpdateState,
	wifiSoftwareUpdateDownloadProgress,
	wifiSoftwareUpdateSecondsSinceGetTargetVersions,
	wifiSoftwareUpdateInfo,
	wifiPoeState,
	wifiPoePowerWatts,
	wifiPoeFaultsFastOvercurrent,
	wifiPoeFaultsSlowOvercurrent,
	wifiPoeFaultsOvervoltage,
	wifiPoeFaultsUndervoltage,
	wifiPoeVsnsVinVolts,
	wifiSetupRequirementState,
	wifiSetupRequirementPauseCountdownSeconds,
}

// Desc is a utility wrapper for prometheus.Desc.
type Desc struct {
	Namespace   string
	Subsystem   string
	Name        string
	Help        string
	Labels      []string
	ConstLabels prometheus.Labels

	fqName string
	desc   *prometheus.Desc
}

// FQName builds the fully-qualified metric name from the name parts.
func (d *Desc) FQName() string {
	if d.fqName == "" {
		d.fqName = prometheus.BuildFQName(d.Namespace, d.Subsystem, d.Name)
	}
	return d.fqName
}

// Desc builds and returns a *prometheus.Desc.
func (d *Desc) Desc() *prometheus.Desc {
	if d.desc == nil {
		d.desc = prometheus.NewDesc(d.FQName(), d.Help, d.Labels, d.ConstLabels)
	}
	return d.desc
}
