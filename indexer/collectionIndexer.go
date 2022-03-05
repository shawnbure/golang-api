package indexer

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/stats/collstats"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/emurmotol/ethconv"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CollectionIndexer struct {
	DeployerAddr string `json:"deployerAddr"`
	ElrondAPI    string `json:"elrondApi"`
	Logger       *log.Logger
	Delay        time.Duration // delay per request in second
}

func NewCollectionIndexer(deployerAddr string, elrondAPI string, delay uint64) (*CollectionIndexer, error) {
	l := log.New(os.Stderr, "", log.LUTC|log.LstdFlags|log.Lshortfile)
	return &CollectionIndexer{DeployerAddr: deployerAddr, ElrondAPI: elrondAPI,
		Delay:  time.Duration(delay),
		Logger: l}, nil
}

func (ci *CollectionIndexer) StartWorker() {
	logErr := ci.Logger
	resultNotAvailableCount := 0
	var colsToCheck []dtos.CollectionToCheck
	for {
	deployLoop:
		var foundDeployedContracts uint64 = 0
		logErr.Println("collection indexer loop")
		time.Sleep(time.Second * ci.Delay)
		deployerStat, err := storage.GetDeployerStat(ci.DeployerAddr)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				deployerStat, err = storage.CreateDeployerStat(ci.DeployerAddr)
				if err != nil {
					logErr.Println(err.Error())
					logErr.Println("something went wrong creating marketstat")
				}
			}
		}
		url := fmt.Sprintf("%s/accounts/%s/transactions?from=%d&withScResults=true&withLogs=false&order=asc",
			ci.ElrondAPI,
			ci.DeployerAddr,
			deployerStat.LastIndex)
		res, err := services.GetResponse(url)

		if err != nil {
			logErr.Println(err.Error())
			logErr.Println(url)
			continue
		}
		var ColResults []map[string]interface{}
		err = json.Unmarshal(res, &ColResults)
		if err != nil {
			logErr.Println(err.Error())
			logErr.Println("error unmarshal nfts deployer")
			continue
		}
		if len(ColResults) == 0 {
			goto colLoop
		}
		foundDeployedContracts += uint64(len(ColResults))
		for _, colR := range ColResults {
			if _, ok := colR["action"]; !ok {
				continue
			}
			name := (colR["action"].(map[string]interface{}))["name"].(string)
			if name == "deployNFTTemplateContract" && colR["status"] != "fail" {
				if _, ok := colR["results"]; !ok {
					goto deployLoop
				}
				mainDataStr := colR["data"].(string)
				mainData64Str, _ := base64.StdEncoding.DecodeString(mainDataStr)
				mainDatas := strings.Split(string(mainData64Str), "@")
				tokenIdHex := mainDatas[1]
				tokenIdStr, _ := hex.DecodeString(mainDatas[1])
				imageLink, _ := hex.DecodeString(mainDatas[4])
				metaLink, _ := hex.DecodeString(mainDatas[9])
				results := (colR["results"].([]interface{}))
				result := results[0]
				data := result.(map[string]interface{})["data"].(string)
				decodedData64, _ := base64.StdEncoding.DecodeString(data)
				decodedData := strings.Split(string(decodedData64), "@")
				hexByte, err := hex.DecodeString(decodedData[2])
				if err != nil {
					logErr.Println(err.Error())
					continue
				}
				byte32, err := bech32.ConvertBits(hexByte, 8, 5, true)
				if err != nil {
					logErr.Println(err.Error())
					continue
				}
				bech32Addr, err := bech32.Encode("erd", byte32)
				if err != nil {
					logErr.Println(err.Error())
					continue
				}
				colsToCheck = append(colsToCheck, dtos.CollectionToCheck{CollectionAddr: bech32Addr, TokenID: string(tokenIdStr)})
				tokenId, err := hex.DecodeString(tokenIdHex)
				dbCol, err := storage.GetCollectionByTokenId(string(tokenId))
				if err != nil {
					logErr.Println(err.Error())
					continue
				}
				dbCol.MetaDataBaseURI = string(metaLink)
				dbCol.TokenBaseURI = string(imageLink)
				err = storage.UpdateCollection(dbCol)
				if err != nil {
					logErr.Println(err.Error())
					continue
				}
				// get collection tx and check mint transactions
				_, err = storage.GetCollectionIndexer(bech32Addr)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						_, err = storage.CreateCollectionStat(bech32Addr)
						if err != nil {
							logErr.Println(err.Error())
							continue
						} else {
							continue
						}
					}
				}

			}

		}
	colLoop:
		if len(colsToCheck) == 0 { // TODO remove
			colsToCheck, err = collstats.GetCollectionToCheck()
			fmt.Println("collection to check from cache ", len(colsToCheck))
			if err != nil {
				logErr.Println(err.Error())
				continue
			}
		}
		for _, colObj := range colsToCheck {
		singleColLoop:
			var foundedTxsCount uint64 = 0

			col, err := storage.GetCollectionByTokenId(colObj.TokenID)
			if err != nil {
				logErr.Println("GetCollectionByTokenId", err.Error(), colObj.TokenID)
				continue
			}
			collectionIndexer, err := storage.GetCollectionIndexer(colObj.CollectionAddr)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					_, err = storage.CreateCollectionStat(colObj.CollectionAddr)
					if err != nil {
						logErr.Println(err.Error())
						logErr.Println("error create colleciton indexer")
						continue
					}
				} else {
					logErr.Println(err.Error())
					logErr.Println("error getting collection indexer")
					continue
				}
			}
			res, err := services.GetResponse(
				fmt.Sprintf("%s/accounts/%s/transactions?from=%d&withScResults=true&withLogs=false&order=asc",
					ci.ElrondAPI,
					collectionIndexer.CollectionAddr,
					collectionIndexer.LastIndex))
			if err != nil {
				logErr.Println(err.Error())
				logErr.Println("error creating request for get nfts deployer")
				if strings.Contains(err.Error(), "429") {
					time.Sleep(time.Second * 10)
					goto singleColLoop
				}
			}

			var ColResults []map[string]interface{}
			err = json.Unmarshal(res, &ColResults)
			if err != nil {
				logErr.Println(err.Error())
				logErr.Println("error unmarshal nfts deployer")
				continue
			}
			if len(ColResults) == 0 {
				continue
			}

			foundedTxsCount += uint64(len(ColResults))

			for _, colR := range ColResults {
				name := (colR["action"].(map[string]interface{}))["name"].(string)
				if name == "mintTokensThroughMarketplace" && colR["status"] != "fail" {
					if _, ok := colR["results"]; !ok {
						logErr.Println("results not available!")
						if colR["status"] == "pending" {
							time.Sleep(time.Second * 2)
							goto colLoop
						} else {
							resultNotAvailableCount++
							if resultNotAvailableCount > 3 {
								resultNotAvailableCount = 0
								continue
							}
						}
					}
					if colR["status"] == "pending" {
						time.Sleep(time.Second * 2)
						goto colLoop
					}
					results := (colR["results"].([]interface{}))
					if len(results) < 2 {
						logErr.Println("this tx wasn't good!", colR["originalTxHash"])
						continue
					}
					for _, r := range results {
						rMap := r.(map[string]interface{})
						rMapData, err := base64.StdEncoding.DecodeString(rMap["data"].(string))
						if err != nil {
							logErr.Println(err.Error())
							continue
						}
						if !strings.Contains(string(rMapData), "ESDTNFTTransfer") {
							continue
						}
						result := r
						nonceRes := result.(map[string]interface{})["data"]
						fmt.Println(nonceRes)
						nonceResBytes, err := base64.StdEncoding.DecodeString(nonceRes.(string))
						if err != nil {
							logErr.Println(err.Error())
							continue
						}
						nonceResArr := strings.Split(string(nonceResBytes), "@")
						nonceResString, err := hex.DecodeString(nonceResArr[2])
						if err != nil {
							logErr.Println(err.Error())
							continue
						}
						countRes, err := hex.DecodeString(nonceResArr[3])
						if err != nil {
							logErr.Println(err.Error())
							continue
						}
						// nonce, err := strconv.ParseUint(string(nonceResString[0]), 10, 64)
						nonce := uint64(nonceResString[0])
						count := int(countRes[0])

						resultMap := result.(map[string]interface{})
						decodedData, err := base64.StdEncoding.DecodeString(resultMap["data"].(string))
						if err != nil {
							logErr.Println(err.Error())
							continue
						}
						transferData := strings.Split(string(decodedData), "@")
						tokenIdByte, err := hex.DecodeString(transferData[1])
						if err != nil {
							logErr.Println(err.Error())
							continue
						}
						for i := 0; i < count; i++ {
							nonceStr := strconv.FormatInt(int64(nonce), 10)
							tokenId := string(tokenIdByte) + "-" + nonceStr
							_, err := storage.GetTokenByTokenIdAndNonceStr(tokenId, nonceStr)
							if err != nil {
								if err == gorm.ErrRecordNotFound {
									acc, err := storage.GetAccountByAddress(rMap["receiver"].(string))
									if err != nil {
										logErr.Println(err.Error())
										continue
									}
									price := colR["value"].(string)
									bigPrice, ok := big.NewInt(0).SetString(price, 10)
									if !ok {
										logErr.Println("conversion to bigInt failed for price")
										continue
									}
									fprice, err := ethconv.FromWei(bigPrice, ethconv.Ether)
									if err != nil {
										logErr.Println(err.Error())
										continue
									}
									priceFloat, _ := fprice.Float64()
									metaURI := col.MetaDataBaseURI
									imageURI := (col.TokenBaseURI)
									if strings.Contains(metaURI, ".json") {
										metaURI = strings.ReplaceAll(metaURI, ".json", "")
									}
									if !strings.Contains(metaURI, "https") {
										logErr.Println("old link", metaURI)
										b, _ := hex.DecodeString(metaURI)
										metaURI = string(b)
									}
									if !strings.Contains(imageURI, "https") {
										logErr.Println("old link", imageURI)
										b, _ := hex.DecodeString(imageURI)
										imageURI = string(b)
									}
									if strings.Contains(ci.ElrondAPI, "devnet") {
										imageURI = strings.Replace(imageURI, "https://gateway.pinata.cloud/ipfs/", "https://devnet-media.elrond.com/nfts/asset/", 1)
									} else {
										imageURI = strings.Replace(imageURI, "https://gateway.pinata.cloud/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
									}
									youbeiMeta := strings.Replace(metaURI, "https://gateway.pinata.cloud/ipfs/", "https://media.youbei.io/ipfs/", 1)
									url := fmt.Sprintf("%s/%s.json", youbeiMeta, nonceStr)
									tokenName := fmt.Sprintf("%s #%d", col.Name, int64(nonce))
									attrbs, err := services.GetResponse(url)
									if err != nil {
										logErr.Println(err.Error())
										if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "deadline") {
											err = storage.AddToken(&entities.Token{
												TokenID:      string(tokenIdByte),
												MintTxHash:   colR["txHash"].(string),
												CollectionID: col.ID,
												Nonce:        nonce,
												NonceStr:     nonceResArr[2],
												MetadataLink: string(youbeiMeta) + "/" + nonceStr + ".json",
												ImageLink:    string(imageURI) + "/" + nonceStr + ".png",
												TokenName:    tokenName,
												Attributes:   []byte{},
												OwnerId:      acc.ID,
												OnSale:       false,
												PriceString:  price,
												PriceNominal: priceFloat,
											})
											if err != nil {
												logErr.Println(err.Error())
												continue
											}
											continue
										}
										logErr.Println(err.Error(), url, col.MetaDataBaseURI, col.TokenBaseURI, col.ID)
										continue
									}
									metadataJSON := make(map[string]interface{})
									err = json.Unmarshal(attrbs, &metadataJSON)
									if err != nil {
										logErr.Println(err.Error(), string(url))
										continue
									}
									var attributes datatypes.JSON
									attributesBytes, err := json.Marshal(metadataJSON["attributes"])
									if err != nil {
										logErr.Println(err.Error())
										continue
									}
									err = json.Unmarshal(attributesBytes, &attributes)
									if err != nil {
										logErr.Println(err.Error())
										continue
									}
									err = storage.AddToken(&entities.Token{
										TokenID:      string(tokenIdByte),
										MintTxHash:   colR["txHash"].(string),
										CollectionID: col.ID,
										Nonce:        nonce,
										NonceStr:     nonceResArr[2],
										MetadataLink: string(youbeiMeta) + "/" + nonceStr + ".json",
										ImageLink:    string(imageURI) + "/" + nonceStr + ".png",
										TokenName:    tokenName,
										Attributes:   attributes,
										OwnerId:      acc.ID,
										OnSale:       false,
										PriceString:  price,
										PriceNominal: priceFloat,
									})
									if err != nil {
										logErr.Println(err.Error())
										continue
									}
								} else {
									logErr.Println(err.Error())
									continue
								}
							}
						}
					}
				}
			}
			collstats.RemoveCollectionToCheck(colsToCheck[0])
			collectionIndexer.LastIndex += foundedTxsCount
			_, err = storage.UpdateCollectionIndexer(collectionIndexer.LastIndex, collectionIndexer.CollectionAddr)
			if err != nil {
				_, err := storage.UpdateCollectionIndexer(collectionIndexer.LastIndex, collectionIndexer.CollectionAddr)
				if err != nil {
					logErr.Println(err.Error())
					continue
				}
			}
			goto singleColLoop
		}
		colsToCheck = []dtos.CollectionToCheck{}
		newStat, err := storage.UpdateDeployerIndexer(deployerStat.LastIndex+foundDeployedContracts, ci.DeployerAddr)
		if err != nil {
			logErr.Println(err.Error())
			logErr.Println("error update deployer index nfts ")
			continue
		}
		if newStat.LastIndex < deployerStat.LastIndex {
			logErr.Println("error something went wrong updating last index of deployer  ")
			continue
		}
	}
}
