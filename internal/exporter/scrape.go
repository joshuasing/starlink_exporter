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
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/joshuasing/starlink_exporter/internal/spacex/api/device"
)

// scrape scrapes metrics from the Starlink Dishy.
func (e *Exporter) scrape(ch chan<- prometheus.Metric) bool {
	e.totalScrapes.Inc()
	start := time.Now()
	defer func() {
		e.scrapeDurationSeconds.Set(time.Since(start).Seconds())
	}()

	return runScrapers(ch, e.scrapeDishStatus)
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
		slog.Error("Failed to scrape dish context", slog.Any("err", err))
		return false
	}

	dishStatus := res.GetDishGetStatus()
	deviceInfo := dishStatus.GetDeviceInfo()
	deviceState := dishStatus.GetDeviceState()
	obstructionStats := dishStatus.GetObstructionStats()

	// starlink_dish_info
	ch <- prometheus.MustNewConstMetric(
		dishInfo, prometheus.GaugeValue, 1,
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
		dishUptimeSeconds, prometheus.GaugeValue,
		float64(deviceState.GetUptimeS()),
	)

	// starlink_dish_snr_above_noise_floor
	ch <- prometheus.MustNewConstMetric(
		dishSnrAboveNoiseFloor, prometheus.GaugeValue,
		btof(dishStatus.GetIsSnrAboveNoiseFloor()),
	)

	// starlink_dish_snr_persistently_low
	ch <- prometheus.MustNewConstMetric(
		dishSnrPersistentlyLow, prometheus.GaugeValue,
		btof(dishStatus.GetIsSnrPersistentlyLow()),
	)

	// starlink_dish_uplink_throughput_bps
	ch <- prometheus.MustNewConstMetric(
		dishUplinkThroughputBps, prometheus.GaugeValue,
		float64(dishStatus.GetUplinkThroughputBps()),
	)

	// starlink_dish_downlink_throughput_bps
	ch <- prometheus.MustNewConstMetric(
		dishDownlinkThroughputBps, prometheus.GaugeValue,
		float64(dishStatus.GetDownlinkThroughputBps()),
	)

	// starlink_dish_pop_ping_drop_ratio
	ch <- prometheus.MustNewConstMetric(
		dishPopPingDropRatio, prometheus.GaugeValue,
		float64(dishStatus.GetPopPingDropRate()),
	)

	// starlink_dish_pop_ping_latency_seconds
	ch <- prometheus.MustNewConstMetric(
		dishPopPingLatencySeconds, prometheus.GaugeValue,
		float64(dishStatus.GetPopPingLatencyMs()/1000),
	)

	// starlink_dish_software_update_state
	ch <- prometheus.MustNewConstMetric(
		dishSoftwareUpdateState, prometheus.GaugeValue,
		float64(dishStatus.GetSoftwareUpdateState()),
	)

	// starlink_dish_software_update_reboot_ready
	ch <- prometheus.MustNewConstMetric(
		dishSoftwareUpdateRebootReady, prometheus.GaugeValue,
		btof(dishStatus.GetSwupdateRebootReady()),
	)

	// starlink_dish_boresight_azimuth_deg
	ch <- prometheus.MustNewConstMetric(
		dishBoresightAzimuthDeg, prometheus.GaugeValue,
		float64(dishStatus.GetBoresightAzimuthDeg()),
	)

	// starlink_dish_boresight_elevation_deg
	ch <- prometheus.MustNewConstMetric(
		dishBoresightElevationDeg, prometheus.GaugeValue,
		float64(dishStatus.GetBoresightElevationDeg()),
	)

	// starlink_dish_currently_obstructed
	ch <- prometheus.MustNewConstMetric(
		dishCurrentlyObstructed, prometheus.GaugeValue,
		btof(obstructionStats.GetCurrentlyObstructed()),
	)

	// starlink_dish_fraction_obstruction_ratio
	ch <- prometheus.MustNewConstMetric(
		dishFractionObstructionRatio, prometheus.GaugeValue,
		float64(obstructionStats.GetFractionObstructed()),
	)

	// starlink_dish_last_24h_obstructed_seconds
	ch <- prometheus.MustNewConstMetric(
		dishLast24HoursObstructedSeconds, prometheus.GaugeValue,
		float64(obstructionStats.GetTimeObstructed()),
	)

	// starlink_dish_prolonged_obstruction_duration_seconds
	ch <- prometheus.MustNewConstMetric(
		dishProlongedObstructionDurationSeconds, prometheus.GaugeValue,
		float64(obstructionStats.GetAvgProlongedObstructionDurationS()),
	)

	return true
}