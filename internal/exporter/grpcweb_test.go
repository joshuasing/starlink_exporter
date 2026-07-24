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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/joshuasing/starlink_exporter/internal/spacex_api/device"
)

// wedgeServer is a fake dish gRPC-Web endpoint that can be flipped "wedged"
// (returns 503, like the dish's ~8.5s unresponsive windows) or healthy
// (returns a valid framed Response with a known uptime).
type wedgeServer struct {
	wedged atomic.Bool
	hits   atomic.Int64
	uptime uint64

	// grpcStatus, when non-empty, makes the handler return a trailers-only
	// gRPC-Web error with this status code (e.g. "12" = Unimplemented).
	grpcStatus  string
	grpcMessage string

	// failLocation, when true, returns Unimplemented for get_location requests
	// while still serving get_status/get_history normally (mimics a dish whose
	// firmware dropped GPS).
	failLocation bool
}

func (w *wedgeServer) handler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		w.hits.Add(1)
		if w.wedged.Load() {
			http.Error(rw, "unavailable", http.StatusServiceUnavailable)
			return
		}
		if w.grpcStatus != "" {
			rw.Header().Set("Content-Type", grpcWebContentType)
			rw.Header().Set("Grpc-Status", w.grpcStatus)
			rw.Header().Set("Grpc-Message", w.grpcMessage)
			rw.WriteHeader(http.StatusOK)
			return
		}

		body, _ := io.ReadAll(r.Body)
		req := decodeRequest(body)

		// Optionally fail only get_location (field 1017).
		if w.failLocation && req.GetGetLocation() != nil {
			rw.Header().Set("Content-Type", grpcWebContentType)
			rw.Header().Set("Grpc-Status", "12") // Unimplemented
			rw.Header().Set("Grpc-Message", "get_location not supported")
			rw.WriteHeader(http.StatusOK)
			return
		}

		resp := w.responseFor(req)
		out, err := proto.Marshal(resp)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", grpcWebContentType)
		rw.Header().Set("Grpc-Status", "0")
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(frameMessage(out))
	}
}

// decodeRequest unframes a gRPC-Web request body and unmarshals the Request.
func decodeRequest(framed []byte) *device.Request {
	req := &device.Request{}
	if len(framed) >= 5 {
		_ = proto.Unmarshal(framed[5:], req) // skip the 5-byte gRPC-Web frame header
	}
	return req
}

// responseFor returns a minimal valid Response matching the request oneof.
func (w *wedgeServer) responseFor(req *device.Request) *device.Response {
	switch {
	case req.GetGetHistory() != nil:
		return &device.Response{Response: &device.Response_DishGetHistory{
			DishGetHistory: &device.DishGetHistoryResponse{},
		}}
	case req.GetGetLocation() != nil:
		return &device.Response{Response: &device.Response_GetLocation{
			GetLocation: &device.GetLocationResponse{},
		}}
	default: // get_status
		return &device.Response{Response: &device.Response_DishGetStatus{
			DishGetStatus: &device.DishGetStatusResponse{
				DeviceState: &device.DeviceState{UptimeS: w.uptime},
			},
		}}
	}
}

// newTestClient builds a grpcWebClient pointed at srv with fast timing so the
// ride-through loop runs quickly under test.
func newTestClient(t *testing.T, srv *httptest.Server) *grpcWebClient {
	t.Helper()
	c := newGRPCWebClient(strings.TrimPrefix(srv.URL, "http://"))
	c.perAttemptTimeout = 100 * time.Millisecond
	c.retryGap = 20 * time.Millisecond
	c.cacheTTL = 2 * time.Second
	c.http = srv.Client() // no real proxy/dialer; talk straight to httptest
	return c
}

func getStatusReq() *device.Request {
	return &device.Request{Request: new(device.Request_GetStatus)}
}

// TestHandle_RidesThroughTransientWedge: a wedge that clears within the context
// budget should be ridden through — Handle returns the fresh response, no error.
func TestHandle_RidesThroughTransientWedge(t *testing.T) {
	ws := &wedgeServer{uptime: 111}
	ws.wedged.Store(true)
	srv := httptest.NewServer(ws.handler())
	defer srv.Close()
	c := newTestClient(t, srv)

	// Clear the wedge after 300ms (several retries in).
	go func() { time.Sleep(300 * time.Millisecond); ws.wedged.Store(false) }()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := c.Handle(ctx, getStatusReq())
	if err != nil {
		t.Fatalf("expected ride-through success, got error: %v", err)
	}
	if got := resp.GetDishGetStatus().GetDeviceState().GetUptimeS(); got != 111 {
		t.Fatalf("uptime = %d, want 111", got)
	}
	if ws.hits.Load() < 2 {
		t.Fatalf("expected multiple attempts (retries), got %d", ws.hits.Load())
	}
}

// TestHandle_CacheFallbackWhenWedgeOutlastsBudget: after one good call primes
// the cache, a wedge that lasts the whole budget should fall back to the cached
// response rather than erroring.
func TestHandle_CacheFallbackWhenWedgeOutlastsBudget(t *testing.T) {
	ws := &wedgeServer{uptime: 222}
	srv := httptest.NewServer(ws.handler())
	defer srv.Close()
	c := newTestClient(t, srv)

	// Prime the cache with a good call.
	if _, err := c.Handle(context.Background(), getStatusReq()); err != nil {
		t.Fatalf("priming call failed: %v", err)
	}

	// Now wedge for longer than the context budget.
	ws.wedged.Store(true)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	resp, err := c.Handle(ctx, getStatusReq())
	if err != nil {
		t.Fatalf("expected cache fallback (nil error), got: %v", err)
	}
	if got := resp.GetDishGetStatus().GetDeviceState().GetUptimeS(); got != 222 {
		t.Fatalf("cached uptime = %d, want 222", got)
	}
	if c.lastOK.Load() {
		t.Fatalf("lastOK should be false on a cache-served (live-failed) call")
	}
}

