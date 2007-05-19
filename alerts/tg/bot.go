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
}

type telegramBot struct {
	bot       *tb.Bot
	recipient *recipient
}

func NewTelegramBot(cfg config.BotConfig) (*telegramBot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	return &telegramBot{
		bot:       b,
		recipient: &recipient{id: cfg.RecID},
	}, nil
}

func (tgb *telegramBot) Start() {
	go tgb.bot.Start()
}

func (tgb *telegramBot) ObserverDownAlert(lastBlockHash string) {
	msg := fmt.Sprintf("🚨🚨🚨\n\nObserver seems down. Last block hash received: %s", lastBlockHash)
	_, _ = tgb.bot.Send(tgb.recipient, msg)
}

func (tgb *telegramBot) registerListeners() {
	tgb.bot.Handle("/hello", func(m *tb.Message) {
		r := &recipient{}

		resp := fmt.Sprintf("Hello %s \U0001F976", m.Sender.FirstName+m.Sender.LastName)
		_, _ = tgb.bot.Send(r, resp)
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
