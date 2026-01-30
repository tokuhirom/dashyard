package model

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestPanelYAMLUnmarshal(t *testing.T) {
	input := `
title: "CPU Usage"
type: "graph"
query: "rate(cpu[5m])"
unit: "percent"
`
	var p Panel
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Title != "CPU Usage" {
		t.Errorf("expected title 'CPU Usage', got %q", p.Title)
	}
	if p.Type != "graph" {
		t.Errorf("expected type 'graph', got %q", p.Type)
	}
	if p.Query != "rate(cpu[5m])" {
		t.Errorf("expected query 'rate(cpu[5m])', got %q", p.Query)
	}
	if p.Unit != "percent" {
		t.Errorf("expected unit 'percent', got %q", p.Unit)
	}
}

func TestPanelMarkdownYAML(t *testing.T) {
	input := `
title: "Notes"
type: "markdown"
content: "## Hello\nWorld"
`
	var p Panel
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Type != "markdown" {
		t.Errorf("expected type 'markdown', got %q", p.Type)
	}
	if p.Content != "## Hello\nWorld" {
		t.Errorf("unexpected content: %q", p.Content)
	}
}

func TestDashboardYAMLUnmarshal(t *testing.T) {
	input := `
title: "Overview"
rows:
  - title: "CPU"
    panels:
      - title: "CPU Usage"
        type: "graph"
        query: "rate(cpu[5m])"
        unit: "percent"
`
	var d Dashboard
	if err := yaml.Unmarshal([]byte(input), &d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Title != "Overview" {
		t.Errorf("expected title 'Overview', got %q", d.Title)
	}
	if len(d.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(d.Rows))
	}
	if d.Rows[0].Title != "CPU" {
		t.Errorf("expected row title 'CPU', got %q", d.Rows[0].Title)
	}
	if len(d.Rows[0].Panels) != 1 {
		t.Fatalf("expected 1 panel, got %d", len(d.Rows[0].Panels))
	}
}

func TestDashboardJSON(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Path:  "test/dash",
		Rows: []Row{
			{
				Title: "Row1",
				Panels: []Panel{
					{Title: "P1", Type: "graph", Query: "up", Unit: "count"},
				},
			},
		},
	}
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decoded["path"] != "test/dash" {
		t.Errorf("expected path 'test/dash', got %v", decoded["path"])
	}
}

func TestDashboardTreeNodeJSON(t *testing.T) {
	node := DashboardTreeNode{
		Name: "infra",
		Children: []*DashboardTreeNode{
			{Name: "network", Path: "infra/network"},
		},
	}
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decoded["name"] != "infra" {
		t.Errorf("expected name 'infra', got %v", decoded["name"])
	}
	children := decoded["children"].([]interface{})
	if len(children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(children))
	}
}
