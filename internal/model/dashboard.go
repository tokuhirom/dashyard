package model

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
	Unit       string      `yaml:"unit,omitempty" json:"unit,omitempty"` // "bytes", "percent", "count"
	YMin       *float64    `yaml:"y_min,omitempty" json:"y_min,omitempty"`
	YMax       *float64    `yaml:"y_max,omitempty" json:"y_max,omitempty"`
	Legend     string      `yaml:"legend,omitempty" json:"legend,omitempty"`
	Thresholds []Threshold `yaml:"thresholds,omitempty" json:"thresholds,omitempty"`
	Stacked    bool        `yaml:"stacked,omitempty" json:"stacked,omitempty"`
	Content    string      `yaml:"content,omitempty" json:"content,omitempty"`
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

// DashboardTreeNode represents a node in the hierarchical dashboard navigation tree.
type DashboardTreeNode struct {
	Name      string               `json:"name"`
	Path      string               `json:"path,omitempty"`      // Only set for leaf nodes (actual dashboards)
	Children  []*DashboardTreeNode `json:"children,omitempty"`  // Only set for directory nodes
}
