package dashboard

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/tokuhirom/dashyard/internal/model"
	"gopkg.in/yaml.v3"
)

var validPathRe = regexp.MustCompile(`^[a-zA-Z0-9_\-/]+$`)

// Store holds loaded dashboards and provides lookup by path.
type Store struct {
	dashboards map[string]*model.Dashboard
	sources    map[string]string
	list       []*model.Dashboard
	tree       []*model.DashboardTreeNode
}

// LoadDir recursively loads all .yaml files from the given directory
// and returns a Store for looking up dashboards.
func LoadDir(dir string) (*Store, error) {
	store := &Store{
		dashboards: make(map[string]*model.Dashboard),
		sources:    make(map[string]string),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("computing relative path: %w", err)
		}

		// Normalize path: strip extension, use forward slashes
		dashPath := strings.TrimSuffix(relPath, ext)
		dashPath = filepath.ToSlash(dashPath)

		if err := validatePath(dashPath); err != nil {
			return fmt.Errorf("invalid dashboard path %q: %w", dashPath, err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		var d model.Dashboard
		if err := yaml.Unmarshal(data, &d); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
		if err := d.Validate(); err != nil {
			return fmt.Errorf("validating %s: %w", path, err)
		}
		d.Path = dashPath

		store.dashboards[dashPath] = &d
		store.sources[dashPath] = string(data)
		store.list = append(store.list, &d)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("loading dashboards from %s: %w", dir, err)
	}

	// Sort list by path for consistent ordering
	sort.Slice(store.list, func(i, j int) bool {
		return store.list[i].Path < store.list[j].Path
	})

	store.tree = buildTree(store.list)
	return store, nil
}

// Get returns a dashboard by its path, or nil if not found.
func (s *Store) Get(path string) *model.Dashboard {
	return s.dashboards[path]
}

// List returns all dashboards sorted by path.
func (s *Store) List() []*model.Dashboard {
	return s.list
}

// Tree returns the hierarchical dashboard navigation tree.
func (s *Store) Tree() []*model.DashboardTreeNode {
	return s.tree
}

// GetSource returns the raw YAML source for a dashboard by path.
func (s *Store) GetSource(path string) (string, bool) {
	src, ok := s.sources[path]
	return src, ok
}

func validatePath(path string) error {
	if strings.Contains(path, "..") {
		return fmt.Errorf("path must not contain '..'")
	}
	if !validPathRe.MatchString(path) {
		return fmt.Errorf("path contains invalid characters")
	}
	return nil
}

func buildTree(dashboards []*model.Dashboard) []*model.DashboardTreeNode {
	root := &model.DashboardTreeNode{}

	for _, d := range dashboards {
		parts := strings.Split(d.Path, "/")
		current := root

		for i, part := range parts {
			if i == len(parts)-1 {
				// Leaf node
				current.Children = append(current.Children, &model.DashboardTreeNode{
					Name: part,
					Path: d.Path,
				})
			} else {
				// Find or create directory node
				var found *model.DashboardTreeNode
				for _, child := range current.Children {
					if child.Name == part && child.Path == "" {
						found = child
						break
					}
				}
				if found == nil {
					found = &model.DashboardTreeNode{Name: part}
					current.Children = append(current.Children, found)
				}
				current = found
			}
		}
	}

	return root.Children
}