// TestHandle_ErrorWhenWedgedAndNoCache: wedged with no usable cache must surface
// an error (so dish_up goes 0).
func TestHandle_ErrorWhenWedgedAndNoCache(t *testing.T) {
	ws := &wedgeServer{uptime: 333}
	ws.wedged.Store(true)
	srv := httptest.NewServer(ws.handler())
	defer srv.Close()
	c := newTestClient(t, srv)

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	if _, err := c.Handle(ctx, getStatusReq()); err == nil {
		t.Fatal("expected error when wedged with no cache, got nil")
	}
}

// TestHandle_ExpiredCacheNotServed: a cache older than cacheTTL must not be
// served; the call should error instead.
func TestHandle_ExpiredCacheNotServed(t *testing.T) {
	ws := &wedgeServer{uptime: 444}
	srv := httptest.NewServer(ws.handler())
	defer srv.Close()
	c := newTestClient(t, srv)
	c.cacheTTL = 50 * time.Millisecond

	if _, err := c.Handle(context.Background(), getStatusReq()); err != nil {
		t.Fatalf("priming call failed: %v", err)
	}
	ws.wedged.Store(true)
	time.Sleep(80 * time.Millisecond) // let the cache entry expire

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	if _, err := c.Handle(ctx, getStatusReq()); err == nil {
		t.Fatal("expected error with expired cache, got nil (stale served)")
	}
}

// TestHandle_GrpcStatusIsParseable: a non-zero gRPC-Web status (e.g.
// Unimplemented=12, returned by dishes whose firmware dropped GPS for
// get_location) must surface as a real google.golang.org/grpc/status error so
// scrapeLocation's status.FromError() can treat it as "location disabled"
// instead of failing the whole scrape (dish_up=0). This is the regression the
// gRPC-Web transport introduced over raw gRPC: a plain fmt.Errorf string is not
// parseable by status.FromError().
func TestHandle_GrpcStatusIsParseable(t *testing.T) {
	ws := &wedgeServer{grpcStatus: "12", grpcMessage: "get_location unimplemented"}
	srv := httptest.NewServer(ws.handler())
	defer srv.Close()
	c := newTestClient(t, srv)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := c.Handle(ctx, &device.Request{Request: new(device.Request_GetLocation)})
	if err == nil {
		t.Fatal("expected an error for grpc-status 12, got nil")
	}
	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error is not a grpc status error (status.FromError failed): %v", err)
	}
	if s.Code() != codes.Unimplemented {
		t.Fatalf("code = %v, want Unimplemented", s.Code())
	}
	if isTransient(err) {
		t.Fatal("grpc-status error should not be transient")
	}
}

// newTestExporterWithClient builds an Exporter wired to a test gRPC-Web client,
// with the self-metric gauges initialised so scrape() can run.
func newTestExporterWithClient(c *grpcWebClient) *Exporter {
	return &Exporter{
		client:                c,
		up:                    prometheus.NewGauge(prometheus.GaugeOpts{Name: "t_up"}),
		totalScrapes:          prometheus.NewCounter(prometheus.CounterOpts{Name: "t_scrapes"}),
		scrapeDurationSeconds: prometheus.NewGauge(prometheus.GaugeOpts{Name: "t_dur"}),
	}
}

// TestScrape_LocationFailureDoesNotMarkDown is the core regression guard for the
// "9 dishes offline" incident: a dish that serves status+history but whose
// firmware dropped get_location (returns Unimplemented) must still report up=1.
// Location is optional telemetry and must never drive dish_up.
func TestScrape_LocationFailureDoesNotMarkDown(t *testing.T) {
	ws := &wedgeServer{uptime: 123, failLocation: true}
	srv := httptest.NewServer(ws.handler())
	defer srv.Close()
	ex := newTestExporterWithClient(newTestClient(t, srv))

	ch := make(chan prometheus.Metric, 500)
	up := ex.scrape(ch)
	close(ch)

	if !up {
		t.Fatal("scrape() returned up=false when only get_location failed; location must not affect dish_up")
	}
	// Sanity: location metrics absent, status metric present.
	var sawStatus, sawLocation bool
	for m := range ch {
		d := m.Desc().String()
		if strings.Contains(d, "starlink_dish_uptime_seconds") {
			sawStatus = true
		}
		if strings.Contains(d, "starlink_dish_location") {
			sawLocation = true
		}
	}
	if !sawStatus {
		t.Fatal("expected status metrics to be emitted")
	}
	if sawLocation {
		t.Fatal("did not expect location metrics when get_location failed")
	}
}

// TestScrape_AllGoodIsUp: when status, history and location all succeed, up=1.
func TestScrape_AllGoodIsUp(t *testing.T) {
	ws := &wedgeServer{uptime: 123}
	srv := httptest.NewServer(ws.handler())
	defer srv.Close()
	ex := newTestExporterWithClient(newTestClient(t, srv))

	ch := make(chan prometheus.Metric, 500)
	up := ex.scrape(ch)
	close(ch)
	for range ch {
	}
	if !up {
		t.Fatal("scrape() returned up=false when all RPCs succeeded")
	}
}
