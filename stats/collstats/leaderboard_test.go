package collstats

import (
	"testing"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/stretchr/testify/require"
)

func Test_GetLeaderboardEntries(t *testing.T) {
	connectToDb()
	cache.InitCacher(cfg)
	defer cache.CloseCacher()

	entries, err := GetLeaderboardEntries("itemsTotal", 0, 10, false)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(entries), 1)
}
