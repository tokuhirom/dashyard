package model

import "fmt"

// Threshold represents a horizontal reference line on a graph panel.
type Threshold struct {
	Value float64 `yaml:"value" json:"value"`
	Color string  `yaml:"color,omitempty" json:"color,omitempty"`
	Label string  `yaml:"label,omitempty" json:"label,omitempty"`
}

// Panel represents a single visualization panel within a dashboard row.
type Panel struct {
	Title      string      `yaml:"title" json:"title"`
	Type       string      `yaml:"type" json:"type"`       // "graph" or "markdown"
	ChartType  string      `yaml:"chart_type,omitempty" json:"chart_type,omitempty"`
	Query      string      `yaml:"query,omitempty" json:"query,omitempty"`
	Unit       string      `yaml:"unit,omitempty" json:"unit,omitempty"` // "bytes", "percent", "count", "seconds"
	YMin       *float64    `yaml:"y_min,omitempty" json:"y_min,omitempty"`
	YMax       *float64    `yaml:"y_max,omitempty" json:"y_max,omitempty"`
	Legend     string      `yaml:"legend,omitempty" json:"legend,omitempty"`
	Thresholds []Threshold `yaml:"thresholds,omitempty" json:"thresholds,omitempty"`
	Stacked    bool        `yaml:"stacked,omitempty" json:"stacked,omitempty"`
	Content    string      `yaml:"content,omitempty" json:"content,omitempty"`
	FullWidth  bool        `yaml:"full_width,omitempty" json:"full_width,omitempty"`
}

// Row represents a horizontal row of panels in a dashboard.
type Row struct {
	Title  string  `yaml:"title" json:"title"`
	Repeat string  `yaml:"repeat,omitempty" json:"repeat,omitempty"`
	Panels []Panel `yaml:"panels" json:"panels"`
}

// Variable represents a dashboard-level template variable populated from Prometheus label values.
type Variable struct {
	Name  string `yaml:"name" json:"name"`
	Label string `yaml:"label,omitempty" json:"label,omitempty"`
	Query string `yaml:"query" json:"query"`
}

// Dashboard represents a single dashboard definition loaded from YAML.
type Dashboard struct {
	Title     string     `yaml:"title" json:"title"`
	Variables []Variable `yaml:"variables,omitempty" json:"variables,omitempty"`
	Rows      []Row      `yaml:"rows" json:"rows"`
	Path      string     `yaml:"-" json:"path"` // Set by loader, not from YAML
}

var validChartTypes = map[string]bool{
	"line": true, "bar": true, "area": true,
	"scatter": true,
}

var validUnits = map[string]bool{
	"bytes": true, "percent": true, "count": true, "seconds": true,
}

// Validate checks the dashboard for semantic correctness.
func (d *Dashboard) Validate() error {
	if d.Title == "" {
		return fmt.Errorf("dashboard title must not be empty")
	}
	if len(d.Rows) == 0 {
		return fmt.Errorf("dashboard %q must have at least one row", d.Title)
	}

	// Build variable name set for repeat validation
	varNames := make(map[string]bool, len(d.Variables))
	for i, v := range d.Variables {
		if v.Name == "" {
			return fmt.Errorf("variable[%d] name must not be empty in dashboard %q", i, d.Title)
		}
		if v.Query == "" {
			return fmt.Errorf("variable %q query must not be empty in dashboard %q", v.Name, d.Title)
		}
		varNames[v.Name] = true
	}

	for i, row := range d.Rows {
		if row.Title == "" {
			return fmt.Errorf("row[%d] title must not be empty in dashboard %q", i, d.Title)
		}
		if len(row.Panels) == 0 {
			return fmt.Errorf("row %q must have at least one panel in dashboard %q", row.Title, d.Title)
		}
		if row.Repeat != "" && !varNames[row.Repeat] {
			return fmt.Errorf("row %q repeat variable %q is not defined in dashboard %q", row.Title, row.Repeat, d.Title)
		}

		for j, panel := range row.Panels {
			switch panel.Type {
			case "graph":
				if panel.Query == "" {
					return fmt.Errorf("graph panel[%d] %q in row %q must have a query in dashboard %q", j, panel.Title, row.Title, d.Title)
				}
				if panel.ChartType != "" && !validChartTypes[panel.ChartType] {
					return fmt.Errorf("graph panel[%d] %q in row %q has invalid chart_type %q in dashboard %q", j, panel.Title, row.Title, panel.ChartType, d.Title)
				}
				if panel.Unit != "" && !validUnits[panel.Unit] {
					return fmt.Errorf("graph panel[%d] %q in row %q has invalid unit %q in dashboard %q", j, panel.Title, row.Title, panel.Unit, d.Title)
				}
			case "markdown":
				if panel.Content == "" {
					return fmt.Errorf("markdown panel[%d] %q in row %q must have content in dashboard %q", j, panel.Title, row.Title, d.Title)
				}
			default:
				return fmt.Errorf("panel[%d] %q in row %q has invalid type %q in dashboard %q", j, panel.Title, row.Title, panel.Type, d.Title)
			}
		}
	}

	return nil
}

// DashboardTreeNode represents a node in the hierarchical dashboard navigation tree.
type DashboardTreeNode struct {
	Name      string               `json:"name"`
	Path      string               `json:"path,omitempty"`      // Only set for leaf nodes (actual dashboards)
	Children  []*DashboardTreeNode `json:"children,omitempty"`  // Only set for directory nodes
}
