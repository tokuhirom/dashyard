package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/tokuhirom/dashyard/internal/prometheus"
	"github.com/tokuhirom/dashyard/internal/prompt"
)

type GenPromptCmd struct {
	URL       string        `arg:"" help:"Prometheus server URL."`
	BearerToken string      `help:"Bearer token for authentication." env:"PROMETHEUS_BEARER_TOKEN"`
	Match     string        `help:"Regex to filter metric names." default:""`
	Timeout   time.Duration `help:"HTTP timeout." default:"30s"`
	OutputDir  string        `help:"Output directory for prompt.md and prompt-metrics.md (default: stdout)." short:"o" default:""`
	Overwrite bool          `help:"Overwrite all write-once files (prompt.md, README.md, config.yaml)." default:"false"`
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

	// Generate output
	if cmd.OutputDir != "" {
		promptFile := filepath.Join(cmd.OutputDir, "prompt.md")
		metricsFile := filepath.Join(cmd.OutputDir, "prompt-metrics.md")
		readmeFile := filepath.Join(cmd.OutputDir, "README.md")
		configFile := filepath.Join(cmd.OutputDir, "config.yaml")

		if err := os.MkdirAll(cmd.OutputDir, 0755); err != nil {
			return fmt.Errorf("creating output directory: %w", err)
		}

		// Write prompt.md only if it doesn't exist (user-editable template), unless --force-prompt
		if cmd.Overwrite {
			promptDoc := generatePromptDoc()
			if err := os.WriteFile(promptFile, []byte(promptDoc), 0644); err != nil {
				return fmt.Errorf("writing prompt file: %w", err)
			}
			slog.Info("wrote prompt file (forced)", "file", promptFile)
		} else if _, err := os.Stat(promptFile); os.IsNotExist(err) {
			promptDoc := generatePromptDoc()
			if err := os.WriteFile(promptFile, []byte(promptDoc), 0644); err != nil {
				return fmt.Errorf("writing prompt file: %w", err)
			}
			slog.Info("wrote prompt file", "file", promptFile)
		} else {
			slog.Info("prompt file already exists, skipping", "file", promptFile)
		}

		// Write README.md (write-once, unless --force-prompt)
		writeOnceFile(readmeFile, generateREADME(), "README", cmd.Overwrite)

		// Write config.yaml (write-once, unless --force-prompt)
		configContent, err := generateConfig(cmd.URL)
		if err != nil {
			return fmt.Errorf("generating config: %w", err)
		}
		writeOnceFile(configFile, configContent, "config", cmd.Overwrite)

		// Always overwrite prompt-metrics.md
		metricsDoc := generateMetricsDoc(metrics)
		if err := os.WriteFile(metricsFile, []byte(metricsDoc), 0644); err != nil {
			return fmt.Errorf("writing metrics file: %w", err)
		}
		slog.Info("wrote metrics file", "file", metricsFile)
	} else {
		// stdout: output everything in one stream
		fmt.Print(generateREADME())
		fmt.Print("\n---\n\n")
		configContent, err := generateConfig(cmd.URL)
		if err != nil {
			return fmt.Errorf("generating config: %w", err)
		}
		fmt.Print(configContent)
		fmt.Print("\n---\n\n")
		fmt.Print(generatePromptDoc())
		fmt.Print("\n---\n\n")
		fmt.Print(generateMetricsDoc(metrics))
	}

	return nil
}

// writeOnceFile writes content to a file only if it doesn't exist, unless force is true.
func writeOnceFile(path, content, label string, force bool) {
	if force {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			slog.Error("writing file", "file", path, "error", err)
			return
		}
		slog.Info("wrote file (forced)", "type", label, "file", path)
	} else if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			slog.Error("writing file", "file", path, "error", err)
			return
		}
		slog.Info("wrote file", "type", label, "file", path)
	} else {
		slog.Info("file already exists, skipping", "type", label, "file", path)
	}
}

// generateREADME returns the README.md content.
func generateREADME() string {
	return prompt.ReadmeTemplate
}

// generateConfig generates config.yaml content with the Prometheus URL embedded.
func generateConfig(prometheusURL string) (string, error) {
	tmpl, err := template.New("config").Parse(prompt.ConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing config template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct{ PrometheusURL string }{PrometheusURL: prometheusURL}); err != nil {
		return "", fmt.Errorf("executing config template: %w", err)
	}
	return buf.String(), nil
}

// generatePromptDoc generates the static prompt template (guidelines + format reference).
func generatePromptDoc() string {
	var sb strings.Builder
	sb.WriteString(prompt.DefaultGuidelines)
	sb.WriteString("\n\n")
	sb.WriteString(prompt.FormatReference)
	sb.WriteString("\n")
	return sb.String()
}

// generateMetricsDoc generates the metrics file (metric listing + label values).
func generateMetricsDoc(metrics []prometheus.MetricInfo) string {
	var sb strings.Builder

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

	// Label values detail
	hasValues := false
	for _, m := range metrics {
		if len(m.LabelValues) > 0 {
			hasValues = true
			break
		}
	}
	if hasValues {
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
	}

	return sb.String()
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
