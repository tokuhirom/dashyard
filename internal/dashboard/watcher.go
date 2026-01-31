package dashboard

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches a dashboard directory for file changes and hot-reloads
// the Store when dashboards are added, modified, or removed.
type Watcher struct {
	dir      string
	holder   *StoreHolder
	debounce time.Duration
}

// NewWatcher creates a Watcher that will reload dashboards from dir into holder.
func NewWatcher(dir string, holder *StoreHolder) *Watcher {
	return &Watcher{
		dir:      dir,
		holder:   holder,
		debounce: 500 * time.Millisecond,
	}
}

// Watch blocks until ctx is cancelled, reloading dashboards on file changes.
func (w *Watcher) Watch(ctx context.Context) error {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer func() { _ = fsw.Close() }()

	// Add the root dir and all subdirectories.
	if err := w.addDirs(fsw, w.dir); err != nil {
		return err
	}

	slog.Info("watching dashboards for changes", "dir", w.dir)

	var timer *time.Timer
	var timerC <-chan time.Time

	for {
		select {
		case <-ctx.Done():
			if timer != nil {
				timer.Stop()
			}
			return nil

		case event, ok := <-fsw.Events:
			if !ok {
				return nil
			}

			if !w.relevant(event) {
				continue
			}

			// If a new directory was created, start watching it.
			if event.Has(fsnotify.Create) {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					_ = fsw.Add(event.Name)
				}
			}

			// Reset debounce timer.
			if timer == nil {
				timer = time.NewTimer(w.debounce)
				timerC = timer.C
			} else {
				timer.Reset(w.debounce)
			}

		case err, ok := <-fsw.Errors:
			if !ok {
				return nil
			}
			slog.Error("filesystem watcher error", "error", err)

		case <-timerC:
			timer = nil
			timerC = nil
			w.reload()
		}
	}
}

// relevant returns true if the event is for a YAML file or a directory.
func (w *Watcher) relevant(event fsnotify.Event) bool {
	// Always handle directory events (Create for new subdirs).
	if event.Has(fsnotify.Create) {
		if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
			return true
		}
	}

	ext := strings.ToLower(filepath.Ext(event.Name))
	return ext == ".yaml" || ext == ".yml"
}

// reload calls LoadDir and swaps the store on success.
func (w *Watcher) reload() {
	store, err := LoadDir(w.dir)
	if err != nil {
		slog.Error("failed to reload dashboards, keeping old data", "error", err)
		return
	}
	w.holder.Replace(store)
	slog.Info("reloaded dashboards", "count", len(store.List()))
}

// addDirs recursively adds dir and all subdirectories to the watcher.
func (w *Watcher) addDirs(fsw *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return fsw.Add(path)
		}
		return nil
	})
}
