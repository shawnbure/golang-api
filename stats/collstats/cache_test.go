package collstats

import (
	"testing"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
)

var cfg = config.CacheConfig{
	Url: "redis://localhost:6379",
}

func Test_AddGetBolt(t *testing.T) {
	cache.InitCacher(cfg)
	defer cache.CloseCacher()

	err := AddCollection(12, "name", "token")
	require.Nil(t, err)

	coll, err := GetCollection("token")
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

	err = AddCollection(collection.ID, collection.Name, collection.TokenID)
	require.Nil(t, err)

	asset := entities.Asset{
		TokenID:      collection.TokenID,
		CollectionID: collection.ID,
		PriceNominal: float64(11),
		OwnerId:      0,
	}
	err = storage.AddAsset(&asset)
	require.Nil(t, err)

	tx := entities.Transaction{
		AssetID:      asset.ID,
		CollectionID: collection.ID,
		PriceNominal: asset.PriceNominal,
		Type:         entities.BuyAsset,
	}
	err = storage.AddTransaction(&tx)
	require.Nil(t, err)

	stats, err := GetStatisticsForTokenId(collection.TokenID)
	require.Nil(t, err)
	require.Equal(t, stats.FloorPrice, asset.PriceNominal)
}

func connectToDb() {
	storage.Connect(config.DatabaseConfig{
		Dialect:       "postgres",
		Host:          "localhost",
		Port:          5432,
		DbName:        "erdsea_db_test",
		User:          "postgres",
		Password:      "root",
		SslMode:       "disable",
		MaxOpenConns:  50,
		MaxIdleConns:  10,
		ShouldMigrate: true,
	})
}
