package indexer

import (
	"encoding/base64"
	"encoding/hex"
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
	"github.com/btcsuite/btcutil/bech32"
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
	log := ci.Logger
	var colsToCheck []struct {
		Addr    string
		TokenID string
	}
	for {
		var foundMintedTokens uint64 = 0
		var foundDeployedContracts uint64 = 0
		log.Println("collection indexer loop")
		time.Sleep(time.Second * ci.Delay)
		deployerStat, err := storage.GetDeployerStat(ci.DeployerAddr)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				deployerStat, err = storage.CreateDeployerStat(ci.DeployerAddr)
				if err != nil {
					log.Println(err.Error())
					log.Println("something went wrong creating marketstat")
				}
			}
		}
		url := fmt.Sprintf("%s/accounts/%s/transactions?from=%d&withScResults=true&withLogs=false&order=asc",
			ci.ElrondAPI,
			ci.DeployerAddr,
			deployerStat.LastIndex)
		res, err := services.GetResponse(url)

		if err != nil {
			log.Println(err.Error())
			log.Println(url)
			continue
		}
		var ColResults []map[string]interface{}
		err = json.Unmarshal(res, &ColResults)
		if err != nil {
			log.Println(err.Error())
			log.Println("error unmarshal nfts deployer")
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
					log.Println(err.Error())
					continue
				}
				byte32, err := bech32.ConvertBits(hexByte, 8, 5, true)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				bech32Addr, err := bech32.Encode("erd", byte32)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				colsToCheck = append(colsToCheck, struct {
					Addr    string
					TokenID string
				}{bech32Addr, string(tokenIdStr)})
				tokenId, err := hex.DecodeString(tokenIdHex)
				dbCol, err := storage.GetCollectionByTokenId(string(tokenId))
				if err != nil {
					log.Println(err.Error())
					continue
				}
				dbCol.MetaDataBaseURI = string(metaLink)
				dbCol.TokenBaseURI = string(imageLink)
				err = storage.UpdateCollection(dbCol)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				// get collection tx and check mint transactions
				_, err = storage.GetCollectionIndexer(bech32Addr)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						_, err = storage.CreateCollectionStat(bech32Addr)
						if err != nil {
							log.Println(err.Error())
							continue
						} else {
							continue
						}
					}
				}

			}

		}
	colLoop:
		// collections, err := storage.GetAllCollections()
		// if err != nil {
		// 	log.Println(err.Error())
		// 	log.Println("error running request get nfts deployer")
		// 	continue
		// }
		for _, colObj := range colsToCheck {
			col, err := storage.GetCollectionByTokenId(colObj.TokenID)
			if err != nil {
				log.Println("GetCollectionByTokenId", err.Error(), colObj.TokenID)
				continue
			}
			collectionIndexer, err := storage.GetCollectionIndexer(colObj.Addr)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					_, err = storage.CreateCollectionStat(colObj.Addr)
					if err != nil {
						log.Println(err.Error())
						log.Println("error create colleciton indexer")
						continue
					}
				} else {
					log.Println(err.Error())
					log.Println("error getting collection indexer")
					continue
				}
			}
			res, err := services.GetResponse(
				fmt.Sprintf("%s/accounts/%s/transactions?from=%d&withScResults=true&withLogs=false&order=asc",
					ci.ElrondAPI,
					collectionIndexer.CollectionAddr,
					collectionIndexer.LastIndex))
			if err != nil {
				log.Println(err.Error())
				log.Println("error creating request for get nfts deployer")
				continue
			}

			var ColResults []map[string]interface{}
			err = json.Unmarshal(res, &ColResults)
			if err != nil {
				log.Println(err.Error())
				log.Println("error unmarshal nfts deployer")
				continue
			}
			if len(ColResults) == 0 {
				continue
			}
			foundMintedTokens += uint64(len(ColResults))

			for _, colR := range ColResults {
				name := (colR["action"].(map[string]interface{}))["name"].(string)
				if name == "mintTokensThroughMarketplace" && colR["status"] != "fail" {
					results := (colR["results"].([]interface{}))
					if len(results) < 2 {
						log.Println("this tx wasn't good!", colR["originalTxHash"])
						continue
					}
					for _, r := range results {
						rMap := r.(map[string]interface{})
						rMapData, err := base64.StdEncoding.DecodeString(rMap["data"].(string))
						if err != nil {
							log.Println(err.Error())
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
							log.Println(err.Error())
							continue
						}
						nonceResArr := strings.Split(string(nonceResBytes), "@")
						nonceResString, err := hex.DecodeString(nonceResArr[2])
						if err != nil {
							log.Println(err.Error())
							continue
						}
						countRes, err := hex.DecodeString(nonceResArr[3])
						if err != nil {
							log.Println(err.Error())
							continue
						}
						// nonce, err := strconv.ParseUint(string(nonceResString[0]), 10, 64)
						nonce := uint64(nonceResString[0])
						count := int(countRes[0])

						resultMap := result.(map[string]interface{})
						decodedData, err := base64.StdEncoding.DecodeString(resultMap["data"].(string))
						if err != nil {
							log.Println(err.Error())
							continue
						}
						transferData := strings.Split(string(decodedData), "@")
						tokenIdByte, err := hex.DecodeString(transferData[1])
						if err != nil {
							log.Println(err.Error())
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
										log.Println(err.Error())
										continue
									}
									price := colR["value"].(string)
									priceFloat, err := strconv.ParseFloat(price, 64)
									metaURI := col.MetaDataBaseURI
									imageURI := (col.TokenBaseURI)
									if !strings.Contains("https", metaURI) {
										b, _ := hex.DecodeString(metaURI)
										metaURI = string(b)
									}
									if !strings.Contains("https", imageURI) {
										b, _ := hex.DecodeString(imageURI)
										imageURI = string(b)
									}

									url = string(metaURI) + "/" + nonceStr + ".json"
									attrbs, err := services.GetResponse(url)
									if err != nil {
										log.Println(err.Error(), url, col.MetaDataBaseURI, col.TokenBaseURI, col.ID)
										continue
									}
									metadataJSON := make(map[string]interface{})
									err = json.Unmarshal(attrbs, &metadataJSON)
									if err != nil {
										log.Println(err.Error())
										continue
									}
									var attributes datatypes.JSON
									attributesBytes, err := json.Marshal(metadataJSON["attributes"])
									if err != nil {
										log.Println(err.Error())
										continue
									}
									err = json.Unmarshal(attributesBytes, &attributes)
									if err != nil {
										log.Println(err.Error())
										continue
									}
									err = storage.AddToken(&entities.Token{
										TokenID:      string(tokenIdByte),
										MintTxHash:   colR["txHash"].(string),
										CollectionID: col.ID,
										Nonce:        nonce,
										NonceStr:     nonceResArr[2],
										MetadataLink: string(metaURI) + "/" + nonceStr + ".json",
										ImageLink:    string(imageURI) + "/" + nonceStr + ".png",
										TokenName:    col.Name,
										Attributes:   attributes,
										OwnerId:      acc.ID,
										OnSale:       false,
										PriceString:  price,
										PriceNominal: priceFloat,
									})
									if err != nil {
										log.Println(err.Error())
										continue
									}
								} else {
									log.Println(err.Error())
									continue
								}
							}
						}
					}
				}
			}
			collectionIndexer.LastIndex += foundMintedTokens
			_, err = storage.UpdateCollectionIndexer(collectionIndexer.LastIndex, collectionIndexer.CollectionAddr)
			if err != nil {
				_, err := storage.UpdateCollectionIndexer(collectionIndexer.LastIndex, collectionIndexer.CollectionAddr)
				if err != nil {
					log.Println(err.Error())
					continue
				}
			}
		}

		newStat, err := storage.UpdateDeployerIndexer(deployerStat.LastIndex+foundDeployedContracts, ci.DeployerAddr)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error update deployer index nfts ")
			continue
		}
		if newStat.LastIndex < deployerStat.LastIndex {
			fmt.Println("error something went wrong updating last index of deployer  ")
			continue
		}
	}
}
