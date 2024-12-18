// Copyright (c) 2024 Joshua Sing <joshua@joshuasing.dev>
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

package exporter

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace         = "starlink"
	dishSubsystem     = "dish"
	exporterSubsystem = "exporter"
)

var (
	// Exporter
	exporterScrapesTotal = Desc{
		Namespace: namespace,
		Subsystem: exporterSubsystem,
		Name:      "scrapes_total",
		Help:      "Total number of Starlink dish scrapes",
	}
	exporterScrapeDurationSeconds = Desc{
		Namespace: namespace,
		Subsystem: exporterSubsystem,
		Name:      "scrape_duration_seconds",
		Help:      "Time taken to scrape metrics from the Starlink dish",
	}

	// Informational
	dishUp = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "up",
		Help:      "Whether scraping metrics from the Starlink dish was successful",
	}
	dishInfo = Desc{
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
		},
	}
	dishUptimeSeconds = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "uptime_seconds",
		Help:      "Starlink dish uptime in seconds",
	}

	// Signal-to-noise ratio
	dishSnrAboveNoiseFloor = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "snr_above_noise_floor",
		Help:      "Whether Starlink dish signal-to-noise ratio is above noise floor",
	}
	dishSnrPersistentlyLow = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "snr_persistently_low",
		Help:      "Whether Starlink dish signal-to-noise ratio is persistently low",
	}

	// Throughput
	dishUplinkThroughputBps = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "uplink_throughput_bps",
		Help:      "Starlink dish uplink throughput in bits/sec",
	}
	dishDownlinkThroughputBps = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "downlink_throughput_bps",
		Help:      "Starlink dish downlink throughput in bit/sec",
	}
	dishDownlinkThroughputHistogram = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "downlink_throughput_bps_histogram",
		Help:      "Histogram of Starlink dish downlink throughput over last 15 minutes",
	}
	dishUplinkThroughputHistogram = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "uplink_throughput_bps_histogram",
		Help:      "Histogram of Starlink dish uplink throughput in bits/sec over last 15 minutes",
	}

	// PoP ping
	dishPopPingDropRatio = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "pop_ping_drop_ratio",
		Help:      "Starlink PoP ping drop ratio",
	}
	dishPopPingLatencySeconds = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "pop_ping_latency_seconds",
		Help:      "Starlink PoP ping latency in seconds",
	}
	dishPopPingLatencyHistogram = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "pop_ping_latency_seconds_histogram",
		Help:      "Histogram of Starlink dish PoP ping latency in seconds over last 15 minutes",
	}

	// Power In
	dishPowerInputHistogram = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "power_input_watts_histogram",
		Help:      "Histogram of Starlink dish power input in watts over last 15 minutes",
	}
	dishPowerInput = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "power_input_watts",
		Help:      "Current power input for the Starlink dish",
	}

	// Software update
	dishSoftwareUpdateState = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "software_update_state",
		Help:      "Starlink dish update state",
	}
	dishSoftwareUpdateRebootReady = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "software_update_reboot_ready",
		Help:      "Whether the Starlink dish is ready to reboot to apply a software update",
	}

	// Boresight
	dishBoresightAzimuthDeg = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "boresight_azimuth_deg",
		Help:      "Starlink dish boresight azimuth in degrees",
	}
	dishBoresightElevationDeg = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "boresight_elevation_deg",
		Help:      "Starlink dish boresight elevation in degrees",
	}

	// Obstruction
	dishCurrentlyObstructed = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "currently_obstructed",
		Help:      "Whether the Starlink dish is currently obstructed",
	}
	dishFractionObstructionRatio = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "fraction_obstruction_ratio",
		Help:      "Fraction of Starlink dish that is obstructed",
	}
	dishLast24HoursObstructedSeconds = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "last_24h_obstructed_seconds",
		Help:      "Number of seconds the Starlink dish was obstructed in the past 24 hours",
	}
	dishProlongedObstructionDurationSeconds = Desc{
		Namespace: namespace,
		Subsystem: dishSubsystem,
		Name:      "prolonged_obstruction_duration_seconds",
		Help:      "Average prolonged obstruction duration in seconds",
	}
)

var Descs = []Desc{
	exporterScrapesTotal,
	exporterScrapeDurationSeconds,
	dishUp,
	dishInfo,
	dishUptimeSeconds,
	dishSnrAboveNoiseFloor,
	dishSnrPersistentlyLow,
	dishUplinkThroughputBps,
	dishDownlinkThroughputBps,
	dishDownlinkThroughputHistogram,
	dishUplinkThroughputHistogram,
	dishPopPingDropRatio,
	dishPopPingLatencySeconds,
	dishPopPingLatencyHistogram,
	dishSoftwareUpdateRebootReady,
	dishBoresightAzimuthDeg,
	dishBoresightElevationDeg,
	dishCurrentlyObstructed,
	dishFractionObstructionRatio,
	dishLast24HoursObstructedSeconds,
	dishProlongedObstructionDurationSeconds,
	dishPowerInputHistogram,
	dishPowerInput,
}

type Desc struct {
	Namespace   string
	Subsystem   string
	Name        string
	Help        string
	Labels      []string
	ConstLabels prometheus.Labels
	desc        *prometheus.Desc
}

func (d Desc) FQName() string {
	return prometheus.BuildFQName(d.Namespace, d.Subsystem, d.Name)
}

func (d Desc) Desc() *prometheus.Desc {
	if d.desc == nil {
		d.desc = prometheus.NewDesc(d.FQName(), d.Help, d.Labels, d.ConstLabels)
	}
	return d.desc
}
