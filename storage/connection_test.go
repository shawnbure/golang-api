package storage

import (
	"testing"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func Test_Connection(t *testing.T) {
	connectToTestDb()
}

func Test_BasicWrite(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	tx := GetDB().Create(&collection)
	require.Nil(t, tx.Error)
}

func Test_BasicWriteRead(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	tx := GetDB().Create(&collection)
	require.Nil(t, tx.Error)

	var collectionRead data.Collection
	txRead := GetDB().Last(&collectionRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, collectionRead, collection)
}

func connectToTestDb() {
	Connect(config.DatabaseConfig{
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
