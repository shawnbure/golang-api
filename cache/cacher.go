package cache

import (
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type BaseCacher struct {
	cacher *cache.Cache
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
		cacher: cacher,
	}
}
