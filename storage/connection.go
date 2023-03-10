package storage

import (
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
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
		dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", cfg.User, cfg.Password, cfg.DbName, "/cloudsql", cfg.ConnectionName)

		storeDSN := fmt.Sprintf("user=%s host=%s port=%d database=%s password=%s sslmode=disable TimeZone=Etc/UTC",
			cfg.User,
			cfg.Host,
			cfg.Port,
			cfg.DbName,
			cfg.Password)

		if cfg.ConnectionName != "" {
			storeDSN = dbURI
		}
		var err error
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
		zlog.Error("account migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.Token{})
	if err != nil {
		zlog.Error("Token migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.Transaction{})
	if err != nil {
		zlog.Error("Transaction migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.Collection{})
	if err != nil {
		zlog.Error("Collection migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.Offer{})
	if err != nil {
		zlog.Error("Offer migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.Bid{})
	if err != nil {
		zlog.Error("Bid migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.Whitelist{})
	if err != nil {
		zlog.Error("Whitelist migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.SessionState{})
	if err != nil {
		zlog.Error("SessionState migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.MarketPlaceStat{})
	if err != nil {
		zlog.Error("MarketPlaceStat migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.DeployerStat{})
	if err != nil {
		zlog.Error("DeployerStat migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.CollectionIndexer{})
	if err != nil {
		zlog.Error("CollectionIndexer migration", zap.Error(err))
	}

	err = db.AutoMigrate(&entities.AggregatedVolumePerHour{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&entities.AggregatedVolumePerCollectionPerHour{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&entities.UserOrders{})
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
