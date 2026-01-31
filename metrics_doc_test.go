package main

import (
	"strings"
	"testing"

	"github.com/tokuhirom/dashyard/internal/prometheus"
)

func TestGenerateMetricsDoc(t *testing.T) {
	metrics := []prometheus.MetricInfo{
		{
			Name:   "go_gc_duration_seconds",
			Type:   "summary",
			Help:   "A summary of the pause duration of garbage collection cycles.",
			Labels: []string{"quantile"},
		},
		{
			Name:   "node_cpu_seconds_total",
			Type:   "counter",
			Help:   "Seconds the CPUs spent in each mode.",
			Labels: []string{"cpu", "mode"},
		},
		{
			Name:   "node_memory_MemTotal_bytes",
			Type:   "gauge",
			Help:   "Memory information field MemTotal_bytes.",
			Labels: []string{"instance", "job"},
		},
		{
			Name: "up",
			Type: "gauge",
			Help: "Whether the target is up.",
		},
	}

	doc := generateMetricsDoc(metrics)

	// Check header
	if !strings.Contains(doc, "# Prometheus Metrics Reference for Dashyard") {
		t.Error("missing header")
	}

	// Check dashboard YAML format section
	if !strings.Contains(doc, "## Dashboard YAML Format") {
		t.Error("missing dashboard YAML format section")
	}
	if !strings.Contains(doc, "chart_type: line") {
		t.Error("missing chart_type reference in YAML format")
	}

	// Check PromQL guidelines
	if !strings.Contains(doc, "## PromQL Guidelines") {
		t.Error("missing PromQL guidelines section")
	}
	if !strings.Contains(doc, "rate(") {
		t.Error("missing rate() guidance")
	}

	// Check unit selection
	if !strings.Contains(doc, "## Unit Selection") {
		t.Error("missing unit selection section")
	}

	// Check metrics are present
	if !strings.Contains(doc, "**`go_gc_duration_seconds`** (summary)") {
		t.Error("missing go_gc_duration_seconds metric")
	}
	if !strings.Contains(doc, "A summary of the pause duration") {
		t.Error("missing go_gc_duration_seconds help text")
	}
	if !strings.Contains(doc, "Labels: `quantile`") {
		t.Error("missing go_gc_duration_seconds labels")
	}

	if !strings.Contains(doc, "**`node_cpu_seconds_total`** (counter)") {
		t.Error("missing node_cpu_seconds_total metric")
	}
	if !strings.Contains(doc, "Labels: `cpu`, `mode`") {
		t.Error("missing node_cpu_seconds_total labels")
	}

	// Check metrics with no labels
	if !strings.Contains(doc, "**`up`** (gauge)") {
		t.Error("missing up metric")
	}

	// Check grouping headers
	if !strings.Contains(doc, "### go_gc") {
		t.Error("missing go_gc group header")
	}
	if !strings.Contains(doc, "### node_cpu") {
		t.Error("missing node_cpu group header")
	}
	if !strings.Contains(doc, "### node_memory") {
		t.Error("missing node_memory group header")
	}
}

func TestGenerateMetricsDocEmpty(t *testing.T) {
	doc := generateMetricsDoc(nil)
	if !strings.Contains(doc, "No metrics found.") {
		t.Error("expected 'No metrics found.' for empty metrics")
	}
}

func TestGroupMetricsByPrefix(t *testing.T) {
	metrics := []prometheus.MetricInfo{
		{Name: "go_gc_duration_seconds"},
		{Name: "go_gc_heap_allocs_total"},
		{Name: "node_cpu_seconds_total"},
		{Name: "node_memory_MemTotal_bytes"},
		{Name: "up"},
	}

	groups := groupMetricsByPrefix(metrics)

	if len(groups) != 4 {
		t.Errorf("expected 4 groups, got %d", len(groups))
	}

	if len(groups["go_gc"]) != 2 {
		t.Errorf("expected 2 metrics in go_gc, got %d", len(groups["go_gc"]))
	}
	if len(groups["node_cpu"]) != 1 {
		t.Errorf("expected 1 metric in node_cpu, got %d", len(groups["node_cpu"]))
	}
	if len(groups["node_memory"]) != 1 {
		t.Errorf("expected 1 metric in node_memory, got %d", len(groups["node_memory"]))
	}
	if len(groups["up"]) != 1 {
		t.Errorf("expected 1 metric in up, got %d", len(groups["up"]))
	}
}

func TestMetricPrefix(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"go_gc_duration_seconds", "go_gc"},
		{"node_cpu_seconds_total", "node_cpu"},
		{"up", "up"},
		{"single_part", "single_part"},
	}

	for _, tt := range tests {
		got := metricPrefix(tt.name)
		if got != tt.expected {
			t.Errorf("metricPrefix(%q) = %q, want %q", tt.name, got, tt.expected)
		}
	}
}
