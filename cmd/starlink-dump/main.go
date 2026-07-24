// Command starlink-dump captures EVERYTHING the Starlink dish exposes over its
// local gRPC-Web Device API, for troubleshooting and Starlink support cases
// (e.g. the lower_signal_than_predicted alert firing while the Starlink app
// shows a healthy dish).
//
// It issues every known READ-ONLY Device/Handle RPC (get_status, get_history,
// get_diagnostics, dish_get_obstruction_map, get_radio_stats, ...) and writes,
// into one timestamped output directory:
//
//	NN_<rpc>.json        full response decoded to JSON (protojson, all fields)
//	NN_<rpc>.bin         raw protobuf wire bytes (lossless — preserves fields
//	                     newer than our vendored protos; SpaceX can decode)
//	obstruction_map.png  the SNR obstruction map rendered as an image
//	summary.txt          human-readable headline info + ACTIVE ALERTS
//	meta.json            per-RPC outcome: gRPC code, attempts, timing, sha256
//	samples/             optional repeated get_status snapshots (-samples N)
//
// plus a <dir>.tar.gz of the whole directory, ready to attach to a case.
//
// Mutating / disruptive RPCs (reboot, stow, self tests, speed tests, RSSI scan
// *activation*, factory reset, ...) are deliberately NOT issued. RPCs that are
// unimplemented or permission-denied on a given firmware are recorded as such;
// that outcome is itself diagnostic information.
//
// Transport: dish firmware 2026.06.x serves the API as gRPC-Web over HTTP/1.1
// on 192.168.100.1:9201 only (raw gRPC :9200 is gone), and the embedded server
// periodically wedges for ~8.5s. Like the patched exporter (grpcweb.go), each
// RPC retries with a short per-attempt timeout within a total budget — but
// there is NO cache here: every byte written was really read from the dish at
// the recorded time.
//
// Build: compiled from the patched joshuasing/starlink_exporter tree by
// exporter/build.sh (the generated SpaceX protobufs live there). Static
// linux/amd64, std-lib only. Run it ON the location server (the dish is not
// reachable from anywhere else).
package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/joshuasing/starlink_exporter/internal/spacex_api/device"
)

// buildVersion is stamped by build.sh via -ldflags "-X main.buildVersion=...".
var buildVersion = "dev"

const (
	protoBase        = "joshuasing/starlink_exporter v0.9.2 vendored SpaceX protos"
	deviceHandlePath = "/SpaceX.API.Device.Device/Handle"
	grpcWebProto     = "application/grpc-web+proto"
)

// rpcSpec describes one read-only Device/Handle request to capture. Core RPCs
// drive the exit code; optional ones are expected to fail on some firmware
// (Unimplemented / PermissionDenied) and only get recorded.
type rpcSpec struct {
	name string
	core bool
	req  func() *device.Request
}

