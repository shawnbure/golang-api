package storage

import (
	"gorm.io/gorm"

	"github.com/erdsea/erdsea-api/data/entities"
)

func AddProffer(p *entities.Proffer) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&p)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
