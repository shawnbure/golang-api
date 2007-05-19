package cache

import (
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
)

const (
	defaultCostFuncFlag = 1

	numCounters = 1e7
	maxCost     = 1 << 30
	bufferItems = 64
	metrics     = false
)

var (
	ErrCouldNotSetWithTTL = errors.New("could not set entry with TTL")
	ErrNoEntryFoundForKey = errors.New("no entry found for key")
	ErrCouldNotSet        = errors.New("could not set entry")
)

type LocalCacher struct {
	cache *ristretto.Cache
}

func NewLocalCacher() (*LocalCacher, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: bufferItems,
		Metrics:     metrics,
	})
	if err != nil {
		return nil, err
	}

	return &LocalCacher{
		cache: cache,
	}, nil
}

func (lc *LocalCacher) SetWithTTLSync(key string, value interface{}, ttl time.Duration) error {
	err := lc.SetWithTTL(key, value, ttl)
	if err != nil {
		return err
	}

	lc.cache.Wait()
	return nil
}

func (lc *LocalCacher) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	ok := lc.cache.SetWithTTL(key, value, defaultCostFuncFlag, ttl)
	if !ok {
		return ErrCouldNotSetWithTTL
	}

	return nil
}

func (lc *LocalCacher) SetSync(key string, value interface{}) error {
	err := lc.Set(key, value)
	if err != nil {
		return err
	}

	lc.cache.Wait()
	return nil
}

func (lc *LocalCacher) Set(key string, value interface{}) error {
	ok := lc.cache.Set(key, value, defaultCostFuncFlag)
	if !ok {
		return ErrCouldNotSet
	}

	return nil
}

func (lc *LocalCacher) Get(key string) (interface{}, error) {
	val, found := lc.cache.Get(key)
	if !found {
		return nil, ErrNoEntryFoundForKey
	}

	return val, nil
}

func (lc *LocalCacher) Del(key string) error {
	lc.cache.Del(key)

	return nil
}

func (lc *LocalCacher) DelMany(keys []string) error {
	for _, key := range keys {
		_ = lc.Del(key)
	}

	return nil
}

func (lc *LocalCacher) Close() error {
	lc.cache.Clear()
	lc.cache.Close()

	return nil
}
