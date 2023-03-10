package storage

import "github.com/ENFT-DAO/youbei-api/data/entities"

func GetMarketPlaceIndexer() (entities.MarketPlaceStat, error) {
	var stat entities.MarketPlaceStat

	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}
	err = database.
		Model(&entities.MarketPlaceStat{}).
		Order("updated_at DESC").
		First(&stat).Error
	if err != nil {
		return stat, err
	}
	return stat, nil
}

func CreateMarketPlaceStat() (stat entities.MarketPlaceStat, err error) {
	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}
	err = database.Create(&stat).Error
	if err != nil {
		return stat, err
	}
	return stat, nil
}
func UpdateMarketPlaceHash(hash string) (entities.MarketPlaceStat, error) {
	var stat entities.MarketPlaceStat

	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}

	err = database.
		Model(&entities.MarketPlaceStat{}).
		Order("updated_at DESC").
		First(&stat).Error
	if err != nil {
		return stat, err
	}
	stat.LastHash = hash
	err = database.Updates(stat).Where("id=?", stat.ID).Error
	if err != nil {
		return stat, err
	}
	return stat, nil
}
func UpdateMarketPlaceIndexerTimestamp(timestamp uint64) (entities.MarketPlaceStat, error) {
	var stat entities.MarketPlaceStat

	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}

	err = database.
		Model(&entities.MarketPlaceStat{}).
		Order("updated_at DESC").
		First(&stat).Error
	if err != nil {
		return stat, err
	}
	if stat.LastTimestamp < timestamp {
		stat.LastTimestamp = timestamp
		err = database.Updates(stat).Where("id=?", stat.ID).Error
		if err != nil {
			return stat, err
		}
	}
	return stat, nil
}
