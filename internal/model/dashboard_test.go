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

func TestPanelLegendYAMLUnmarshal(t *testing.T) {
	input := `
title: "Bytes Received"
type: "graph"
query: 'rate(system_network_io_bytes_total{direction="receive"}[5m])'
unit: "bytes"
legend: "{device} {direction}"
`
	var p Panel
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Legend != "{device} {direction}" {
		t.Errorf("expected legend '{device} {direction}', got %q", p.Legend)
	}
}

func TestPanelLegendOmittedYAML(t *testing.T) {
	input := `
title: "CPU Usage"
type: "graph"
query: "rate(cpu[5m])"
`
	var p Panel
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Legend != "" {
		t.Errorf("expected empty legend when omitted, got %q", p.Legend)
	}
}

func TestPanelLegendJSON(t *testing.T) {
	p := Panel{Title: "Test", Type: "graph", Query: "up", Legend: "{instance}"}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decoded["legend"] != "{instance}" {
		t.Errorf("expected legend '{instance}', got %v", decoded["legend"])
	}

	// Verify legend is omitted from JSON when empty
	p2 := Panel{Title: "Test", Type: "graph", Query: "up"}
	data2, err := json.Marshal(p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded2 map[string]interface{}
	if err := json.Unmarshal(data2, &decoded2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := decoded2["legend"]; ok {
		t.Errorf("expected legend to be omitted from JSON when empty")
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

func TestPanelChartTypeYAMLUnmarshal(t *testing.T) {
	input := `
title: "CPU Bar"
type: "graph"
chart_type: "bar"
query: "rate(cpu[5m])"
unit: "percent"
`
	var p Panel
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ChartType != "bar" {
		t.Errorf("expected chart_type 'bar', got %q", p.ChartType)
	}
}

func TestPanelChartTypeOmittedYAML(t *testing.T) {
	input := `
title: "CPU Usage"
type: "graph"
query: "rate(cpu[5m])"
`
	var p Panel
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ChartType != "" {
		t.Errorf("expected empty chart_type when omitted, got %q", p.ChartType)
	}
}

func TestPanelChartTypeJSON(t *testing.T) {
	p := Panel{Title: "Test", Type: "graph", Query: "up", ChartType: "scatter"}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decoded["chart_type"] != "scatter" {
		t.Errorf("expected chart_type 'scatter', got %v", decoded["chart_type"])
	}

	// Verify chart_type is omitted from JSON when empty
	p2 := Panel{Title: "Test", Type: "graph", Query: "up"}
	data2, err := json.Marshal(p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded2 map[string]interface{}
	if err := json.Unmarshal(data2, &decoded2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := decoded2["chart_type"]; ok {
		t.Errorf("expected chart_type to be omitted from JSON when empty")
	}
}

func TestDashboardWithVariablesYAML(t *testing.T) {
	input := `
title: "Network by Interface"
variables:
  - name: device
    label: "Network Device"
    query: "label_values(system_network_io_bytes_total, device)"
rows:
  - title: "Traffic for $device"
    repeat: device
    panels:
      - title: "Bytes Received"
        type: "graph"
        query: 'rate(system_network_io_bytes_total{device="$device"}[5m])'
        unit: "bytes"
`
	var d Dashboard
	if err := yaml.Unmarshal([]byte(input), &d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.Variables) != 1 {
		t.Fatalf("expected 1 variable, got %d", len(d.Variables))
	}
	v := d.Variables[0]
	if v.Name != "device" {
		t.Errorf("expected variable name 'device', got %q", v.Name)
	}
	if v.Label != "Network Device" {
		t.Errorf("expected variable label 'Network Device', got %q", v.Label)
	}
	if v.Query != "label_values(system_network_io_bytes_total, device)" {
		t.Errorf("expected variable query, got %q", v.Query)
	}
	if d.Rows[0].Repeat != "device" {
		t.Errorf("expected row repeat 'device', got %q", d.Rows[0].Repeat)
	}
}

func TestDashboardVariablesJSON(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Variables: []Variable{
			{Name: "device", Label: "Device", Query: "label_values(m, device)"},
		},
		Rows: []Row{
			{Title: "Row1", Repeat: "device", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up"}}},
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
	vars, ok := decoded["variables"].([]interface{})
	if !ok || len(vars) != 1 {
		t.Fatalf("expected 1 variable in JSON, got %v", decoded["variables"])
	}
	varMap := vars[0].(map[string]interface{})
	if varMap["name"] != "device" {
		t.Errorf("expected variable name 'device', got %v", varMap["name"])
	}
	rows := decoded["rows"].([]interface{})
	rowMap := rows[0].(map[string]interface{})
	if rowMap["repeat"] != "device" {
		t.Errorf("expected row repeat 'device', got %v", rowMap["repeat"])
	}
}

func TestDashboardVariablesOmittedJSON(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Rows: []Row{
			{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up"}}},
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
	if _, ok := decoded["variables"]; ok {
		t.Errorf("expected variables to be omitted from JSON when nil")
	}
	rows := decoded["rows"].([]interface{})
	rowMap := rows[0].(map[string]interface{})
	if _, ok := rowMap["repeat"]; ok {
		t.Errorf("expected repeat to be omitted from JSON when empty")
	}
}

func TestPanelThresholdsYAMLUnmarshal(t *testing.T) {
	input := `
title: "CPU Usage"
type: "graph"
query: "rate(cpu[5m])"
unit: "percent"
thresholds:
  - value: 80
    color: "#f59e0b"
    label: "Warning"
  - value: 95
    color: "#ef4444"
    label: "Critical"
`
	var p Panel
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Thresholds) != 2 {
		t.Fatalf("expected 2 thresholds, got %d", len(p.Thresholds))
	}
	if p.Thresholds[0].Value != 80 {
		t.Errorf("expected threshold value 80, got %v", p.Thresholds[0].Value)
	}
	if p.Thresholds[0].Color != "#f59e0b" {
		t.Errorf("expected threshold color '#f59e0b', got %q", p.Thresholds[0].Color)
	}
	if p.Thresholds[0].Label != "Warning" {
		t.Errorf("expected threshold label 'Warning', got %q", p.Thresholds[0].Label)
	}
	if p.Thresholds[1].Value != 95 {
		t.Errorf("expected threshold value 95, got %v", p.Thresholds[1].Value)
	}
}

func TestPanelThresholdValueOnlyYAML(t *testing.T) {
	input := `
title: "CPU Usage"
type: "graph"
query: "rate(cpu[5m])"
thresholds:
  - value: 50
`
	var p Panel
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Thresholds) != 1 {
		t.Fatalf("expected 1 threshold, got %d", len(p.Thresholds))
	}
	if p.Thresholds[0].Value != 50 {
		t.Errorf("expected threshold value 50, got %v", p.Thresholds[0].Value)
	}
	if p.Thresholds[0].Color != "" {
		t.Errorf("expected empty color when omitted, got %q", p.Thresholds[0].Color)
	}
	if p.Thresholds[0].Label != "" {
		t.Errorf("expected empty label when omitted, got %q", p.Thresholds[0].Label)
	}
}

func TestPanelThresholdsJSON(t *testing.T) {
	p := Panel{
		Title: "Test",
		Type:  "graph",
		Query: "up",
		Thresholds: []Threshold{
			{Value: 80, Color: "#f59e0b", Label: "Warning"},
		},
	}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	thresholds, ok := decoded["thresholds"].([]interface{})
	if !ok || len(thresholds) != 1 {
		t.Fatalf("expected 1 threshold in JSON, got %v", decoded["thresholds"])
	}
	th := thresholds[0].(map[string]interface{})
	if th["value"] != float64(80) {
		t.Errorf("expected threshold value 80, got %v", th["value"])
	}
	if th["color"] != "#f59e0b" {
		t.Errorf("expected threshold color '#f59e0b', got %v", th["color"])
	}

	// Verify thresholds is omitted from JSON when nil
	p2 := Panel{Title: "Test", Type: "graph", Query: "up"}
	data2, err := json.Marshal(p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded2 map[string]interface{}
	if err := json.Unmarshal(data2, &decoded2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := decoded2["thresholds"]; ok {
		t.Errorf("expected thresholds to be omitted from JSON when nil")
	}
}

func TestPanelStackedYAMLUnmarshal(t *testing.T) {
	input := `
title: "Memory Stacked"
type: "graph"
chart_type: "area"
stacked: true
query: "system_memory_usage_bytes"
unit: "bytes"
`
	var p Panel
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.Stacked {
		t.Errorf("expected stacked to be true")
	}
}

func TestPanelStackedOmittedYAML(t *testing.T) {
	input := `
title: "CPU Usage"
type: "graph"
query: "rate(cpu[5m])"
`
	var p Panel
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Stacked {
		t.Errorf("expected stacked to be false when omitted")
	}
}

func TestPanelStackedJSON(t *testing.T) {
	p := Panel{Title: "Test", Type: "graph", Query: "up", Stacked: true}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decoded["stacked"] != true {
		t.Errorf("expected stacked true, got %v", decoded["stacked"])
	}

	// Verify stacked is omitted from JSON when false
	p2 := Panel{Title: "Test", Type: "graph", Query: "up"}
	data2, err := json.Marshal(p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded2 map[string]interface{}
	if err := json.Unmarshal(data2, &decoded2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := decoded2["stacked"]; ok {
		t.Errorf("expected stacked to be omitted from JSON when false")
	}
}

func TestValidateValidDashboard(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Rows: []Row{
			{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up"}}},
		},
	}
	if err := d.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateEmptyTitle(t *testing.T) {
	d := Dashboard{
		Title: "",
		Rows:  []Row{{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up"}}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for empty title")
	}
}

func TestValidateNoRows(t *testing.T) {
	d := Dashboard{Title: "Test", Rows: []Row{}}
	if err := d.Validate(); err == nil {
		t.Error("expected error for no rows")
	}
}

func TestValidateRowEmptyTitle(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Rows:  []Row{{Title: "", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up"}}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for empty row title")
	}
}

func TestValidateRowNoPanels(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Rows:  []Row{{Title: "Row1", Panels: []Panel{}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for row with no panels")
	}
}

func TestValidateInvalidPanelType(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Rows:  []Row{{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "unknown"}}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for invalid panel type")
	}
}

func TestValidateGraphPanelNoQuery(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Rows:  []Row{{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "graph", Query: ""}}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for graph panel without query")
	}
}

func TestValidateGraphPanelInvalidChartType(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Rows:  []Row{{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up", ChartType: "invalid"}}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for invalid chart_type")
	}
}

func TestValidateGraphPanelValidChartTypes(t *testing.T) {
	for _, ct := range []string{"line", "bar", "area", "scatter", "pie", "doughnut"} {
		d := Dashboard{
			Title: "Test",
			Rows:  []Row{{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up", ChartType: ct}}}},
		}
		if err := d.Validate(); err != nil {
			t.Errorf("expected no error for chart_type %q, got %v", ct, err)
		}
	}
}

func TestValidateGraphPanelInvalidUnit(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Rows:  []Row{{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up", Unit: "invalid"}}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for invalid unit")
	}
}

func TestValidateGraphPanelValidUnits(t *testing.T) {
	for _, u := range []string{"bytes", "percent", "count"} {
		d := Dashboard{
			Title: "Test",
			Rows:  []Row{{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up", Unit: u}}}},
		}
		if err := d.Validate(); err != nil {
			t.Errorf("expected no error for unit %q, got %v", u, err)
		}
	}
}

func TestValidateMarkdownPanelNoContent(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Rows:  []Row{{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "markdown", Content: ""}}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for markdown panel without content")
	}
}

func TestValidateVariableEmptyName(t *testing.T) {
	d := Dashboard{
		Title:     "Test",
		Variables: []Variable{{Name: "", Query: "label_values(m, x)"}},
		Rows:      []Row{{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up"}}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for variable with empty name")
	}
}

func TestValidateVariableEmptyQuery(t *testing.T) {
	d := Dashboard{
		Title:     "Test",
		Variables: []Variable{{Name: "x", Query: ""}},
		Rows:      []Row{{Title: "Row1", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up"}}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for variable with empty query")
	}
}

func TestValidateRepeatUndefinedVariable(t *testing.T) {
	d := Dashboard{
		Title: "Test",
		Rows:  []Row{{Title: "Row1", Repeat: "missing", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up"}}}},
	}
	if err := d.Validate(); err == nil {
		t.Error("expected error for repeat referencing undefined variable")
	}
}

func TestValidateRepeatDefinedVariable(t *testing.T) {
	d := Dashboard{
		Title:     "Test",
		Variables: []Variable{{Name: "device", Query: "label_values(m, device)"}},
		Rows:      []Row{{Title: "Row1", Repeat: "device", Panels: []Panel{{Title: "P1", Type: "graph", Query: "up"}}}},
	}
	if err := d.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
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
