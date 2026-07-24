// gRPC-Web transport for the Starlink dish.
//
// Starlink dish firmware 2026.06.x removed the legacy raw-gRPC endpoint on
// 192.168.100.1:9200. The Device API now lives on :9201 and is served as
// gRPC-Web over HTTP/1.1 ONLY (the same transport the dish's local web UI
// uses). The standard google.golang.org/grpc client requires HTTP/2 and so
// fails with "http2: frame too large, ... looked like an HTTP/1.1 header".
//
// This file implements just enough gRPC-Web to satisfy the single unary call
// the exporter makes (Device/Handle). It marshals the existing generated
// *device.Request, frames it as gRPC-Web, POSTs over HTTP/1.1, unframes the
// reply, and proto.Unmarshals into the existing *device.Response. Nothing in
// scrape.go or the generated protobufs changes.

package exporter

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/joshuasing/starlink_exporter/internal/spacex_api/device"
)

const (
	// deviceHandlePath is the gRPC-Web request path for the unary Handle RPC.
	deviceHandlePath = "/SpaceX.API.Device.Device/Handle"

	// grpcWebContentType is the binary (proto) gRPC-Web framing.
	grpcWebContentType = "application/grpc-web+proto"
)

// Default ride-through / cache tuning. The dish's gRPC-Web endpoint goes fully
// unresponsive in ~8.5s windows (all RPCs, periodically). Handle retries with a
// short per-attempt timeout until the caller's context budget (callTimeout in
// scrape.go, ~12s) is spent, which clears nearly all wedges; the cache covers
// the rest. cacheTTL is set above the wedge and a few scrape intervals, but low
// enough that a genuinely dead dish goes stale (-> dish_up 0) reasonably fast.
// These are struct fields (defaulted here) so tests can shrink the timing.
const (
	defaultCacheTTL          = 90 * time.Second
	defaultPerAttemptTimeout = 1500 * time.Millisecond
	defaultRetryGap          = 500 * time.Millisecond
)

type cachedResponse struct {
	resp *device.Response
	at   time.Time
}

// grpcWebClient speaks gRPC-Web over HTTP/1.1 to the dish. It implements the
// subset of device.DeviceClient that the exporter actually uses (Handle); the
// streaming Stream method is intentionally unimplemented.
type grpcWebClient struct {
	baseURL string // e.g. http://192.168.100.1:9201
	http    *http.Client

	// ride-through / cache tuning (defaulted in newGRPCWebClient).
	cacheTTL          time.Duration
	perAttemptTimeout time.Duration
	retryGap          time.Duration

	// lastOK tracks whether the most recent call succeeded, surfaced via
	// ConnState() for the /health endpoint (mirrors the old gRPC conn state).
	lastOK atomic.Bool

	// lastGood caches the most recent successful response per request type
	// (keyed by the Request oneof type name), for the short-cache fallback.
	mu       sync.Mutex
	lastGood map[string]cachedResponse
}

// newGRPCWebClient builds a gRPC-Web client for the given dish address
// (host:port, no scheme). Plain HTTP/1.1; the dish has no TLS on the LAN.
func newGRPCWebClient(address string) *grpcWebClient {
	return &grpcWebClient{
		baseURL:           "http://" + address,
		lastGood:          make(map[string]cachedResponse),
		cacheTTL:          defaultCacheTTL,
		perAttemptTimeout: defaultPerAttemptTimeout,
		retryGap:          defaultRetryGap,
		http: &http.Client{
			// Hard safety net above the per-call context (callTimeout in
			// scrape.go). The context normally bounds each request; this just
			// guarantees no request can ever hang the collector indefinitely.
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				// Honor HTTP_PROXY/NO_PROXY (the dish is on-LAN in production,
				// but this allows reaching it via a proxy when the exporter
				// runs off the dish network, e.g. during testing).
				Proxy: http.ProxyFromEnvironment,
				// Force HTTP/1.1: the dish gRPC-Web endpoint does NOT speak
				// HTTP/2 (h2c). Leaving ForceAttemptHTTP2 default could let
				// the client try h2c on a cleartext POST and get "empty reply".
				ForceAttemptHTTP2: false,
				// Disable HTTP keep-alive: one fresh connection per request.
				// The dish's embedded gRPC-Web server drops idle connections,
				// so a reused-but-dead pooled connection would hang the next
				// POST while "awaiting headers". The dish is on-LAN and scrapes
				// are infrequent (~15s), so per-request connections are cheap.
				DisableKeepAlives: true,
				// Bound connection setup so a network blip fails fast instead
				// of consuming the whole scrape context.
				DialContext: (&net.Dialer{
					Timeout:   3 * time.Second,
					KeepAlive: -1,
				}).DialContext,
				// NOTE: the dish firmware periodically wedges a single
				// get_status request entirely (no response for 15s+), clearing
				// by the next request. We can't make a wedged request return;
				// the per-call context (callTimeout) abandons it quickly and
				// the next scrape recovers. See Handle's one-shot retry.
				ResponseHeaderTimeout: 5 * time.Second,
			},
		},
	}
}

