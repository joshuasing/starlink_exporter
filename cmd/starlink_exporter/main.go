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

// Package main is a Prometheus Exporter for Starlink Dishy metrics.
package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	_ "net/http/pprof" //nolint: gosec // pprof is exposed intentionally.
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	versionCollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"google.golang.org/grpc/connectivity"

	"github.com/joshuasing/starlink_exporter/internal/exporter"
)

const defaultListenAddress = ":9451"

var (
	listenAddress = flag.String("listen", defaultListenAddress, "Listen address")
	dishAddress   = flag.String("dish", exporter.DefaultDishAddress, "Dish address")
)

func main() {
	flag.Parse()
	os.Exit(run())
}

func run() int {
	slog.Info("Starting Starlink exporter", slog.String("version", version.Info()))
	slog.Info("Build context", slog.String("build_context", version.BuildContext()))

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
	r.MustRegister(versionCollector.NewCollector("starlink_exporter"))

	// Health check handler.
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		connState := ex.ConnState()
		switch connState {
		case connectivity.Ready, connectivity.Idle:
			w.WriteHeader(http.StatusOK)
		case connectivity.Connecting, connectivity.TransientFailure:
			w.WriteHeader(http.StatusServiceUnavailable)
		case connectivity.Shutdown:
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = w.Write([]byte(strings.ToLower(connState.String())))
	})

	// Metrics handler.
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	// Landing page.
	landingPage, err := web.NewLandingPage(web.LandingConfig{
		Name:        "Starlink Exporter",
		Description: "A simple Prometheus exporter for Starlink",
		Version:     version.Info(),
		Links: []web.LandingLinks{
			{Address: "/metrics", Text: "Metrics"},
			{Address: "/health", Text: "Health"},
			{Address: "https://github.com/joshuasing/starlink_exporter", Text: "GitHub"},
		},
	})
	if err != nil {
		slog.Error("Failed to create landing page", slog.String("err", err.Error()))
		return 1
	}
	http.Handle("/", landingPage)

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
