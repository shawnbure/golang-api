package tg

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTelegramBot(t *testing.T) {
	b, err := NewTelegramBot()
	require.Nil(t, err)

	b.ObserverDownAlert("abcdef")
}

func TestNewTelegramBot_Start(t *testing.T) {
	b, err := NewTelegramBot()
	require.Nil(t, err)

	b.Start()

	ch := make(chan bool)
	<-ch
}
