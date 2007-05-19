package cache

import (
	"math/big"
	"testing"
	"time"

	"github.com/erdsea/erdsea-api/config"
	"github.com/stretchr/testify/require"
)

var cfg = config.CacheConfig{
	Addrs: []string{":6379"},
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

	cacher := NewBaseCacher(cfg)

	require.NotNil(t, cacher)
}

func TestBaseCacher_SetThenGetShouldWork(t *testing.T) {
	t.Parallel()

	cacher := NewBaseCacher(cfg)

	k := "test-key"

	err := cacher.Set(k, defaultStruct, 1*time.Minute)
	require.Nil(t, err)

	var res TestStruct

	err = cacher.Get(k, &res)

	require.Nil(t, err)
	require.Equal(t, res, defaultStruct)
}
