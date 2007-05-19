package cache

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/erdsea/erdsea-api/config"
	"github.com/stretchr/testify/require"
)

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