var rpcs = []rpcSpec{
	{"device_info", true, func() *device.Request {
		return &device.Request{Request: &device.Request_GetDeviceInfo{GetDeviceInfo: &device.GetDeviceInfoRequest{}}}
	}},
	{"status", true, func() *device.Request {
		return &device.Request{Request: &device.Request_GetStatus{GetStatus: &device.GetStatusRequest{}}}
	}},
	{"diagnostics", false, func() *device.Request {
		return &device.Request{Request: &device.Request_GetDiagnostics{GetDiagnostics: &device.GetDiagnosticsRequest{}}}
	}},
	{"history", true, func() *device.Request {
		return &device.Request{Request: &device.Request_GetHistory{GetHistory: &device.GetHistoryRequest{}}}
	}},
	{"obstruction_map", false, func() *device.Request {
		return &device.Request{Request: &device.Request_DishGetObstructionMap{DishGetObstructionMap: &device.DishGetObstructionMapRequest{}}}
	}},
	{"context", false, func() *device.Request {
		return &device.Request{Request: &device.Request_DishGetContext{DishGetContext: &device.DishGetContextRequest{}}}
	}},
	{"config", false, func() *device.Request {
		return &device.Request{Request: &device.Request_DishGetConfig{DishGetConfig: &device.DishGetConfigRequest{}}}
	}},
	{"location", false, func() *device.Request {
		return &device.Request{Request: &device.Request_GetLocation{GetLocation: &device.GetLocationRequest{}}}
	}},
	{"radio_stats", false, func() *device.Request {
		return &device.Request{Request: &device.Request_GetRadioStats{GetRadioStats: &device.GetRadioStatsRequest{}}}
	}},
	{"persistent_stats", false, func() *device.Request {
		return &device.Request{Request: &device.Request_GetPersistentStats{GetPersistentStats: &device.GetPersistentStatsRequest{}}}
	}},
	{"network_interfaces", false, func() *device.Request {
		return &device.Request{Request: &device.Request_GetNetworkInterfaces{GetNetworkInterfaces: &device.GetNetworkInterfacesRequest{}}}
	}},
	{"ping", false, func() *device.Request {
		return &device.Request{Request: &device.Request_GetPing{GetPing: &device.GetPingRequest{}}}
	}},
	{"connections", false, func() *device.Request {
		return &device.Request{Request: &device.Request_GetConnections{GetConnections: &device.GetConnectionsRequest{}}}
	}},
	{"time", false, func() *device.Request {
		return &device.Request{Request: &device.Request_Time{Time: &device.GetTimeRequest{}}}
	}},
	{"gnss_measurement", false, func() *device.Request {
		return &device.Request{Request: &device.Request_GetGnssMeasurement{GetGnssMeasurement: &device.GetGnssMeasurementRequest{}}}
	}},
	{"rssi_scan_result", false, func() *device.Request {
		return &device.Request{Request: &device.Request_DishGetRssiScanResult{DishGetRssiScanResult: &device.DishGetRssiScanResultRequest{}}}
	}},
	{"emc", false, func() *device.Request {
		return &device.Request{Request: &device.Request_DishGetEmc{DishGetEmc: &device.DishGetEmcRequest{}}}
	}},
	{"transceiver_status", false, func() *device.Request {
		return &device.Request{Request: &device.Request_TransceiverGetStatus{TransceiverGetStatus: &device.TransceiverGetStatusRequest{}}}
	}},
	{"transceiver_telemetry", false, func() *device.Request {
		return &device.Request{Request: &device.Request_TransceiverGetTelemetry{TransceiverGetTelemetry: &device.TransceiverGetTelemetryRequest{}}}
	}},
	{"log", false, func() *device.Request {
		return &device.Request{Request: &device.Request_GetLog{GetLog: &device.GetLogRequest{}}}
	}},
}

// rpcResult is the per-RPC record written to meta.json.
type rpcResult struct {
	Name        string `json:"name"`
	Core        bool   `json:"core"`
	Request     string `json:"request"`
	At          string `json:"at"`
	OK          bool   `json:"ok"`
	GrpcCode    string `json:"grpc_code,omitempty"`
	Error       string `json:"error,omitempty"`
	Attempts    int    `json:"attempts"`
	DurationMs  int64  `json:"duration_ms"`
	ResponseLen int    `json:"response_bytes,omitempty"`
	Sha256      string `json:"sha256,omitempty"`
	JSONFile    string `json:"json_file,omitempty"`
	BinFile     string `json:"bin_file,omitempty"`
}

type metaDoc struct {
	Tool       string      `json:"tool"`
	Version    string      `json:"version"`
	ProtoBase  string      `json:"proto_base"`
	Hostname   string      `json:"hostname"`
	Dish       string      `json:"dish"`
	StartedAt  string      `json:"started_at"`
	FinishedAt string      `json:"finished_at"`
	Results    []rpcResult `json:"results"`
}

