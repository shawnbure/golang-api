package storage

import (
	"testing"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/entities"
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

	var collectionRead entities.Collection
	txRead := GetDB().Last(&collectionRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, collectionRead, collection)
}

func Test_InitDb(t *testing.T) {
	connectToTestDb()

	account := entities.Account{
		ID: 0,
	}
	tx := GetDB().Find(&account)
	require.Nil(t, tx.Error)
	require.Equal(t, tx.RowsAffected, int64(1))

	collection := entities.Collection{
		ID: 0,
	}
	tx = GetDB().Find(&collection)
	require.Nil(t, tx.Error)
	require.Equal(t, tx.RowsAffected, int64(1))
}

func connectToTestDb() {
	Connect(config.DatabaseConfig{
		Dialect:       "postgres",
		Host:          "localhost",
		Port:          5432,
		DbName:        "youbei_db_test",
		User:          "postgres",
		Password:      "postgres",
		SslMode:       "disable",
		MaxOpenConns:  50,
		MaxIdleConns:  10,
		ShouldMigrate: true,
	})
}
