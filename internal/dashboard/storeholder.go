package dashboard

import "sync/atomic"

// StoreHolder provides lock-free concurrent access to a Store that can be
// atomically replaced (e.g. on dashboard file changes).
type StoreHolder struct {
	p atomic.Pointer[Store]
}

// NewStoreHolder creates a StoreHolder initialised with the given Store.
func NewStoreHolder(s *Store) *StoreHolder {
	h := &StoreHolder{}
	h.p.Store(s)
	return h
}

// Store returns the current Store. Safe for concurrent use.
func (h *StoreHolder) Store() *Store {
	return h.p.Load()
}

// Replace atomically swaps the current Store with a new one.
func (h *StoreHolder) Replace(s *Store) {
	h.p.Store(s)
}
