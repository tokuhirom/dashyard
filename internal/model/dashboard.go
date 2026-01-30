package model

// Panel represents a single visualization panel within a dashboard row.
type Panel struct {
	Title   string `yaml:"title" json:"title"`
	Type    string `yaml:"type" json:"type"`       // "graph" or "markdown"
	Query   string `yaml:"query,omitempty" json:"query,omitempty"`
	Unit    string `yaml:"unit,omitempty" json:"unit,omitempty"` // "bytes", "percent", "count"
	Content string `yaml:"content,omitempty" json:"content,omitempty"`
}

// Row represents a horizontal row of panels in a dashboard.
type Row struct {
	Title  string  `yaml:"title" json:"title"`
	Panels []Panel `yaml:"panels" json:"panels"`
}

// Dashboard represents a single dashboard definition loaded from YAML.
type Dashboard struct {
	Title string `yaml:"title" json:"title"`
	Rows  []Row  `yaml:"rows" json:"rows"`
	Path  string `yaml:"-" json:"path"` // Set by loader, not from YAML
}

// DashboardTreeNode represents a node in the hierarchical dashboard navigation tree.
type DashboardTreeNode struct {
	Name      string               `json:"name"`
	Path      string               `json:"path,omitempty"`      // Only set for leaf nodes (actual dashboards)
	Children  []*DashboardTreeNode `json:"children,omitempty"`  // Only set for directory nodes
}
