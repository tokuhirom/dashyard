package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
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
	Output      string        `help:"Output file (default: stdout). A labels file is also written alongside." short:"o" default:""`
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

		// Fetch label names (soft failure per metric)
		labels, err := client.MetricLabels(ctx, name)
		if err != nil {
			slog.Warn("could not fetch labels", "metric", name, "error", err)
		} else {
			info.Labels = labels
		}

		// Fetch label values for each label (soft failure per label)
		if len(info.Labels) > 0 {
			info.LabelValues = make(map[string][]string)
			for _, label := range info.Labels {
				values, err := client.MetricLabelValues(ctx, name, label)
				if err != nil {
					slog.Warn("could not fetch label values", "metric", name, "label", label, "error", err)
				} else {
					info.LabelValues[label] = values
				}
			}
		}

		metrics = append(metrics, info)
	}

	// Generate output
	if cmd.Output != "" {
		labelsFile := labelsFilePath(cmd.Output)
		labelsBaseName := filepath.Base(labelsFile)

		mainDoc := generateMetricsDoc(metrics, labelsBaseName)
		labelsDoc := generateLabelsDoc(metrics)

		if err := os.WriteFile(cmd.Output, []byte(mainDoc), 0644); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
		slog.Info("wrote prompt file", "file", cmd.Output)

		if err := os.WriteFile(labelsFile, []byte(labelsDoc), 0644); err != nil {
			return fmt.Errorf("writing labels file: %w", err)
		}
		slog.Info("wrote labels file", "file", labelsFile)
	} else {
		// stdout: output everything in one stream
		mainDoc := generateMetricsDoc(metrics, "")
		fmt.Print(mainDoc)
		labelsDoc := generateLabelsDoc(metrics)
		if labelsDoc != "" {
			fmt.Print("\n---\n\n")
			fmt.Print(labelsDoc)
		}
	}

	return nil
}

