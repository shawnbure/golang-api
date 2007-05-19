package services

import (
	"gorm.io/datatypes"
	"strconv"
	"testing"
	"time"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
)

func Test_BasicProxyRequest(T *testing.T) {
	var resp ProxyRegisteredNFTsResponse
	err := HttpGet("https://devnet-gateway.elrond.com/address/erd17s2pz8qrds6ake3qwheezgy48wzf7dr5nhdpuu2h4rr4mt5rt9ussj7xzh/registered-nfts", &resp)

	require.Nil(T, err)
	require.Equal(T, resp.Code, "successful")
}

func Test_CreateCollection(T *testing.T) {
	connectToDb()
	cache.InitCacher(config.CacheConfig{Url: "redis://localhost:6379"})

	request := &CreateCollectionRequest{
		UserAddress:   "erd17s2pz8qrds6ake3qwheezgy48wzf7dr5nhdpuu2h4rr4mt5rt9ussj7xzh",
		Name:          "this name is unique",
		TokenId:       "MYNFT-55f092",
		Description:   "this description is flawless",
		Website:       "this is a website",
		DiscordLink:   "this is a discord link",
		TwitterLink:   "this is a twitter link",
		InstagramLink: "this is an instagram link",
		TelegramLink:  "this is a telegram link",
	}

	proxy := "https://devnet-gateway.elrond.com"
	_, err := CreateCollection(request, proxy)
	require.Nil(T, err)

	_, err = storage.GetCollectionByName("this name is unique")
	require.Nil(T, err)
}

func Test_GetCollectionStatistics(T *testing.T) {
	connectToDb()
	cache.InitCacher(config.CacheConfig{Url: "redis://localhost:6379"})

	stats, err := GetStatisticsForCollection(1)
	require.Nil(T, err)
	require.GreaterOrEqual(T, stats.FloorPrice, float64(1))
	require.GreaterOrEqual(T, stats.ItemsCount, uint64(1))
	require.GreaterOrEqual(T, stats.OwnersCount, uint64(1))
	require.GreaterOrEqual(T, stats.VolumeTraded, float64(1))

	stats, err = GetStatisticsForCollection(1)
	require.Nil(T, err)
	hits := cache.GetCacher().GetStats().Hits
	require.GreaterOrEqual(T, hits.Load(), int64(1))
}

func Test_SearchCollection(T *testing.T) {
	connectToDb()
	cache.InitCacher(config.CacheConfig{Url: "redis://localhost:6379"})

	coll := &entities.Collection{
		Name: "this name is uniquee",
	}

	coll.ID = 0
	err := storage.AddCollection(coll)
	require.Nil(T, err)

	coll.ID = 0
	err = storage.AddCollection(coll)
	require.Nil(T, err)

	coll.ID = 0
	err = storage.AddCollection(coll)
	require.Nil(T, err)

	coll.ID = 0
	err = storage.AddCollection(coll)
	require.Nil(T, err)

	coll.ID = 0
	err = storage.AddCollection(coll)
	require.Nil(T, err)

	coll.ID = 0
	err = storage.AddCollection(coll)
	require.Nil(T, err)

	colls, err := GetCollectionsWithNameAlike("uniquee", 5)
	require.Nil(T, err)
	require.Equal(T, len(colls), 5)
	require.Equal(T, colls[0].Name, "this name is uniquee")
	require.Equal(T, colls[1].Name, "this name is uniquee")
	require.Equal(T, colls[2].Name, "this name is uniquee")
	require.Equal(T, colls[3].Name, "this name is uniquee")
	require.Equal(T, colls[4].Name, "this name is uniquee")

	colls, err = GetCollectionsWithNameAlike("uniquee", 5)
	require.Nil(T, err)
	require.Equal(T, len(colls), 5)
	require.Equal(T, colls[0].Name, "this name is uniquee")
	require.Equal(T, colls[1].Name, "this name is uniquee")
	require.Equal(T, colls[2].Name, "this name is uniquee")
	require.Equal(T, colls[3].Name, "this name is uniquee")
	require.Equal(T, colls[4].Name, "this name is uniquee")

	hits := cache.GetCacher().GetStats().Hits
	require.GreaterOrEqual(T, hits.Load(), int64(1))
}

func Test_GetCollectionMetadata(t *testing.T) {
	connectToDb()

	coll := entities.Collection{
		Name: strconv.Itoa(int(time.Now().Unix())),
	}
	err := storage.AddCollection(&coll)
	require.Nil(t, err)

	asset1 := entities.Asset{
		CollectionID: coll.ID,
		Listed:       true,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "red", "background": "dark"}`),
	}
	err = storage.AddAsset(&asset1)
	require.Nil(t, err)

	asset2 := entities.Asset{
		CollectionID: coll.ID,
		Listed:       true,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "green", "background": "dark"}`),
	}
	err = storage.AddAsset(&asset2)
	require.Nil(t, err)

	asset3 := entities.Asset{
		CollectionID: coll.ID,
		Listed:       true,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "blue", "background": "dark"}`),
	}
	err = storage.AddAsset(&asset3)
	require.Nil(t, err)

	asset4 := entities.Asset{
		CollectionID: coll.ID,
		Listed:       true,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{}`),
	}
	err = storage.AddAsset(&asset4)
	require.Nil(t, err)

	asset5 := entities.Asset{
		CollectionID: coll.ID,
		Listed:       true,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "green"}`),
	}
	err = storage.AddAsset(&asset5)
	require.Nil(t, err)

	asset6 := entities.Asset{
		CollectionID: coll.ID,
		Listed:       true,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"background": "dark"}`),
	}
	err = storage.AddAsset(&asset6)
	require.Nil(t, err)

	asset7 := entities.Asset{
		CollectionID: coll.ID,
		Listed:       true,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "yellow", "background": "dark"}`),
	}
	err = storage.AddAsset(&asset7)
	require.Nil(t, err)

	asset8 := entities.Asset{
		CollectionID: coll.ID,
		Listed:       true,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "white", "background": "dark"}`),
	}
	err = storage.AddAsset(&asset8)
	require.Nil(t, err)

	asset9 := entities.Asset{
		CollectionID: coll.ID,
		Listed:       true,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "white", "background": "dark"}`),
	}
	err = storage.AddAsset(&asset9)
	require.Nil(t, err)

	asset10 := entities.Asset{
		CollectionID: coll.ID,
		Listed:       true,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"something_else": "yea"}`),
	}
	err = storage.AddAsset(&asset10)
	require.Nil(t, err)

	collStats, err := computeCollectionMetadata(coll.ID)
	require.Nil(t, err)

	expected := CollectionMetadata{
		NumItems: 10,
		Owners:   map[uint64]bool{1: true},
		AttrStats: map[string]map[string]int{
			"hair": {
				"white":  2,
				"red":    1,
				"green":  2,
				"blue":   1,
				"yellow": 1,
			},
			"background": {
				"dark": 7,
			},
			"something_else": {
				"yea": 1,
			},
		},
	}
	require.Equal(t, expected, *collStats)
}
