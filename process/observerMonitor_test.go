package process

import (
	"context"
	"fmt"
	"github.com/erdsea/erdsea-api/alerts/tg"
	"github.com/erdsea/erdsea-api/config"
	"github.com/stretchr/testify/require"
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

func (m *mockBot) StoreBlockHash(lastBlockHash string) {
	fmt.Println("got new block hash", lastBlockHash)
}

func TestObserverMonitor_Monitor(t *testing.T) {
	t.Parallel()

	bot := &mockBot{}

	monit := NewObserverMonitor(bot, context.Background(), true)

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

func TestObserverMonitor_NotifyBotOnBlockHash(t *testing.T) {
	t.Parallel()

	cfg := config.BotConfig{
		Token: "2010065738:AAH0J6N2meI7Wj2c_5AsnMAXXGNcj2YTYPk",
		RecID: "-1001435005959",
	}

	bot, err := tg.NewTelegramBot(cfg)
	require.Nil(t, err)

	bot.Start()

	monit := NewObserverMonitor(bot, context.Background(), true)
	monit.LivenessChan() <- "test-hash-cool-cool"

	ch := make(chan bool)
	<-ch
}
