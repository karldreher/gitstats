package persistence

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const redisHashKey = "gitstats:commits"

type redisStore struct {
	client *redis.Client
}

// NewRedisStore returns a StateStore backed by a Redis hash at host, authenticated with password.
func NewRedisStore(host, password string) (StateStore, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
	})
	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &redisStore{client: c}, nil
}

func (s *redisStore) Load() (map[string]float64, error) {
	vals, err := s.client.HGetAll(context.Background(), redisHashKey).Result()
	if err != nil {
		return nil, err
	}
	out := make(map[string]float64, len(vals))
	for k, v := range vals {
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			out[k] = f
		}
	}
	return out, nil
}

// Increment atomically increments the key using HINCRBY — immediately durable.
func (s *redisStore) Increment(key string) error {
	return s.client.HIncrBy(context.Background(), redisHashKey, key, 1).Err()
}

// Set stores a value using HSET — immediately durable.
func (s *redisStore) Set(key string, value float64) error {
	return s.client.HSet(context.Background(), redisHashKey, key, fmt.Sprintf("%g", value)).Err()
}

// Flush is a no-op for Redis; all writes are immediately durable.
func (s *redisStore) Flush() error { return nil }