// dishClient is a minimal gRPC-Web (HTTP/1.1) client for Device/Handle,
// mirroring the transport quirks handled in the exporter patch (grpcweb.go):
// HTTP/1.1 only, no keep-alive, short attempt timeout + retry inside a budget
// to ride through the dish's periodic ~8.5s unresponsive windows. No caching.
type dishClient struct {
	baseURL        string
	http           *http.Client
	budget         time.Duration
	attemptTimeout time.Duration
	retryGap       time.Duration
}

func newDishClient(address string, budget, attemptTimeout, retryGap time.Duration) *dishClient {
	return &dishClient{
		baseURL:        "http://" + address,
		budget:         budget,
		attemptTimeout: attemptTimeout,
		retryGap:       retryGap,
		http: &http.Client{
			Transport: &http.Transport{
				Proxy:             http.ProxyFromEnvironment,
				ForceAttemptHTTP2: false, // dish is HTTP/1.1 only
				DisableKeepAlives: true,  // dish drops idle conns; fresh conn per request
				DialContext: (&net.Dialer{
					Timeout:   3 * time.Second,
					KeepAlive: -1,
				}).DialContext,
			},
		},
	}
}

// transportError marks retryable failures (connect/timeout/5xx/429 — the
// dish's wedge windows). gRPC status errors are deterministic and final.
type transportError struct{ err error }

func (t *transportError) Error() string { return t.err.Error() }
func (t *transportError) Unwrap() error { return t.err }

func isTransient(err error) bool {
	var te *transportError
	return errors.As(err, &te)
}

// call runs one RPC with ride-through retries. Returns the raw Response wire
// bytes, the decoded Response, and the number of attempts made.
func (c *dishClient) call(req *device.Request) (raw []byte, resp *device.Response, attempts int, err error) {
	deadline := time.Now().Add(c.budget)
	for {
		attempts++
		ctx, cancel := context.WithTimeout(context.Background(), c.attemptTimeout)
		raw, err = c.post(ctx, req)
		cancel()

		if err == nil {
			resp = new(device.Response)
			if uerr := proto.Unmarshal(raw, resp); uerr != nil {
				return raw, nil, attempts, fmt.Errorf("unmarshal response: %w", uerr)
			}
			return raw, resp, attempts, nil
		}
		if !isTransient(err) || time.Now().Add(c.retryGap).After(deadline) {
			return nil, nil, attempts, err
		}
		time.Sleep(c.retryGap)
	}
}

// post performs a single gRPC-Web POST and returns the response message bytes.
func (c *dishClient) post(ctx context.Context, req *device.Request) ([]byte, error) {
	payload, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+deviceHandlePath, bytes.NewReader(frameMessage(payload)))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", grpcWebProto)
	httpReq.Header.Set("Accept", grpcWebProto)
	httpReq.Header.Set("X-Grpc-Web", "1")

	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, &transportError{fmt.Errorf("grpc-web request: %w", err)}
	}
	defer func() { _, _ = io.Copy(io.Discard, httpResp.Body); _ = httpResp.Body.Close() }()

	if httpResp.StatusCode != http.StatusOK {
		if httpResp.StatusCode >= 500 || httpResp.StatusCode == http.StatusTooManyRequests {
			return nil, &transportError{fmt.Errorf("HTTP status %d", httpResp.StatusCode)}
		}
		return nil, fmt.Errorf("HTTP status %d", httpResp.StatusCode)
	}
	if s := httpResp.Header.Get("Grpc-Status"); s != "" && s != "0" {
		return nil, grpcStatusErr(s, httpResp.Header.Get("Grpc-Message"))
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, &transportError{fmt.Errorf("read body: %w", err)}
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
	return msg, nil
}

func grpcStatusErr(codeStr, msg string) error {
	n, err := strconv.Atoi(strings.TrimSpace(codeStr))
	if err != nil {
		return fmt.Errorf("grpc status %s: %s", codeStr, msg)
	}
	return status.Error(codes.Code(n), msg)
}

