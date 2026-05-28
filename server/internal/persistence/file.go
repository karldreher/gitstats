package persistence

import (
	"encoding/json"
	"os"
	"sync"
)

type fileStore struct {
	path string
	mu   sync.Mutex
	data map[string]float64
}

func NewFileStore(path string) (StateStore, error) {
	s := &fileStore{path: path, data: make(map[string]float64)}
	if err := s.loadFromDisk(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

func (s *fileStore) Load() (map[string]float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]float64, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out, nil
}

func (s *fileStore) Increment(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key]++
	return s.writeToDisk()
}

func (s *fileStore) loadFromDisk() error {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}

func (s *fileStore) writeToDisk() error {
	b, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0644)
}
