package process

import (
	"context"
	"fmt"
	"github.com/erdsea/erdsea-api/alerts/tg"
	"testing"
	"time"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/stretchr/testify/require"
)

func TestEventProcessor_OnEvents(t *testing.T) {
	t.Parallel()

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

	cacher, _ := cache.NewLocalCacher()

	monitor := NewObserverMonitor(&tg.DisabledBot{}, context.Background(), false)
	proc := NewEventProcessor(addresses, identifiers, cacher, monitor)
	proc.OnEvents(blockEvents)
}

func TestNewEventProcessor_OnEventsWithFinalized(t *testing.T) {
	t.Parallel()

	cacher, err := cache.NewLocalCacher()
	require.Nil(t, err)

	addresses := []string{"erd1", "erd2", "erd3"}
	identifiers := []string{"putNftForSale", "buyNft", "withdrawNft"}

	monitor := NewObserverMonitor(&tg.DisabledBot{}, context.Background(), false)
	proc := NewEventProcessor(addresses, identifiers, cacher, monitor)

	events := []entities.Event{
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
	}

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)
	finalIdx := 0
	offset := 0
	cnt := 0

	var hashes []string

	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				hash := fmt.Sprintf("hash-%03d", offset+1)
				hashes = append(hashes, hash)
				offset++

				proc.OnEvents(entities.BlockEvents{
					Hash:   hash,
					Events: events,
				})

				if cnt == 3 {
					proc.OnFinalizedEvent(entities.FinalizedBlock{
						Hash: hashes[finalIdx],
					})
					finalIdx++
					cnt = 0
				} else {
					cnt++
				}
			}
		}
	}()

	time.Sleep(5 * time.Minute)
	ticker.Stop()
	done <- true
}
