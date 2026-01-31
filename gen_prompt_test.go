package main

import (
	"strings"
	"testing"

	"github.com/tokuhirom/dashyard/internal/prometheus"
	"github.com/tokuhirom/dashyard/internal/prompt"
)

func TestGenerateMetricsDoc(t *testing.T) {
	metrics := []prometheus.MetricInfo{
		{
			Name:   "go_gc_duration_seconds",
			Type:   "summary",
			Help:   "A summary of the pause duration of garbage collection cycles.",
			Labels: []string{"quantile"},
			LabelValues: map[string][]string{
				"quantile": {"0.5", "0.75", "0.99"},
			},
		},
		{
			Name:   "node_cpu_seconds_total",
			Type:   "counter",
			Help:   "Seconds the CPUs spent in each mode.",
			Labels: []string{"cpu", "mode"},
			LabelValues: map[string][]string{
				"cpu":  {"0", "1", "2", "3"},
				"mode": {"idle", "system", "user"},
			},
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

	doc := generateMetricsDoc(metrics, "metrics-labels.md", prompt.DefaultGuidelines)

	// Check LLM role instruction (from default guidelines)
	if !strings.Contains(doc, "You are a Dashyard dashboard generator") {
		t.Error("missing LLM role instruction")
	}

	// Check multi-file generation instruction
	if !strings.Contains(doc, "generate one or more YAML files") {
		t.Error("missing multi-file generation instruction")
	}
	if !strings.Contains(doc, "# File:") {
		t.Error("missing file path comment example")
	}

	// Check file organization section (from guidelines)
	if !strings.Contains(doc, "# File Organization") {
		t.Error("missing file organization section")
	}

	// Check dashboard YAML format section (from format reference)
	if !strings.Contains(doc, "# Dashboard YAML Format") {
		t.Error("missing dashboard YAML format section")
	}

	// Check rules (from format reference)
	if !strings.Contains(doc, "# Rules") {
		t.Error("missing rules section")
	}
	if !strings.Contains(doc, "rate(") {
		t.Error("missing rate() guidance")
	}

	// Check rate/unit correspondence
	if !strings.Contains(doc, "bytes/sec") {
		t.Error("missing rate/unit correspondence for bytes")
	}

	// Check stacked chart guidance
	if !strings.Contains(doc, "Stacked Charts") {
		t.Error("missing stacked chart guidance")
	}

	// Check ratio/percent guidance
	if !strings.Contains(doc, "Ratio and Percent") {
		t.Error("missing ratio/percent guidance")
	}

	// Check default behavior
	if !strings.Contains(doc, "# Default Behavior") {
		t.Error("missing default behavior section")
	}

	// Check complete example
	if !strings.Contains(doc, "# Complete Example") {
		t.Error("missing complete example")
	}

	// Check dashyard validate reference
	if !strings.Contains(doc, "dashyard validate") {
		t.Error("missing dashyard validate reference")
	}

	// Check label details reference
	if !strings.Contains(doc, "metrics-labels.md") {
		t.Error("missing labels file reference")
	}

	// Check metrics with label value counts
	if !strings.Contains(doc, "cpu (4 values)") {
		t.Error("missing cpu label value count")
	}
	if !strings.Contains(doc, "mode (3 values)") {
		t.Error("missing mode label value count")
	}

	// Check metrics without label values still show label names
	if !strings.Contains(doc, "Labels: instance, job") {
		t.Error("missing labels without counts for node_memory_MemTotal_bytes")
	}

	// Check grouping headers
	if !strings.Contains(doc, "## go_gc") {
		t.Error("missing go_gc group header")
	}
	if !strings.Contains(doc, "## node_cpu") {
		t.Error("missing node_cpu group header")
	}
}

func TestGenerateMetricsDocCustomGuidelines(t *testing.T) {
	metrics := []prometheus.MetricInfo{
		{Name: "up", Type: "gauge"},
	}
	customGuidelines := "You are a custom dashboard bot. Follow my rules."
	doc := generateMetricsDoc(metrics, "", customGuidelines)

	// Custom guidelines should be present
	if !strings.Contains(doc, "You are a custom dashboard bot") {
		t.Error("missing custom guidelines")
	}
	// Default guidelines should NOT be present
	if strings.Contains(doc, "You are a Dashyard dashboard generator") {
		t.Error("default guidelines should not be present when custom guidelines are used")
	}
	// Format reference should still be present
	if !strings.Contains(doc, "# Dashboard YAML Format") {
		t.Error("format reference should always be present")
	}
}

func TestGenerateMetricsDocNoLabelsFileName(t *testing.T) {
	metrics := []prometheus.MetricInfo{
		{Name: "up", Type: "gauge"},
	}
	doc := generateMetricsDoc(metrics, "", prompt.DefaultGuidelines)
	if strings.Contains(doc, "# Label Details") {
		t.Error("should not contain Label Details section when labelsFileName is empty")
	}
}

func TestGenerateMetricsDocEmpty(t *testing.T) {
	doc := generateMetricsDoc(nil, "", prompt.DefaultGuidelines)
	if !strings.Contains(doc, "No metrics available.") {
		t.Error("expected 'No metrics available.' for empty metrics")
	}
}

func TestGenerateLabelsDoc(t *testing.T) {
	metrics := []prometheus.MetricInfo{
		{
			Name:   "node_cpu_seconds_total",
			Labels: []string{"cpu", "mode"},
			LabelValues: map[string][]string{
				"cpu":  {"0", "1"},
				"mode": {"idle", "system", "user"},
			},
		},
		{
			Name: "up",
		},
	}

	doc := generateLabelsDoc(metrics)

	if !strings.Contains(doc, "# Label Values") {
		t.Error("missing header")
	}
	if !strings.Contains(doc, "## node_cpu_seconds_total") {
		t.Error("missing metric name")
	}
	if !strings.Contains(doc, "**cpu**: 0, 1") {
		t.Error("missing cpu values")
	}
	if !strings.Contains(doc, "**mode**: idle, system, user") {
		t.Error("missing mode values")
	}
	// "up" has no label values, should not appear
	if strings.Contains(doc, "## up") {
		t.Error("up should not appear in labels doc")
	}
}

func TestGenerateLabelsDocEmpty(t *testing.T) {
	metrics := []prometheus.MetricInfo{
		{Name: "up"},
	}
	doc := generateLabelsDoc(metrics)
	if doc != "" {
		t.Errorf("expected empty string for metrics without label values, got %q", doc)
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

