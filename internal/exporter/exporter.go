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

// DefaultDishAddress is the default address of the Starlink Dishy's gRPC server.
const DefaultDishAddress = "192.168.100.1:9200"

// Exporter is a Starlink Dishy metrics exporter.
type Exporter struct {
	mx     sync.Mutex
	conn   *grpc.ClientConn    // Starlink Dishy gRPC connection
	client device.DeviceClient // Starlink Dishy gRPC client

	up                    prometheus.Gauge   // starlink_dish_up
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
			Name: dishUp.FQName(),
			Help: dishUp.Help,
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: exporterScrapesTotal.FQName(),
			Help: exporterScrapesTotal.Help,
		}),
		scrapeDurationSeconds: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: exporterScrapeDurationSeconds.FQName(),
			Help: exporterScrapeDurationSeconds.Help,
		}),
	}, nil
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
