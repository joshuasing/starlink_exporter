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

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/joshuasing/starlink_exporter/internal/spacex/api/device"
)

const (
	namespace         = "starlink"
	dishSubsystem     = "dish"
	exporterSubsystem = "exporter"
)

const DefaultDishAddress = "192.168.100.1:9200"

var (
	// Informational
	dishInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "info"),
		"Starlink dish software information",
		[]string{
			"device_id",
			"hardware_version",
			"board_rev",
			"software_version",
			"manufactured_version",
			"generation_number",
			"country_code",
			"utc_offset",
			"boot_count",
		}, nil,
	)
	dishUptimeSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "uptime_seconds"),
		"Starlink dish uptime in seconds",
		nil, nil,
	)

	// Signal-to-noise ratio
	dishSnrAboveNoiseFloor = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "snr_above_noise_floor"),
		"Whether Starlink dish signal-to-noise ratio is above noise floor",
		nil, nil,
	)
	dishSnrPersistentlyLow = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "snr_persistently_low"),
		"Whether Starlink dish signal-to-noise ratio is persistently low",
		nil, nil,
	)

	// Throughput
	dishUplinkThroughputBps = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "uplink_throughput_bps"),
		"Starlink dish uplink throughput in bits/sec",
		nil, nil,
	)
	dishDownlinkThroughputBps = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "downlink_throughput_bps"),
		"Starlink dish downlink throughput in bit/sec",
		nil, nil,
	)
	dishDownlinkThroughputHistogram = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "downlink_throughput_histogram"),
		"Histogram of Starlink dish downlink throughput over last 15 minutes",
		nil, nil,
	)
	dishUplinkThroughputHistogram = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "uplink_throughput_histogram"),
		"Histogram of Starlink dish uplink throughput over last 15 minutes",
		nil, nil,
	)

	// PoP ping
	dishPopPingDropRatio = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "pop_ping_drop_ratio"),
		"Starlink PoP ping drop ratio",
		nil, nil,
	)
	dishPopPingLatencySeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "pop_ping_latency_seconds"),
		"Starlink PoP ping latency in seconds",
		nil, nil,
	)
	dishPopPingLatencyHistogram = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "pop_ping_latency_histogram"),
		"Histogram of Starlink dish pop ping latency over last 15 minutes",
		nil, nil,
	)

	// Power In
	dishPowerInputHistogram = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "power_input_histogram"),
		"Histogram of Starlink dish power input over last 15 minutes",
		nil, nil,
	)
	dishPowerInput = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "power_input"),
		"Current power input for the Starlink dish",
		nil, nil,
	)

	// Software update
	dishSoftwareUpdateState = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "software_update_state"),
		"Starlink dish update state",
		nil, nil,
	)
	dishSoftwareUpdateRebootReady = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "software_update_reboot_ready"),
		"Whether the Starlink dish is ready to reboot to apply a software update",
		nil, nil,
	)

	// Boresight
	dishBoresightAzimuthDeg = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "boresight_azimuth_deg"),
		"Starlink dish boresight azimuth in degrees",
		nil, nil,
	)
	dishBoresightElevationDeg = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "boresight_elevation_deg"),
		"Starlink dish boresight elevation in degrees",
		nil, nil,
	)

	// Obstruction
	dishCurrentlyObstructed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "currently_obstructed"),
		"Whether the Starlink dish is currently obstructed",
		nil, nil,
	)
	dishFractionObstructionRatio = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "fraction_obstruction_ratio"),
		"Fraction of Starlink dish that is obstructed",
		nil, nil,
	)
	dishLast24HoursObstructedSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "last_24h_obstructed_seconds"),
		"Number of seconds the Starlink dish was obstructed in the past 24 hours",
		nil, nil,
	)
	dishProlongedObstructionDurationSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, dishSubsystem, "prolonged_obstruction_duration_seconds"),
		"Average prolonged obstruction duration in seconds",
		nil, nil,
	)
)

// Exporter is a Starlink Dishy metrics exporter.
type Exporter struct {
	mx     sync.Mutex
	conn   *grpc.ClientConn    // Starlink Dishy gRPC connection
	client device.DeviceClient // Starlink Dishy gRPC client

	up                    prometheus.Gauge   // starlink_up
	totalScrapes          prometheus.Counter // starlink_exporter_scrapes_total
	scrapeDurationSeconds prometheus.Gauge   // starlink_exporter_scrape_duration_seconds
}

var _ prometheus.Collector = (*Exporter)(nil)

// NewExporter returns a new exporter that connects to the Starlink Dishy at
// the given address.
func NewExporter(address string) (*Exporter, error) {
	slog.Info("Connecting to Starlink Dishy", slog.String("address", address))
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	client := device.NewDeviceClient(conn)
	return &Exporter{
		conn:   conn,
		client: client,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: dishSubsystem,
			Name:      "up",
			Help:      "Whether scraping metrics from the Starlink dish was successful",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "exporter",
			Name:      "scrapes_total",
			Help:      "Total number of Starlink dish scrapes",
		}),
		scrapeDurationSeconds: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporterSubsystem,
			Name:      "scrape_duration_seconds",
			Help:      "Time taken to scrape metrics from the Starlink dish",
		}),
	}, nil
}

// ConnState returns the gRPC connection state.
func (e *Exporter) ConnState() connectivity.State {
	return e.conn.GetState()
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up.Desc()
	ch <- e.totalScrapes.Desc()
	ch <- e.scrapeDurationSeconds.Desc()

	// Scraped by scrapeDishStatus
	ch <- dishInfo
	ch <- dishUptimeSeconds
	ch <- dishSnrAboveNoiseFloor
	ch <- dishSnrPersistentlyLow
	ch <- dishUplinkThroughputBps
	ch <- dishDownlinkThroughputBps
	ch <- dishDownlinkThroughputHistogram
	ch <- dishUplinkThroughputHistogram
	ch <- dishPopPingDropRatio
	ch <- dishPopPingLatencySeconds
	ch <- dishPopPingLatencyHistogram
	ch <- dishSoftwareUpdateRebootReady
	ch <- dishBoresightAzimuthDeg
	ch <- dishBoresightElevationDeg
	ch <- dishCurrentlyObstructed
	ch <- dishFractionObstructionRatio
	ch <- dishLast24HoursObstructedSeconds
	ch <- dishProlongedObstructionDurationSeconds
	ch <- dishPowerInputHistogram
	ch <- dishPowerInput
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mx.Lock()
	defer e.mx.Unlock()

	up := e.scrape(ch)
	e.up.Set(btof(up))

	ch <- e.up
	ch <- e.totalScrapes
	ch <- e.scrapeDurationSeconds
}

// Close closes the gRPC connection.
func (e *Exporter) Close() {
	if err := e.conn.Close(); err != nil {
		slog.Error("An error occurred while closing gRPC connection", slog.Any("err", err))
	}
}
