package persistence

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const redisHashKey = "gitstats:commits"

type redisStore struct {
	client *redis.Client
}

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

func (s *redisStore) Increment(key string) error {
	return s.client.HIncrBy(context.Background(), redisHashKey, key, 1).Err()
}
