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
	"github.com/tokuhirom/dashyard/internal/prompt"
)

type GenPromptCmd struct {
	URL           string        `arg:"" help:"Prometheus server URL."`
	BearerToken   string        `help:"Bearer token for authentication." env:"PROMETHEUS_BEARER_TOKEN"`
	Match         string        `help:"Regex to filter metric names." default:""`
	Timeout       time.Duration `help:"HTTP timeout." default:"30s"`
	OutputDir     string        `help:"Output directory for prompt.md and prompt-labels.md (default: stdout)." short:"o" default:""`
	Guidelines    string        `help:"Custom guidelines markdown file to replace default guidelines." default:""`
}

func (cmd *GenPromptCmd) Run() error {
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

	// Load guidelines
	guidelines := prompt.DefaultGuidelines
	if cmd.Guidelines != "" {
		data, err := os.ReadFile(cmd.Guidelines)
		if err != nil {
			return fmt.Errorf("reading guidelines file: %w", err)
		}
		guidelines = string(data)
		slog.Info("using custom guidelines", "file", cmd.Guidelines)
	}

	// Generate output
	if cmd.OutputDir != "" {
		promptFile := filepath.Join(cmd.OutputDir, "prompt.md")
		labelsFile := filepath.Join(cmd.OutputDir, "prompt-labels.md")

		mainDoc := generateMetricsDoc(metrics, "prompt-labels.md", guidelines)
		labelsDoc := generateLabelsDoc(metrics)

		if err := os.MkdirAll(cmd.OutputDir, 0755); err != nil {
			return fmt.Errorf("creating output directory: %w", err)
		}
		if err := os.WriteFile(promptFile, []byte(mainDoc), 0644); err != nil {
			return fmt.Errorf("writing prompt file: %w", err)
		}
		slog.Info("wrote prompt file", "file", promptFile)

		if err := os.WriteFile(labelsFile, []byte(labelsDoc), 0644); err != nil {
			return fmt.Errorf("writing labels file: %w", err)
		}
		slog.Info("wrote labels file", "file", labelsFile)
	} else {
		// stdout: output everything in one stream
		mainDoc := generateMetricsDoc(metrics, "", guidelines)
		fmt.Print(mainDoc)
		labelsDoc := generateLabelsDoc(metrics)
		if labelsDoc != "" {
			fmt.Print("\n---\n\n")
			fmt.Print(labelsDoc)
		}
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

// generateMetricsDoc generates the main prompt file.
// labelsFileName is the basename of the labels file (empty if stdout mode).
// guidelines is the customizable guidelines content.
func generateMetricsDoc(metrics []prometheus.MetricInfo, labelsFileName string, guidelines string) string {
	var sb strings.Builder

	// 1. Guidelines (customizable)
	sb.WriteString(guidelines)
	sb.WriteString("\n\n")

	// 2. Format reference (universal, always included)
	sb.WriteString(prompt.FormatReference)
	sb.WriteString("\n\n")

	// 3. Labels file reference
	if labelsFileName != "" {
		sb.WriteString(fmt.Sprintf("# Label Details\n\nThe full list of label values for each metric is available in `%s`. Refer to it when you need to know the exact values of a label (e.g. to enumerate devices, CPU cores, or states).\n\n", labelsFileName))
	}

	// 5. Metrics listing
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
