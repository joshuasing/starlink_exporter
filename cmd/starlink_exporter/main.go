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

package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/connectivity"

	"github.com/joshuasing/starlink_exporter/internal/exporter"
)

const (
	defaultListenAddress = ":9451"
	defaultDishAddress   = "192.168.100.1:9200"
)

var (
	listenAddress = flag.String("listen", defaultListenAddress, "Listen address")
	dishAddress   = flag.String("dish", defaultDishAddress, "Dish address")
)

func main() {
	flag.Parse()
	os.Exit(run())
}

func run() int {
	slog.Info("Starting Starlink exporter")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ex, err := exporter.NewExporter(*dishAddress)
	if err != nil {
		slog.Error("Failed to create exporter", slog.Any("err", err))
		return 1
	}
	defer ex.Close()

	// Prometheus registry
	r := prometheus.NewRegistry()
	r.MustRegister(ex)

	// Health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		switch ex.ConnState() {
		case connectivity.Ready, connectivity.Idle:
			w.WriteHeader(http.StatusOK)
		case connectivity.Connecting, connectivity.TransientFailure:
			w.WriteHeader(http.StatusServiceUnavailable)
		case connectivity.Shutdown:
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	// Metrics handler
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	// Run HTTP server in a goroutine
	httpErr := make(chan error)
	go func() {
		srv := &http.Server{
			Addr:         *listenAddress,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		}
		slog.Info("HTTP server listening", slog.String("address", srv.Addr))
		httpErr <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
	case err := <-httpErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start HTTP server", slog.Any("err", err))
			return 1
		}
	}

	return 0
}
