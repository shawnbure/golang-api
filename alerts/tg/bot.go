package tg

import (
	"fmt"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot interface {
	Start()
	ObserverDownAlert(lastBlockHash string)
}

type TelegramBot struct {
	bot       *tb.Bot
	recipient *recipient
}

func NewTelegramBot() (*TelegramBot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  "2010065738:AAH0J6N2meI7Wj2c_5AsnMAXXGNcj2YTYPk",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	b.Handle("/hello", func(m *tb.Message) {
		r := &recipient{}

		resp := fmt.Sprintf("Hello %s \U0001F976", m.Sender.FirstName+m.Sender.LastName)
		_, innerErr := b.Send(r, resp)

		log.Println(innerErr)
	})

	return &TelegramBot{
		bot:       b,
		recipient: &recipient{},
	}, nil
}

func (tgb *TelegramBot) Start() {
	go tgb.bot.Start()
}

func (tgb *TelegramBot) ObserverDownAlert(lastBlockHash string) {
	msg := fmt.Sprintf("🚨🚨🚨\n\nObserver seems down. Last block hash received: %s", lastBlockHash)
	_, _ = tgb.bot.Send(tgb.recipient, msg)
}

type recipient struct {
}

func (r *recipient) Recipient() string {
	return "-1001435005959"
}
