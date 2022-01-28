package storage

import (
	"database/sql"
	"errors"
	"fmt"
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
		dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", cfg.User, cfg.Password, cfg.DbName, "/cloudsql", cfg.ConnectionName)
		storeDSN := fmt.Sprintf("user=%s host=%s port=%d database=%s password=%s sslmode=disable TimeZone=Etc/UTC",
			cfg.User,
			cfg.Host,
			cfg.Port,
			cfg.DbName,
			cfg.Password)
		sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)

		if cfg.ConnectionName != "" {
			storeDSN = dbURI
		}
		db, err = gorm.Open(postgres.Open(storeDSN), &gorm.Config{})
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
