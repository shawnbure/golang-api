package gatherer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/ENFT-DAO/youbei-api/utils"
	"go.uber.org/zap"
)

const (
	RarityUpdaterDurationMilli                = 300
	RarityUpdaterAllCollectionDurationMinutes = 30
	RarityUpdaterTokenDurationMilli           = 50
)

func syncRarityRunner(cha chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-cha:
			ticker.Stop()
			return
		case <-ticker.C:
			//getMissedRarity()
			// computeRarityScorePreCollection()
			// computeRanks()
		}
	}
}

func getMissedRarity() {
	for {
		// Get the last unsynced value from database
		success := true
		tokenInstance, err := storage.GetOldTokenWithZeroRarity()
		if err == nil {
			// token with this condition exist
			if tokenInstance.MetadataLink != "" {
				success = false

				content, err := OnePage(tokenInstance.MetadataLink)
				if err == nil {
					var metadataJSON map[string]interface{}
					err1 := json.Unmarshal([]byte(content), &metadataJSON)
					if err1 == nil {
						if _, ok := metadataJSON["rarity"]; ok {
							rarityBody, err2 := json.Marshal(metadataJSON["rarity"].(map[string]interface{}))
							if err2 == nil {
								var rarity entities.TokenRarity
								if err := json.Unmarshal(rarityBody, &rarity); err == nil {
									tokenInstance.RarityScoreNorm = rarity.RarityScoreNormed
									tokenInstance.RarityUsedTraitCount = uint(rarity.UsedTraitsCount)
									tokenInstance.RarityScore = rarity.RarityScore
									tokenInstance.IsRarityInserted = true
									success = true
								}
							}
						}
					}
				}
			} else {
				logInstance.Error("Cannot get metadata from link ", err)
			}

			if !success {
				tokenInstance.IsRarityInserted = false
			}

			err3 := storage.UpdateToken(&tokenInstance)
			if err3 != nil {
				logInstance.Error("Cannot update token info ", err3)
			}
		}

		time.Sleep(RarityUpdaterDurationMilli * time.Millisecond)
	}
}

func OnePage(link string) (string, error) {
	res, err := http.Get(link)
	if err != nil {
		return "", err
	}

	if res.StatusCode == 200 {
		content, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		return string(content), nil
	}

	return "", errors.New(res.Status)
}

func computeRarityScorePreCollection() {
	// Get all collections from database
	collections, err := storage.GetAllCollections()
	if err == nil {
		for _, col := range collections {
			// get tokens associated with collection
			tokens, err := storage.GetTokensByCollectionIdNotRanked(col.ID)
			if err == nil {
				traits := make(map[string]map[string]int)
				traitsInTokens := make(map[uint64][]string)
				traitsRank := make(map[string]map[string]float64)

				totalTokens := len(tokens)

				for _, token := range tokens {
					traitsInToken := []string{}

					v := []map[string]interface{}{}
					bytes, _ := token.Attributes.MarshalJSON()
					err1 := json.Unmarshal(bytes, &v)
					if err1 == nil {
						for _, item := range v {
							var key1, key2 string
							if _, ok := item["trait_type"]; !ok {
								for k, v := range item {
									key1 = k
									key2 = v.(string)
								}
							} else {
								key1 = item["trait_type"].(string)
								key2 = item["value"].(string)
							}
							if val, ok := traits[key1]; ok {
								if val2, ok2 := val[key2]; ok2 {
									traits[key1][key2] = val2 + 1
								} else {
									traits[key1][key2] = 1
								}
							} else {
								traits[key1] = map[string]int{}
								traits[key1][key2] = 1
							}

							key := fmt.Sprintf("%v$$$$$%v", key1, key2)
							traitsInToken = append(traitsInToken, key)
						}
						if len(traitsInTokens[token.ID]) == 0 {
							traitsInTokens[token.ID] = traitsInToken
						}
					}
				}

				for key, val := range traits {
					traitsRank[key] = make(map[string]float64)
					for key2, val2 := range val {
						traitsRank[key][key2] = float64(float64(val2) / float64(totalTokens))
					}
				}

				for _, token := range tokens {
					localTraits := traitsInTokens[token.ID]

					totalRank := float64(0)
					for key, _ := range traits {
						index := utils.IndexInArray(localTraits, key)
						if index >= 0 {
							splittedKeys := strings.Split(localTraits[index], "$$$$$")
							key2 := splittedKeys[1]

							totalRank += 1 / traitsRank[key][key2]
						}
					}

					token.RarityScoreNorm = 0
					token.RarityUsedTraitCount = uint(len(traitsInTokens[token.ID]))
					token.RarityScore = totalRank
					token.IsRarityInserted = true

					err3 := storage.UpdateToken(&token)
					if err3 != nil {
						logInstance.Error("Cannot update token info ", err3)
					}
					time.Sleep(RarityUpdaterTokenDurationMilli * time.Millisecond)
				}
			}
		}
	}
}

func computeRanks() {
	// Get all collections from database
	collections, err := storage.GetAllCollections()
	if err == nil {
		for _, col := range collections {
			// get tokens associated with collection
			count, err := storage.GetTokensWithNoRankCount(col.ID)
			if err != nil {
				zlog.Error("critical", zap.Error(err))
			}
			if count == 0 {
				continue
			}
			tokens, err := storage.GetTokensByCollectionIdOrderedByRarityScore(col.ID, "DESC")
			if len(tokens) == 0 {
				continue
			}
			if err == nil {
				lastToken := tokens[0]
				for i, token := range tokens {
					if token.RarityScore < lastToken.RarityScore {
						token.Rank = uint(i + 1)
					} else {
						if lastToken.Rank == 0 {
							token.Rank = uint(i + 1)
						} else {
							token.Rank = lastToken.Rank
						}
					}
					lastToken = token
					storage.UpdateToken(&token)
				}
			}
		}
	}
}
