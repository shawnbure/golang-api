package services

import (
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

const ConfigTestFilePath = "../config/config_test.toml"

func Test_CreateNewCollection(t *testing.T) {
	connectToDb(t)

	CreateNewCollection("ownerAddress", "tokenId", "collectionName", "collectionDescription")

	_, err := storage.GetAccountByAddress("ownerAddress")
	require.Nil(t, err)

	_, err = storage.GetCollectionByName("collectionName")
	require.Nil(t, err)
}

func connectToDb(t *testing.T) {
	cfg, err := config.LoadConfig(ConfigTestFilePath)
	require.Nil(t, err)
	storage.Connect(cfg.Database)
}
