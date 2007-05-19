package tg

import (
	"fmt"
	"time"

	"github.com/erdsea/erdsea-api/config"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot interface {
	Start()
	ObserverDownAlert(lastBlockHash string)
	StoreBlockHash(lastBlockHash string)
}

type telegramBot struct {
	bot       *tb.Bot
	recipient *recipient

	lastBlockHash string
}

func NewTelegramBot(cfg config.BotConfig) (*telegramBot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	tgb := &telegramBot{
		bot:       b,
		recipient: &recipient{id: cfg.RecID},
	}

	tgb.registerListeners()

	return tgb, nil
}

func (tgb *telegramBot) Start() {
	go tgb.bot.Start()
}

func (tgb *telegramBot) ObserverDownAlert(lastBlockHash string) {
	msg := fmt.Sprintf("🚨🚨🚨\n\nObserver seems down. Last block hash received: %s", lastBlockHash)
	_, _ = tgb.bot.Send(tgb.recipient, msg)
}

func (tgb *telegramBot) StoreBlockHash(lastBlockHash string) {
	tgb.lastBlockHash = lastBlockHash
}

func (tgb *telegramBot) registerListeners() {
	tgb.bot.Handle("/hash", func(m *tb.Message) {
		var msg string

		if tgb.lastBlockHash != "" {
			msg = fmt.Sprintf("last block hash 👇\n\n%s", tgb.lastBlockHash)
		} else {
			msg = "no block hash found in storage 🤔"
		}

		_, _ = tgb.bot.Send(tgb.recipient, msg)
	})
}

type recipient struct {
	id string
}

func (r *recipient) Recipient() string {
	return r.id
}

type DisabledBot struct{}

func (db *DisabledBot) Start() {}

func (db *DisabledBot) ObserverDownAlert(_ string) {}

func (db *DisabledBot) StoreBlockHash(_ string) {}
