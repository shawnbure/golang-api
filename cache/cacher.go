package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type BaseCacher struct {
	cache *cache.Cache

	ctx context.Context
}

func NewBaseCacher() *BaseCacher {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{},
	})

	cacher := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Second),
	})

	return &BaseCacher{
		cache: cacher,
		ctx:   context.Background(),
	}
}

func (c *BaseCacher) Set(k string, v interface{}, ttl time.Duration) error {
	return c.cache.Set(&cache.Item{
		Ctx:   c.ctx,
		Key:   k,
		Value: v,
		TTL:   ttl,
	})
}

func (c *BaseCacher) Get(k string, v interface{}) (interface{}, error) {
	err := c.cache.Get(c.ctx, k, &v)

	return v, err
}
