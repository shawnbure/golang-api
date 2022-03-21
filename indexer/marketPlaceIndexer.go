package indexer

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/emurmotol/ethconv"
	"gorm.io/gorm"
)

type MarketPlaceIndexer struct {
	MarketPlaceAddr string `json:"marketPlaceAddr"`
	ElrondAPI       string `json:"elrondAPI"`
	Logger          *log.Logger
	Delay           time.Duration // delay between each call
}

func NewMarketPlaceIndexer(marketPlaceAddr string, elrondAPI string, delay uint64) (*MarketPlaceIndexer, error) {
	lerr := log.New(os.Stderr, "", log.LUTC|log.LstdFlags|log.Lshortfile)
	return &MarketPlaceIndexer{MarketPlaceAddr: marketPlaceAddr, ElrondAPI: elrondAPI, Logger: lerr, Delay: time.Duration(delay)}, nil
}

func (mpi *MarketPlaceIndexer) StartWorker() {
	lerr := mpi.Logger
	lastHashMet := false
	lastIndex := 0
	for {
		var foundResults uint64 = 0
		time.Sleep(time.Second * mpi.Delay)
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
		body, err := services.GetResponse(fmt.Sprintf("%s/accounts/%s/sc-results?from=%d&size=100", // sc-result endpoint doesn't have order!
			mpi.ElrondAPI,
			mpi.MarketPlaceAddr,
			lastIndex,
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
		foundResults += uint64(len(txResult))
		for _, tx := range txResult {
		txloop:
			orgtxByte, err := services.GetResponse(fmt.Sprintf("%s/transactions/%s", mpi.ElrondAPI, tx.OriginalTxHash))
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			var orgTx entities.TransactionBC
			err = json.Unmarshal(orgtxByte, &orgTx)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			if orgTx.TxHash == marketStat.LastHash {
				lastHashMet = true
				lastIndex = 0
			} else {
				marketStat.LastHash = orgTx.TxHash
			}
			if orgTx.Status == "fail" {
				continue
			}
			if orgTx.Status == "pending" {
				goto txloop
			}
			token := &entities.Token{}
			data, err := base64.StdEncoding.DecodeString(tx.Data)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			orgData, err := base64.StdEncoding.DecodeString(orgTx.Data)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			isWithdrawn := strings.Contains(string(orgData), "withdrawNft")
			isOnSale := strings.Contains(string(data), "putNftForSale")
			isBuyNft := strings.Contains(string(orgData), "buyNft")
			if !isOnSale && !isBuyNft && !isWithdrawn {
				continue
			}
			body, err := services.GetResponse(fmt.Sprintf("%s/transactions?hashes=%s&order=asc", mpi.ElrondAPI, tx.OriginalTxHash))
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			var txBody []entities.TransactionBC
			err = json.Unmarshal(body, &txBody)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			mainTxDataStr := txBody[0].Data
			mainTxData, err := base64.StdEncoding.DecodeString(mainTxDataStr)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			mainDataParts := strings.Split(string(mainTxData), "@")
			hexTokenId := mainDataParts[1]
			tokenId, err := hex.DecodeString(hexTokenId)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			hexNonce := mainDataParts[2]
			// nonce, err := hex.DecodeString(hexNonce)
			data, err = base64.StdEncoding.DecodeString(tx.Data)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}
			dataStr := string(data)
			dataParts := strings.Split(dataStr, "@")
			price, ok := big.NewInt(0).SetString(dataParts[1], 16)
			if !ok {
				lerr.Println("can not convert price", price, dataParts[1])
				continue
			}
			fprice, err := ethconv.FromWei(price, ethconv.Ether)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}

			txTimestamp := orgTx.Timestamp
			token, err = storage.GetTokenByTokenIdAndNonceStr(string(tokenId), hexNonce)
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					lerr.Println(err.Error())
					continue
				} else {
					goto txloop
				}
			}
			if token.LastMarketTimestamp < txTimestamp {
				if isOnSale {
					token.OnSale = true
					token.Status = "List"
					token.PriceString = fprice.String()
					token.PriceNominal, _ = fprice.Float64()
				} else if isWithdrawn {
					token.OnSale = false
					token.Status = "Withdrawn"
				} else if isBuyNft {
					token.OnSale = false
					token.Status = "Bought"
				}
				token.LastMarketTimestamp = txTimestamp
				err = storage.UpdateTokenWhere(token, map[string]interface{}{
					"OnSale":              token.OnSale,
					"Status":              token.Status,
					"PriceString":         token.PriceString,
					"PriceNominal":        token.PriceNominal,
					"LastMarketTimestamp": txTimestamp,
				}, "token_id=? AND nonce_str=?", tokenId, hexNonce)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						continue
					}
					lerr.Println(err.Error())
					lerr.Println("error updating token ", fmt.Sprintf("tokenID %d", token.ID))
					continue
				}
			}
		}
		// marketStat, err = storage.UpdateMarketPlaceIndexer(marketStat.LastIndex + foundResults)
		// if err != nil {
		// 	lerr.Println(err.Error())
		// 	lerr.Println("error update marketplace index nfts ")
		// 	continue
		// }
		marketStat, err = storage.UpdateMarketPlaceHash(marketStat.LastHash)
		if err != nil {
			lerr.Println(err.Error())
			lerr.Println("error update marketplace index nfts ")
			continue
		}
		// if newStat.LastIndex <= marketStat.LastIndex {
		// 	lerr.Println("error something went wrong updating last index of marketplace  ")
		// 	continue
		// }
		if !lastHashMet {
			lastIndex += 100
		}
	}
}
