package persistence

import (
	"testing"
)

func TestFileStore_NewFromNonexistentFile(t *testing.T) {
	store, err := NewFileStore(t.TempDir() + "/state.json")
	if err != nil {
		t.Fatalf("NewFileStore on nonexistent file: %v", err)
	}
	data, _ := store.Load()
	if len(data) != 0 {
		t.Errorf("new store should be empty, got %v", data)
	}
}

func TestFileStore_IncrementAndFlush(t *testing.T) {
	store, _ := NewFileStore(t.TempDir() + "/state.json")

	_ = store.Increment("key1")
	_ = store.Increment("key1")
	_ = store.Increment("key2")

	// In-memory state should be updated immediately
	data, _ := store.Load()
	if data["key1"] != 2 {
		t.Errorf("in-memory key1 = %v, want 2", data["key1"])
	}
	if data["key2"] != 1 {
		t.Errorf("in-memory key2 = %v, want 1", data["key2"])
	}
}

func TestFileStore_FlushPersists(t *testing.T) {
	path := t.TempDir() + "/state.json"
	store, _ := NewFileStore(path)

	_ = store.Increment("k1")
	_ = store.Increment("k1")
	_ = store.Set(KeyLastPolledAt, 1700000000.0)
	_ = store.Flush()

	// Open a fresh store from the same file — simulates restart
	store2, err := NewFileStore(path)
	if err != nil {
		t.Fatalf("reopening store: %v", err)
	}
	data, _ := store2.Load()

	if data["k1"] != 2 {
		t.Errorf("after restart k1 = %v, want 2", data["k1"])
	}
	if data[KeyLastPolledAt] != 1700000000.0 {
		t.Errorf("after restart %s = %v, want 1700000000", KeyLastPolledAt, data[KeyLastPolledAt])
	}
}

func TestFileStore_NoWriteBeforeFlush(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/state.json"
	store, _ := NewFileStore(path)

	_ = store.Increment("k1")

	// File should not exist yet (no Flush called)
	store2, err := NewFileStore(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := store2.Load()
	if data["k1"] != 0 {
		t.Errorf("before flush, new store should see nothing; got %v", data["k1"])
	}

	_ = store.Flush()

	store3, _ := NewFileStore(path)
	data3, _ := store3.Load()
	if data3["k1"] != 1 {
		t.Errorf("after flush k1 = %v, want 1", data3["k1"])
	}
}
