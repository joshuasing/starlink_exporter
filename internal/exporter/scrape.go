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
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/joshuasing/starlink_exporter/internal/spacex_api/device"
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
	dishPopPingDropRatioHistOpts = prometheus.HistogramOpts{
		Name:    dishPopPingDropRatioHistogram.FQName(),
		Help:    dishPopPingDropRatioHistogram.Help,
		Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.25, 0.5, 0.75, 1},
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
	out := make(chan bool, len(scrapers))

	for _, s := range scrapers {
		wg.Go(func() {
			out <- s(ctx, ch)
		})
	}

	wg.Wait()
	close(out)

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
	alerts := dishStatus.GetAlerts()

	// starlink_dish_info
	ch <- metric(dishInfo, prometheus.GaugeValue, 1,
		deviceInfo.GetId(),
		deviceInfo.GetHardwareVersion(),
		itos(deviceInfo.GetBoardRev()),
		deviceInfo.GetSoftwareVersion(),
		deviceInfo.GetManufacturedVersion(),
		itos(deviceInfo.GetGenerationNumber()),
		deviceInfo.GetCountryCode(),
		itos(deviceInfo.GetUtcOffsetS()),
		itos(deviceInfo.GetBootcount()),
		dishStatus.GetMobilityClass().String(),
	)

	// starlink_dish_uptime_seconds
	ch <- metric(dishUptimeSeconds, prometheus.GaugeValue,
		deviceState.GetUptimeS())

	// starlink_dish_mobility_class
	ch <- metric(dishMobilityClass, prometheus.GaugeValue,
		dishStatus.GetMobilityClass())

	// starlink_dish_pop_ping_latency_seconds
	ch <- metric(dishPopPingLatencySeconds, prometheus.GaugeValue,
		dishStatus.GetPopPingLatencyMs()/1000)

	// starlink_dish_gps_valid
	ch <- metric(dishGPSValid, prometheus.GaugeValue,
		btof(dishStatus.GetGpsStats().GetGpsValid()))

	// starlink_dish_gps_satellites
	ch <- metric(dishGPSSatellites, prometheus.GaugeValue,
		dishStatus.GetGpsStats().GetGpsSats())

	// starlink_dish_snr_above_noise_floor
	ch <- metric(dishSnrAboveNoiseFloor, prometheus.GaugeValue,
		btof(dishStatus.GetIsSnrAboveNoiseFloor()))

	// starlink_dish_snr_persistently_low
	ch <- metric(dishSnrPersistentlyLow, prometheus.GaugeValue,
		btof(dishStatus.GetIsSnrPersistentlyLow()))

	// starlink_dish_pop_ping_drop_ratio
	ch <- metric(dishPopPingDropRatio, prometheus.GaugeValue,
		dishStatus.GetPopPingDropRate())

	// starlink_dish_software_update_state
	ch <- metric(dishSoftwareUpdateState, prometheus.GaugeValue,
		dishStatus.GetSoftwareUpdateState())

	// starlink_dish_software_update_reboot_ready
	ch <- metric(dishSoftwareUpdateRebootReady, prometheus.GaugeValue,
		btof(dishStatus.GetSwupdateRebootReady()))

	// starlink_dish_tilt_angle_deg
	ch <- metric(dishTiltAngleDeg, prometheus.GaugeValue,
		dishStatus.GetAlignmentStats().GetTiltAngleDeg())

	// starlink_dish_boresight_azimuth_deg
	ch <- metric(dishBoresightAzimuthDeg, prometheus.GaugeValue,
		dishStatus.GetAlignmentStats().GetBoresightAzimuthDeg())

	// starlink_dish_boresight_elevation_deg
	ch <- metric(dishBoresightElevationDeg, prometheus.GaugeValue,
		dishStatus.GetAlignmentStats().GetBoresightElevationDeg())

	// starlink_dish_desired_boresight_azimuth_deg
	ch <- metric(dishDesiredBoresightAzimuthDeg, prometheus.GaugeValue,
		dishStatus.GetAlignmentStats().GetDesiredBoresightAzimuthDeg())

	// starlink_dish_desired_boresight_elevation_deg
	ch <- metric(dishDesiredBoresightElevationDeg, prometheus.GaugeValue,
		dishStatus.GetAlignmentStats().GetDesiredBoresightElevationDeg())

	// starlink_dish_currently_obstructed
	ch <- metric(dishCurrentlyObstructed, prometheus.GaugeValue,
		btof(obstructionStats.GetCurrentlyObstructed()))

	// starlink_dish_fraction_obstruction_ratio
	ch <- metric(dishFractionObstructionRatio, prometheus.GaugeValue,
		obstructionStats.GetFractionObstructed())

	// starlink_dish_last_24h_obstructed_seconds
	ch <- metric(dishLast24HoursObstructedSeconds, prometheus.GaugeValue,
		obstructionStats.GetTimeObstructed())

	// starlink_dish_alert_unexpected_location
	ch <- metric(dishAlertUnexpectedLocation, prometheus.GaugeValue,
		btof(alerts.GetUnexpectedLocation()))

	// starlink_dish_alert_install_pending
	ch <- metric(dishAlertInstallPending, prometheus.GaugeValue,
		btof(alerts.GetInstallPending()))

	// starlink_dish_alert_is_heating
	ch <- metric(dishAlertIsHeating, prometheus.GaugeValue,
		btof(alerts.GetIsHeating()))

	// starlink_dish_alert_is_power_save_idle
	ch <- metric(dishAlertIsPowerSaveIdle, prometheus.GaugeValue,
		btof(alerts.GetIsPowerSaveIdle()))

	// starlink_dish_alert_lower_than_predicted
	ch <- metric(dishAlertSignalLowerThanPredicted, prometheus.GaugeValue,
		btof(alerts.GetLowerSignalThanPredicted()))

	// Additional alerts
	ch <- metric(dishAlertMotorsStuck, prometheus.GaugeValue,
		btof(alerts.GetMotorsStuck()))
	ch <- metric(dishAlertThermalThrottle, prometheus.GaugeValue,
		btof(alerts.GetThermalThrottle()))
	ch <- metric(dishAlertThermalShutdown, prometheus.GaugeValue,
		btof(alerts.GetThermalShutdown()))
	ch <- metric(dishAlertMastNotNearVertical, prometheus.GaugeValue,
		btof(alerts.GetMastNotNearVertical()))
	ch <- metric(dishAlertSlowEthernetSpeeds, prometheus.GaugeValue,
		btof(alerts.GetSlowEthernetSpeeds()))
	ch <- metric(dishAlertSlowEthernetSpeeds100, prometheus.GaugeValue,
		btof(alerts.GetSlowEthernetSpeeds_100()))
	ch <- metric(dishAlertRoaming, prometheus.GaugeValue,
		btof(alerts.GetRoaming()))
	ch <- metric(dishAlertPowerSupplyThermalThrottle, prometheus.GaugeValue,
		btof(alerts.GetPowerSupplyThermalThrottle()))
	ch <- metric(dishAlertDbfTelemStale, prometheus.GaugeValue,
		btof(alerts.GetDbfTelemStale()))
	ch <- metric(dishAlertLowMotorCurrent, prometheus.GaugeValue,
		btof(alerts.GetLowMotorCurrent()))
	ch <- metric(dishAlertObstructionMapReset, prometheus.GaugeValue,
		btof(alerts.GetObstructionMapReset()))
	ch <- metric(dishAlertDishWaterDetected, prometheus.GaugeValue,
		btof(alerts.GetDishWaterDetected()))
	ch <- metric(dishAlertRouterWaterDetected, prometheus.GaugeValue,
		btof(alerts.GetRouterWaterDetected()))
	ch <- metric(dishAlertUpsuRouterPortSlow, prometheus.GaugeValue,
		btof(alerts.GetUpsuRouterPortSlow()))
	ch <- metric(dishAlertNoEthernetLink, prometheus.GaugeValue,
		btof(alerts.GetNoEthernetLink()))

	// Top-level status scalars
	ch <- metric(dishSecondsToFirstNonemptySlot, prometheus.GaugeValue,
		dishStatus.GetSecondsToFirstNonemptySlot())
	ch <- metric(dishStowRequested, prometheus.GaugeValue,
		btof(dishStatus.GetStowRequested()))
	ch <- metric(dishEthSpeedMbps, prometheus.GaugeValue,
		dishStatus.GetEthSpeedMbps())
	ch <- metric(dishClassOfService, prometheus.GaugeValue,
		dishStatus.GetClassOfService())
	ch <- metric(dishRebootReason, prometheus.GaugeValue,
		dishStatus.GetRebootReason())
	ch <- metric(dishDisablementCode, prometheus.GaugeValue,
		dishStatus.GetDisablementCode())
	ch <- metric(dishDlBandwidthRestrictedReason, prometheus.GaugeValue,
		dishStatus.GetDlBandwidthRestrictedReason())
	ch <- metric(dishUlBandwidthRestrictedReason, prometheus.GaugeValue,
		dishStatus.GetUlBandwidthRestrictedReason())
	ch <- metric(dishIsCellDisabled, prometheus.GaugeValue,
		btof(dishStatus.GetIsCellDisabled()))
	ch <- metric(dishSecondsUntilSwupdateRebootPossible, prometheus.GaugeValue,
		dishStatus.GetSecondsUntilSwupdateRebootPossible())
	ch <- metric(dishHighPowerTestMode, prometheus.GaugeValue,
		btof(dishStatus.GetHighPowerTestMode()))
	ch <- metric(dishIsMovingFastPersisted, prometheus.GaugeValue,
		btof(dishStatus.GetIsMovingFastPersisted()))
	ch <- metric(dishMacFlag, prometheus.GaugeValue,
		btof(dishStatus.GetMacFlag()))
	ch <- metric(dishNatFlag, prometheus.GaugeValue,
		dishStatus.GetNatFlag())
	ch <- metric(dishAccountShard, prometheus.GaugeValue,
		dishStatus.GetAccountShard())
	ch <- metric(dishConnectedRoutersCount, prometheus.GaugeValue,
		len(dishStatus.GetConnectedRouters()))
	ch <- metric(dishHasActuators, prometheus.GaugeValue,
		dishStatus.GetHasActuators())

	// NED-to-dish orientation quaternion
	quat := dishStatus.GetNed2DishQuaternion()
	ch <- metric(dishNed2DishQuaternionW, prometheus.GaugeValue, quat.GetQScalar())
	ch <- metric(dishNed2DishQuaternionX, prometheus.GaugeValue, quat.GetQX())
	ch <- metric(dishNed2DishQuaternionY, prometheus.GaugeValue, quat.GetQY())
	ch <- metric(dishNed2DishQuaternionZ, prometheus.GaugeValue, quat.GetQZ())

	// Current outage
	if outage := dishStatus.GetOutage(); outage != nil {
		ch <- metric(dishOutageInfo, prometheus.GaugeValue, 1,
			outage.GetCause().String(),
			strconv.FormatBool(outage.GetDidSwitch()),
		)
		ch <- metric(dishOutageStartTimestampSeconds, prometheus.GaugeValue,
			float64(outage.GetStartTimestampNs())/1e9)
		ch <- metric(dishOutageDurationSeconds, prometheus.GaugeValue,
			float64(outage.GetDurationNs())/1e9)
	}

	// GPS (additional)
	ch <- metric(dishGPSNoSatsAfterTtff, prometheus.GaugeValue,
		btof(dishStatus.GetGpsStats().GetNoSatsAfterTtff()))
	ch <- metric(dishGPSInhibit, prometheus.GaugeValue,
		btof(dishStatus.GetGpsStats().GetInhibitGps()))

	// Obstruction (additional)
	ch <- metric(dishObstructionValidSeconds, prometheus.GaugeValue,
		obstructionStats.GetValidS())
	ch <- metric(dishObstructionPatchesValid, prometheus.GaugeValue,
		obstructionStats.GetPatchesValid())
	ch <- metric(dishAvgProlongedObstructionDurationSeconds, prometheus.GaugeValue,
		obstructionStats.GetAvgProlongedObstructionDurationS())
	ch <- metric(dishAvgProlongedObstructionIntervalSeconds, prometheus.GaugeValue,
		obstructionStats.GetAvgProlongedObstructionIntervalS())
	ch <- metric(dishAvgProlongedObstructionValid, prometheus.GaugeValue,
		btof(obstructionStats.GetAvgProlongedObstructionValid()))

	// Ready states
	ready := dishStatus.GetReadyStates()
	ch <- metric(dishReadyStateCady, prometheus.GaugeValue, btof(ready.GetCady()))
	ch <- metric(dishReadyStateScp, prometheus.GaugeValue, btof(ready.GetScp()))
	ch <- metric(dishReadyStateL1L2, prometheus.GaugeValue, btof(ready.GetL1L2()))
	ch <- metric(dishReadyStateXphy, prometheus.GaugeValue, btof(ready.GetXphy()))
	ch <- metric(dishReadyStateAap, prometheus.GaugeValue, btof(ready.GetAap()))
	ch <- metric(dishReadyStateRf, prometheus.GaugeValue, btof(ready.GetRf()))

	// Software update stats
	swStats := dishStatus.GetSoftwareUpdateStats()
	ch <- metric(dishSoftwareUpdateProgress, prometheus.GaugeValue,
		swStats.GetSoftwareUpdateProgress())
	ch <- metric(dishSoftwareUpdateRequiresReboot, prometheus.GaugeValue,
		btof(swStats.GetUpdateRequiresReboot()))
	ch <- metric(dishSoftwareUpdateRebootScheduledUtcSeconds, prometheus.GaugeValue,
		swStats.GetRebootScheduledUtcTime())

	// Alignment (additional)
	alignment := dishStatus.GetAlignmentStats()
	ch <- metric(dishActuatorState, prometheus.GaugeValue,
		alignment.GetActuatorState())
	ch <- metric(dishAttitudeEstimationState, prometheus.GaugeValue,
		alignment.GetAttitudeEstimationState())
	ch <- metric(dishAttitudeUncertaintyDeg, prometheus.GaugeValue,
		alignment.GetAttitudeUncertaintyDeg())

	// Initialization durations
	init := dishStatus.GetInitializationDurationSeconds()
	ch <- metric(dishInitAttitudeSeconds, prometheus.GaugeValue,
		init.GetAttitudeInitialization())
	ch <- metric(dishInitBurstDetectedSeconds, prometheus.GaugeValue,
		init.GetBurstDetected())
	ch <- metric(dishInitEkfConvergedSeconds, prometheus.GaugeValue,
		init.GetEkfConverged())
	ch <- metric(dishInitFirstCplaneSeconds, prometheus.GaugeValue,
		init.GetFirstCplane())
	ch <- metric(dishInitFirstPopPingSeconds, prometheus.GaugeValue,
		init.GetFirstPopPing())
	ch <- metric(dishInitGpsValidSeconds, prometheus.GaugeValue,
		init.GetGpsValid())
	ch <- metric(dishInitInitialNetworkEntrySeconds, prometheus.GaugeValue,
		init.GetInitialNetworkEntry())
	ch <- metric(dishInitNetworkScheduleSeconds, prometheus.GaugeValue,
		init.GetNetworkSchedule())
	ch <- metric(dishInitRfReadySeconds, prometheus.GaugeValue,
		init.GetRfReady())
	ch <- metric(dishInitStableConnectionSeconds, prometheus.GaugeValue,
		init.GetStableConnection())

	// PLC (Mini battery)
	plc := dishStatus.GetPlcStats()
	ch <- metric(dishPlcReceiving, prometheus.GaugeValue,
		btof(plc.GetReceivingPlc()))
	ch <- metric(dishPlcAverageTimeToEmptySeconds, prometheus.GaugeValue,
		plc.GetAverageTimeToEmpty())
	ch <- metric(dishPlcAverageTimeToFullSeconds, prometheus.GaugeValue,
		plc.GetAverageTimeToFull())
	ch <- metric(dishPlcBatteryHealth, prometheus.GaugeValue,
		plc.GetBatteryHealth())
	ch <- metric(dishPlcPermanentFailure, prometheus.GaugeValue,
		btof(plc.GetPermanentFailure()))
	ch <- metric(dishPlcSafetyModeActive, prometheus.GaugeValue,
		btof(plc.GetSafetyModeActive()))
	ch <- metric(dishPlcStateOfChargePercent, prometheus.GaugeValue,
		plc.GetStateOfCharge())
	ch <- metric(dishPlcThermalThrottleLevel, prometheus.GaugeValue,
		plc.GetThermalThrottleLevel())
	ch <- metric(dishPlcRevision, prometheus.GaugeValue,
		plc.GetPlcRevision())

	// UPSU
	upsu := dishStatus.GetUpsuStats()
	ch <- metric(dishUpsuDishPowerWatts, prometheus.GaugeValue,
		upsu.GetDishPower())
	ch <- metric(dishUpsuRouterPowerWatts, prometheus.GaugeValue,
		upsu.GetRouterPower())
	ch <- metric(dishUpsuUptimeSeconds, prometheus.GaugeValue,
		upsu.GetUptime())
	ch <- metric(dishUpsuBoardRev, prometheus.GaugeValue,
		upsu.GetBoardRev())

	// APS
	aps := dishStatus.GetApsStats()
	ch <- metric(dishApsDishPowerWatts, prometheus.GaugeValue,
		aps.GetDishPower())
	ch <- metric(dishApsUptimeSeconds, prometheus.GaugeValue,
		aps.GetUptime())
	ch <- metric(dishApsBoardRev, prometheus.GaugeValue,
		aps.GetBoardRev())

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
		ch <- metric(dishDownlinkThroughputBps, prometheus.GaugeValue,
			downlinkData[len(downlinkData)-1])
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
		ch <- metric(dishUplinkThroughputBps, prometheus.GaugeValue,
			uplinkData[len(uplinkData)-1])
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
		ch <- metric(dishPowerInput, prometheus.GaugeValue,
			powerData[len(powerData)-1])
	}

	// starlink_dish_pop_ping_drop_ratio_histogram
	dropData := parseRingBuffer(dishHistory.GetPopPingDropRate(), dishHistory.GetCurrent())
	dropHist := prometheus.NewHistogram(dishPopPingDropRatioHistOpts)
	for _, drop := range dropData {
		dropHist.Observe(float64(drop))
	}
	ch <- dropHist

	// starlink_dish_history_outages_count
	ch <- metric(dishHistoryOutagesCount, prometheus.GaugeValue,
		len(dishHistory.GetOutages()))

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
	ch <- metric(dishLocationInfo, prometheus.GaugeValue, 1,
		loc.GetSource().String(),
		ftos(lla.GetLat()),
		ftos(lla.GetLon()),
		ftos(lla.GetAlt()),
	)

	// starlink_dish_location_latitude_deg
	ch <- metric(dishLocationLatitude, prometheus.GaugeValue, lla.GetLat())

	// starlink_dish_location_longitude_deg
	ch <- metric(dishLocationLongitude, prometheus.GaugeValue, lla.GetLon())

	// starlink_dish_location_altitude_meters
	ch <- metric(dishLocationAltitude, prometheus.GaugeValue, lla.GetAlt())

	// starlink_dish_location_uncertainty_meters
	ch <- metric(dishLocationUncertaintyMeters, prometheus.GaugeValue, loc.GetSigmaM())

	// starlink_dish_horizontal_speed_mps
	ch <- metric(dishHorizontalSpeedMps, prometheus.GaugeValue, loc.GetHorizontalSpeedMps())

	// starlink_dish_vertical_speed_mps
	ch <- metric(dishVerticalSpeedMps, prometheus.GaugeValue, loc.GetVerticalSpeedMps())

	return true
}
