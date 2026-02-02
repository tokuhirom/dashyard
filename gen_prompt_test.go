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

	// Check legend candidates
	if !strings.Contains(doc, "Legend candidates: mode, cpu") {
		t.Errorf("missing or wrong legend candidates for node_cpu_seconds_total, got:\n%s", doc)
	}
	if !strings.Contains(doc, "Legend candidates: quantile") {
		t.Error("missing legend candidates for go_gc_duration_seconds")
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

	// Verify new file names are referenced
	if !strings.Contains(readme, "prompt-system.md") {
		t.Error("README should reference prompt-system.md")
	}
	if !strings.Contains(readme, "prompt-user.md") {
		t.Error("README should reference prompt-user.md")
	}
	if !strings.Contains(readme, "DASHYARD-USAGE.md") {
		t.Error("README should reference DASHYARD-USAGE.md")
	}

	// Verify old file names and --overwrite are not referenced
	if strings.Contains(readme, "--overwrite") {
		t.Error("README should not reference --overwrite flag")
	}
}

func TestPromptUserTemplate(t *testing.T) {
	tmpl := prompt.PromptUserTemplate
	if tmpl == "" {
		t.Fatal("PromptUserTemplate is empty")
	}
	if !strings.Contains(tmpl, "# Custom Prompt") {
		t.Error("PromptUserTemplate should contain '# Custom Prompt' heading")
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

	// Check users section with password hash
	if !strings.Contains(config, "users:") {
		t.Error("config should contain users section")
	}
	if !strings.Contains(config, `id: "admin"`) {
		t.Error("config should contain admin user")
	}
	if !strings.Contains(config, "$6$") {
		t.Error("config should contain SHA-512 crypt password hash")
	}
	// Check that plain password is shown in comment
	if !strings.Contains(config, "# Default user: admin /") {
		t.Error("config should contain plain password in comment")
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

func TestGenerateConfigRandomPassword(t *testing.T) {
	config1, err := generateConfig("http://localhost:9090")
	if err != nil {
		t.Fatalf("generateConfig returned error: %v", err)
	}
	config2, err := generateConfig("http://localhost:9090")
	if err != nil {
		t.Fatalf("generateConfig returned error: %v", err)
	}

	if config1 == config2 {
		t.Error("each call should generate a different random password")
	}
}

func TestClassifyLabels(t *testing.T) {
	tests := []struct {
		name             string
		metric           prometheus.MetricInfo
		wantFixed        []string
		wantVariable     []string
		wantLegend       []string
	}{
		{
			name: "mixed labels",
			metric: prometheus.MetricInfo{
				Name:   "lb_server_up_ratio",
				Labels: []string{"sakuracloud_variant", "sakuracloud_publisher", "status", "server_id"},
				LabelValues: map[string][]string{
					"sakuracloud_variant":   {"lb_metrics"},                                     // 1 value → fixed
					"sakuracloud_publisher": {"apprun-dedicated"},                               // 1 value → fixed
					"status":               {"up", "down", "maintenance"},                       // 3 values → legend
					"server_id":            {"s1", "s2", "s3", "s4", "s5", "s6", "s7"},         // 7 values → variable
				},
			},
			wantFixed:    []string{"sakuracloud_variant", "sakuracloud_publisher"},
			wantVariable: []string{"server_id"},
			wantLegend:   []string{"status"},
		},
		{
			name: "legend sorted by value count ascending",
			metric: prometheus.MetricInfo{
				Name:   "http_requests_total",
				Labels: []string{"method", "status", "path"},
				LabelValues: map[string][]string{
					"method": {"GET", "POST", "PUT", "DELETE"}, // 4 values
					"status": {"200", "404"},                   // 2 values → first
					"path":   {"a", "b", "c"},                  // 3 values → second
				},
			},
			wantFixed:    nil,
			wantVariable: nil,
			wantLegend:   []string{"status", "path", "method"},
		},
		{
			name: "no label values known",
			metric: prometheus.MetricInfo{
				Name:   "some_metric",
				Labels: []string{"instance", "job"},
			},
			wantFixed:    []string{"instance", "job"},
			wantVariable: nil,
			wantLegend:   nil,
		},
		{
			name: "all fixed",
			metric: prometheus.MetricInfo{
				Name:   "singleton_metric",
				Labels: []string{"env"},
				LabelValues: map[string][]string{
					"env": {"production"},
				},
			},
			wantFixed:    []string{"env"},
			wantVariable: nil,
			wantLegend:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixed, variable, legend := classifyLabels(tt.metric)
			if !slicesEqual(fixed, tt.wantFixed) {
				t.Errorf("fixed: got %v, want %v", fixed, tt.wantFixed)
			}
			if !slicesEqual(variable, tt.wantVariable) {
				t.Errorf("variable: got %v, want %v", variable, tt.wantVariable)
			}
			if !slicesEqual(legend, tt.wantLegend) {
				t.Errorf("legend: got %v, want %v", legend, tt.wantLegend)
			}
		})
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGenerateMetricsDocLabelClassification(t *testing.T) {
	metrics := []prometheus.MetricInfo{
		{
			Name:   "lb_server_up_ratio",
			Type:   "gauge",
			Help:   "Server up ratio.",
			Labels: []string{"variant", "publisher", "status"},
			LabelValues: map[string][]string{
				"variant":   {"lb_metrics"},
				"publisher": {"apprun-dedicated"},
				"status":    {"up", "down"},
			},
		},
	}

	doc := generateMetricsDoc(metrics)

	if !strings.Contains(doc, `Fixed: variant="lb_metrics", publisher="apprun-dedicated"`) {
		t.Errorf("missing fixed labels with values, got:\n%s", doc)
	}
	if !strings.Contains(doc, "Legend candidates: status") {
		t.Errorf("missing legend candidates, got:\n%s", doc)
	}
	if strings.Contains(doc, "Variable candidates:") {
		t.Error("should not have variable candidates")
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
