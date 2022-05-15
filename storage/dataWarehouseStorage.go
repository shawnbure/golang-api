package storage

import (
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"gorm.io/gorm"
)

func AddOrUpdateAggregatedVolumePerHour(record *entities.AggregatedVolumePerHour) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	recordCount := int64(0)
	err = db.Model(&entities.AggregatedVolumePerHour{}).
		Where("hour=?", record.Hour).
		Count(&recordCount).
		Error

	if recordCount == 0 {
		//if the record does not exist in the db create it return error
		txCreate := database.Create(&record)
		if txCreate.Error != nil {
			return txCreate.Error
		}
		if txCreate.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
	} else {
		txCreate := database.Model(record).Where("hour=?", record.Hour).Updates(record)
		if txCreate.Error != nil {
			return txCreate.Error
		}
		if txCreate.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
	}

	return nil
}

func GetAllAggregatedVolumePerHourInRange(lowBoundary, highBoundary int64) ([]entities.AggregatedVolumePerHour, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	var records []entities.AggregatedVolumePerHour

	txRead := database.
		Where("hour>=? and hour <?", lowBoundary, highBoundary).
		Find(&records)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return records, nil
}

func GetOneAggregatedVolumePerHour(hour int64) (entities.AggregatedVolumePerHour, error) {
	var record entities.AggregatedVolumePerHour

	database, err := GetDBOrError()
	if err != nil {
		return record, err
	}

	txRead := database.
		Where("hour=?", hour).
		First(&record)

	if txRead.Error != nil {
		return record, txRead.Error
	}

	return record, nil
}
