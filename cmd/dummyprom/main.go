// dummyprom is a fake Prometheus query_range API server that returns
// synthetic host-metrics-style data (similar to OpenTelemetry Collector's
// hostmetricsreceiver). Useful for demoing Dashyard without a real Prometheus.
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := "9090"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/api/v1/query_range", handleQueryRange)
	http.HandleFunc("/api/v1/label/", handleLabelValues)

	slog.Info("dummy prometheus server starting", "port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

type promResponse struct {
	Status string   `json:"status"`
	Data   promData `json:"data"`
}

type promData struct {
	ResultType string       `json:"resultType"`
	Result     []promResult `json:"result"`
}

type promResult struct {
	Metric map[string]string `json:"metric"`
	Values [][2]interface{}  `json:"values"`
}

func handleQueryRange(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	stepStr := r.URL.Query().Get("step")

	slog.Info("query_range", "query", query, "start", startStr, "end", endStr, "step", stepStr)

	start, _ := strconv.ParseFloat(startStr, 64)
	end, _ := strconv.ParseFloat(endStr, 64)
	step := parseStep(stepStr)
	if step == 0 {
		step = 15
	}

	results := generateData(query, start, end, step)

	resp := promResponse{
		Status: "success",
		Data: promData{
			ResultType: "matrix",
			Result:     results,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode query_range response", "error", err)
	}
}

func parseStep(s string) float64 {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "s") {
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "s"), 64)
		return v
	}
	if strings.HasSuffix(s, "m") {
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "m"), 64)
		return v * 60
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func generateData(query string, start, end, step float64) []promResult {
	switch {
	case strings.Contains(query, "cpu_utilization"):
		return generateCPUUtilization(start, end, step)
	case strings.Contains(query, "cpu_load_average"):
		return generateLoadAverage(start, end, step)
	case strings.Contains(query, "memory_usage"):
		return generateMemoryUsage(start, end, step)
	case strings.Contains(query, "network_io") && strings.Contains(query, "receive"):
		return generateNetworkIO(start, end, step, "receive")
	case strings.Contains(query, "network_io") && strings.Contains(query, "transmit"):
		return generateNetworkIO(start, end, step, "transmit")
	case strings.Contains(query, "disk_io") && strings.Contains(query, "read"):
		return generateDiskIO(start, end, step, "read")
	case strings.Contains(query, "disk_io") && strings.Contains(query, "write"):
		return generateDiskIO(start, end, step, "write")
	default:
		return generateGeneric(query, start, end, step)
	}
}

// noise returns a deterministic pseudo-random value in [0,1) derived from the
// timestamp t and an optional seed. Using a pure function of t means the same
// query with the same time range always produces identical data, which makes
// screenshots stable and easier to review.
func noise(t float64, seed uint64) float64 {
	// Mix timestamp bits with the seed using a simple hash (splitmix64-style).
	x := uint64(t*1000) + seed
	x ^= x >> 30
	x *= 0xbf58476d1ce4e5b9
	x ^= x >> 27
	x *= 0x94d049bb133111eb
	x ^= x >> 31
	return float64(x>>11) / float64(1<<53) // [0, 1)
}

func generateCPUUtilization(start, end, step float64) []promResult {
	cpus := []string{"cpu0", "cpu1", "cpu2", "cpu3"}
	baseVals := []float64{20.0, 25.0, 30.0, 18.0}
	var results []promResult
	for i, cpu := range cpus {
		seed := uint64(i)
		base := baseVals[i]
		values := generateTimeSeries(start, end, step, func(t float64) float64 {
			// Simulate CPU with some periodic pattern + noise
			return math.Max(0, math.Min(100,
				base+15*math.Sin(t/600)+noise(t, seed)*10-5))
		})
		results = append(results, promResult{
			Metric: map[string]string{
				"__name__": "system_cpu_utilization_ratio",
				"cpu":      cpu,
			},
			Values: values,
		})
	}
	return results
}

func generateLoadAverage(start, end, step float64) []promResult {
	values := generateTimeSeries(start, end, step, func(t float64) float64 {
		return math.Max(0, 1.5+0.8*math.Sin(t/1200)+noise(t, 100)*0.3)
	})
	return []promResult{
		{
			Metric: map[string]string{
				"__name__": "system_cpu_load_average_1m_ratio",
			},
			Values: values,
		},
	}
}

