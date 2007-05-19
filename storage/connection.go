package storage

import (
	"database/sql"
	"errors"
	"github.com/erdsea/erdsea-api/data/entities"
	images2 "github.com/erdsea/erdsea-api/data/entities/images"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"

	"github.com/erdsea/erdsea-api/config"
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
	err := db.AutoMigrate(&entities.Account{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&entities.Asset{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&entities.Transaction{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&entities.Collection{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&images2.AccountCoverImage{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&images2.AccountProfileImage{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&images2.CollectionCoverImage{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&images2.CollectionProfileImage{})
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
