package collstats

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/stats"
	"github.com/ENFT-DAO/youbei-api/storage"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/boltdb/bolt"
)

var (
	redisCollectionStatsKeyFormat = "CollStats:%s"
	redisSetNXCollectionKeyFormat = "SetNxColl:%s"
	redisCollectionChange         = "coll_change_queue"
	redisProccesingCollection     = "coll_change_queue-processing"
	redisSetNXCollectionExpire    = 1 * time.Second
	tokenIdToCollectionCacheInfo  = []byte("tokenToColl")

	log = logger.GetOrCreate("stats")
)

func AddCollectionToCheck(col dtos.CollectionToCheck) error {
	redis := cache.GetRedis()
	redisCtx := cache.GetContext()
	cmd := redis.SAdd(redisCtx, redisCollectionChange, fmt.Sprintf("%s,%s,%d", col.CollectionAddr, col.TokenID, col.Counter))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
func RemoveCollectionToCheck(col dtos.CollectionToCheck) error {
	redis := cache.GetRedis()
	redisCtx := cache.GetContext()
	cmd := redis.SRem(redisCtx, redisCollectionChange, fmt.Sprintf("%s,%s,%d", col.CollectionAddr, col.TokenID, col.Counter))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
func GetCollectionToCheck() ([]dtos.CollectionToCheck, error) {
	redis := cache.GetRedis()
	redisCtx := cache.GetContext()
	var colToCheck []dtos.CollectionToCheck
	cmd := redis.SMembers(redisCtx, redisCollectionChange)
	cols, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	// colData := strings.Split(cols, ",")
	// col = dtos.CollectionToCheck{CollectionAddr: colData[0], TokenID: colData[1]}
	for _, col := range cols {
		colData := strings.Split(col, ",")
		ci, err := strconv.Atoi(colData[2])
		if err != nil {
			return nil, err
		}
		colToCheck = append(colToCheck, dtos.CollectionToCheck{CollectionAddr: colData[0], TokenID: colData[1], Counter: ci})
	}
	return colToCheck, nil
}
func ClearCollectionToCheck() ([]dtos.CollectionToCheck, error) {
	redis := cache.GetRedis()
	redisCtx := cache.GetContext()
	cmdCard := redis.SCard(redisCtx, redisCollectionChange)
	size, err := cmdCard.Result()
	if err != nil {
		return nil, err
	}
	cmd := redis.SPopN(redisCtx, redisCollectionChange, size)
	cols, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	var colToCheck []dtos.CollectionToCheck
	for _, col := range cols {
		colData := strings.Split(col, ",")
		colToCheck = append(colToCheck, dtos.CollectionToCheck{CollectionAddr: colData[0], TokenID: colData[1]})
	}
	return colToCheck, nil
}
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
		statistics, innerErr := setStatisticsRaw(tokenId)
		if innerErr != nil {
			_, _ = redis.Del(redisCtx, redisExpireKey).Result()
			return nil, innerErr
		}
		return statistics, nil
	} else {
		return getStatisticsRaw(tokenId)
	}
}

func AddCollectionToCache(collectionId uint64, collectionName string, collectionFlags datatypes.JSON, tokenId string) (*dtos.CollectionCacheInfo, error) {
	db := cache.GetBolt()

	cacheInfo := dtos.CollectionCacheInfo{
		CollectionId:    collectionId,
		CollectionName:  collectionName,
		CollectionFlags: collectionFlags,
	}

	entryBytes, err := json.Marshal(&cacheInfo)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, innerErr := tx.CreateBucketIfNotExists(tokenIdToCollectionCacheInfo)
		if innerErr != nil {
			return innerErr
		}

		innerErr = bucket.Put([]byte(tokenId), entryBytes)
		return innerErr
	})

	return &cacheInfo, err
}

func GetCollectionCacheInfo(tokenId string) (*dtos.CollectionCacheInfo, error) {
	db := cache.GetBolt()

	var bytes []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tokenIdToCollectionCacheInfo)
		if bucket == nil {
			return errors.New("no bucket for collection cache")
		}

		bytes = bucket.Get([]byte(tokenId))
		return nil
	})
	if err != nil {
		return nil, err
	}

	var cacheInfo dtos.CollectionCacheInfo
	err = json.Unmarshal(bytes, &cacheInfo)
	if err != nil {
		return nil, err
	}

	return &cacheInfo, nil
}

func GetOrAddCollectionCacheInfo(tokenId string) (*dtos.CollectionCacheInfo, error) {
	cacheInfo, err := GetCollectionCacheInfo(tokenId)
	if err != nil {
		coll, innerErr := storage.GetCollectionByTokenId(tokenId)
		if innerErr != nil {
			return nil, innerErr
		}

		cacheInfo, innerErr = AddCollectionToCache(coll.ID, coll.Name, coll.Flags, coll.TokenID)
		if innerErr != nil {
			return nil, innerErr
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
