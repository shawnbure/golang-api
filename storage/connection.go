package storage

import (
	"database/sql"
	"errors"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/ENFT-DAO/youbei-api/data/entities"
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

	err = db.AutoMigrate(&entities.Token{})
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

	err = db.AutoMigrate(&entities.Offer{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&entities.Bid{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&entities.Whitelist{})
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
