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

func UpdateMarketPlaceIndexer(lastIndex uint64) (entities.MarketPlaceStat, error) {
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
	stat.LastIndex = stat.LastIndex + 1
	err = database.Updates(stat).Where("id=?", stat.ID).Error
	if err != nil {
		return stat, err
	}
	return stat, nil
}
