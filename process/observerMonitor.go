package process

import (
	"context"
	"time"
)

const (
	tickInterval = 20 * time.Second
)

type observerMonitor struct {
	watchDog bool

	ticker *time.Ticker

	livenessChan chan string
	alertBot     interface{}

	lastBlockHash string

	ctx context.Context
}

func NewObserverMonitor(alertBot interface{}, ctx context.Context) *observerMonitor {
	om := &observerMonitor{
		watchDog:     false,
		ticker:       time.NewTicker(tickInterval),
		livenessChan: make(chan string),
		alertBot:     alertBot,
		ctx:          ctx,
	}

	go om.Monitor()

	return om
}

func (om *observerMonitor) Monitor() {
	for {
		select {
		case hash := <-om.livenessChan:
			om.lastBlockHash = hash
			om.watchDog = true
		case _ = <-om.ticker.C:
			_ = om.alertBot
		case <-om.ctx.Done():
			om.ticker.Stop()
			return
		}
	}
}

func (om *observerMonitor) WatchDogChan() <-chan string {
	return om.livenessChan
}
