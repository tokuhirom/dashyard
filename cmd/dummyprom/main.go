// dummyprom is a fake Prometheus query_range API server that returns
// synthetic host-metrics-style data (similar to OpenTelemetry Collector's
// hostmetricsreceiver). Useful for demoing Dashyard without a real Prometheus.
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
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
	json.NewEncoder(w).Encode(resp)
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

func generateCPUUtilization(start, end, step float64) []promResult {
	cpus := []string{"cpu0", "cpu1", "cpu2", "cpu3"}
	var results []promResult
	for _, cpu := range cpus {
		baseVal := 15.0 + rand.Float64()*20
		values := generateTimeSeries(start, end, step, func(t float64) float64 {
			// Simulate CPU with some periodic pattern + noise
			return math.Max(0, math.Min(100,
				baseVal+15*math.Sin(t/600)+rand.Float64()*10-5))
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
		return math.Max(0, 1.5+0.8*math.Sin(t/1200)+rand.Float64()*0.3)
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
	baseBytes := 4.0 * 1024 * 1024 * 1024 // 4 GB base
	values := generateTimeSeries(start, end, step, func(t float64) float64 {
		return baseBytes + 512*1024*1024*math.Sin(t/1800) + rand.Float64()*100*1024*1024
	})
	return []promResult{
		{
			Metric: map[string]string{
				"__name__": "system_memory_usage_bytes",
				"state":    "used",
			},
			Values: values,
		},
	}
}

func generateNetworkIO(start, end, step float64, direction string) []promResult {
	devices := []string{"eth0", "eth1"}
	var results []promResult
	for _, dev := range devices {
		baseRate := 1024.0 * 1024 * (5 + rand.Float64()*10) // 5-15 MB/s
		values := generateTimeSeries(start, end, step, func(t float64) float64 {
			return math.Max(0, baseRate+baseRate*0.3*math.Sin(t/900)+rand.Float64()*1024*512)
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
		return math.Max(0, baseRate+baseRate*0.4*math.Sin(t/600)+rand.Float64()*1024*1024*10)
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
		return math.Max(0, 50+30*math.Sin(t/600)+rand.Float64()*10)
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
		values = append(values, [2]interface{}{
			t,
			fmt.Sprintf("%.6f", fn(t)),
		})
	}
	return values
}
