// Package persistence defines the StateStore interface and its file and Redis implementations.
package persistence

import (
	"errors"
	"os"
)

// KeyLastPolledAt is the special store key tracking the last successful poll timestamp.
const KeyLastPolledAt = "__last_polled_at"

// StateStore persists and restores counter state across restarts.
type StateStore interface {
	// Load returns all saved key → value pairs (counters + internal keys).
	Load() (map[string]float64, error)
	// Increment adds 1 to the given key. Memory-only for file; durable for Redis.
	Increment(key string) error
	// Set stores an arbitrary value. Memory-only for file; durable for Redis.
	Set(key string, value float64) error
	// Flush writes buffered state to durable storage. No-op for Redis.
	Flush() error
}

// LabelKey builds a deterministic composite key from label values.
func LabelKey(repo, author, commitType, conventional string) string {
	return repo + "|" + author + "|" + commitType + "|" + conventional
}

// FromEnv constructs the appropriate StateStore from environment variables.
// Returns nil, nil if no persistence is configured.
func FromEnv() (StateStore, error) {
	redisHost := os.Getenv("PERSISTENCE_REDIS_HOST")
	filePath := os.Getenv("PERSISTENCE_FILE")

	if redisHost != "" && filePath != "" {
		return nil, errors.New("only one persistence backend may be configured: unset PERSISTENCE_REDIS_HOST or PERSISTENCE_FILE")
	}
	if redisHost != "" {
		pass := os.Getenv("PERSISTENCE_REDIS_PASS")
		if pass == "" {
			return nil, errors.New("PERSISTENCE_REDIS_PASS is required when PERSISTENCE_REDIS_HOST is set")
		}
		return NewRedisStore(redisHost, pass)
	}
	if filePath != "" {
		return NewFileStore(filePath)
	}
	return nil, nil
}
