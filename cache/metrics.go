package cache

import (
	"context"

	"github.com/go-redis/cache/v8"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
)

var meter = metric.Must(global.Meter("cache.metrics"))

func MonitorCache(cache *cache.Cache) {
	var hitsCounter, missesCounter metric.Int64SumObserver

	batchObserver := meter.NewBatchObserver(
		func(ctx context.Context, result metric.BatchObserverResult) {
			stats := cache.Stats()

			result.Observe(nil,
				hitsCounter.Observation(int64(stats.Hits)),
				missesCounter.Observation(int64(stats.Misses)),
			)
		})

	hitsCounter = batchObserver.NewInt64SumObserver("cache.hits")
	missesCounter = batchObserver.NewInt64SumObserver("cache.misses")
}
