package storage

import (
	"github.com/erdsea/erdsea-api/data"
	"sync"

	"database/sql"
	"github.com/erdsea/erdsea-api/config"
	_ "github.com/lib/pq"
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

		if cfg.ShouldMigrate {
			err = TryMigrate()
		}
		if err != nil {
			panic(err)
		}
	})
}

func TryMigrate() error {
	err := db.AutoMigrate(&data.Account{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&data.Asset{})
	if err != nil {
		return err
	}


	err = db.AutoMigrate(&data.Transaction{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&data.Collection{})
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *gorm.DB {
	return db
}
