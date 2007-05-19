package cache

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"time"

	"github.com/erdsea/erdsea-api/config"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type Stats struct {
	Hits   atomic.Int64
	Misses atomic.Int64
}

type BaseCacher struct {
	cache *cache.Cache
	stats Stats
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
		stats: Stats{},
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
	err := c.cache.Get(c.ctx, k, v)

	if err == nil {
		c.stats.Hits.Add(1)
	} else {
		c.stats.Misses.Add(1)
	}

	return err
}
