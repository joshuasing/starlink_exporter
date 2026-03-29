// Copyright (c) 2026 Joshua Sing <joshua@joshuasing.dev>
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

package exporter_test

import (
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/connectivity"

	"github.com/joshuasing/starlink_exporter/internal/exporter"
)

func dishAddress(t *testing.T) string {
	t.Helper()
	addr := os.Getenv("STARLINK_DISH_ADDR")
	if addr == "" {
		addr = exporter.DefaultDishAddress
	}
	return addr
}

func newTestExporter(t *testing.T) *exporter.Exporter {
	t.Helper()
	if os.Getenv("STARLINK_INTEGRATION") == "" {
		t.Skip("set STARLINK_INTEGRATION=1 to run integration tests")
	}
	ex, err := exporter.NewExporter(dishAddress(t))
	if err != nil {
		t.Fatalf("NewExporter: %v", err)
	}
	t.Cleanup(ex.Close)
	return ex
}

func TestIntegrationConnState(t *testing.T) {
	t.Parallel()
	ex := newTestExporter(t)
	state := ex.ConnState()
	switch state {
	case connectivity.Ready, connectivity.Idle, connectivity.Connecting:
		// OK
	default:
		t.Errorf("unexpected connection state: %v", state)
	}
}

func TestIntegrationCollect(t *testing.T) {
	t.Parallel()
	ex := newTestExporter(t)

	ch := make(chan prometheus.Metric, 200)
	go func() {
		ex.Collect(ch)
		close(ch)
	}()

	var metrics []prometheus.Metric
	for m := range ch {
		metrics = append(metrics, m)
	}

	if len(metrics) == 0 {
		t.Fatal("Collect produced no metrics")
	}
	t.Logf("Collect produced %d metrics", len(metrics))
}

func TestIntegrationDescribe(t *testing.T) {
	t.Parallel()
	ex := newTestExporter(t)

	ch := make(chan *prometheus.Desc, 200)
	go func() {
		ex.Describe(ch)
		close(ch)
	}()

	var descs []*prometheus.Desc
	for d := range ch {
		descs = append(descs, d)
	}

	if len(descs) != len(exporter.Descs) {
		t.Errorf("Describe produced %d descs, want %d", len(descs), len(exporter.Descs))
	}
}

func TestIntegrationRegistration(t *testing.T) {
	t.Parallel()
	ex := newTestExporter(t)

	r := prometheus.NewRegistry()
	if err := r.Register(ex); err != nil {
		t.Fatalf("Register: %v", err)
	}

	families, err := r.Gather()
	if err != nil {
		t.Fatalf("Gather: %v", err)
	}

	if len(families) == 0 {
		t.Fatal("Gather produced no metric families")
	}
	t.Logf("Gather produced %d metric families", len(families))

	// Verify some key metrics are present.
	want := map[string]bool{
		"starlink_dish_up":             false,
		"starlink_dish_uptime_seconds": false,
		"starlink_dish_info":           false,
	}
	for _, f := range families {
		if _, ok := want[f.GetName()]; ok {
			want[f.GetName()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("expected metric %q not found", name)
		}
	}
}
