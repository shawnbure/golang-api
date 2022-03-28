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

		var deployerTxs []entities.TransactionBC
		err = json.Unmarshal(res, &deployerTxs)
		if err != nil {
			logErr.Println(err.Error())
			logErr.Println("error unmarshal nfts deployer")
			continue
		}
		if len(deployerTxs) == 0 {
			goto colLoop
		}
		foundDeployedContracts += uint64(len(deployerTxs))
		for _, colR := range deployerTxs {
			if colR.Action.Name == "" {
				continue
			}
			name := colR.Action.Name
			if name == "deployNFTTemplateContract" && colR.Status != "fail" {
				if len(colR.Results) == 0 {
					goto deployLoop
				}
				mainDataStr := colR.Data
				mainData64Str, _ := base64.StdEncoding.DecodeString(mainDataStr)
				mainDatas := strings.Split(string(mainData64Str), "@")
				tokenIdHex := mainDatas[1]
				tokenIdStr, _ := hex.DecodeString(mainDatas[1])
				imageLink, _ := hex.DecodeString(mainDatas[4])
				metaLink, _ := hex.DecodeString(mainDatas[9])
				results := colR.Results
				result := results[0]
				data := result.Data
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
				if err != nil {
					logErr.Println(err.Error())
					continue
				}

				dbCol, err := storage.GetCollectionByTokenId(string(tokenId))
				if err != nil {
					logErr.Println(err.Error())
					continue
				}

				dbCol.MetaDataBaseURI = string(metaLink)
				dbCol.TokenBaseURI = string(imageLink)
				metaInfoByte, err := services.GetResponse(dbCol.MetaDataBaseURI + "/1.json")

				if err != nil {
					logErr.Println(err.Error())
					continue
				}

				metaInfo := map[string]interface{}{}
				err = json.Unmarshal(metaInfoByte, &metaInfo)
				if err != nil {
					logErr.Println(err.Error())
					continue
				}

				dbCol.Description = metaInfo["description"].(string)
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
		cols, err := storage.GetAllCollections()
		if err != nil {
			logErr.Println(err.Error())
			continue
		}
		for _, colObj := range cols {
		singleColLoop:
			var foundedTxsCount uint64 = 0
			var minted uint64 = 0

			col, err := storage.GetCollectionByTokenId(colObj.TokenID)
			if err != nil {
				logErr.Println("GetCollectionByTokenId", err.Error(), colObj.TokenID)
				continue
			}
			if colObj.ContractAddress == "" {
				colDetail, err := services.GetCollectionDetailBC(col.TokenID, ci.ElrondAPI)
				if err != nil {
					continue
				}
				var address string
				for _, role := range colDetail.Roles {
					rolesStr, ok := role["roles"].([]interface{})
					if ok {
						for _, roleStr := range rolesStr {
							if strings.EqualFold(roleStr.(string), "ESDTRoleNFTCreate") {
								address = role["address"].(string)
							}
						}
					}
				}
				colObj.ContractAddress = address
				colObj.Name = colDetail.Name
				err = services.UpdateCollectionWithAddress(&colObj, map[string]interface{}{
					"Name":            colObj.Name,
					"ContractAddress": colObj.ContractAddress,
				})
				if err != nil {
					continue
				}
			}
			collectionIndexer, err := storage.GetCollectionIndexer(colObj.ContractAddress)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					collectionIndexer, err = storage.CreateCollectionStat(colObj.ContractAddress)
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
				if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "deadline") {
					time.Sleep(time.Second * 10)
					goto singleColLoop
				}
			}

			var collectionTxs []entities.TransactionBC
			err = json.Unmarshal(res, &collectionTxs)
			if err != nil {
				logErr.Println(err.Error())
				logErr.Println("error unmarshal nfts deployer")
				continue
			}
			if len(collectionTxs) == 0 {
				logErr.Println("ColResults no collection related tx found")
				continue
			}

			foundedTxsCount += uint64(len(collectionTxs))
			fmt.Println("entering colResults loop")

			for _, colR := range collectionTxs {

				dataString := colR.Data
				dataByte, err := base64.StdEncoding.DecodeString(dataString)
				if err != nil {
					logErr.Println(err.Error())
					continue
				}
				dataStr := string(dataByte)
				if strings.Contains(dataStr, "mintTokensThroughMarketplace") && colR.Status != "fail" {
					fmt.Println("mint found")
					if len(colR.Results) == 0 {
						logErr.Println("results not available!")
						if colR.Status == "pending" {
							goto singleColLoop
						} else {
							resultNotAvailableCount++
							if resultNotAvailableCount > 3 {
								resultNotAvailableCount = 0
							}
						}
						// goto singleColLoop
						continue
					}
					if colR.Status == "pending" {
						if colR.Timestamp-uint64(time.Now().UTC().Unix()) < 10*60 { // i've seen stuff on devnet pending for hours
							goto singleColLoop
						}
					}
					if colR.Results == nil {
						logErr.Println("CRITICAL  result was nil")
						time.Sleep(time.Second * 2)
						continue
					}
					results := colR.Results
					if len(results) < 2 {
						logErr.Println("CRITICAL this tx wasn't good!", colR.TxHash)
						continue
					}
					dataParts := strings.Split(dataStr, "@")
					mintCount, err := strconv.ParseInt(dataParts[1], 16, 64)
					if err != nil {
						logErr.Println("CRITICAL  mint count conversion failed")
						continue
					}
					for _, r := range results {
						decodedData, err := base64.StdEncoding.DecodeString(r.Data)
						if err != nil {
							logErr.Println("CRITICAL ", err.Error())
							continue
						}
						r.Data = string(decodedData)
						if !strings.Contains(r.Data, "ESDTNFTTransfer") {
							continue
						}
						transferData := strings.Split(r.Data, "@")
						nonce, err := strconv.ParseInt(transferData[2], 16, 64)
						if err != nil {
							logErr.Println(err.Error())
							continue
						}
						tokenIdByte, err := hex.DecodeString(transferData[1])
						if err != nil {
							logErr.Println(err.Error())
							continue
						}
						for i := 0; i < 1; i++ {
							nonceStr := transferData[2] //strconv.FormatInt(int64(nonce), 10)
							tokenId := string(tokenIdByte) + "-" + nonceStr
							_, err := storage.GetTokenByTokenIdAndNonceStr(tokenId, nonceStr)
							if err != nil {
								if err == gorm.ErrRecordNotFound {
									acc, err := storage.GetAccountByAddress(r.Receiver)
									if err != nil {
										if err == gorm.ErrRecordNotFound {
											acc = &entities.Account{}
											acc.Address = r.Receiver
											acc.Name = services.RandomName()
											err := storage.AddAccount(acc)
											if err != nil {
												logErr.Println("CRITICAL ", "fatal ", err.Error())
												continue
											}
										} else {
											logErr.Println(err.Error())
											continue
										}
									}
									price := colR.Value
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
									bigPriceStr := bigPrice.Div(bigPrice, big.NewInt(mintCount)).String()
									priceFraction := float64(1 / float64(mintCount))
									priceBigFloat := fprice.Mul(fprice, big.NewFloat(priceFraction))
									priceFloat, _ := priceBigFloat.Float64()

									// priceFloat, _ := fprice.Float64()
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
										imageURI = strings.Replace(imageURI, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
									}
									youbeiMeta := strings.Replace(metaURI, "https://gateway.pinata.cloud/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
									youbeiMeta = strings.Replace(youbeiMeta, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
									nonceStr = strconv.FormatInt(int64(nonce), 10)

									url := fmt.Sprintf("%s/%s.json", youbeiMeta, nonceStr)
									tokenDetail, err := services.GetResponse(fmt.Sprintf(`%s/nfts/%s`, ci.ElrondAPI, tokenId))
									if err != nil {
										logErr.Println(err.Error())
										continue
									}
									var tokenDetailObj entities.TokenBC
									err = json.Unmarshal(tokenDetail, &tokenDetailObj)
									if err != nil {
										logErr.Println(err.Error())
										continue
									}
									attrbs, err := services.GetResponse(url)
									if err != nil {
										logErr.Println(err.Error())
										if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "deadline") {
											err = storage.AddToken(&entities.Token{
												TokenID:      string(tokenIdByte),
												MintTxHash:   colR.TxHash,
												CollectionID: col.ID,
												Nonce:        uint64(nonce),
												NonceStr:     transferData[2],
												MetadataLink: string(youbeiMeta) + "/" + nonceStr + ".json",
												ImageLink:    string(imageURI) + "/" + nonceStr + ".png",
												TokenName:    tokenDetailObj.Name,
												Attributes:   []byte{},
												OwnerId:      acc.ID,
												OnSale:       false,
												PriceString:  bigPriceStr,
												PriceNominal: priceFloat,
											})
											if err != nil {
												logErr.Println(err.Error())
												continue
											}
											minted++
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
										MintTxHash:   colR.TxHash,
										CollectionID: col.ID,
										Nonce:        uint64(nonce),
										NonceStr:     transferData[2],
										MetadataLink: string(youbeiMeta) + "/" + nonceStr + ".json",
										ImageLink:    string(imageURI) + "/" + nonceStr + ".png",
										TokenName:    tokenDetailObj.Name,
										Attributes:   attributes,
										OwnerId:      acc.ID,
										OnSale:       false,
										PriceString:  bigPriceStr,
										PriceNominal: priceFloat,
									})
									if err != nil {
										logErr.Println(err.Error())
										continue
									}
									minted++
								} else {
									logErr.Println(err.Error())
									continue
								}
							}
						}
					}
				}
			}
			// collstats.RemoveCollectionToCheck(colObj) TODO
			if foundedTxsCount != minted {
				fmt.Println("CRITICAL", "mint!=founded", collectionIndexer.CollectionAddr, collectionIndexer.ID)
			}
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
