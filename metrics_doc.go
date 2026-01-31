package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/tokuhirom/dashyard/internal/prometheus"
)

type MetricsDocCmd struct {
	URL         string        `arg:"" help:"Prometheus server URL."`
	BearerToken string        `help:"Bearer token for authentication." env:"PROMETHEUS_BEARER_TOKEN"`
	Match       string        `help:"Regex to filter metric names." default:""`
	Timeout     time.Duration `help:"HTTP timeout." default:"30s"`
	Output      string        `help:"Output file (default: stdout)." short:"o" default:""`
}

func (cmd *MetricsDocCmd) Run() error {
	ctx := context.Background()

	var opts []prometheus.ClientOption
	if cmd.BearerToken != "" {
		opts = append(opts, prometheus.WithBearerToken(cmd.BearerToken))
	}
	client := prometheus.NewClient(cmd.URL, cmd.Timeout, opts...)

	// Fetch all metric names
	slog.Info("fetching metric names", "url", cmd.URL)
	names, err := client.MetricNames(ctx)
	if err != nil {
		return fmt.Errorf("fetching metric names: %w", err)
	}
	slog.Info("discovered metrics", "count", len(names))

	// Filter by --match regex
	if cmd.Match != "" {
		re, err := regexp.Compile(cmd.Match)
		if err != nil {
			return fmt.Errorf("invalid --match regex: %w", err)
		}
		var filtered []string
		for _, name := range names {
			if re.MatchString(name) {
				filtered = append(filtered, name)
			}
		}
		slog.Info("filtered metrics", "count", len(filtered), "pattern", cmd.Match)
		names = filtered
	}

	// Fetch metadata (soft failure)
	slog.Info("fetching metric metadata")
	metadata, err := client.MetricMetadata(ctx)
	if err != nil {
		slog.Warn("could not fetch metadata (some Prometheus-compatible systems don't support this endpoint)", "error", err)
		metadata = nil
	}

	// Build MetricInfo for each name
	metrics := make([]prometheus.MetricInfo, 0, len(names))
	for _, name := range names {
		info := prometheus.MetricInfo{Name: name}

		// Fill metadata if available
		if metadata != nil {
			if entries, ok := metadata[name]; ok && len(entries) > 0 {
				info.Type = entries[0].Type
				info.Help = entries[0].Help
				info.Unit = entries[0].Unit
			}
		}

		// Fetch labels (soft failure per metric)
		labels, err := client.MetricLabels(ctx, name)
		if err != nil {
			slog.Warn("could not fetch labels", "metric", name, "error", err)
		} else {
			info.Labels = labels
		}

		metrics = append(metrics, info)
	}

	// Generate markdown
	doc := generateMetricsDoc(metrics)

	// Write output
	if cmd.Output != "" {
		if err := os.WriteFile(cmd.Output, []byte(doc), 0644); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
		slog.Info("wrote metrics documentation", "file", cmd.Output)
	} else {
		fmt.Print(doc)
	}

	return nil
}

// groupMetricsByPrefix groups metrics by their common prefix (first two underscore-separated segments).
func groupMetricsByPrefix(metrics []prometheus.MetricInfo) map[string][]prometheus.MetricInfo {
	groups := make(map[string][]prometheus.MetricInfo)
	for _, m := range metrics {
		prefix := metricPrefix(m.Name)
		groups[prefix] = append(groups[prefix], m)
	}
	return groups
}

// metricPrefix extracts the first two underscore-separated segments, or the full name if fewer.
func metricPrefix(name string) string {
	parts := strings.SplitN(name, "_", 3)
	if len(parts) >= 2 {
		return parts[0] + "_" + parts[1]
	}
	return name
}

func generateMetricsDoc(metrics []prometheus.MetricInfo) string {
	var sb strings.Builder

	sb.WriteString("# Prometheus Metrics Reference for Dashyard\n\n")
	sb.WriteString("This document lists all available Prometheus metrics. Use it as context for generating Dashyard dashboard YAML files.\n\n")

	// Dashboard YAML format reference
	sb.WriteString("## Dashboard YAML Format\n\n")
	sb.WriteString("```yaml\n")
	sb.WriteString(`title: "Dashboard Title"
variables:                          # optional
  - name: device                    # variable name used as $device in queries
    label: "Network Device"         # display label
    query: "label_values(metric_name, label_name)"
rows:
  - title: "Row Title"
    repeat: device                  # optional: repeat row for each variable value
    panels:
      - title: "Panel Title"
        type: graph                 # "graph" or "markdown"
        query: 'promql_expression'  # required for graph
        unit: bytes                 # bytes, percent, count, seconds
        chart_type: line            # line, bar, area, scatter, pie, doughnut
        legend: "{label_name}"      # legend template
        y_min: 0                    # optional y-axis bounds
        y_max: 100
        stacked: false              # stack series
        thresholds:                 # optional reference lines
          - value: 80
            color: orange
            label: "Warning"
      - title: "Notes"
        type: markdown
        content: |
          Markdown content here.
`)
	sb.WriteString("```\n\n")

	// Guidelines
	sb.WriteString("## PromQL Guidelines\n\n")
	sb.WriteString("- **counter** metrics: always wrap with `rate(...[5m])` or `increase(...[5m])`\n")
	sb.WriteString("- **histogram** `_bucket` metrics: use `histogram_quantile(0.99, rate(...[5m]))`\n")
	sb.WriteString("- **gauge** metrics: use directly, or apply aggregation like `avg()`, `sum()`\n")
	sb.WriteString("- **summary** `_sum`/`_count` metrics: `rate(sum[5m]) / rate(count[5m])` for average\n")
	sb.WriteString("- Use `{label=\"value\"}` for filtering, `by (label)` for grouping\n")
	sb.WriteString("- Use `$variable` syntax to reference dashboard variables in queries\n\n")

	sb.WriteString("## Unit Selection\n\n")
	sb.WriteString("- `bytes` — metrics measuring bytes (memory, disk, network I/O)\n")
	sb.WriteString("- `percent` — ratios and utilization (0-100 scale)\n")
	sb.WriteString("- `seconds` — durations and latencies\n")
	sb.WriteString("- `count` — counts, rates, and dimensionless values\n\n")

	// Metrics listing
	sb.WriteString("## Available Metrics\n\n")

	if len(metrics) == 0 {
		sb.WriteString("No metrics found.\n")
		return sb.String()
	}

	groups := groupMetricsByPrefix(metrics)

	// Sort group keys
	groupKeys := make([]string, 0, len(groups))
	for k := range groups {
		groupKeys = append(groupKeys, k)
	}
	sort.Strings(groupKeys)

	for _, prefix := range groupKeys {
		group := groups[prefix]
		sb.WriteString(fmt.Sprintf("### %s\n\n", prefix))

		for _, m := range group {
			sb.WriteString(fmt.Sprintf("**`%s`**", m.Name))
			if m.Type != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", m.Type))
			}
			sb.WriteString("\n")

			if m.Help != "" {
				sb.WriteString(fmt.Sprintf("  %s\n", m.Help))
			}
			if m.Unit != "" {
				sb.WriteString(fmt.Sprintf("  Unit: %s\n", m.Unit))
			}
			if len(m.Labels) > 0 {
				sb.WriteString(fmt.Sprintf("  Labels: `%s`\n", strings.Join(m.Labels, "`, `")))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