func generateMemoryUsage(start, end, step float64) []promResult {
	// Simulate 16 GB total memory split into used/cached/free/buffers.
	// Each state varies over time but the total stays around 16 GB.
	type memState struct {
		name string
		base float64 // base bytes
		seed uint64
	}
	states := []memState{
		{"used", 4.0 * 1024 * 1024 * 1024, 200},
		{"cached", 6.0 * 1024 * 1024 * 1024, 201},
		{"free", 4.5 * 1024 * 1024 * 1024, 202},
		{"buffers", 1.5 * 1024 * 1024 * 1024, 203},
	}
	var results []promResult
	for _, s := range states {
		base := s.base
		seed := s.seed
		values := generateTimeSeries(start, end, step, func(t float64) float64 {
			return math.Max(0, base+512*1024*1024*math.Sin(t/1800)+noise(t, seed)*100*1024*1024)
		})
		results = append(results, promResult{
			Metric: map[string]string{
				"__name__": "system_memory_usage_bytes",
				"state":    s.name,
			},
			Values: values,
		})
	}
	return results
}

func generateNetworkIO(start, end, step float64, direction string) []promResult {
	devices := []string{"eth0", "eth1"}
	baseRates := []float64{1024.0 * 1024 * 8, 1024.0 * 1024 * 12} // 8 MB/s, 12 MB/s
	var results []promResult
	for i, dev := range devices {
		baseRate := baseRates[i]
		seed := uint64(300 + i)
		values := generateTimeSeries(start, end, step, func(t float64) float64 {
			return math.Max(0, baseRate+baseRate*0.3*math.Sin(t/900)+noise(t, seed)*1024*512)
		})
		results = append(results, promResult{
			Metric: map[string]string{
				"__name__":  "system_network_io_bytes_total",
				"device":    dev,
				"direction": direction,
			},
			Values: values,
		})
	}
	return results
}

func generateDiskIO(start, end, step float64, direction string) []promResult {
	baseRate := 1024.0 * 1024 * 50 // 50 MB/s
	values := generateTimeSeries(start, end, step, func(t float64) float64 {
		return math.Max(0, baseRate+baseRate*0.4*math.Sin(t/600)+noise(t, 400)*1024*1024*10)
	})
	return []promResult{
		{
			Metric: map[string]string{
				"__name__":  "system_disk_io_bytes_total",
				"device":    "sda",
				"direction": direction,
			},
			Values: values,
		},
	}
}

func generateGeneric(query string, start, end, step float64) []promResult {
	values := generateTimeSeries(start, end, step, func(t float64) float64 {
		return math.Max(0, 50+30*math.Sin(t/600)+noise(t, 500)*10)
	})
	return []promResult{
		{
			Metric: map[string]string{
				"__name__": query,
			},
			Values: values,
		},
	}
}

func generateTimeSeries(start, end, step float64, fn func(float64) float64) [][2]interface{} {
	var values [][2]interface{}
	for t := start; t <= end; t += step {
		// Pass relative time (offset from start) so data shape is
		// independent of when the query is made.
		values = append(values, [2]interface{}{
			t,
			fmt.Sprintf("%.6f", fn(t-start)),
		})
	}
	return values
}

// labelRegistry maps metric names to their label name -> values.
var labelRegistry = map[string]map[string][]string{
	"system_network_io_bytes_total": {
		"device":    {"eth0", "eth1"},
		"direction": {"receive", "transmit"},
	},
	"system_cpu_utilization_ratio": {
		"cpu": {"cpu0", "cpu1", "cpu2", "cpu3"},
	},
	"system_memory_usage_bytes": {
		"state": {"used", "cached", "free", "buffers"},
	},
	"system_disk_io_bytes_total": {
		"device":    {"sda"},
		"direction": {"read", "write"},
	},
}

type labelValuesResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

func handleLabelValues(w http.ResponseWriter, r *http.Request) {
	// Path: /api/v1/label/{label}/values
	path := r.URL.Path
	const prefix = "/api/v1/label/"
	const suffix = "/values"
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		http.NotFound(w, r)
		return
	}
	label := path[len(prefix) : len(path)-len(suffix)]
	if label == "" {
		http.NotFound(w, r)
		return
	}

	match := r.URL.Query().Get("match[]")

	slog.Info("label_values", "label", label, "match", match)

	values := collectLabelValues(label, match)

	resp := labelValuesResponse{
		Status: "success",
		Data:   values,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode label_values response", "error", err)
	}
}

func collectLabelValues(label, match string) []string {
	seen := map[string]bool{}
	var values []string

	for metric, labels := range labelRegistry {
		if match != "" && metric != match {
			continue
		}
		if vals, ok := labels[label]; ok {
			for _, v := range vals {
				if !seen[v] {
					seen[v] = true
					values = append(values, v)
				}
			}
		}
	}

	return values
}
