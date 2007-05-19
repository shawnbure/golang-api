package process

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type mockBot struct {
}

func (m *mockBot) Start() {
}

func (m *mockBot) ObserverDownAlert(lastBlockHash string) {
	fmt.Println("observer down. last block hash", lastBlockHash)
}

func TestObserverMonitor_Monitor(t *testing.T) {
	t.Parallel()

	bot := &mockBot{}

	monit := NewObserverMonitor(bot, context.Background())

	ticker := time.NewTicker(3 * time.Second)
	done := make(chan bool)

	shouldNotBeZero := 2

	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				if shouldNotBeZero != 0 {
					hash := fmt.Sprintf("hash-%d", shouldNotBeZero)
					monit.LivenessChan() <- hash
					shouldNotBeZero--
				}
			}
		}
	}()

	time.Sleep(10 * time.Minute)
	ticker.Stop()
	done <- true
}