// Handle performs a unary Device/Handle call over gRPC-Web, riding through the
// dish's periodic ~8.5s unresponsive windows.
//
// The dish's gRPC-Web endpoint goes fully unresponsive (every RPC) for ~8.5s at
// a time, periodically. Within the caller's context budget we retry with a
// short per-attempt timeout (perAttemptTimeout) every retryGap, which clears
// nearly all wedges. If the whole budget is exhausted we fall back to the last
// good response for this request type if it's still within cacheTTL — so a blip
// produces stale-but-present metrics rather than a gap. Only when there's no
// fresh-enough cache either do we surface the error (-> dish_up 0).
//
// A successful or non-retryable (valid gRPC status) response returns
// immediately. The web UI hides these windows by polling many times/second;
// this gives the once-per-scrape exporter the same resilience.
func (c *grpcWebClient) Handle(ctx context.Context, in *device.Request) (*device.Response, error) {
	key := requestKey(in)
	var lastErr error

	for {
		attemptCtx, cancel := context.WithTimeout(ctx, c.perAttemptTimeout)
		out, err := c.doHandle(attemptCtx, in)
		cancel()

		if err == nil {
			c.store(key, out)
			c.lastOK.Store(true)
			return out, nil
		}
		lastErr = err

		// A real gRPC error status is deterministic — don't retry or mask it.
		if !isTransient(err) {
			c.lastOK.Store(false)
			return nil, err
		}

		// Out of budget? Stop retrying and try the cache.
		if ctx.Err() != nil {
			break
		}
		select {
		case <-ctx.Done():
		case <-time.After(c.retryGap):
		}
		if ctx.Err() != nil {
			break
		}
	}

	// Ride-through exhausted: serve a fresh-enough cached response if we have one.
	if cached, ok := c.cached(key); ok {
		slog.Warn("dish unresponsive; serving cached response",
			slog.String("request", key), slog.Any("err", lastErr))
		c.lastOK.Store(false) // a cache hit still means the live fetch failed
		return cached, nil
	}

	c.lastOK.Store(false)
	return nil, lastErr
}

// requestKey identifies the Request oneof variant (e.g. "*device.Request_GetStatus")
// so responses are cached per request type.
func requestKey(in *device.Request) string {
	if in == nil || in.Request == nil {
		return "nil"
	}
	return fmt.Sprintf("%T", in.Request)
}

func (c *grpcWebClient) store(key string, resp *device.Response) {
	c.mu.Lock()
	c.lastGood[key] = cachedResponse{resp: resp, at: time.Now()}
	c.mu.Unlock()
}

// cached returns the last good response for key if it's within cacheTTL.
func (c *grpcWebClient) cached(key string) (*device.Response, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.lastGood[key]
	if !ok || time.Since(e.at) > c.cacheTTL {
		return nil, false
	}
	return e.resp, true
}

// grpcStatusErr builds a proper google.golang.org/grpc/status error from the
// gRPC-Web wire status code + message. Returning a real status error (rather
// than a plain fmt.Errorf string) lets callers use status.FromError(err) —
// e.g. scrapeLocation treats Unimplemented/PermissionDenied as "location
// disabled" instead of a hard failure. With the old raw-gRPC transport these
// were already status errors; this preserves that behaviour over gRPC-Web.
func grpcStatusErr(codeStr, msg string) error {
	n, err := strconv.Atoi(strings.TrimSpace(codeStr))
	if err != nil {
		return fmt.Errorf("grpc status %s: %s", codeStr, msg)
	}
	return status.Error(codes.Code(n), msg)
}

