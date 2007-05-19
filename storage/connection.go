package storage

import (
	"sync"

	"database/sql"
	"github.com/erdsea/erdsea-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once sync.Once
	db   *gorm.DB
)

func Connect(cfg config.DatabaseConfig) {
	once.Do(func() {
		sqlDb, err := sql.Open(cfg.Dialect, cfg.Url())
		if err != nil {
			panic(err)
		}

		sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)

		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDb,
		}))
		if err != nil {
			panic(err)
		}
	})
}

func GetDB() *gorm.DB {
	return db
}
