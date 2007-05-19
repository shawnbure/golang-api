package cache

import (
	"context"
	"github.com/uptrace/uptrace-go/uptrace"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func GetSet(c *BaseCacher){
	_ = c.Set("test-key", "test-value", 1*time.Second)

	var s string
	_ = c.Get("test-key", &s)
}

func TestMonitorCache(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN("https://Cto_ZSq3mcbp5d7V7ILe5g@api.uptrace.dev/215"),
		uptrace.WithServiceName("pula-n-pizda"),
		uptrace.WithServiceVersion("1.0.0"),
	)

	defer func(ctx context.Context) {
		err := uptrace.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}(ctx)

	cacher := NewBaseCacher(cfg)

	go MonitorCache(cacher.cache)
	go GetSet(cacher)

	ch := make(chan os.Signal, 3)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-ch
}