// isTransient reports whether an error is worth retrying. Everything wrapped
// with the "grpc-web request:" prefix is transient: transport failures from
// c.http.Do (timeouts, resets, refused connections) and HTTP 5xx/429 from the
// dish's wedge windows. Deterministic protocol errors (other 4xx, non-zero
// grpc-status, unmarshal failures) are not.
func isTransient(err error) bool {
	return strings.HasPrefix(err.Error(), "grpc-web request:")
}

// doHandle performs a single unary Device/Handle call over gRPC-Web.
func (c *grpcWebClient) doHandle(ctx context.Context, in *device.Request) (*device.Response, error) {
	reqBytes, err := proto.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+deviceHandlePath, bytes.NewReader(frameMessage(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", grpcWebContentType)
	httpReq.Header.Set("Accept", grpcWebContentType)
	httpReq.Header.Set("X-Grpc-Web", "1")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("grpc-web request: %w", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		// 5xx / 429 are "temporarily unavailable" — exactly the dish's wedge
		// windows (it returns errors, not just connection resets, while
		// unresponsive). Tag them transient so Handle rides through / caches.
		// Other 4xx are deterministic and not retried.
		if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
			return nil, fmt.Errorf("grpc-web request: HTTP status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("grpc-web HTTP status %d", resp.StatusCode)
	}

	// gRPC status may arrive in headers (trailers-only response) or in a
	// trailer frame appended to the body. Check the header form first.
	if s := resp.Header.Get("Grpc-Status"); s != "" && s != "0" {
		return nil, grpcStatusErr(s, resp.Header.Get("Grpc-Message"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// A read error mid-body is a transport failure; tag it so it retries.
		return nil, fmt.Errorf("grpc-web request: read body: %w", err)
	}

	msg, trailer, err := deframeMessage(body)
	if err != nil {
		return nil, err
	}
	if code := trailer["grpc-status"]; code != "" && code != "0" {
		return nil, grpcStatusErr(code, trailer["grpc-message"])
	}
	if msg == nil {
		return nil, fmt.Errorf("no message frame in gRPC-Web response")
	}

	out := new(device.Response)
	if err := proto.Unmarshal(msg, out); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return out, nil
}

// ConnState reports a synthetic connectivity state for the /health endpoint,
// based on whether the last call succeeded. There is no persistent connection
// in gRPC-Web (one HTTP/1.1 request per call), so this is best-effort.
func (c *grpcWebClient) ConnState() connectivity.State {
	if c.lastOK.Load() {
		return connectivity.Ready
	}
	return connectivity.TransientFailure
}

// frameMessage wraps a protobuf payload in a gRPC-Web data frame:
// 1 flag byte (0x00 = data, uncompressed) + 4-byte big-endian length + payload.
func frameMessage(payload []byte) []byte {
	buf := make([]byte, 5+len(payload))
	buf[0] = 0x00
	binary.BigEndian.PutUint32(buf[1:5], uint32(len(payload)))
	copy(buf[5:], payload)
	return buf
}

// deframeMessage walks the gRPC-Web framed response body. It returns the first
// data frame (the message) and the parsed trailer frame (flag bit 0x80) as a
// lower-cased key/value map. Either may be empty.
func deframeMessage(body []byte) (msg []byte, trailer map[string]string, err error) {
	trailer = map[string]string{}
	for off := 0; off+5 <= len(body); {
		flag := body[off]
		n := binary.BigEndian.Uint32(body[off+1 : off+5])
		off += 5
		if off+int(n) > len(body) {
			return nil, trailer, fmt.Errorf("truncated gRPC-Web frame (need %d, have %d)", n, len(body)-off)
		}
		frame := body[off : off+int(n)]
		off += int(n)

		if flag&0x80 != 0 {
			// Trailer frame: HTTP/1.1-style "key: value" lines.
			for _, line := range strings.Split(string(frame), "\r\n") {
				if line == "" {
					continue
				}
				if i := strings.IndexByte(line, ':'); i >= 0 {
					k := strings.ToLower(strings.TrimSpace(line[:i]))
					v := strings.TrimSpace(line[i+1:])
					trailer[k] = v
				}
			}
			continue
		}
		if msg == nil {
			msg = frame
		}
	}
	return msg, trailer, nil
}
