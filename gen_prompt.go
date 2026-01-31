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
	Output        string        `help:"Output file (default: stdout). A labels file is also written alongside." short:"o" default:""`
	Guidelines    string        `help:"Custom guidelines markdown file to replace default guidelines." default:""`
	DashboardsDir string        `help:"Directory of existing dashboard YAML files to include as context." name:"dashboards-dir" default:""`
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

	// Load existing dashboards
	var existingDashboards string
	if cmd.DashboardsDir != "" {
		content, err := loadExistingDashboards(cmd.DashboardsDir)
		if err != nil {
			return fmt.Errorf("loading existing dashboards: %w", err)
		}
		existingDashboards = content
		slog.Info("loaded existing dashboards", "dir", cmd.DashboardsDir)
	}

	// Generate output
	if cmd.Output != "" {
		labelsFile := labelsFilePath(cmd.Output)
		labelsBaseName := filepath.Base(labelsFile)

		mainDoc := generateMetricsDoc(metrics, labelsBaseName, guidelines, existingDashboards)
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
		mainDoc := generateMetricsDoc(metrics, "", guidelines, existingDashboards)
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

// loadExistingDashboards reads all YAML files from a directory and returns their content
// formatted for inclusion in the prompt.
func loadExistingDashboards(dir string) (string, error) {
	var sb strings.Builder
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = filepath.Base(path)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			slog.Warn("could not read dashboard file", "path", path, "error", err)
			return nil
		}

		sb.WriteString(fmt.Sprintf("## %s\n\n```yaml\n%s```\n\n", relPath, string(data)))
		return nil
	})
	if err != nil {
		return "", err
	}
	return sb.String(), nil
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
// existingDashboards is the formatted content of existing dashboard files (empty if not provided).
func generateMetricsDoc(metrics []prometheus.MetricInfo, labelsFileName string, guidelines string, existingDashboards string) string {
	var sb strings.Builder

	// 1. Guidelines (customizable)
	sb.WriteString(guidelines)
	sb.WriteString("\n\n")

	// 2. Format reference (universal, always included)
	sb.WriteString(prompt.FormatReference)
	sb.WriteString("\n\n")

	// 3. Existing dashboards (optional)
	if existingDashboards != "" {
		sb.WriteString("# Existing Dashboards\n\n")
		sb.WriteString("The following dashboards already exist. When adding new metrics, add panels to the appropriate existing dashboard or create a new file if the domain is new. Before modifying an existing file, check `git log -p <file>` for manual edits and ask the user before overwriting those.\n\n")
		sb.WriteString(existingDashboards)
	}

	// 4. Labels file reference
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
