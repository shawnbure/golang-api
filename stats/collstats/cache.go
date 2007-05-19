package collstats

import (
	"encoding/json"
	"fmt"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/boltdb/bolt"
	"github.com/erdsea/erdsea-api/cache"
	"github.com/erdsea/erdsea-api/data/dtos"
	"github.com/erdsea/erdsea-api/stats"
	"github.com/erdsea/erdsea-api/storage"
)

type CacheInfo struct {
	CollectionId   uint64
	CollectionName string
}

var (
	redisCollectionStatsKeyFormat = "CollStats:%s"
	redisSetNXCollectionKeyFormat = "SetNxColl:%s"
	redisSetNXCollectionExpire    = 15 * time.Minute
	tokenIdToCollectionCacheInfo  = []byte("tokenToColl")

	log = logger.GetOrCreate("stats")
)

func GetStatisticsForTokenId(tokenId string) (*dtos.CollectionStatistics, error) {
	redis := cache.GetRedis()
	redisCtx := cache.GetContext()

	redisExpireKey := fmt.Sprintf(redisSetNXCollectionKeyFormat, tokenId)
	ok, err := redis.SetNX(redisCtx, redisExpireKey, true, redisSetNXCollectionExpire).Result()
	if err != nil {
		log.Debug("set nx resulted in error", err)
	}

	shouldComputeStats := ok == true && err == nil
	if shouldComputeStats {
		return setStatisticsRaw(tokenId)
	} else {
		return getStatisticsRaw(tokenId)
	}
}

func AddCollectionToCache(collectionId uint64, collectionName string, tokenId string) (*CacheInfo, error) {
	db := cache.GetBolt()

	cacheInfo := CacheInfo{
		CollectionId:   collectionId,
		CollectionName: collectionName,
	}

	entryBytes, err := json.Marshal(&cacheInfo)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, innerErr := tx.CreateBucketIfNotExists(tokenIdToCollectionCacheInfo)
		if innerErr != nil {
			return innerErr
		}

		innerErr = tx.Bucket(tokenIdToCollectionCacheInfo).Put([]byte(tokenId), entryBytes)
		return innerErr
	})

	return &cacheInfo, err
}

func GetCollectionCacheInfo(tokenId string) (*CacheInfo, error) {
	db := cache.GetBolt()

	var bytes []byte
	err := db.View(func(tx *bolt.Tx) error {
		bytes = tx.Bucket(tokenIdToCollectionCacheInfo).Get([]byte(tokenId))
		return nil
	})
	if err != nil {
		return nil, err
	}

	var cacheInfo CacheInfo
	err = json.Unmarshal(bytes, &cacheInfo)
	if err != nil {
		return nil, err
	}

	return &cacheInfo, nil
}

func GetOrAddCollectionCacheInfo(tokenId string) (*CacheInfo, error) {
	cacheInfo, err := GetCollectionCacheInfo(tokenId)
	if err != nil {
		coll, innerErr := storage.GetCollectionByTokenId(tokenId)
		if innerErr != nil {
			return nil, err
		}

		cacheInfo, innerErr = AddCollectionToCache(coll.ID, coll.Name, coll.TokenID)
		if innerErr != nil {
			return nil, err
		}
	}

	return cacheInfo, nil
}

func getStatisticsRaw(tokenId string) (*dtos.CollectionStatistics, error) {
	redis := cache.GetRedis()
	redisCtx := cache.GetContext()

	redisKey := fmt.Sprintf(redisCollectionStatsKeyFormat, tokenId)
	statsStr, err := redis.Get(redisCtx, redisKey).Result()
	if err != nil {
		log.Debug("could not get from redis")
		return nil, err
	}

	var cacheStats dtos.CollectionStatistics
	err = json.Unmarshal([]byte(statsStr), &cacheStats)
	if err != nil {
		log.Debug("could not unmarshal")
		return nil, err
	}

	return &cacheStats, nil
}

func setStatisticsRaw(tokenId string) (*dtos.CollectionStatistics, error) {
	redis := cache.GetRedis()
	redisCtx := cache.GetContext()

	cacheInfo, err := GetOrAddCollectionCacheInfo(tokenId)
	if err != nil {
		log.Debug("get collection from bolt failed", err)
		return nil, err
	}

	cacheStats, err := stats.ComputeStatisticsForCollection(cacheInfo.CollectionId)
	if err != nil {
		log.Debug("could not compute cacheStats", err)
		return nil, err
	}

	err = updateLeaderboardTables(tokenId, cacheStats)
	if err != nil {
		log.Debug("could not update leaderboard table")
	}

	bytes, err := json.Marshal(cacheStats)
	if err != nil {
		log.Debug("could not marshal", err)
		return nil, err
	}

	redisKey := fmt.Sprintf(redisCollectionStatsKeyFormat, tokenId)
	err = redis.Set(redisCtx, redisKey, bytes, 0).Err()
	if err != nil {
		log.Debug("could not set to redis")
	}

	return cacheStats, nil
}
