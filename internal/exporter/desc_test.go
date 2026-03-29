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
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestDescFQName(t *testing.T) {
	t.Parallel()
	d := &Desc{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "test_metric",
	}
	want := "starlink_dish_test_metric"
	if got := d.FQName(); got != want {
		t.Errorf("FQName() = %q, want %q", got, want)
	}
}

func TestDescFQNameCaching(t *testing.T) {
	t.Parallel()
	d := &Desc{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "test_metric",
	}
	first := d.FQName()
	second := d.FQName()
	if first != second {
		t.Errorf("FQName() not stable: %q != %q", first, second)
	}
	if d.fqName == "" {
		t.Error("FQName() did not cache the result")
	}
}

func TestDescDesc(t *testing.T) {
	t.Parallel()
	d := &Desc{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "test_metric",
		Help:      "A test metric",
	}
	desc := d.Desc()
	if desc == nil {
		t.Fatal("Desc() returned nil")
	}
	if d.desc == nil {
		t.Error("Desc() did not cache the result")
	}
	if got := d.Desc(); got != desc {
		t.Error("Desc() returned different pointer on second call")
	}
}

func TestDescWithLabels(t *testing.T) {
	t.Parallel()
	d := &Desc{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "info",
		Help:      "info metric",
		Labels:    []string{"device_id", "version"},
	}
	desc := d.Desc()
	if desc == nil {
		t.Fatal("Desc() returned nil")
	}
	// Verify the desc works with the right number of labels.
	m := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1, "id1", "v1")
	if m == nil {
		t.Fatal("failed to create metric with labels")
	}
}

func TestDescsComplete(t *testing.T) {
	t.Parallel()
	// Verify all descriptors in the Descs slice are valid.
	seen := make(map[string]bool)
	for _, d := range Descs {
		fqName := d.FQName()
		if fqName == "" {
			t.Error("Desc with empty FQName")
		}
		if seen[fqName] {
			t.Errorf("duplicate descriptor: %s", fqName)
		}
		seen[fqName] = true

		if d.Help == "" {
			t.Errorf("descriptor %s has empty help", fqName)
		}
		if d.Desc() == nil {
			t.Errorf("descriptor %s returned nil Desc()", fqName)
		}
	}
}
