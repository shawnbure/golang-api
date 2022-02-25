package indexer

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/btcsuite/btcutil/bech32"
	"gorm.io/gorm"
)

type CollectionIndexer struct {
	DeployerAddr string `json:"deployerAddr"`
	ElrondAPI    string `json:"elrondApi"`
}

func NewCollectionIndexer(deployerAddr string, elrondAPI string) (*CollectionIndexer, error) {
	return &CollectionIndexer{DeployerAddr: deployerAddr, ElrondAPI: elrondAPI}, nil
}

func (ci *CollectionIndexer) StartWorker() {
	client := &http.Client{}
	for {
		log.Println("collection indexer loop")
		time.Sleep(time.Second * 2)
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
				fmt.Sprintf("%s/accounts/%s/transactions?from=%d&withScResults=true&withLogs=false",
					ci.ElrondAPI,
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
		var ColResults []map[string]interface{}
		err = json.Unmarshal(body, &ColResults)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error unmarshal nfts deployer")
			continue
		}
		if len(ColResults) == 0 {
			goto colLoop
		}
		for _, colR := range ColResults {
			name := (colR["action"].(map[string]interface{}))["name"].(string)
			if name == "deployNFTTemplateContract" {
				mainDataStr := colR["data"].(string)
				mainData64Str, _ := base64.StdEncoding.DecodeString(mainDataStr)
				mainDatas := strings.Split(string(mainData64Str), "@")
				tokenIdHex := mainDatas[1]
				imageLink, _ := hex.DecodeString(mainDatas[4])
				metaLink, _ := hex.DecodeString(mainDatas[9])
				results := (colR["results"].([]interface{}))
				result := results[0]
				data := result.(map[string]interface{})["data"].(string)
				decodedData64, _ := base64.StdEncoding.DecodeString(data)
				decodedData := strings.Split(string(decodedData64), "@")
				hexByte, err := hex.DecodeString(decodedData[2])
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				byte32, err := bech32.ConvertBits(hexByte, 8, 5, true)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				bech32Addr, err := bech32.Encode("erd", byte32)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				tokenId, err := hex.DecodeString(tokenIdHex)
				dbCol, err := storage.GetCollectionByTokenId(string(tokenId))
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				dbCol.MetaDataBaseURI = string(metaLink)
				dbCol.TokenBaseURI = string(imageLink)
				err = storage.UpdateCollection(dbCol)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				// get collection tx and check mint transactions
				_, err = storage.GetCollectionIndexer(bech32Addr)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						_, err = storage.CreateCollectionStat(bech32Addr)
						if err != nil {
							fmt.Println(err.Error())
							continue
						} else {
							continue
						}
					}
				}

			}

		}
	colLoop:
		collections, err := storage.GetAllCollections()
		if err != nil {
			fmt.Errorf(err.Error())
			fmt.Errorf("error running request get nfts deployer")
			continue
		}
		for _, col := range collections {

			collectionIndexer, err := storage.GetCollectionIndexer(col.ContractAddress)
			if err != nil {

				if err == gorm.ErrRecordNotFound {
					_, err = storage.CreateCollectionStat(col.ContractAddress)
					if err != nil {
						fmt.Errorf(err.Error())
						fmt.Errorf("error running request get nfts deployer")
						continue
					}
				} else {
					fmt.Errorf(err.Error())
					fmt.Errorf("error running request get nfts deployer")
					continue
				}
			}

			req, err := http.
				NewRequest("GET",
					fmt.Sprintf("%s/accounts/%s/transactions?from=%d&withScResults=true&withLogs=false",
						ci.ElrondAPI,
						collectionIndexer.CollectionAddr,
						collectionIndexer.LastIndex),
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
				fmt.Println("response not successful  get nfts deployer", resp.Status)
				continue
			}

			var ColResults []map[string]interface{}
			err = json.Unmarshal(body, &ColResults)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("error unmarshal nfts deployer")
				continue
			}
			if len(ColResults) == 0 {
				continue
			}
			for _, colR := range ColResults {
				name := (colR["action"].(map[string]interface{}))["name"].(string)
				if name == "mintTokensThroughMarketplace" {
					results := (colR["results"].([]interface{}))
					if len(results) < 2 {
						fmt.Println("this tx wasn't good!", colR["originalTxHash"])
						continue
					}
					for _, r := range results {
						rMap := r.(map[string]interface{})
						rMapData, err := base64.StdEncoding.DecodeString(rMap["data"].(string))
						if err != nil {
							fmt.Println(err.Error())
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
							fmt.Println(err.Error())
							continue
						}
						nonceResArr := strings.Split(string(nonceResBytes), "@")
						nonceResString, err := hex.DecodeString(nonceResArr[2])
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						countRes, err := hex.DecodeString(nonceResArr[3])
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						// nonce, err := strconv.ParseUint(string(nonceResString[0]), 10, 64)
						nonce := uint64(nonceResString[0])
						count := int(countRes[0])
						// if err != nil {
						// 	fmt.Println(err.Error())
						// 	continue
						// }
						resultMap := result.(map[string]interface{})
						decodedData, err := base64.StdEncoding.DecodeString(resultMap["data"].(string))
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						transferData := strings.Split(string(decodedData), "@")
						tokenIdByte, err := hex.DecodeString(transferData[1])
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						for i := 0; i < count; i++ {
							nonceStr := strconv.FormatInt(int64(nonce)+int64(count), 10)
							tokenId := string(tokenIdByte) + "-" + nonceStr
							_, err := storage.GetTokenByTokenIdAndNonceStr(tokenId, nonceStr)
							if err != nil {
								if err == gorm.ErrRecordNotFound {
									// lastToken, err := storage.GetLastNonceTokenByCollectionId(col.ID)
									// if err != nil {
									// 	if err == gorm.ErrRecordNotFound {
									// 		lastToken.Nonce = 1
									// 	} else {
									// 		fmt.Println(err.Error())
									// 		continue
									// 	}
									// } else {
									// 	//if this is a collection with tokens increment the nince
									// 	lastToken.Nonce = lastToken.Nonce + 1
									// }
									acc, err := storage.GetAccountByAddress(rMap["receiver"].(string))
									if err != nil {
										fmt.Println(err.Error())
										continue
									}
									price := colR["value"].(string)
									priceFloat, err := strconv.ParseFloat(price, 64)
									metaURI, err := hex.DecodeString(col.MetaDataBaseURI)
									imageURI, err := hex.DecodeString(col.TokenBaseURI)

									err = storage.AddToken(&entities.Token{
										TokenID:      string(tokenIdByte),
										MintTxHash:   colR["txHash"].(string),
										CollectionID: col.ID,
										Nonce:        nonce,
										NonceStr:     nonceResArr[2],
										MetadataLink: string(metaURI) + "/" + nonceStr + ".json",
										ImageLink:    string(imageURI) + "/" + nonceStr + ".png",
										TokenName:    col.Name,
										OwnerId:      acc.ID,
										OnSale:       false,
										PriceString:  price,
										PriceNominal: priceFloat,
									})
									if err != nil {
										fmt.Errorf(err.Error())
										continue
									}
								} else {
									fmt.Println(err.Error())
									continue
								}
							}
						}

					}

				}
			}
			collectionIndexer.LastIndex += 1
			_, err = storage.UpdateCollectionIndexer(collectionIndexer.LastIndex, collectionIndexer.CollectionAddr)
			if err != nil {
				_, err := storage.UpdateCollectionIndexer(collectionIndexer.LastIndex, collectionIndexer.CollectionAddr)
				if err != nil {
					fmt.Errorf(err.Error())
					continue
				}
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
