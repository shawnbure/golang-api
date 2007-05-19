package process

import (
	"context"
	"time"

	"github.com/erdsea/erdsea-api/alerts/tg"
)

const (
	tickInterval = 20 * time.Second
)

type observerMonitor struct {
	watchDog bool

	ticker *time.Ticker

	livenessChan chan string
	alertBot     tg.Bot

	lastBlockHash string

	ctx context.Context
}

func NewObserverMonitor(alertBot tg.Bot, ctx context.Context) *observerMonitor {
	om := &observerMonitor{
		watchDog:     false,
		ticker:       time.NewTicker(tickInterval),
		livenessChan: make(chan string),
		alertBot:     alertBot,
		ctx:          ctx,
	}

	go om.monitor()

	return om
}

func (om *observerMonitor) WatchDogChan() chan string {
	return om.livenessChan
}

func (om *observerMonitor) monitor() {
	for {
		select {
		case hash := <-om.livenessChan:
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
