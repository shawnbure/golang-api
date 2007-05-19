package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/erdsea/erdsea-api/config"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type BaseCacher struct {
	cache *cache.Cache

	ctx context.Context
}

func NewBaseCacher(cfg config.CacheConfig) *BaseCacher {
	addrs := map[string]string{}
	for i, addr := range cfg.Addrs {
		k := fmt.Sprintf("cache-%d", i)
		addrs[k] = addr
	}

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: addrs,
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

func (c *BaseCacher) Get(k string, v interface{}) error {
	return c.cache.Get(c.ctx, k, v)
}
