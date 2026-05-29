package persistence

import (
	"testing"
)

func TestLabelKey(t *testing.T) {
	got := LabelKey("myorg/repo", "alice", "feat", "true")
	want := "myorg/repo|alice|feat|true"
	if got != want {
		t.Errorf("LabelKey = %q, want %q", got, want)
	}
}

func TestFromEnv_NoPersistence(t *testing.T) {
	t.Setenv("PERSISTENCE_REDIS_HOST", "")
	t.Setenv("PERSISTENCE_FILE", "")
	store, err := FromEnv()
	if err != nil {
		t.Fatalf("no persistence vars: unexpected error: %v", err)
	}
	if store != nil {
		t.Errorf("no persistence vars: want nil store, got %T", store)
	}
}

func TestFromEnv_BothBackends(t *testing.T) {
	t.Setenv("PERSISTENCE_REDIS_HOST", "localhost:6379")
	t.Setenv("PERSISTENCE_FILE", "/tmp/test.json")
	_, err := FromEnv()
	if err == nil {
		t.Error("both backends configured: expected error")
	}
}

func TestFromEnv_RedisWithoutPassword(t *testing.T) {
	t.Setenv("PERSISTENCE_REDIS_HOST", "localhost:6379")
	t.Setenv("PERSISTENCE_REDIS_PASS", "")
	_, err := FromEnv()
	if err == nil {
		t.Error("redis without password: expected error")
	}
}

func TestFromEnv_FileStore(t *testing.T) {
	t.Setenv("PERSISTENCE_FILE", t.TempDir()+"/state.json")
	t.Setenv("PERSISTENCE_REDIS_HOST", "")
	store, err := FromEnv()
	if err != nil {
		t.Fatalf("file store: unexpected error: %v", err)
	}
	if store == nil {
		t.Fatal("file store: expected non-nil store")
	}
}
