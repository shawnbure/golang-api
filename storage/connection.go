package storage

import (
	"database/sql"
	"errors"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/erdsea/erdsea-api/config"
	"github.com/erdsea/erdsea-api/data/entities"
	"github.com/erdsea/erdsea-api/data/entities/images"
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

		err = createDefaultEntitiesIfNotExist()
		if err != nil {
			panic(err)
		}
	})
}

func createDefaultEntitiesIfNotExist() error {
	err := createDefaultAccountIfNotExist()
	if err != nil {
		return err
	}

	err = createDefaultCollectionIfNotExist()
	if err != nil {
		return err
	}

	return nil
}

func createDefaultAccountIfNotExist() error {
	account := entities.Account{}
	tx := db.Where("id = 0").Find(&account)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 1 {
		return nil
	}

	tx = db.Create(&account)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("could not create account")
	}

	tx = db.Table("accounts").Where("id = ?", account.ID).Update("id", uint64(0))
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("could not update new account to id = 0")
	}

	return nil
}

func createDefaultCollectionIfNotExist() error {
	collection := entities.Collection{}
	tx := db.Where("id = 0").Find(&collection)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 1 {
		return nil
	}

	tx = db.Create(&collection)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("could not create collection")
	}

	tx = db.Table("collections").Where("id = ?", collection.ID).Update("id", uint64(0))
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("could not update new collection to id = 0")
	}

	return nil
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
