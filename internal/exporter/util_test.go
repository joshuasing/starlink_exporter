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
)

func TestBtof(t *testing.T) {
	t.Parallel()
	if got := btof(true); got != 1 {
		t.Errorf("btof(true) = %v, want 1", got)
	}
	if got := btof(false); got != 0 {
		t.Errorf("btof(false) = %v, want 0", got)
	}
}

func TestItos(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   int32
		want string
	}{
		{0, "0"},
		{1, "1"},
		{-1, "-1"},
		{42, "42"},
	}
	for _, tt := range tests {
		if got := itos(tt.in); got != tt.want {
			t.Errorf("itos(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestFtos(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{1.5, "1.5"},
		{-3.14, "-3.14"},
		{100, "100"},
	}
	for _, tt := range tests {
		if got := ftos(tt.in); got != tt.want {
			t.Errorf("ftos(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParseRingBuffer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		data    []float32
		current uint64
		want    []float32
	}{
		{
			name:    "empty",
			data:    []float32{},
			current: 0,
			want:    []float32{},
		},
		{
			name:    "not wrapped",
			data:    []float32{1, 2, 3, 0, 0},
			current: 3,
			want:    []float32{1, 2, 3},
		},
		{
			name:    "exactly full",
			data:    []float32{1, 2, 3, 4, 5},
			current: 5,
			want:    []float32{1, 2, 3, 4, 5},
		},
		{
			name:    "wrapped once",
			data:    []float32{6, 7, 3, 4, 5},
			current: 7,
			want:    []float32{3, 4, 5, 6, 7},
		},
		{
			name:    "wrapped multiple times",
			data:    []float32{11, 12, 13, 9, 10},
			current: 13,
			want:    []float32{9, 10, 11, 12, 13},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parseRingBuffer(tt.data, tt.current)
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("index %d: got %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
