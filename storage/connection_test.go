package storage

import (
	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	_ "github.com/lib/pq"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const ConfigTestFilePath = "../config/config_test.toml"

func Test_Connection(t *testing.T) {
	connectToDb(t)
}

func Test_BasicWrite(t *testing.T) {
	connectToDb(t)

	collection := defaultCollection()
	tx := GetDB().Create(&collection)
	require.Nil(t, tx.Error)
}

func Test_BasicWriteRead(t *testing.T) {
	connectToDb(t)

	collection := defaultCollection()
	tx := GetDB().Create(&collection)
	require.Nil(t, tx.Error)

	var collectionRead data.Collection
	txRead := GetDB().Last(&collectionRead)

	require.Nil(t, txRead.Error)
	assert.Equal(t, collectionRead, collection)
}

func connectToDb(t *testing.T) {
	cfg, err := config.LoadConfig(ConfigTestFilePath)
	require.Nil(t, err)
	Connect(cfg.Database)
}
