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
	"log/slog"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/connectivity"
)

// DefaultDishAddress is the default address of the Starlink Dishy's gRPC-Web
// server. Firmware 2026.06.x moved the Device API from the legacy raw-gRPC
// port 9200 (now closed) to gRPC-Web over HTTP/1.1 on port 9201.
const DefaultDishAddress = "192.168.100.1:9201"

// Exporter is a Starlink Dishy metrics exporter.
type Exporter struct {
	mx     sync.Mutex
	client *grpcWebClient // Starlink Dishy gRPC-Web client (HTTP/1.1)

	up                    prometheus.Gauge   // starlink_dish_up
	totalScrapes          prometheus.Counter // starlink_exporter_scrapes_total
	scrapeDurationSeconds prometheus.Gauge   // starlink_exporter_scrape_duration_seconds
}

var _ prometheus.Collector = (*Exporter)(nil)

// NewExporter returns a new exporter that connects to the Starlink Dishy at
// the given address over gRPC-Web (HTTP/1.1).
func NewExporter(address string) (*Exporter, error) {
	slog.Info("Connecting to Starlink Dishy", slog.String("address", address), slog.String("transport", "grpc-web"))
	client := newGRPCWebClient(address)
	return &Exporter{
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

// ConnState returns a synthetic connection state derived from the last
// gRPC-Web call. gRPC-Web is stateless (one HTTP/1.1 request per call), so
// there is no persistent connection to query.
func (e *Exporter) ConnState() connectivity.State {
	return e.client.ConnState()
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

// Close is a no-op for the gRPC-Web transport (no persistent connection).
func (e *Exporter) Close() {}
