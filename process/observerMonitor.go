package process

import (
	"context"
	"time"

	"github.com/ENFT-DAO/youbei-api/alerts/tg"
)

const (
	tickInterval = 20 * time.Second
)

type observerMonitor struct {
	watchDog bool

	ticker *time.Ticker

	livenessChan chan string
	alertBot     tg.Bot
	isEnabled    bool

	lastBlockHash string

	ctx context.Context
}

func NewObserverMonitor(alertBot tg.Bot, ctx context.Context, isEnabled bool) *observerMonitor {
	om := &observerMonitor{
		watchDog:     false,
		ticker:       time.NewTicker(tickInterval),
		livenessChan: make(chan string),
		alertBot:     alertBot,
		isEnabled:    isEnabled,
		ctx:          ctx,
	}

	if isEnabled {
		go om.monitor()
	}

	return om
}

func (om *observerMonitor) LivenessChan() chan string {
	return om.livenessChan
}

func (om *observerMonitor) IsEnabled() bool {
	return om.isEnabled
}

func (om *observerMonitor) monitor() {
	for {
		select {
		case hash := <-om.livenessChan:
			om.alertBot.StoreBlockHash(hash)
			om.lastBlockHash = hash
			om.watchDog = true
		case _ = <-om.ticker.C:
			if !om.watchDog {
				om.alertBot.ObserverDownAlert(om.lastBlockHash)
			}
			om.watchDog = false
		case <-om.ctx.Done():
			om.ticker.Stop()
			return
		}
	}
}
