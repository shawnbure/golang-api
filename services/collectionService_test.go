package services

import (
	"strconv"
	"testing"
	"time"

	"gorm.io/datatypes"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/interaction"
	"github.com/ENFT-DAO/youbei-api/stats"
	"github.com/ENFT-DAO/youbei-api/storage"
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

	collectionStats, err := stats.ComputeStatisticsForCollection(1)
	require.Nil(T, err)
	require.GreaterOrEqual(T, collectionStats.FloorPrice, float64(1))
	require.GreaterOrEqual(T, collectionStats.ItemsTotal, uint64(1))
	require.GreaterOrEqual(T, collectionStats.OwnersTotal, uint64(1))
	require.GreaterOrEqual(T, collectionStats.VolumeTraded, float64(1))
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

	token1 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "red", "background": "dark"}`),
	}
	err = storage.AddToken(&token1)
	require.Nil(t, err)

	token2 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "green", "background": "dark"}`),
	}
	err = storage.AddToken(&token2)
	require.Nil(t, err)

	token3 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "blue", "background": "dark"}`),
	}
	err = storage.AddToken(&token3)
	require.Nil(t, err)

	token4 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{}`),
	}
	err = storage.AddToken(&token4)
	require.Nil(t, err)

	token5 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "green", "background": "dark"}`),
	}
	err = storage.AddToken(&token5)
	require.Nil(t, err)

	token6 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"background": "dark"}`),
	}
	err = storage.AddToken(&token6)
	require.Nil(t, err)

	token7 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "yellow", "background": "dark"}`),
	}
	err = storage.AddToken(&token7)
	require.Nil(t, err)

	token8 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "white", "background": "dark"}`),
	}
	err = storage.AddToken(&token8)
	require.Nil(t, err)

	token9 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "white", "background": "dark"}`),
	}
	err = storage.AddToken(&token9)
	require.Nil(t, err)

	token10 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"something_else": "yea"}`),
	}
	err = storage.AddToken(&token10)
	require.Nil(t, err)

	collStats, err := stats.ComputeCollectionMetadata(coll.ID)
	require.Nil(t, err)

	expected := stats.CollectionMetadata{
		NumItems: 10,
		Owners:   map[uint64]bool{1: true},
		AttrStats: []dtos.AttributeStat{{
			TraitType: "hair",
			Value:     "red",
			Total:     1,
		}, {
			TraitType: "background",
			Value:     "dark",
			Total:     8,
		}, {
			TraitType: "hair",
			Value:     "green",
			Total:     2,
		}, {
			TraitType: "hair",
			Value:     "blue",
			Total:     1,
		}, {
			TraitType: "hair",
			Value:     "yellow",
			Total:     1,
		}, {
			TraitType: "hair",
			Value:     "white",
			Total:     2,
		}, {
			TraitType: "something_else",
			Value:     "yea",
			Total:     1,
		},
		},
	}
	require.Equal(t, expected, *collStats)
}

func Test_GetMintInfoFromContract(t *testing.T) {
	cache.InitCacher(config.CacheConfig{Url: "redis://localhost:6379"})
	defer cache.CloseCacher()

	cfg := config.BlockchainConfig{
		ProxyUrl: "https://devnet-gateway.elrond.com",
		ChainID:  "D",
	}

	interaction.InitBlockchainInteractor(cfg)
	bi := interaction.GetBlockchainInteractor()
	require.NotNil(t, bi)

	info, err := GetMintInfoForContract("erd1qqqqqqqqqqqqqpgq3uvfynvpvcs8aldhuyrseuyepmp0cj7at9usgefv56")
	require.Nil(t, err)
	require.True(t, info.MaxSupply > 0)
	require.True(t, info.TotalSold > 0)

	info, err = GetMintInfoForContract("erd1qqqqqqqqqqqqqpgq3uvfynvpvcs8aldhuyrseuyepmp0cj7at9usgefv56")
	require.Nil(t, err)
	require.True(t, info.MaxSupply > 0)
	require.True(t, info.TotalSold > 0)
}

func Test_StandardizeName(t *testing.T) {
	name1 := "\n  \t    Name       1  \t   \n   \t"
	require.Equal(t, "Name 1", standardizeName(name1))
}

func Test_CollectionFlags(t *testing.T) {
	connectToDb()

	collection := entities.Collection{
		Flags: datatypes.JSON(`["k1", "k2"]`),
	}
	err := storage.AddCollection(&collection)
	require.Nil(t, err)

	db := storage.GetDB()
	var read entities.Collection
	tx := db.Where(datatypes.JSONQuery("flags").HasKey("k1")).Find(&read)
	require.Nil(t, tx.Error)
	require.NotZero(t, tx.RowsAffected)
}

func Test_ValidFlags(t *testing.T) {
	flags := []string{"k1", "k2"}
	err := CheckValidFlags(flags)
	require.Nil(t, err)

	flags = []string{"k1", "k2", ";"}
	err = CheckValidFlags(flags)
	require.NotNil(t, err)
}

func Test_FlagsFilter(t *testing.T) {
	connectToDb()

	collection := entities.Collection{
		Flags: datatypes.JSON(`["k1", "k2"]`),
	}
	err := storage.AddCollection(&collection)
	require.Nil(t, err)

	colls, err := storage.GetCollectionsWithOffsetLimit(0, 10, []string{"k1", "k2"})
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(colls), 1)

	colls, err = storage.GetCollectionsWithOffsetLimit(0, 10, []string{"k1", "bkqbdwkqehwleqh"})
	require.Nil(t, err)
	require.Zero(t, len(colls))
}