// frameMessage wraps a protobuf payload in a gRPC-Web data frame.
func frameMessage(payload []byte) []byte {
	buf := make([]byte, 5+len(payload))
	buf[0] = 0x00
	binary.BigEndian.PutUint32(buf[1:5], uint32(len(payload)))
	copy(buf[5:], payload)
	return buf
}

// deframeMessage walks the framed body: returns the first data frame and the
// trailer frame (flag 0x80) parsed into a lower-cased map.
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
			for _, line := range strings.Split(string(frame), "\r\n") {
				if i := strings.IndexByte(line, ':'); i >= 0 {
					trailer[strings.ToLower(strings.TrimSpace(line[:i]))] = strings.TrimSpace(line[i+1:])
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

var (
	jsonFull = protojson.MarshalOptions{Multiline: true, Indent: "  ", EmitUnpopulated: true}
	jsonLine = protojson.MarshalOptions{}
)

func main() {
	var (
		dish           = flag.String("dish", "192.168.100.1:9201", "dish gRPC-Web address (host:port)")
		outDir         = flag.String("out", "", "output directory (default ./starlink-dump_<host>_<utc-timestamp>)")
		budget         = flag.Duration("timeout", 20*time.Second, "total per-RPC budget incl. retries (rides through the dish's ~8.5s wedge windows)")
		attemptTimeout = flag.Duration("attempt-timeout", 4*time.Second, "single attempt timeout (raise if 'log' truncates)")
		retryGap       = flag.Duration("retry-gap", 500*time.Millisecond, "pause between retries")
		samples        = flag.Int("samples", 1, "number of get_status snapshots (extra ones land in samples/, spaced by -sample-interval)")
		sampleInterval = flag.Duration("sample-interval", 2*time.Second, "interval between extra get_status samples")
		only           = flag.String("only", "", "comma-separated RPC names to run (default all)")
		skip           = flag.String("skip", "", "comma-separated RPC names to skip")
		noTar          = flag.Bool("no-tar", false, "do not create <out>.tar.gz")
		noPng          = flag.Bool("no-png", false, "do not render obstruction_map.png")
		showVersion    = flag.Bool("version", false, "print version and exit")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `starlink-dump %s — full read-only dump of the Starlink dish gRPC-Web API
(for troubleshooting / Starlink support cases; run ON the location server)

Usage: starlink-dump [flags]

RPCs captured: `, buildVersion)
		names := make([]string, len(rpcs))
		for i, r := range rpcs {
			names[i] = r.name
		}
		fmt.Fprintf(os.Stderr, "%s\n\nFlags:\n", strings.Join(names, ", "))
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExit codes: 0 = all core RPCs (device_info,status,history) OK; 1 = some core failed; 2 = dish unreachable (no RPC succeeded)\n")
	}
	flag.Parse()

	if *showVersion {
		fmt.Printf("starlink-dump %s (%s)\n", buildVersion, protoBase)
		return
	}

	hostname, _ := os.Hostname()
	start := time.Now().UTC()

	dir := *outDir
	if dir == "" {
		dir = fmt.Sprintf("starlink-dump_%s_%s", hostname, start.Format("20060102-150405Z"))
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fatal("create output dir: %v", err)
	}

	onlySet := nameSet(*only)
	skipSet := nameSet(*skip)

	client := newDishClient(*dish, *budget, *attemptTimeout, *retryGap)
	meta := &metaDoc{
		Tool:      "starlink-dump",
		Version:   buildVersion,
		ProtoBase: protoBase,
		Hostname:  hostname,
		Dish:      *dish,
		StartedAt: start.Format(time.RFC3339),
	}
	responses := map[string]*device.Response{}

	fmt.Fprintf(os.Stderr, "starlink-dump %s | dish %s | host %s | out %s\n",
		buildVersion, *dish, hostname, dir)

	anyOK := false
	for i, spec := range rpcs {
		if len(onlySet) > 0 && !onlySet[spec.name] {
			continue
		}
		if skipSet[spec.name] {
			continue
		}
		res := runRPC(client, dir, fmt.Sprintf("%02d_%s", i+1, spec.name), spec)
		meta.Results = append(meta.Results, res)
		if res.OK {
			anyOK = true
			responses[spec.name] = lastResp
		}
	}

	// Extra get_status samples to catch flapping (e.g. an alert that is set in
	// one poll and clear the next while the app shows "everything OK").
	if *samples > 1 && (len(onlySet) == 0 || onlySet["status"]) && !skipSet["status"] {
		samplesDir := filepath.Join(dir, "samples")
		if err := os.MkdirAll(samplesDir, 0o755); err != nil {
			fatal("create samples dir: %v", err)
		}
		var statusSpec rpcSpec
		for _, spec := range rpcs {
			if spec.name == "status" {
				statusSpec = rpcSpec{name: "status", core: false, req: spec.req}
				break
			}
		}
		for s := 2; s <= *samples; s++ {
			time.Sleep(*sampleInterval)
			res := runRPC(client, samplesDir, fmt.Sprintf("status_s%02d", s), statusSpec)
			res.Name = fmt.Sprintf("status(sample %d/%d)", s, *samples)
			meta.Results = append(meta.Results, res)
		}
	}

	// Obstruction map PNG.
	if !*noPng {
		if om := responses["obstruction_map"].GetDishGetObstructionMap(); om != nil && len(om.GetSnr()) > 0 {
			if err := writeObstructionPNG(filepath.Join(dir, "obstruction_map.png"), om); err != nil {
				fmt.Fprintf(os.Stderr, "  ! obstruction_map.png: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "  wrote obstruction_map.png (%dx%d)\n", om.GetNumCols(), om.GetNumRows())
			}
		}
	}

	meta.FinishedAt = time.Now().UTC().Format(time.RFC3339)

	// summary.txt + meta.json
	if err := os.WriteFile(filepath.Join(dir, "summary.txt"),
		[]byte(buildSummary(meta, responses)), 0o644); err != nil {
		fatal("write summary.txt: %v", err)
	}
	mj, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, "meta.json"), append(mj, '\n'), 0o644); err != nil {
		fatal("write meta.json: %v", err)
	}

	if !*noTar {
		tgz := dir + ".tar.gz"
		if err := makeTarGz(dir, tgz); err != nil {
			fmt.Fprintf(os.Stderr, "  ! tar.gz: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "\narchive: %s\n", tgz)
		}
	}

	coreFailed := 0
	for _, r := range meta.Results {
		if r.Core && !r.OK {
			coreFailed++
		}
	}
	fmt.Fprintf(os.Stderr, "done: %s (core failures: %d)\n", dir, coreFailed)
	switch {
	case !anyOK:
		os.Exit(2)
	case coreFailed > 0:
		os.Exit(1)
	}
}

// lastResp holds the decoded response of the most recent successful runRPC
// call (single-threaded main loop; avoids re-decoding for summary/PNG).
var lastResp *device.Response

// runRPC executes one RPC and writes <prefix>.json/.bin into dir.
func runRPC(client *dishClient, dir, prefix string, spec rpcSpec) rpcResult {
	req := spec.req()
	reqJSON, _ := jsonLine.Marshal(req)
	res := rpcResult{
		Name:    spec.name,
		Core:    spec.core,
		Request: string(reqJSON),
		At:      time.Now().UTC().Format(time.RFC3339),
	}

	fmt.Fprintf(os.Stderr, "  %-22s ", spec.name)
	t0 := time.Now()
	raw, resp, attempts, err := client.call(req)
	res.Attempts = attempts
	res.DurationMs = time.Since(t0).Milliseconds()

	if err != nil {
		res.Error = err.Error()
		if s, ok := status.FromError(err); ok && !isTransient(err) {
			res.GrpcCode = s.Code().String()
		}
		fmt.Fprintf(os.Stderr, "FAILED after %d attempt(s) in %dms: %v\n", attempts, res.DurationMs, err)
		return res
	}

	res.OK = true
	res.ResponseLen = len(raw)
	sum := sha256.Sum256(raw)
	res.Sha256 = hex.EncodeToString(sum[:])
	lastResp = resp

	binFile := prefix + ".bin"
	if werr := os.WriteFile(filepath.Join(dir, binFile), raw, 0o644); werr != nil {
		fatal("write %s: %v", binFile, werr)
	}
	res.BinFile = binFile

	jsonFile := prefix + ".json"
	jb, jerr := jsonFull.Marshal(resp)
	if jerr != nil {
		// Keep the .bin — the evidence survives even if protojson chokes.
		res.Error = fmt.Sprintf("protojson: %v (raw .bin still written)", jerr)
		fmt.Fprintf(os.Stderr, "OK (bin only, json failed: %v)\n", jerr)
		return res
	}
	if werr := os.WriteFile(filepath.Join(dir, jsonFile), append(jb, '\n'), 0o644); werr != nil {
		fatal("write %s: %v", jsonFile, werr)
	}
	res.JSONFile = jsonFile

	fmt.Fprintf(os.Stderr, "OK  %6s  %d attempt(s)  %dms\n", byteSize(len(raw)), attempts, res.DurationMs)
	return res
}

// buildSummary renders the human-readable summary.txt.
func buildSummary(meta *metaDoc, responses map[string]*device.Response) string {
	var b strings.Builder
	w := func(format string, a ...any) { fmt.Fprintf(&b, format+"\n", a...) }

	w("starlink-dump %s (%s)", meta.Version, meta.ProtoBase)
	w("host: %s   dish: %s", meta.Hostname, meta.Dish)
	w("started: %s   finished: %s", meta.StartedAt, meta.FinishedAt)

	if di := responses["device_info"].GetGetDeviceInfo().GetDeviceInfo(); di != nil {
		w("")
		w("== device ==")
		w("id:                %s", di.GetId())
		w("hardware:          %s (board rev %d, generation %d)", di.GetHardwareVersion(), di.GetBoardRev(), di.GetGenerationNumber())
		w("software:          %s", di.GetSoftwareVersion())
		w("build id:          %s", di.GetBuildId())
		w("country:           %s   bootcount: %d   is_dev: %v", di.GetCountryCode(), di.GetBootcount(), di.GetIsDev())
	}

	if st := responses["status"].GetDishGetStatus(); st != nil {
		w("")
		w("== status ==")
		w("uptime_s:          %d", st.GetDeviceState().GetUptimeS())
		w("sw update state:   %s   reboot_reason: %s", st.GetSoftwareUpdateState(), st.GetRebootReason())
		w("boresight:         az %.2f deg / el %.2f deg", st.GetBoresightAzimuthDeg(), st.GetBoresightElevationDeg())
		w("throughput:        down %.0f bps / up %.0f bps", st.GetDownlinkThroughputBps(), st.GetUplinkThroughputBps())
		w("pop ping:          drop %.4f / latency %.1f ms", st.GetPopPingDropRate(), st.GetPopPingLatencyMs())
		w("SNR flags:         is_snr_above_noise_floor=%v  is_snr_persistently_low=%v",
			st.GetIsSnrAboveNoiseFloor(), st.GetIsSnrPersistentlyLow())
		w("eth speed:         %d Mbps   mobility_class: %s   class_of_service: %s",
			st.GetEthSpeedMbps(), st.GetMobilityClass(), st.GetClassOfService())
		if obs := st.GetObstructionStats(); obs != nil {
			w("obstruction:       %s", compactJSON(obs))
		}
		if as := st.GetAlignmentStats(); as != nil {
			w("alignment:         %s", compactJSON(as))
		}
		if gs := st.GetGpsStats(); gs != nil {
			w("gps:               %s", compactJSON(gs))
		}
		if o := st.GetOutage(); o != nil {
			w("current outage:    %s", compactJSON(o))
		}

		w("")
		w("== ACTIVE ALERTS ==")
		alerts := activeAlerts(st.GetAlerts())
		if len(alerts) == 0 {
			w("(none)")
		} else {
			for _, a := range alerts {
				w("%s", a)
			}
		}
	}

	w("")
	w("== RPC results ==")
	for _, r := range meta.Results {
		state := "OK"
		if !r.OK {
			state = "FAILED"
			if r.GrpcCode != "" {
				state = r.GrpcCode
			}
		}
		w("%-24s %-16s %8s  %d attempt(s)  %dms", r.Name, state, byteSize(r.ResponseLen), r.Attempts, r.DurationMs)
	}
	w("")
	w("Raw wire bytes of every response are in the *.bin files (protobuf,")
	w("lossless even for fields newer than this tool's protos); decoded JSON in")
	w("*.json; per-RPC gRPC codes/timing/sha256 in meta.json.")
	return b.String()
}

// activeAlerts lists the names of all boolean alert fields that are true,
// via protobuf reflection (so new alert fields show up without code changes).
func activeAlerts(alerts proto.Message) []string {
	var out []string
	if alerts == nil {
		return out
	}
	alerts.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.Kind() == protoreflect.BoolKind && v.Bool() {
			out = append(out, string(fd.Name()))
		}
		return true
	})
	sort.Strings(out)
	return out
}

