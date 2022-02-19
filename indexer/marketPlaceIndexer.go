package indexer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	"gorm.io/gorm"
)

type MarketPlaceIndexer struct {
	MarketPlaceAddr string `json:"marketPlaceAddr"`
}

func NewMarketPlaceIndexer(marketPlaceAddr string) (*MarketPlaceIndexer, error) {
	return &MarketPlaceIndexer{MarketPlaceAddr: marketPlaceAddr}, nil
}

func (mpi *MarketPlaceIndexer) StartWorker() {
	client := &http.Client{}
	for {
		marketStat, err := storage.GetMarketPlaceIndexer()
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				marketStat, err = storage.CreateMarketPlaceStat()
				if err != nil {
					fmt.Println(err.Error())
					fmt.Println("something went wrong creating marketstat")
				}
			}
		}
		req, err := http.
			NewRequest("GET",
				fmt.Sprintf("https://devnet-api.elrond.com/accounts/%s/nfts?from=%d",
					mpi.MarketPlaceAddr,
					marketStat.LastIndex),
				nil)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error creating request for get nfts marketplace")
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error running request get nfts marketplace")
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			resp.Body.Close()
			fmt.Println("error readall response get nfts marketplace")
			continue
		}
		resp.Body.Close()
		if resp.Status != "200 OK" {
			fmt.Println("response not successful  get nfts marketplace", resp.Status, req.URL.RawPath)
			continue
		}
		var NFTs []entities.MarketPlaceNFT
		err = json.Unmarshal(body, &NFTs)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error unmarshal nfts marketplace")
			continue
		}
		if len(NFTs) == 0 {
			time.Sleep(time.Second * 10)
			continue
		}
		for _, nft := range NFTs {
			var token entities.Token
			token.TokenID = nft.Collection
			token.Nonce = nft.Nonce
			token.OnSale = true
			err := storage.UpdateTokenWhere(&token, "token_id=? AND nonce=?", nft.Collection, nft.Nonce)
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
					// continue
				}
				fmt.Println(err.Error())
				fmt.Println("error updating token ", fmt.Sprintf("tokenID %d", token.ID))
				continue
			}
		}
		newStat, err := storage.UpdateMarketPlaceIndexer(marketStat.LastIndex + 1)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error update marketplace index nfts ")
			continue
		}
		if newStat.LastIndex <= marketStat.LastIndex {
			fmt.Println("error something went wrong updating last index of marketplace  ")
			continue
		}
	}
}
