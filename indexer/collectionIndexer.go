package indexer

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/btcsuite/btcutil/bech32"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var getCollectionNFTSAPI string = "%s/collections/%s/nfts?from=%d&withOwner=true"
var getCollectionNFTSCountsAPI string = "%s/collections/%s/nfts/count?withOwner=true"

type CollectionIndexer struct {
	DeployerAddr string `json:"deployerAddr"`
	ElrondAPI    string `json:"elrondApi"`
	ElrondAPISec string `json:"elrondApiSec"`
	Logger       *log.Logger
	Delay        time.Duration // delay per request in second
}

func NewCollectionIndexer(deployerAddr string, elrondAPI string, elrondAPISec string, delay uint64) (*CollectionIndexer, error) {
	l := log.New(os.Stderr, "", log.LUTC|log.LstdFlags|log.Lshortfile)
	return &CollectionIndexer{
		DeployerAddr: deployerAddr,
		ElrondAPI:    elrondAPI,
		ElrondAPISec: elrondAPISec,
		Delay:        time.Duration(delay),
		Logger:       l}, nil
}
func (ci *CollectionIndexer) CorrectIfAddressIsEmpty(colObj *entities.Collection, blockchainApi string) error {
	if colObj.ContractAddress == "" {
		colDetail, err := services.GetCollectionDetailBC(colObj.TokenID, blockchainApi)
		if err != nil {
			return err
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
		err = services.UpdateCollectionWithAddress(colObj, map[string]interface{}{
			"Name":            colObj.Name,
			"ContractAddress": colObj.ContractAddress,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func (ci *CollectionIndexer) StartWorker() {
	logErr := ci.Logger
	var colsToCheck []dtos.CollectionToCheck
	api := ci.ElrondAPI
	if ci.ElrondAPISec != "" {
		api = ci.ElrondAPISec
		getCollectionNFTSAPI = "%s/nftsFromCollection?collection=%s&from=%d&withOwner=true"
		getCollectionNFTSCountsAPI = "%s/nfts/count?collection=%s&withOwner=true"
	}
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
			api,
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
						_, err = storage.CreateCollectionStat(entities.CollectionIndexer{CollectionName: string(tokenIdStr), CollectionAddr: bech32Addr})
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
		cols, err := storage.GetVerifiedCollections()
		if err != nil {
			logErr.Println(err.Error())
			continue
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(cols), func(i, j int) { cols[i], cols[j] = cols[j], cols[i] })
		zlog.Info("going into cols loop")
		for _, colObj := range cols {
			fmt.Println("iterate_col", colObj)
			if err := ci.CorrectIfAddressIsEmpty(&colObj, api); err != nil {
				if err != nil {
					zlog.Error(err.Error())
					continue
				}
			}
			collectionIndexer, err := storage.GetCollectionIndexer(colObj.ContractAddress)
			if err != nil {
				if err == gorm.ErrRecordNotFound { //indexer not found
					collectionIndexer, err = storage.CreateCollectionStat(entities.CollectionIndexer{
						CollectionAddr: colObj.ContractAddress,
						CollectionName: colObj.TokenID,
					})
					if err != nil { // bad error
						zlog.Error("error create colleciton indexer", zap.Error(err))
						continue
					}
				} else { // unknown error
					zlog.Error("error getting colleciton indexer", zap.Error(err))
					continue
				}
			}
			if collectionIndexer.CollectionName == "" { //update collection name inside collection indexer
				err := storage.UpdateCollectionIndexerWhere(&collectionIndexer, map[string]interface{}{"collection_name": colObj.TokenID}, "id=?", collectionIndexer.ID)
				if err != nil {
					zlog.Error("error UpdateCollectionndexerWhere collection indexer", zap.Error(err))
					continue
				}
			}
			countNftRes, _ := services.GetResponse(fmt.Sprintf(getCollectionNFTSCountsAPI, api, collectionIndexer.CollectionName))
			var count uint64
			json.Unmarshal(countNftRes, &count)
			lastIndex := 0
			done := false
			if count <= collectionIndexer.CountIndexed {
				continue
			}
			for !done {
				if lastIndex > 9999 {
					done = true
				}
				// Get NFTS from collection from lastIndex , index can't be higher than 10k as elastic query by default won't support that and api returns error
				url := fmt.Sprintf(getCollectionNFTSAPI+"&size=100",
					api,
					collectionIndexer.CollectionName,
					lastIndex)
				res, err := services.GetResponse(url)
				if err != nil {
					if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "deadline") || strings.Contains(err.Error(), "404") {
						time.Sleep(time.Second * 10)
						continue
					}
					logErr.Println("BADERR", err.Error())
				}

				var tokens []entities.TokenBC
				err = json.Unmarshal(res, &tokens)
				if err != nil {
					logErr.Println("BADERR", err.Error())
					logErr.Println(err.Error(), "collection name ",
						collectionIndexer.CollectionName,
						"lastIndex", lastIndex,
						"url", url,
						"raw data", res)
				}
				if len(tokens) == 0 {
					fmt.Println("DONE")
					done = true
					continue
				}
				tokenCountSuc := 0
				for _, token := range tokens {
				tokenLoop:
					imageURI, attributeURI := services.GetTokenUris(token)
					nonce10Str := strconv.FormatUint(token.Nonce, 10)
					nonceStr := strconv.FormatUint(token.Nonce, 16)
					if len(nonceStr)%2 != 0 {
						nonceStr = "0" + nonceStr
					}
					// Convert URI to elrond url for faster retreive
					if strings.Contains(api, "devnet") {
						imageURI = strings.Replace(imageURI, "https://gateway.pinata.cloud/ipfs/", "https://devnet-media.elrond.com/nfts/asset/", 1)
					} else {
						imageURI = strings.Replace(imageURI, "https://gateway.pinata.cloud/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
						imageURI = strings.Replace(imageURI, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
						imageURI = strings.Replace(imageURI, "ipfs://", "https://media.elrond.com/nfts/asset/", 1)
					}
					if imageURI != "" {
						if string(imageURI[len(imageURI)-1]) == "/" {
							imageURI = imageURI[:len(imageURI)-1]
						}
					}

					youbeiMeta := strings.Replace(attributeURI, "https://gateway.pinata.cloud/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
					youbeiMeta = strings.Replace(youbeiMeta, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
					youbeiMeta = strings.Replace(youbeiMeta, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
					youbeiMeta = strings.Replace(youbeiMeta, "ipfs://", "https://media.elrond.com/nfts/asset/", 1)
					if youbeiMeta != "" {
						if string(youbeiMeta[len(youbeiMeta)-1]) == "/" {
							youbeiMeta = youbeiMeta[:len(youbeiMeta)-1]
						}
					}

					url := fmt.Sprintf("%s/%s.json", youbeiMeta, nonce10Str)
					attrbs, err := services.GetResponse(url)

					metadataJSON := make(map[string]interface{})
					err = json.Unmarshal(attrbs, &metadataJSON)
					if err != nil {
						logErr.Println(err.Error(), string(url), token.Collection, token.Attributes, token.Identifier, token.Media, token.Metadata)
					}
					var attributes datatypes.JSON
					attributesBytes, err := json.Marshal(metadataJSON["attributes"])
					if err != nil {
						logErr.Println(err.Error())
						attributesBytes = []byte{}
					}
					err = json.Unmarshal(attributesBytes, &attributes)
					if err != nil {
						logErr.Println(err.Error())
					}

					//get owner of token from database TODO
					if token.Owner == "" {
						tokenRes, err := services.GetResponse(fmt.Sprintf("%s/nfts/%s", api, token.Identifier))
						if err != nil {
							logErr.Println("CRITICAL can't get nft data", err.Error())
							if strings.Contains(err.Error(), "deadline") {
								goto tokenLoop
							}
							continue
						}
						json.Unmarshal(tokenRes, &token)
					}
					if token.Owner == "" {
						token.Owner = token.Creator
					}
					acc, err := storage.GetAccountByAddress(token.Owner)
					if err != nil {
						if err != gorm.ErrRecordNotFound {
							continue
						} else {
							name := services.RandomName()
							acc = &entities.Account{
								Address: token.Owner,
								Name:    name,
							}
							err := storage.AddAccount(acc)
							if err != nil {
								if !strings.Contains(err.Error(), "duplicate") {
									logErr.Println("CRITICAL can't create user", err.Error())
									continue
								} else {
									acc, err = storage.GetAccountByAddress(token.Owner)
									if err != nil {
										logErr.Println("CRITICAL can't get user", err.Error())
										continue
									}
								}
							}
						}
					}
					//try get token from database TODO
					dbToken, err := storage.GetTokenByTokenIdAndNonce(token.Collection, token.Nonce)
					if err != nil {
						if err != gorm.ErrRecordNotFound {
							logErr.Println(err.Error(), "getTokenByTokenIdAndNonce_error")
							continue
						} else {

						}
					}
					if dbToken == nil {
						dbToken = &entities.Token{}
					}
					err = storage.AddOrUpdateToken(&entities.Token{
						TokenID:      token.Collection,
						MintTxHash:   dbToken.MintTxHash,
						CollectionID: colObj.ID,
						Nonce:        token.Nonce,
						NonceStr:     nonceStr,
						MetadataLink: string(youbeiMeta) + "/" + nonce10Str + ".json",
						ImageLink:    string(imageURI),
						TokenName:    token.Name,
						Attributes:   attributes,
						OwnerID:      acc.ID,
						OnSale:       false,
						PriceString:  dbToken.PriceString,
						PriceNominal: dbToken.PriceNominal,
					})
					if err != nil {
						logErr.Println("BADERR", err.Error())
						if err == gorm.ErrRegistered {
							tokenCountSuc++
						}
					} else {
						tokenCountSuc++
					}
				}
				lastIndex += tokenCountSuc
				countIndexed := collectionIndexer.CountIndexed + uint64(tokenCountSuc)
				if countIndexed > count {
					countIndexed = count
				}
				err = storage.UpdateCollectionIndexerWhere(&collectionIndexer,
					map[string]interface{}{
						"LastIndex":    lastIndex,
						"CountIndexed": countIndexed,
					},
					"id=?",
					collectionIndexer.ID)
				if err != nil {
					logErr.Println("CRITICAL", err.Error())
				}
			}
		}

	}

}
