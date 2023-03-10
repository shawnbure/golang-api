package collstats

import (
	"testing"

	"gorm.io/datatypes"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/stretchr/testify/require"
)

var cfg = config.CacheConfig{
	Url: "redis://localhost:6379",
}

func Test_AddGetBolt(t *testing.T) {
	cache.InitCacher(cfg)
	defer cache.CloseCacher()

	_, err := AddCollectionToCache(12, "name", datatypes.JSON("[]"), "token")
	require.Nil(t, err)

	coll, err := GetCollectionCacheInfo("token")
	require.Nil(t, err)

	require.Equal(t, coll.CollectionId, uint64(12))
	require.Equal(t, coll.CollectionName, "name")
}

func Test_GetStats(t *testing.T) {
	connectToDb()
	cache.InitCacher(cfg)
	defer cache.CloseCacher()

	collection := entities.Collection{
		TokenID: "Token12",
		Name:    "Stats",
	}
	err := storage.AddCollection(&collection)
	require.Nil(t, err)

	_, err = AddCollectionToCache(collection.ID, collection.Name, datatypes.JSON("[]"), collection.TokenID)
	require.Nil(t, err)

	token := entities.Token{
		TokenID:      collection.TokenID,
		CollectionID: collection.ID,
		PriceNominal: float64(11),
		OwnerID:      0,
	}
	err = storage.AddToken(&token)
	require.Nil(t, err)

	tx := entities.Transaction{
		TokenID:      token.ID,
		CollectionID: collection.ID,
		PriceNominal: token.PriceNominal,
		Type:         entities.BuyToken,
	}
	err = storage.AddTransaction(&tx)
	require.Nil(t, err)

	stats, err := GetStatisticsForTokenId(collection.TokenID)
	require.Nil(t, err)
	require.Equal(t, stats.FloorPrice, token.PriceNominal)
}

func connectToDb() {
	storage.Connect(config.DatabaseConfig{
		Dialect:       "postgres",
		Host:          "localhost",
		Port:          5432,
		DbName:        "youbei_dev",
		User:          "postgres",
		Password:      "root",
		SslMode:       "disable",
		MaxOpenConns:  50,
		MaxIdleConns:  10,
		ShouldMigrate: true,
	})
}
