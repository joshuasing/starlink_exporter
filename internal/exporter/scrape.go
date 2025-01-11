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

package exporter

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/joshuasing/starlink_exporter/internal/spacex/api/device"
)

var (
	dishPopPingLatencyHistOpts = prometheus.HistogramOpts{
		Name:    dishPopPingLatencyHistogram.FQName(),
		Help:    dishPopPingLatencyHistogram.Help,
		Buckets: []float64{0.01, 0.02, 0.021, 0.022, 0.023, 0.024, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}
	dishDownlinkThroughputHistOpts = prometheus.HistogramOpts{
		Name:    dishDownlinkThroughputHistogram.FQName(),
		Help:    dishDownlinkThroughputHistogram.Help,
		Buckets: []float64{1e6, 5e6, 10e6, 25e6, 50e6, 100e6, 250e6, 500e6},
	}
	dishUplinkThroughputHistOpts = prometheus.HistogramOpts{
		Name:    dishUplinkThroughputHistogram.FQName(),
		Help:    dishUplinkThroughputHistogram.Help,
		Buckets: []float64{1e6, 5e6, 10e6, 25e6, 50e6, 100e6, 250e6, 500e6},
	}
	dishPowerInputHistOpts = prometheus.HistogramOpts{
		Name:    dishPowerInputHistogram.FQName(),
		Help:    dishPowerInputHistogram.Help,
		Buckets: []float64{20, 25, 30, 40, 50, 75, 100, 150, 200},
	}
)

// scrape scrapes metrics from the Starlink Dishy.
func (e *Exporter) scrape(ch chan<- prometheus.Metric) bool {
	e.totalScrapes.Inc()
	start := time.Now()
	defer func() {
		e.scrapeDurationSeconds.Set(time.Since(start).Seconds())
	}()

	return runScrapers(ch,
		e.scrapeDishStatus,
		e.scrapeDishHistory,
		e.scrapeLocation,
	)
}

type scraper func(ctx context.Context, ch chan<- prometheus.Metric) bool

// runScrapers runs the scrapers in parallel and returns true if all succeed,
// otherwise false is returned.
func runScrapers(ch chan<- prometheus.Metric, scrapers ...scraper) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	out := make(chan bool)

	wg.Add(len(scrapers))
	for _, s := range scrapers {
		go func() {
			defer wg.Done()
			out <- s(ctx, ch)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for s := range out {
		if !s {
			return false
		}
	}
	return true
}

// scrapeDishStatus scrapes metrics from the GetStatus response.
func (e *Exporter) scrapeDishStatus(ctx context.Context, ch chan<- prometheus.Metric) bool {
	res, err := e.client.Handle(ctx, &device.Request{
		Request: new(device.Request_GetStatus),
	})
	if err != nil {
		slog.Error("Failed to scrape dish status", slog.Any("err", err))
		return false
	}

	dishStatus := res.GetDishGetStatus()
	deviceInfo := dishStatus.GetDeviceInfo()
	deviceState := dishStatus.GetDeviceState()
	obstructionStats := dishStatus.GetObstructionStats()

	// starlink_dish_info
	ch <- prometheus.MustNewConstMetric(
		dishInfo.Desc(), prometheus.GaugeValue, 1,
		deviceInfo.GetId(),
		deviceInfo.GetHardwareVersion(),
		itos(deviceInfo.GetBoardRev()),
		deviceInfo.GetSoftwareVersion(),
		deviceInfo.GetManufacturedVersion(),
		itos(deviceInfo.GetGenerationNumber()),
		deviceInfo.GetCountryCode(),
		itos(deviceInfo.GetUtcOffsetS()),
		itos(deviceInfo.GetBootcount()),
	)

	// starlink_dish_uptime_seconds
	ch <- prometheus.MustNewConstMetric(
		dishUptimeSeconds.Desc(), prometheus.GaugeValue,
		float64(deviceState.GetUptimeS()),
	)

	// starlink_dish_pop_ping_latency_seconds
	ch <- prometheus.MustNewConstMetric(
		dishPopPingLatencySeconds.Desc(), prometheus.GaugeValue,
		float64(dishStatus.GetPopPingLatencyMs()/1000),
	)

	// starlink_dish_gps_valid
	ch <- prometheus.MustNewConstMetric(
		dishGPSValid.Desc(), prometheus.GaugeValue,
		btof(dishStatus.GetGpsStats().GetGpsValid()),
	)

	// starlink_dish_gps_satellites
	ch <- prometheus.MustNewConstMetric(
		dishGPSSatellites.Desc(), prometheus.GaugeValue,
		float64(dishStatus.GetGpsStats().GetGpsSats()),
	)

	// starlink_dish_snr_above_noise_floor
	ch <- prometheus.MustNewConstMetric(
		dishSnrAboveNoiseFloor.Desc(), prometheus.GaugeValue,
		btof(dishStatus.GetIsSnrAboveNoiseFloor()),
	)

	// starlink_dish_snr_persistently_low
	ch <- prometheus.MustNewConstMetric(
		dishSnrPersistentlyLow.Desc(), prometheus.GaugeValue,
		btof(dishStatus.GetIsSnrPersistentlyLow()),
	)

	// starlink_dish_pop_ping_drop_ratio
	ch <- prometheus.MustNewConstMetric(
		dishPopPingDropRatio.Desc(), prometheus.GaugeValue,
		float64(dishStatus.GetPopPingDropRate()),
	)

	// starlink_dish_software_update_state
	ch <- prometheus.MustNewConstMetric(
		dishSoftwareUpdateState.Desc(), prometheus.GaugeValue,
		float64(dishStatus.GetSoftwareUpdateState()),
	)

	// starlink_dish_software_update_reboot_ready
	ch <- prometheus.MustNewConstMetric(
		dishSoftwareUpdateRebootReady.Desc(), prometheus.GaugeValue,
		btof(dishStatus.GetSwupdateRebootReady()),
	)

	// starlink_dish_tilt_angle_deg
	ch <- prometheus.MustNewConstMetric(
		dishTiltAngleDeg.Desc(), prometheus.GaugeValue,
		float64(dishStatus.GetAlignmentStats().GetTiltAngleDeg()),
	)

	// starlink_dish_boresight_azimuth_deg
	ch <- prometheus.MustNewConstMetric(
		dishBoresightAzimuthDeg.Desc(), prometheus.GaugeValue,
		float64(dishStatus.GetAlignmentStats().GetBoresightAzimuthDeg()),
	)

	// starlink_dish_boresight_elevation_deg
	ch <- prometheus.MustNewConstMetric(
		dishBoresightElevationDeg.Desc(), prometheus.GaugeValue,
		float64(dishStatus.GetAlignmentStats().GetBoresightElevationDeg()),
	)

	// starlink_dish_desired_boresight_azimuth_deg
	ch <- prometheus.MustNewConstMetric(
		dishDesiredBoresightAzimuthDeg.Desc(), prometheus.GaugeValue,
		float64(dishStatus.GetAlignmentStats().GetDesiredBoresightAzimuthDeg()),
	)

	// starlink_dish_desired_boresight_elevation_deg
	ch <- prometheus.MustNewConstMetric(
		dishDesiredBoresightElevationDeg.Desc(), prometheus.GaugeValue,
		float64(dishStatus.GetAlignmentStats().GetDesiredBoresightElevationDeg()),
	)

	// starlink_dish_currently_obstructed
	ch <- prometheus.MustNewConstMetric(
		dishCurrentlyObstructed.Desc(), prometheus.GaugeValue,
		btof(obstructionStats.GetCurrentlyObstructed()),
	)

	// starlink_dish_fraction_obstruction_ratio
	ch <- prometheus.MustNewConstMetric(
		dishFractionObstructionRatio.Desc(), prometheus.GaugeValue,
		float64(obstructionStats.GetFractionObstructed()),
	)

	// starlink_dish_last_24h_obstructed_seconds
	ch <- prometheus.MustNewConstMetric(
		dishLast24HoursObstructedSeconds.Desc(), prometheus.GaugeValue,
		float64(obstructionStats.GetTimeObstructed()),
	)

	return true
}

func (e *Exporter) scrapeDishHistory(ctx context.Context, ch chan<- prometheus.Metric) bool {
	res, err := e.client.Handle(ctx, &device.Request{
		Request: new(device.Request_GetHistory),
	})
	if err != nil {
		slog.Error("Failed to scrape dish history", slog.Any("err", err))
		return false
	}

	dishHistory := res.GetDishGetHistory()

	// starlink_dish_pop_ping_latency_histogram
	latencyData := parseRingBuffer(dishHistory.GetPopPingLatencyMs(), dishHistory.GetCurrent())
	latencyHist := prometheus.NewHistogram(dishPopPingLatencyHistOpts)
	for _, latency := range latencyData {
		latencyHist.Observe(float64(latency) / 1000) // Convert ms to seconds
	}
	ch <- latencyHist

	// starlink_dish_downlink_throughput_histogram
	downlinkData := parseRingBuffer(dishHistory.GetDownlinkThroughputBps(), dishHistory.GetCurrent())
	downlinkHist := prometheus.NewHistogram(dishDownlinkThroughputHistOpts)
	for _, throughput := range downlinkData {
		downlinkHist.Observe(float64(throughput))
	}
	ch <- downlinkHist

	// starlink_dish_downlink_throughput_bps
	if len(downlinkData) > 0 {
		ch <- prometheus.MustNewConstMetric(
			dishDownlinkThroughputBps.Desc(), prometheus.GaugeValue,
			float64(downlinkData[len(downlinkData)-1]),
		)
	}

	// starlink_dish_uplink_throughput_bps_histogram
	uplinkData := parseRingBuffer(dishHistory.GetUplinkThroughputBps(), dishHistory.GetCurrent())
	uplinkHist := prometheus.NewHistogram(dishUplinkThroughputHistOpts)
	for _, throughput := range uplinkData {
		uplinkHist.Observe(float64(throughput))
	}
	ch <- uplinkHist

	// starlink_dish_uplink_throughput_bps
	if len(uplinkData) > 0 {
		ch <- prometheus.MustNewConstMetric(
			dishUplinkThroughputBps.Desc(), prometheus.GaugeValue,
			float64(uplinkData[len(uplinkData)-1]),
		)
	}

	// starlink_dish_power_input_watts_histogram
	powerData := parseRingBuffer(dishHistory.GetPowerIn(), dishHistory.GetCurrent())
	powerHist := prometheus.NewHistogram(dishPowerInputHistOpts)
	for _, power := range powerData {
		powerHist.Observe(float64(power))
	}
	ch <- powerHist

	// starlink_dish_power_input_watts
	if len(powerData) > 0 {
		ch <- prometheus.MustNewConstMetric(
			dishPowerInput.Desc(), prometheus.GaugeValue,
			float64(powerData[len(powerData)-1]),
		)
	}

	return true
}

func (e *Exporter) scrapeLocation(ctx context.Context, ch chan<- prometheus.Metric) bool {
	res, err := e.client.Handle(ctx, &device.Request{
		Request: new(device.Request_GetLocation),
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.PermissionDenied {
				// Location service may be disabled, ignore.
				return true
			}
		}
		slog.Error("Failed to scrape location", slog.Any("err", err))
		return false
	}
	loc := res.GetGetLocation()
	lla := loc.GetLla()

	// starlink_dish_location_info
	ch <- prometheus.MustNewConstMetric(
		dishLocationInfo.Desc(), prometheus.GaugeValue, 1,
		loc.GetSource().String(),
		ftos(lla.GetLat()),
		ftos(lla.GetLon()),
		ftos(lla.GetAlt()),
	)

	// starlink_dish_location_latitude
	ch <- prometheus.MustNewConstMetric(
		dishLocationLatitude.Desc(), prometheus.GaugeValue, lla.GetLat(),
	)

	// starlink_dish_location_longitude
	ch <- prometheus.MustNewConstMetric(
		dishLocationLongitude.Desc(), prometheus.GaugeValue, lla.GetLon(),
	)

	// starlink_dish_location_altitude
	ch <- prometheus.MustNewConstMetric(
		dishLocationAltitude.Desc(), prometheus.GaugeValue, lla.GetAlt(),
	)

	return true
}
