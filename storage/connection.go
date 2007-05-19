package storage

import (
	"database/sql"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data"
	"github.com/erdsea/erdsea-api/data/images"
	_ "github.com/lib/pq"
)

var NoDBError = errors.New("no DB Connection")

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

	err = db.AutoMigrate(&data.Token{})
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

	err = db.AutoMigrate(&images.AccountCoverImage{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&images.AccountProfileImage{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&images.CollectionCoverImage{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&images.CollectionProfileImage{})
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *gorm.DB {
	return db
}

func GetDBOrError() (*gorm.DB, error) {
	database := GetDB()
	if database == nil {
		return nil, NoDBError
	}

	return database, nil
}
