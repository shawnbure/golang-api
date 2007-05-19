package process

import (
	"context"
	"github.com/erdsea/erdsea-api/alerts/tg"
	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/entities"
	"testing"
)

var cacheCfg = config.CacheConfig{
	Url: "redis://localhost:6379",
}

var blockchainCfg = config.BlockchainConfig{
	ProxyUrl:           "https://devnet-gateway.elrond.com",
	MarketplaceAddress: "erd1qqqqqqqqqqqqqpgqm4dmwyxc5fsj49z3jcu9h08azjrcf60kt9uspxs483",
}

func TestEventProcessor_OnEvents(t *testing.T) {
	t.Parallel()
	cache.InitCacher(cacheCfg)

	addresses := []string{"erd1", "erd2", "erd3"}
	identifiers := []string{"func1", "func2", "func3"}

	blockEvents := entities.BlockEvents{
		Hash: "abcdef",
		Events: []entities.Event{
			{
				Address:    addresses[0],
				Identifier: identifiers[0],
			},
			{
				Address:    addresses[1],
				Identifier: identifiers[1],
			},
			{
				Address:    addresses[2],
				Identifier: identifiers[2],
			},
		},
	}

	monitor := NewObserverMonitor(&tg.DisabledBot{}, context.Background(), false)
	proc := NewEventProcessor(addresses, identifiers, blockchainCfg.ProxyUrl, blockchainCfg.MarketplaceAddress, monitor)
	proc.OnEvents(blockEvents)
}
