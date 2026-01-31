package dashboard

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const validDashboardYAML = `title: Test Dashboard
rows:
  - title: Test Row
    panels:
      - title: Test Panel
        type: markdown
        content: hello
`

// waitFor polls condition every 50ms up to timeout, returning true if met.
func waitFor(t *testing.T, timeout time.Duration, condition func() bool) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

func TestWatcherReloadOnCreate(t *testing.T) {
	dir := t.TempDir()

	store, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	holder := NewStoreHolder(store)

	w := NewWatcher(dir, holder)
	w.debounce = 100 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := w.Watch(ctx); err != nil {
			t.Errorf("watcher error: %v", err)
		}
	}()

	// Give the watcher time to start.
	time.Sleep(200 * time.Millisecond)

	// Create a new dashboard file.
	if err := os.WriteFile(filepath.Join(dir, "new.yaml"), []byte(validDashboardYAML), 0644); err != nil {
		t.Fatal(err)
	}

	if !waitFor(t, 3*time.Second, func() bool {
		return holder.Store().Get("new") != nil
	}) {
		t.Error("expected 'new' dashboard after file creation")
	}
}

func TestWatcherReloadOnDelete(t *testing.T) {
	dir := t.TempDir()

	// Pre-create a dashboard.
	if err := os.WriteFile(filepath.Join(dir, "existing.yaml"), []byte(validDashboardYAML), 0644); err != nil {
		t.Fatal(err)
	}

	store, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if store.Get("existing") == nil {
		t.Fatal("expected 'existing' dashboard")
	}

	holder := NewStoreHolder(store)
	w := NewWatcher(dir, holder)
	w.debounce = 100 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := w.Watch(ctx); err != nil {
			t.Errorf("watcher error: %v", err)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	// Remove the file.
	if err := os.Remove(filepath.Join(dir, "existing.yaml")); err != nil {
		t.Fatal(err)
	}

	if !waitFor(t, 3*time.Second, func() bool {
		return holder.Store().Get("existing") == nil
	}) {
		t.Error("expected 'existing' dashboard to disappear after deletion")
	}
}

func TestWatcherKeepsOldStoreOnError(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "good.yaml"), []byte(validDashboardYAML), 0644); err != nil {
		t.Fatal(err)
	}

	store, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	holder := NewStoreHolder(store)
	w := NewWatcher(dir, holder)
	w.debounce = 100 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := w.Watch(ctx); err != nil {
			t.Errorf("watcher error: %v", err)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	// Write an invalid YAML file.
	if err := os.WriteFile(filepath.Join(dir, "bad.yaml"), []byte("{invalid"), 0644); err != nil {
		t.Fatal(err)
	}

	// Wait for debounce + reload attempt.
	time.Sleep(500 * time.Millisecond)

	// Old data should still be available.
	if holder.Store().Get("good") == nil {
		t.Error("expected 'good' dashboard to survive a failed reload")
	}
}

func TestWatcherSubdirectory(t *testing.T) {
	dir := t.TempDir()

	store, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	holder := NewStoreHolder(store)

	w := NewWatcher(dir, holder)
	w.debounce = 100 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := w.Watch(ctx); err != nil {
			t.Errorf("watcher error: %v", err)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	// Create a subdirectory and add a file to it.
	subdir := filepath.Join(dir, "infra")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	// Small delay to let watcher pick up the new directory.
	time.Sleep(200 * time.Millisecond)

	if err := os.WriteFile(filepath.Join(subdir, "sub.yaml"), []byte(validDashboardYAML), 0644); err != nil {
		t.Fatal(err)
	}

	if !waitFor(t, 3*time.Second, func() bool {
		return holder.Store().Get("infra/sub") != nil
	}) {
		t.Error("expected 'infra/sub' dashboard after subdirectory file creation")
	}
}
