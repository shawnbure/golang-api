package storage

import (
	"database/sql"
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

func AddOrUpdateAggregatedVolumePerCollectionPerHour(record *entities.AggregatedVolumePerCollectionPerHour) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	recordCount := int64(0)
	err = db.Model(&entities.AggregatedVolumePerCollectionPerHour{}).
		Where("hour=? and collection_id=?", record.Hour, record.CollectionId).
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
		txCreate := database.Model(record).Where("hour=? and collection_id=?", record.Hour, record.CollectionId).Updates(record)
		if txCreate.Error != nil {
			return txCreate.Error
		}
		if txCreate.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
	}

	return nil
}

func GetOneAggregatedVolumePerCollectionPerHour(hour int64) ([]entities.AggregatedVolumePerCollectionPerHour, error) {
	var records []entities.AggregatedVolumePerCollectionPerHour

	database, err := GetDBOrError()
	if err != nil {
		return records, err
	}

	txRead := database.
		Where("hour=?", hour).
		First(&records)

	if txRead.Error != nil {
		return records, txRead.Error
	}

	return records, nil
}

func GetAggregatedTradedVolumeHourly(fromDate, toDate string, _type entities.TxType) (float64, error) {
	database, err := GetDBOrError()
	if err != nil {
		return 0, err
	}

	nullFloat := sql.NullFloat64{}

	txRead := database.Table("transactions").
		Select("sum(transactions.price_nominal)").
		Where("date_trunc('hour', to_timestamp(transactions.timestamp))>=? and date_trunc('hour', to_timestamp(transactions.timestamp))<? and transactions.type=?", fromDate, toDate, _type).
		Scan(&nullFloat)

	if txRead.Error != nil {
		return 0, txRead.Error
	}

	if !nullFloat.Valid {
		return 0, nil
	}

	return nullFloat.Float64, nil
}

func GetAggregatedTradedVolumePerCollectionHourly(fromDate, toDate string) ([]entities.GroupAggregatedVolumePerCollection, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	var records []entities.GroupAggregatedVolumePerCollection

	txRead := database.Table("transactions").
		Select("transactions.type as type, transactions.collection_id as collectionId, sum(transactions.price_nominal) as total").
		Where("date_trunc('hour', to_timestamp(transactions.timestamp))>=? and date_trunc('hour', to_timestamp(transactions.timestamp))<?", fromDate, toDate).
		Group("transactions.type").
		Group("transactions.collection_id").
		Scan(&records)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return records, nil
}
