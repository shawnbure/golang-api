package gatherer

import (
	"encoding/json"
	"errors"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	RarityUpdaterDurationMilli = 300
)

func syncRarityRunner(cha chan bool) {
	ticker := time.NewTicker(60 * time.Minute)
	for {
		select {
		case <-cha:
			ticker.Stop()
			return
		case <-ticker.C:
			getMissedRarity()
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
