package services

import (
	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
	"testing"
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
	err := CreateCollection(request, proxy)
	require.Nil(T, err)

	_, err = storage.GetCollectionByName("this name is unique")
	require.Nil(T, err)
}
