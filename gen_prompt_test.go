package main

import (
	"strings"
	"testing"

	"github.com/tokuhirom/dashyard/internal/prometheus"
	"github.com/tokuhirom/dashyard/internal/prompt"
)

func TestGeneratePromptDoc(t *testing.T) {
	doc := generatePromptDoc()

	// Check guidelines are included
	if !strings.Contains(doc, "You are a Dashyard dashboard generator") {
		t.Error("missing LLM role instruction from guidelines")
	}

	// Check format reference is included
	if !strings.Contains(doc, "# Dashboard YAML Format") {
		t.Error("missing dashboard YAML format section")
	}
	if !strings.Contains(doc, "# Rules") {
		t.Error("missing rules section")
	}
}

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

	doc := generateMetricsDoc(metrics)

	// Check metrics listing
	if !strings.Contains(doc, "# Available Metrics") {
		t.Error("missing available metrics header")
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

	// Check label values section is included
	if !strings.Contains(doc, "# Label Values") {
		t.Error("missing label values section")
	}
	if !strings.Contains(doc, "**cpu**: 0, 1, 2, 3") {
		t.Error("missing cpu label values")
	}
	if !strings.Contains(doc, "**mode**: idle, system, user") {
		t.Error("missing mode label values")
	}
}

func TestGenerateMetricsDocEmpty(t *testing.T) {
	doc := generateMetricsDoc(nil)
	if !strings.Contains(doc, "No metrics available.") {
		t.Error("expected 'No metrics available.' for empty metrics")
	}
}

func TestGenerateMetricsDocNoLabelValues(t *testing.T) {
	metrics := []prometheus.MetricInfo{
		{Name: "up", Type: "gauge"},
	}
	doc := generateMetricsDoc(metrics)
	if strings.Contains(doc, "# Label Values") {
		t.Error("should not contain Label Values section when no metrics have label values")
	}
}

func TestPromptDocDoesNotContainMetrics(t *testing.T) {
	doc := generatePromptDoc()
	if strings.Contains(doc, "# Available Metrics") {
		t.Error("prompt doc should not contain metrics listing")
	}
}

func TestMetricsDocDoesNotContainGuidelines(t *testing.T) {
	metrics := []prometheus.MetricInfo{
		{Name: "up", Type: "gauge"},
	}
	doc := generateMetricsDoc(metrics)
	if strings.Contains(doc, "You are a Dashyard dashboard generator") {
		t.Error("metrics doc should not contain guidelines")
	}
}

func TestGeneratePromptDocContainsDefaultGuidelines(t *testing.T) {
	doc := generatePromptDoc()
	// Verify key sections from default guidelines
	if !strings.Contains(doc, prompt.DefaultGuidelines) {
		t.Error("prompt doc should contain default guidelines verbatim")
	}
	if !strings.Contains(doc, prompt.FormatReference) {
		t.Error("prompt doc should contain format reference verbatim")
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

func TestGenerateREADME(t *testing.T) {
	readme := generateREADME()

	// Check that README is non-empty
	if readme == "" {
		t.Fatal("generateREADME returned empty string")
	}

	// Check key section headings
	sections := []string{
		"## Files",
		"## Workflow",
		"## Customizing the Prompt",
		"## Configuring config.yaml",
		"## Multiple Datasources",
		"## Viewing Dashboards",
	}
	for _, section := range sections {
		if !strings.Contains(readme, section) {
			t.Errorf("README missing section: %s", section)
		}
	}
}

func TestGenerateConfig(t *testing.T) {
	config, err := generateConfig("http://localhost:9090")
	if err != nil {
		t.Fatalf("generateConfig returned error: %v", err)
	}

	// Check that URL is embedded
	if !strings.Contains(config, "http://localhost:9090") {
		t.Error("config should contain the Prometheus URL")
	}

	// Check basic structure
	if !strings.Contains(config, "datasources:") {
		t.Error("config should contain datasources key")
	}
	if !strings.Contains(config, "name: default") {
		t.Error("config should contain default datasource name")
	}
	if !strings.Contains(config, "default: true") {
		t.Error("config should mark datasource as default")
	}
}

func TestGenerateConfigDifferentURL(t *testing.T) {
	url := "https://prometheus.example.com:9090"
	config, err := generateConfig(url)
	if err != nil {
		t.Fatalf("generateConfig returned error: %v", err)
	}

	if !strings.Contains(config, url) {
		t.Errorf("config should contain URL %q", url)
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
