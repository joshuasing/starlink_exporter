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

	"github.com/joshuasing/starlink_exporter/internal/spacex_api/device"
)

// DefaultDishAddress is the default address of the Starlink Dishy's gRPC server.
const DefaultDishAddress = "192.168.100.1:9200"

// DefaultRouterAddress is the default address of the Starlink WiFi router's
// gRPC server.
const DefaultRouterAddress = "192.168.1.1:9000"

// Exporter is a Starlink Dishy metrics exporter.
type Exporter struct {
	mx     sync.Mutex
	conn   *grpc.ClientConn    // Starlink Dishy gRPC connection
	client device.DeviceClient // Starlink Dishy gRPC client

	wifiConn   *grpc.ClientConn    // Starlink WiFi router gRPC connection (optional)
	wifiClient device.DeviceClient // Starlink WiFi router gRPC client (optional)

	up                    prometheus.Gauge   // starlink_dish_up
	wifiUp                prometheus.Gauge   // starlink_wifi_up
	totalScrapes          prometheus.Counter // starlink_exporter_scrapes_total
	scrapeDurationSeconds prometheus.Gauge   // starlink_exporter_scrape_duration_seconds
}

var _ prometheus.Collector = (*Exporter)(nil)

// NewExporter returns a new exporter that connects to the Starlink Dishy at
// the given address. If routerAddress is non-empty, the exporter also scrapes
// WiFi router metrics from that address; router scrape failures are logged but
// do not cause the dish scrape to fail.
func NewExporter(address, routerAddress string) (*Exporter, error) {
	slog.Info("Connecting to Starlink Dishy", slog.String("address", address))
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	e := &Exporter{
		conn:   conn,
		client: device.NewDeviceClient(conn),
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: dishUp.FQName(),
			Help: dishUp.Help,
		}),
		wifiUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: wifiUp.FQName(),
			Help: wifiUp.Help,
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: exporterScrapesTotal.FQName(),
			Help: exporterScrapesTotal.Help,
		}),
		scrapeDurationSeconds: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: exporterScrapeDurationSeconds.FQName(),
			Help: exporterScrapeDurationSeconds.Help,
		}),
	}

	if routerAddress != "" {
		slog.Info("Connecting to Starlink WiFi router; scrape failures will be handled gracefully",
			slog.String("address", routerAddress))
		wifiConn, err := grpc.NewClient(routerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				slog.Error("An error occurred while closing dish gRPC connection", slog.Any("err", cerr))
			}
			return nil, fmt.Errorf("new wifi grpc client: %w", err)
		}
		e.wifiConn = wifiConn
		e.wifiClient = device.NewDeviceClient(wifiConn)
	} else {
		slog.Info("WiFi router scraping disabled")
	}

	return e, nil
}

// ConnState returns the gRPC connection state.
func (e *Exporter) ConnState() connectivity.State {
	return e.conn.GetState()
}

// Describe provides all descriptors for metrics provided by the exporter.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range Descs {
		ch <- d.Desc()
	}
}

// Collect collects metrics from the Starlink Dishy.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mx.Lock()
	defer e.mx.Unlock()

	up := e.scrape(ch)
	e.up.Set(btof(up))

	if e.wifiClient != nil {
		e.wifiUp.Set(btof(e.scrapeWifi(ch)))
		ch <- e.wifiUp
	}

	ch <- e.up
	ch <- e.totalScrapes
	ch <- e.scrapeDurationSeconds
}

// Close closes the gRPC connections.
func (e *Exporter) Close() {
	if err := e.conn.Close(); err != nil {
		slog.Error("An error occurred while closing dish gRPC connection", slog.Any("err", err))
	}
	if e.wifiConn != nil {
		if err := e.wifiConn.Close(); err != nil {
			slog.Error("An error occurred while closing WiFi router gRPC connection", slog.Any("err", err))
		}
	}
}
