package dashboard

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDir(t *testing.T) {
	// Use the testdata directory from the project root
	store, err := LoadDir("../../testdata/dashboards")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list := store.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 dashboards, got %d", len(list))
	}

	// Should be sorted by path
	if list[0].Path != "infra/network" {
		t.Errorf("expected first dashboard path 'infra/network', got %q", list[0].Path)
	}
	if list[1].Path != "overview" {
		t.Errorf("expected second dashboard path 'overview', got %q", list[1].Path)
	}

	// Test Get
	d := store.Get("overview")
	if d == nil {
		t.Fatal("expected to find 'overview' dashboard")
	}
	if d.Title != "System Overview" {
		t.Errorf("expected title 'System Overview', got %q", d.Title)
	}
	if len(d.Rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(d.Rows))
	}

	d = store.Get("infra/network")
	if d == nil {
		t.Fatal("expected to find 'infra/network' dashboard")
	}
	if d.Title != "Network" {
		t.Errorf("expected title 'Network', got %q", d.Title)
	}

	// Test Get non-existent
	if store.Get("nonexistent") != nil {
		t.Error("expected nil for nonexistent dashboard")
	}

	// Test Tree
	tree := store.Tree()
	if len(tree) != 2 {
		t.Fatalf("expected 2 top-level tree nodes, got %d", len(tree))
	}
}

func TestLoadDirEmpty(t *testing.T) {
	dir := t.TempDir()
	store, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.List()) != 0 {
		t.Errorf("expected 0 dashboards, got %d", len(store.List()))
	}
}

func TestLoadDirNonExistent(t *testing.T) {
	_, err := LoadDir("/nonexistent/path")
	if err == nil {
		t.Error("expected error for non-existent directory")
	}
}

func TestLoadDirInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bad.yaml"), []byte("{invalid"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadDir(dir)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		path    string
		wantErr bool
	}{
		{"overview", false},
		{"infra/network", false},
		{"my-dash_1", false},
		{"../escape", true},
		{"path/../escape", true},
		{"invalid chars!", true},
		{"spaces are bad", true},
	}

	for _, tt := range tests {
		err := validatePath(tt.path)
		if (err != nil) != tt.wantErr {
			t.Errorf("validatePath(%q): got err=%v, wantErr=%v", tt.path, err, tt.wantErr)
		}
	}
}

func TestBuildTree(t *testing.T) {
	store, err := LoadDir("../../testdata/dashboards")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tree := store.Tree()

	// Should have "infra" directory and "overview" leaf
	var infraNode, overviewNode *struct{ found bool }
	for _, node := range tree {
		switch node.Name {
		case "infra":
			infraNode = &struct{ found bool }{true}
			if node.Path != "" {
				t.Error("directory node should not have a path")
			}
			if len(node.Children) != 1 {
				t.Errorf("expected 1 child under infra, got %d", len(node.Children))
			} else if node.Children[0].Name != "network" {
				t.Errorf("expected child 'network', got %q", node.Children[0].Name)
			} else if node.Children[0].Path != "infra/network" {
				t.Errorf("expected path 'infra/network', got %q", node.Children[0].Path)
			}
		case "overview":
			overviewNode = &struct{ found bool }{true}
			if node.Path != "overview" {
				t.Errorf("expected path 'overview', got %q", node.Path)
			}
		}
	}
	if infraNode == nil {
		t.Error("missing 'infra' node in tree")
	}
	if overviewNode == nil {
		t.Error("missing 'overview' node in tree")
	}
}
