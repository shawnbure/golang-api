package cache

import (
	"bytes"
	"context"
	"crypto/rand"
	"github.com/go-redis/redis/v8"
	"go.uber.org/atomic"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"testing"
	"time"

	"github.com/erdsea/erdsea-api/config"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

var cfg = config.CacheConfig{
	Url: "redis://localhost:6379",
}

type TestStruct struct {
	Str        string
	BytesSlice []byte
	U64        uint64
	BigNum     *big.Int
}

var defaultStruct = TestStruct{
	Str:        "test-this-string",
	BytesSlice: []byte("test-this-slice"),
	U64:        123456789,
	BigNum:     big.NewInt(123456789),
}

func TestNewBaseCacher_ShouldCreate(t *testing.T) {
	t.Parallel()

	InitCacher(cfg)

	require.NotNil(t, cacher)
}

func TestBaseCacher_SetThenGetShouldWork(t *testing.T) {
	t.Parallel()

	InitCacher(cfg)

	k := "test-key"

	err := cacher.Set(k, defaultStruct, 1*time.Minute)
	require.Nil(t, err)

	var res TestStruct

	err = cacher.Get(k, &res)

	require.Nil(t, err)
	require.Equal(t, res, defaultStruct)
}

func TestBaseCacher_MaxEntry(t *testing.T) {
	t.Parallel()

	InitCacher(cfg)

	MaxSize := 100_000
	ObjSize := 4_000

	objects := make([][]byte, MaxSize)
	for i := 0; i < MaxSize; i++ {
		objects[i] = make([]byte, ObjSize)
		_, err := rand.Read(objects[i])
		require.Nil(t, err)
	}

	for i := 0; i < MaxSize; i++ {
		err := cacher.Set(strconv.Itoa(i), objects[i], 10*time.Minute)
		require.Nil(t, err)
	}

	result := make([]byte, ObjSize)
	for i := 0; i < MaxSize; i++ {
		err := cacher.Get(strconv.Itoa(i), &result)
		if err == nil {
			require.True(t, bytes.Equal(result, objects[i]))
		}
	}

	println(cacher.stats.Misses.Load())
	println(cacher.stats.Hits.Load())
}

func TestRedisClient_ReconnectOnHiccup(t *testing.T) {
	t.Parallel()

	url := "redis://localhost:6379"

	opt, err := redis.ParseURL(url)
	if err != nil {
		t.Log(err)
	}

	client := redis.NewClient(opt)

	cancellableCtx, cancel := context.WithCancel(context.Background())

	counter := atomic.NewInt64(0)

	exitCh := make(chan struct{})
	go func(ctx context.Context) {
		for {
			t.Log("looping...")
			time.Sleep(1 * time.Second)

			pong, pingErr := ping(client)
			if pingErr != nil {
				t.Log(pingErr)
			} else {
				t.Log(pong)
				counter.Add(1)
			}
			t.Logf("pong counter: %d", counter.Load())

			select {
			case <-ctx.Done():
				t.Log("will exit soon")
				time.Sleep(100 * time.Millisecond)
				exitCh <- struct{}{}
				return
			default:
			}
		}
	}(cancellableCtx)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()
	<-exitCh
}

func ping(client *redis.Client) (string, error) {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return "", err
	}

	return pong, nil
}
