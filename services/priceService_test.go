package services

import (
	"testing"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/stretchr/testify/require"
)

func TestGetEGLDPrice(t *testing.T) {
	cache.InitCacher(config.CacheConfig{Url: "redis://localhost:6379"})

	price, err := GetEGLDPrice()
	require.Nil(t, err)
	require.Greater(t, price, float64(0))

	var cachePrice float64
	err = cache.GetCacher().Get(EGLDPriceCacheKey, &cachePrice)
	require.Nil(t, err)
	require.Equal(t, price, cachePrice)

	stats := cache.GetCacher().GetStats()
	require.GreaterOrEqual(t, stats.Hits.Load(), int64(1))

	err = cache.GetCacher().Get(EGLDPriceCacheKey, &cachePrice)
	require.Nil(t, err)
	require.Equal(t, price, cachePrice)

	stats = cache.GetCacher().GetStats()
	require.GreaterOrEqual(t, stats.Hits.Load(), int64(2))
}
