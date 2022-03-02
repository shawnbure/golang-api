package cache

import (
	"context"
	"sync"
	"time"

	"github.com/ENFT-DAO/youbei-api/config"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/boltdb/bolt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"go.uber.org/atomic"
)

type Stats struct {
	Hits   atomic.Int64
	Misses atomic.Int64
}

type Cacher struct {
	cache *cache.Cache
	stats Stats
	ctx   context.Context
	redis *redis.ClusterClient
	bolt  *bolt.DB
}

var (
	initOnce    sync.Once
	closeOnce   sync.Once
	cacher      *Cacher
	localCacher *LocalCacher

	BoltDbPath = "/tmp/bolt.db"
	log        = logger.GetOrCreate("cacheLog")
)

func InitCacher(cfg config.CacheConfig) {
	initOnce.Do(func() {

		clusterSlots := func(ctx context.Context) ([]redis.ClusterSlot, error) {
			slots := []redis.ClusterSlot{
				// First node with 1 master and 1 slave.
				{
					Start: 0,
					End:   8191,
					Nodes: []redis.ClusterNode{{
						Addr: cfg.WriteUrl, // master
					}, {
						Addr: cfg.ReadUrl, // 1st slave
					}},
				},
				// Second node with 1 master and 1 slave.
				{
					Start: 8192,
					End:   16383,
					Nodes: []redis.ClusterNode{{
						Addr: cfg.WriteUrl, // master
					}, {
						Addr: cfg.ReadUrl, // 1st slave
					}},
				},
			}
			return slots, nil
		}
		redisdb := redis.NewClusterClient(&redis.ClusterOptions{
			ClusterSlots:  clusterSlots,
			ReadOnly:      true,
			RouteRandomly: true,
		})

		// ReloadState reloads cluster state. It calls ClusterSlots func
		// to get cluster slots information.
		// ctx := context.Background()
		// err := redisdb.ReloadState(ctx)
		// if err != nil {
		// 	panic(err)
		// }

		// opt, err := redis.ParseURL(cfg.WriteUrl)
		// if err != nil {
		// 	panic(err)
		// }

		// redisClient := redis.NewClient(opt)
		newCache := cache.New(&cache.Options{
			Redis:      redisdb,
			LocalCache: cache.NewTinyLFU(1000, time.Second),
		})

		boltDb, err := bolt.Open(BoltDbPath, 0600, nil)
		if err != nil {
			panic(err)
		}

		localCacher, err = NewLocalCacher()
		if err != nil {
			panic(err)
		}

		cacher = &Cacher{
			cache: newCache,
			stats: Stats{},
			ctx:   context.Background(),
			redis: redisdb,
			bolt:  boltDb,
		}
	})
}

func CloseCacher() {
	closeOnce.Do(func() {
		err := cacher.bolt.Close()
		if err != nil {
			log.Error("db close", err)
		}
	})
}

func (c *Cacher) Set(k string, v interface{}, ttl time.Duration) error {
	return c.cache.Set(&cache.Item{
		Ctx:   c.ctx,
		Key:   k,
		Value: v,
		TTL:   ttl,
	})
}

func (c *Cacher) Get(k string, v interface{}) error {
	err := c.cache.Get(c.ctx, k, v)

	if err == nil {
		c.stats.Hits.Add(1)
	} else {
		c.stats.Misses.Add(1)
	}

	return err
}

func (c *Cacher) GetStats() Stats {
	return cacher.stats
}

func GetCacher() *Cacher {
	return cacher
}

func GetRedis() *redis.ClusterClient {
	return cacher.redis
}

func GetContext() context.Context {
	return cacher.ctx
}

func GetBolt() *bolt.DB {
	return cacher.bolt
}

func GetLocalCacher() *LocalCacher {
	return localCacher
}
