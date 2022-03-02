package storage

import "github.com/ENFT-DAO/youbei-api/data/entities"

func GetDeployerStat(deployerAddr string) (entities.DeployerStat, error) {
	var stat entities.DeployerStat

	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}
	err = database.
		Model(&entities.DeployerStat{}).
		Where("deployer_addr = ?", deployerAddr).
		Order("updated_at DESC").
		First(&stat).Error
	if err != nil {
		return stat, err
	}
	return stat, nil
}

func CreateDeployerStat(deployerAddr string) (stat entities.DeployerStat, err error) {
	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}
	stat.DeployerAddr = deployerAddr
	err = database.Create(&stat).Error
	if err != nil {
		return stat, err
	}
	return stat, nil
}

func UpdateDeployerIndexer(lastIndex uint64, deployerAddr string) (entities.DeployerStat, error) {
	var stat entities.DeployerStat

	database, err := GetDBOrError()
	if err != nil {
		return stat, err
	}

	err = database.
		Model(&entities.DeployerStat{}).
		Where("deployer_addr = ? ", deployerAddr).
		Order("updated_at DESC").
		First(&stat).Error
	if err != nil {
		return stat, err
	}
	stat.LastIndex = lastIndex
	err = database.Updates(stat).Where("id=?", stat.ID).Error
	if err != nil {
		return stat, err
	}
	return stat, nil
}
