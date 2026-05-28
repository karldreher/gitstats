package persistence

import (
	"errors"
	"os"
)

// StateStore persists and restores counter state across restarts.
type StateStore interface {
	// Load returns all saved label-key → count pairs.
	Load() (map[string]float64, error)
	// Increment atomically adds 1 to the given key.
	Increment(key string) error
}

// LabelKey builds a deterministic composite key from label values.
func LabelKey(repo, author, commitType, conventional string) string {
	return repo + "|" + author + "|" + commitType + "|" + conventional
}

// FromEnv constructs the appropriate StateStore from environment variables.
// Returns nil, nil if no persistence is configured.
// Returns an error if both backends are configured or Redis is configured without a password.
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