func compactJSON(m proto.Message) string {
	b, err := jsonLine.Marshal(m)
	if err != nil {
		return fmt.Sprintf("<%v>", err)
	}
	return string(b)
}

// writeObstructionPNG renders the SNR obstruction map: unknown cells (-1)
// transparent, otherwise red (0.0) -> yellow -> green (1.0), scaled 4x.
func writeObstructionPNG(path string, om *device.DishGetObstructionMapResponse) error {
	rows, cols := int(om.GetNumRows()), int(om.GetNumCols())
	snr := om.GetSnr()
	if rows*cols != len(snr) {
		return fmt.Errorf("map size mismatch: %dx%d != %d samples", rows, cols, len(snr))
	}
	const scale = 4
	img := image.NewNRGBA(image.Rect(0, 0, cols*scale, rows*scale))
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			v := snr[r*cols+c]
			var col color.NRGBA
			switch {
			case v < 0: // no data
				col = color.NRGBA{0, 0, 0, 0}
			case v < 0.5: // red -> yellow
				col = color.NRGBA{255, uint8(255 * (v * 2)), 0, 255}
			default: // yellow -> green
				col = color.NRGBA{uint8(255 * (2 - v*2)), 255, 0, 255}
			}
			for y := 0; y < scale; y++ {
				for x := 0; x < scale; x++ {
					img.SetNRGBA(c*scale+x, r*scale+y, col)
				}
			}
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// makeTarGz packs dir into tgzPath (entries prefixed with the dir base name).
func makeTarGz(dir, tgzPath string) error {
	f, err := os.Create(tgzPath)
	if err != nil {
		return err
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()

	base := filepath.Base(dir)
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(filepath.Join(base, rel))
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = io.Copy(tw, src)
		return err
	})
}

func nameSet(csv string) map[string]bool {
	set := map[string]bool{}
	for _, n := range strings.Split(csv, ",") {
		if n = strings.TrimSpace(n); n != "" {
			set[n] = true
		}
	}
	return set
}

func byteSize(n int) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1fMB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1fKB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%dB", n)
	}
}

func fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "FATAL: "+format+"\n", a...)
	os.Exit(3)
}