// labelsFilePath derives the labels file path from the main output path.
// e.g., "metrics.md" -> "metrics-labels.md"
func labelsFilePath(mainPath string) string {
	ext := filepath.Ext(mainPath)
	base := strings.TrimSuffix(mainPath, ext)
	return base + "-labels" + ext
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

// generateMetricsDoc generates the main prompt file.
// labelsFileName is the basename of the labels file (empty if stdout mode).
func generateMetricsDoc(metrics []prometheus.MetricInfo, labelsFileName string) string {
	var sb strings.Builder

	// Role and task
	sb.WriteString("You are a Dashyard dashboard generator. Your task is to create Dashyard dashboard YAML files based on user requests.\n\n")
	sb.WriteString("When the user asks for dashboards, generate one or more YAML files. Group metrics by domain — put closely related metrics together in the same dashboard. For example, host metrics (CPU, memory, disk, network) belong in a single dashboard. Separate dashboards are for distinct domains like JVM, HTTP, database, or application-specific metrics.\n\n")
	sb.WriteString("Dashyard loads all YAML files from a dashboard directory. Subdirectories become collapsible groups in the sidebar. Output each file with a comment indicating its path, for example:\n\n")
	sb.WriteString("```\n")
	sb.WriteString("# File: host.yaml\n")
	sb.WriteString("title: \"Host Metrics\"\n")
	sb.WriteString("...\n\n")
	sb.WriteString("# File: jvm.yaml\n")
	sb.WriteString("title: \"JVM\"\n")
	sb.WriteString("...\n")
	sb.WriteString("```\n\n")

	// Dashboard YAML format
	sb.WriteString("# Dashboard YAML Format\n\n")
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

	// Rules
	sb.WriteString("# Rules\n\n")

	sb.WriteString("## PromQL by Metric Type\n\n")
	sb.WriteString("- counter: always wrap with `rate(...[5m])` or `increase(...[5m])`. Never use a raw counter.\n")
	sb.WriteString("- histogram (_bucket): use `histogram_quantile(0.99, rate(...[5m]))`\n")
	sb.WriteString("- gauge: use directly, or apply `avg()`, `sum()`, `min()`, `max()`\n")
	sb.WriteString("- summary (_sum/_count): `rate(sum[5m]) / rate(count[5m])` for average\n\n")

	sb.WriteString("## Filtering and Grouping\n\n")
	sb.WriteString("- Filter with `{label=\"value\"}`, group with `by (label)`\n")
	sb.WriteString("- Use `$variable` to reference dashboard variables in queries\n\n")

	sb.WriteString("## Unit Selection\n\n")
	sb.WriteString("- `bytes` — memory, disk, network I/O metrics\n")
	sb.WriteString("- `percent` — ratios and utilization (0-100 scale)\n")
	sb.WriteString("- `seconds` — durations and latencies\n")
	sb.WriteString("- `count` — counts, rates, and dimensionless values\n\n")

	sb.WriteString("## File Organization\n\n")
	sb.WriteString("- Group metrics by domain into one dashboard (e.g. host metrics: CPU + memory + disk + network in one file)\n")
	sb.WriteString("- Separate dashboards for distinct domains (e.g. `host.yaml`, `jvm.yaml`, `http.yaml`, `database.yaml`)\n")
	sb.WriteString("- Use subdirectories when there are many dashboards (e.g. `app/api.yaml`, `app/workers.yaml`)\n")
	sb.WriteString("- Use rows within a dashboard to separate sub-topics (e.g. CPU row, Memory row, Disk row)\n\n")

	sb.WriteString("## Best Practices\n\n")
	sb.WriteString("- Group related panels into rows with descriptive titles\n")
	sb.WriteString("- When a metric has a label with many values (e.g. device, cpu), use a variable with `label_values()` and `$variable` in queries\n")
	sb.WriteString("- Use `repeat` on a row to auto-expand for each variable value\n")
	sb.WriteString("- Add `thresholds` for metrics with known warning/critical levels\n")
	sb.WriteString("- Add a markdown panel to explain what the dashboard monitors\n")
	sb.WriteString("- Validate generated YAML with `dashyard validate` before deploying\n")
	sb.WriteString("- Output each file starting with `# File: path/name.yaml` followed by the YAML content\n\n")

	// Labels file reference
	if labelsFileName != "" {
		sb.WriteString(fmt.Sprintf("# Label Details\n\nThe full list of label values for each metric is available in `%s`. Refer to it when you need to know the exact values of a label (e.g. to enumerate devices, CPU cores, or states).\n\n", labelsFileName))
	}

	// Metrics listing
	sb.WriteString("# Available Metrics\n\n")

	if len(metrics) == 0 {
		sb.WriteString("No metrics available.\n")
		return sb.String()
	}

	groups := groupMetricsByPrefix(metrics)

	groupKeys := make([]string, 0, len(groups))
	for k := range groups {
		groupKeys = append(groupKeys, k)
	}
	sort.Strings(groupKeys)

	for _, prefix := range groupKeys {
		group := groups[prefix]
		sb.WriteString(fmt.Sprintf("## %s\n\n", prefix))

		for _, m := range group {
			sb.WriteString(fmt.Sprintf("- `%s`", m.Name))
			if m.Type != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", m.Type))
			}
			if m.Help != "" {
				sb.WriteString(fmt.Sprintf(" — %s", m.Help))
			}
			sb.WriteString("\n")
			if m.Unit != "" {
				sb.WriteString(fmt.Sprintf("  Unit: %s\n", m.Unit))
			}
			if len(m.Labels) > 0 {
				var labelParts []string
				for _, l := range m.Labels {
					if vals, ok := m.LabelValues[l]; ok && len(vals) > 0 {
						labelParts = append(labelParts, fmt.Sprintf("%s (%d values)", l, len(vals)))
					} else {
						labelParts = append(labelParts, l)
					}
				}
				sb.WriteString(fmt.Sprintf("  Labels: %s\n", strings.Join(labelParts, ", ")))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// generateLabelsDoc generates the label values detail file.
func generateLabelsDoc(metrics []prometheus.MetricInfo) string {
	// Check if there are any label values to write
	hasValues := false
	for _, m := range metrics {
		if len(m.LabelValues) > 0 {
			hasValues = true
			break
		}
	}
	if !hasValues {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("# Label Values\n\n")
	sb.WriteString("Full label value listings for each metric.\n\n")

	for _, m := range metrics {
		if len(m.LabelValues) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("## %s\n\n", m.Name))
		for _, label := range m.Labels {
			vals, ok := m.LabelValues[label]
			if !ok || len(vals) == 0 {
				continue
			}
			sb.WriteString(fmt.Sprintf("- **%s**: %s\n", label, strings.Join(vals, ", ")))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
