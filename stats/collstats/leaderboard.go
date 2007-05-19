package collstats

import (
	"errors"

	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/go-redis/redis/v8"
)

type LeaderboardEntry struct {
	CollectionId   string  `json:"CollectionId"`
	CollectionName string  `json:"CollectionName"`
	ItemsTotal     uint64  `json:"itemsTotal"`
	OwnersTotal    uint64  `json:"ownersTotal"`
	FloorPrice     float64 `json:"floorPrice"`
	VolumeTraded   float64 `json:"volumeTraded"`
}

const (
	ItemsTotal   = "itemsTotal"
	OwnersTotal  = "ownersTotal"
	FloorPrice   = "floorPrice"
	VolumeTraded = "volumeTraded"
)

func updateLeaderboardTables(tokenId string, stats *dtos.CollectionStatistics) error {
	redisCache := cache.GetRedis()
	redisCtx := cache.GetContext()

	_, err := redisCache.ZAdd(redisCtx, ItemsTotal, &redis.Z{
		Score:  float64(stats.ItemsTotal),
		Member: tokenId,
	}).Result()
	if err != nil {
		log.Debug("sorted set add failed")
	}

	_, err = redisCache.ZAdd(redisCtx, OwnersTotal, &redis.Z{
		Score:  float64(stats.OwnersTotal),
		Member: tokenId,
	}).Result()
	if err != nil {
		log.Debug("sorted set add failed")
	}

	_, err = redisCache.ZAdd(redisCtx, FloorPrice, &redis.Z{
		Score:  stats.FloorPrice,
		Member: tokenId,
	}).Result()
	if err != nil {
		log.Debug("sorted set add failed")
	}

	_, err = redisCache.ZAdd(redisCtx, VolumeTraded, &redis.Z{
		Score:  stats.VolumeTraded,
		Member: tokenId,
	}).Result()
	if err != nil {
		log.Debug("sorted set add failed")
	}

	return nil
}

func GetLeaderboardEntries(table string, start int, stop int, rev bool) ([]LeaderboardEntry, error) {
	redisCache := cache.GetRedis()
	redisCtx := cache.GetContext()

	err := testTableName(table)
	if err != nil {
		return nil, err
	}

	var tokenIds []string
	if rev {
		tokenIds, err = redisCache.ZRevRange(redisCtx, table, int64(start), int64(stop)).Result()
	} else {
		tokenIds, err = redisCache.ZRange(redisCtx, table, int64(start), int64(stop)).Result()
	}
	if err != nil {
		return nil, err
	}

	entries := make([]LeaderboardEntry, len(tokenIds))
	for index, tokenId := range tokenIds {
		collCacheInfo, innerErr := GetOrAddCollectionCacheInfo(tokenId)
		if innerErr != nil {
			log.Debug("could not get from bolt")
			continue
		}

		collStats, innerErr := getStatisticsRaw(tokenId)
		if innerErr != nil {
			log.Debug("could not get from statistics")
			continue
		}

		entries[index] = LeaderboardEntry{
			CollectionId:   tokenId,
			CollectionName: collCacheInfo.CollectionName,
			ItemsTotal:     collStats.ItemsTotal,
			OwnersTotal:    collStats.OwnersTotal,
			FloorPrice:     collStats.FloorPrice,
			VolumeTraded:   collStats.VolumeTraded,
		}
	}

	return entries, nil
}

func testTableName(table string) error {
	switch table {
	case ItemsTotal:
		return nil
	case OwnersTotal:
		return nil
	case FloorPrice:
		return nil
	case VolumeTraded:
		return nil
	default:
		return errors.New("not a valid lb table name")
	}
}
