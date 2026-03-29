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

package exporter

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestRunScrapersAllSucceed(t *testing.T) {
	t.Parallel()
	ch := make(chan prometheus.Metric, 100)
	ok := runScrapers(ch,
		func(_ context.Context, _ chan<- prometheus.Metric) bool { return true },
		func(_ context.Context, _ chan<- prometheus.Metric) bool { return true },
	)
	if !ok {
		t.Error("runScrapers returned false, want true")
	}
}

func TestRunScrapersOneFails(t *testing.T) {
	t.Parallel()
	ch := make(chan prometheus.Metric, 100)
	ok := runScrapers(ch,
		func(_ context.Context, _ chan<- prometheus.Metric) bool { return true },
		func(_ context.Context, _ chan<- prometheus.Metric) bool { return false },
	)
	if ok {
		t.Error("runScrapers returned true, want false")
	}
}

func TestRunScrapersNoScrapers(t *testing.T) {
	t.Parallel()
	ch := make(chan prometheus.Metric, 100)
	ok := runScrapers(ch)
	if !ok {
		t.Error("runScrapers with no scrapers returned false, want true")
	}
}

func TestRunScrapersEmitsMetrics(t *testing.T) {
	t.Parallel()
	ch := make(chan prometheus.Metric, 100)
	d := &Desc{
		Namespace: "test",
		Name:      "gauge",
		Help:      "a test gauge",
	}

	ok := runScrapers(ch,
		func(_ context.Context, ch chan<- prometheus.Metric) bool {
			ch <- metric(d, prometheus.GaugeValue, 42)
			return true
		},
	)
	if !ok {
		t.Fatal("runScrapers returned false")
	}
	if len(ch) != 1 {
		t.Errorf("got %d metrics, want 1", len(ch))
	}
}

func TestRunScrapersReceivesContext(t *testing.T) {
	t.Parallel()
	ch := make(chan prometheus.Metric, 100)
	ok := runScrapers(ch,
		func(ctx context.Context, _ chan<- prometheus.Metric) bool {
			if ctx == nil {
				t.Error("context is nil")
				return false
			}
			if _, ok := ctx.Deadline(); !ok {
				t.Error("context has no deadline")
				return false
			}
			return true
		},
	)
	if !ok {
		t.Error("runScrapers returned false")
	}
}
