package indexer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/storage"
	"gorm.io/gorm"
)

var lerr *log.Logger = log.New(os.Stderr, "", 1)

type MarketPlaceIndexer struct {
	MarketPlaceAddr string `json:"marketPlaceAddr"`
	ElrondAPI       string `json:"elrondAPI"`
}

func NewMarketPlaceIndexer(marketPlaceAddr string, elrondAPI string) (*MarketPlaceIndexer, error) {
	return &MarketPlaceIndexer{MarketPlaceAddr: marketPlaceAddr, ElrondAPI: elrondAPI}, nil
}

func (mpi *MarketPlaceIndexer) StartWorker() {
	for {
		time.Sleep(time.Second * 1)
		marketStat, err := storage.GetMarketPlaceIndexer()
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				marketStat, err = storage.CreateMarketPlaceStat()
				if err != nil {
					lerr.Println(err.Error())
					lerr.Println("something went wrong creating marketstat")
					continue
				}
			}
		}
		body, err := services.GetResponse(fmt.Sprintf("%s/accounts/%s/sc-results?from=%d",
			mpi.ElrondAPI,
			mpi.MarketPlaceAddr,
			marketStat.LastIndex,
		))
		if err != nil {
			lerr.Println(err.Error())
			continue
		}

		var txResult []entities.SCResult
		err = json.Unmarshal(body, &txResult)
		if err != nil {
			lerr.Println(err.Error())
			lerr.Println("error unmarshal nfts marketplace")
			continue
		}
		if len(txResult) == 0 {
			continue
		}
		for _, tx := range txResult {
			var token entities.Token
			data, err := base64.StdEncoding.DecodeString(tx.Data)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			if !strings.Contains(string(data), "putNftForSale") {
				continue
			}
			body, err := services.GetResponse(fmt.Sprintf("%s/transactions?hashes=%s", mpi.ElrondAPI, tx.OriginalTxHash))
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			var txBody entities.TransactionBC
			err = json.Unmarshal(body, &txBody)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			data, err = base64.StdEncoding.DecodeString(tx.Data)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			dataStr := string(data)
			dataParts := strings.Split(dataStr, "@")
			tokenId := dataParts[1]
			nonce, err := strconv.ParseUint(dataParts[2], 10, 64)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			token.TokenID = tokenId
			token.Nonce = nonce
			token.OnSale = true
			err = storage.UpdateTokenWhere(&token, "token_id=? AND nonce=?", tokenId, nonce)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					// col, err := storage.GetCollectionByTokenId(nft.Collection)
					// if err != nil {
					// 	fmt.Println(err.Error())
					// 	fmt.Println("error get collection by tokenid", fmt.Sprintf("tokenID %d", token.ID))
					// 	continue
					// }
					// token.PriceNominal = col.MintPricePerTokenNominal
					// token.PriceString = col.MintPricePerTokenString
					// token.RoyaltiesPercent = float64(nft.Royalties)
					// twoParts := strings.Split(nft.URL, "/")
					// metadataLink, err := hex.DecodeString(twoParts[0])
					// if err != nil {
					// 	fmt.Println(err.Error())
					// 	continue
					// }
					// token.ImageLink = fmt.Sprintf(string(metadataLink)+"/%d.%s", nft.Nonce, "png")     // TODO image type should come from collection probably
					// token.MetadataLink = fmt.Sprintf(string(metadataLink)+"/%d.%s", nft.Nonce, "json") // TODO image type should come from collection probably
					// token.TokenName = nft.Collection
					// token.CollectionID = col.ID
					// // get owner
					// req, err := http.
					// 	NewRequest("GET",
					// 		fmt.Sprintf("https://devnet-api.elrond.com/accounts/%s",
					// 			nft.Creator,
					// 		),
					// 		nil)
					// if err != nil {
					// 	fmt.Println(err.Error())
					// 	fmt.Println("error creating request for get nfts marketplace")
					// 	continue
					// }
					// resp, err := client.Do(req)
					// if err != nil {
					// 	fmt.Println(err.Error())
					// 	fmt.Println("error running request get nfts marketplace")
					// 	continue
					// }

					// body, err := ioutil.ReadAll(resp.Body)
					// if err != nil {
					// 	fmt.Println(err.Error())
					// 	resp.Body.Close()
					// 	fmt.Println("error readall response get nfts marketplace")
					// 	continue
					// }
					// resp.Body.Close()
					// if resp.Status != "200 OK" {
					// 	fmt.Println("response not successful  get nfts marketplace")
					// 	continue
					// }
					// var accountObj map[string]interface{}
					// err = json.Unmarshal(body, &accountObj)
					// if err != nil {
					// 	fmt.Println(err.Error())
					// 	continue
					// }
					// ownerAddr := accountObj["ownerAddress"].(string)
					// acc, err := storage.GetAccountByAddress(ownerAddr)
					// if err != nil {
					// 	fmt.Println(err.Error())
					// 	fmt.Println("error get account by address", fmt.Sprintf("nft create %s", nft.Creator))
					// 	continue
					// }
					// token.OwnerId = acc.ID
					// fmt.Println("token not added to db ", token.TokenID, fmt.Sprintf("nonce %d", token.Nonce))
					continue
				}
				lerr.Println(err.Error())
				lerr.Println("error updating token ", fmt.Sprintf("tokenID %d", token.ID))
				continue
			}
		}
		newStat, err := storage.UpdateMarketPlaceIndexer(marketStat.LastIndex + 1)
		if err != nil {
			lerr.Println(err.Error())
			lerr.Println("error update marketplace index nfts ")
			continue
		}
		if newStat.LastIndex <= marketStat.LastIndex {
			lerr.Println("error something went wrong updating last index of marketplace  ")
			continue
		}
	}
}
