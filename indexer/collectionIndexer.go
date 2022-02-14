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

type CollectionIndexer struct {
	DeployerAddr string `json:"deployerAddr"`
}

func NewCollectionIndexer(deployerAddr string) (*CollectionIndexer, error) {
	return &CollectionIndexer{DeployerAddr: deployerAddr}, nil
}

func (ci *CollectionIndexer) StartWorker() {
	client := &http.Client{}
	for {
		deployerStat, err := storage.GetDeployerStat(ci.DeployerAddr)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				deployerStat, err = storage.CreateDeployerStat(ci.DeployerAddr)
				if err != nil {
					fmt.Println(err.Error())
					fmt.Println("something went wrong creating marketstat")
				}
			}
		}
		req, err := http.
			NewRequest("GET",
				fmt.Sprintf("https://devnet-api.elrond.com/accounts/%s/nfts?from=%d",
					ci.DeployerAddr,
					deployerStat.LastIndex),
				nil)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error creating request for get nfts deployer")
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error running request get nfts deployer")
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			resp.Body.Close()
			fmt.Println("error readall response get nfts deployer")
			continue
		}
		resp.Body.Close()
		if resp.Status != "200 OK" {
			fmt.Println("response not successful  get nfts deployer")
			continue
		}
		var NFTs []entities.MarketPlaceNFT
		err = json.Unmarshal(body, &NFTs)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error unmarshal nfts deployer")
			continue
		}
		if len(NFTs) == 0 {
			time.Sleep(time.Minute * 2)
		}
		for _, nft := range NFTs {
			var token entities.Token
			token.TokenID = nft.Collection
			token.Nonce = nft.Nonce
			token.OnSale = true
			err := storage.UpdateTokenWhere(&token, "token_id=? AND nonce=?", nft.Collection, nft.Nonce)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					fmt.Println("token not added to db ", token.TokenID, fmt.Sprintf("nonce %d", token.Nonce))
					continue
				}
				fmt.Println(err.Error())
				fmt.Println("error updating token ", fmt.Sprintf("tokenID", token.ID))
				continue
			}
		}
		newStat, err := storage.UpdateDeployerIndexer(deployerStat.LastIndex+1, ci.DeployerAddr)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error update deployer index nfts ")
			continue
		}
		if newStat.LastIndex <= deployerStat.LastIndex {
			fmt.Println("error something went wrong updating last index of deployer  ")
			continue
		}
	}
}
