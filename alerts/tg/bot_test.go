package tg

import (
	"testing"

	"github.com/erdsea/erdsea-api/config"
	"github.com/stretchr/testify/require"
)

var cfg = config.BotConfig{
	Token: "2010065738:AAH0J6N2meI7Wj2c_5AsnMAXXGNcj2YTYPk",
	RecID: "-1001435005959",
}

func TestNewTelegramBot(t *testing.T) {
	b, err := NewTelegramBot(cfg)
	require.Nil(t, err)

	b.ObserverDownAlert("abcdef")
}

func TestNewTelegramBot_Start(t *testing.T) {
	b, err := NewTelegramBot(cfg)
	require.Nil(t, err)

	b.Start()

	ch := make(chan bool)
	<-ch
}
