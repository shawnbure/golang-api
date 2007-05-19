package services

import (
	"testing"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/config"
)

var cfg = config.CacheConfig{
	Url: "redis://localhost:6379",
}

func Test_UpdateDeposit(t *testing.T) {
	connectToDb()
	cache.InitCacher(cfg)

	//TODO: make test after sc deploy
}
