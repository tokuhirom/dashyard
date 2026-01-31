package dashboard

import (
	"testing"
)

func TestStoreHolderSwap(t *testing.T) {
	store1, err := LoadDir("../../testdata/dashboards")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	holder := NewStoreHolder(store1)

	// Verify initial store is accessible.
	got := holder.Store()
	if got != store1 {
		t.Fatal("expected holder to return store1")
	}
	if got.Get("overview") == nil {
		t.Fatal("expected 'overview' dashboard in store1")
	}

	// Create a second (empty) store and swap it in.
	store2, err := LoadDir(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	holder.Replace(store2)

	got = holder.Store()
	if got != store2 {
		t.Fatal("expected holder to return store2 after Replace")
	}
	if got.Get("overview") != nil {
		t.Error("expected nil for 'overview' in empty store2")
	}
}
