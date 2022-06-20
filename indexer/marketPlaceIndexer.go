package indexer

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/emurmotol/ethconv"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MarketPlaceIndexer struct {
	MarketPlaceAddr string `json:"marketPlaceAddr"`
	ElrondAPI       string `json:"elrondAPI"`
	ElrondAPISec    string `json:"elrondAPISec"`
	Logger          *log.Logger
	Delay           time.Duration // delay between each call
}

func NewMarketPlaceIndexer(marketPlaceAddr string, elrondAPI string, elrondAPISec string, delay uint64) (*MarketPlaceIndexer, error) {
	lerr := log.New(os.Stderr, "", log.LUTC|log.LstdFlags|log.Lshortfile)
	return &MarketPlaceIndexer{
		MarketPlaceAddr: marketPlaceAddr,
		ElrondAPI:       elrondAPI,
		ElrondAPISec:    elrondAPISec,
		Logger:          lerr,
		Delay:           time.Duration(delay)}, nil
}

func (mpi *MarketPlaceIndexer) StartWorker() {
	lerr := mpi.Logger

	api := mpi.ElrondAPI
	if mpi.ElrondAPISec != "" {
		api = mpi.ElrondAPISec
	}
	var lastProcessed entities.TransactionBC

	for {
	mainLoop:
		var foundTxs uint64 = 0
		time.Sleep(time.Second * mpi.Delay)
		marketStat, err := storage.GetMarketPlaceIndexer()
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				marketStat, err = storage.CreateMarketPlaceStat()
				if err != nil {
					lerr.Println("something went wrong creating marketstat", err.Error())
					continue
				}
			}
		}
		reqUrl := fmt.Sprintf("%s/accounts/%s/transactions?after=%d&size=50&order=desc&withScResults=true&withLogs=true",
			api,
			mpi.MarketPlaceAddr,
			marketStat.LastTimestamp)
		body, err := services.GetResponse(reqUrl)
		if err != nil {
			lerr.Println(err.Error())
			continue
		}
		var txs []entities.TransactionBC
		err = json.Unmarshal(body, &txs)
		if err != nil {
			lerr.Println("error unmarshal nfts marketplace", err.Error())
			continue
		}
		if len(txs) == 0 {
			goto mainLoop
		}
		if marketStat.LastTimestamp == txs[0].Timestamp {
			goto mainLoop
		}
		foundTxs += uint64(len(txs))

		for _, tx := range txs {
			if len(tx.Results) == 0 {
				continue
			}
			// txCounter := 0
		txloop:
			var orgTx entities.TransactionBC = tx
			if orgTx.Status == string(transaction.TxStatusPending) {
				lerr.Println("REPEAT", "no final state of tx")
				goto txloop
			}
			if (orgTx.Status == string(transaction.TxStatusSuccess) ||
				orgTx.Status == string(transaction.TxStatusFail) ||
				orgTx.Status == string(transaction.TxStatusInvalid)) &&
				!orgTx.PendingResults {
			} else {
				lerr.Println("REPEAT", "no final state of tx")
				goto txloop
			}

			results := tx.Results
			for _, result := range results {
				orgDataHex, err := base64.StdEncoding.DecodeString(result.Data)
				if err != nil {
					lerr.Println("BADERR", err.Error())
					continue
				}
				orgDataHexParts := strings.Split(string(orgDataHex), "@")
				orgDataHexStr := strings.Join(orgDataHexParts[1:], "")
				orgData, err := hex.DecodeString(orgDataHexStr)
				if err != nil {
					lerr.Println("BADERR", err.Error())
					continue
				}
				orgData = []byte(orgDataHexParts[0] + string(orgData))
				var actions map[string]bool = make(map[string]bool)
				if strings.Contains(orgDataHexParts[0], "upgradeContract") {
					continue
				}
				actions["isWithdrawn"] = strings.EqualFold(orgTx.Function, "withdrawNft")
				actions["isOnSale"] = strings.Contains(string(orgData), "putNftForSale")
				actions["isOnAuction"] = strings.Contains(string(orgData), "startAuction")
				actions["isBuyNft"] = strings.Contains(string(orgData), "buyNft")
				actions["isOffer"] = strings.Contains(string(orgData), "makeOffer")
				actions["isCancelOffer"] = strings.Contains(string(orgData), "cancelOffer")
				actions["isAcceptOffer"] = strings.Contains(string(orgData), "acceptOffer")
				actions["isBid"] = strings.Contains(string(orgData), "placeBid")
				actions["isEndAuction"] = strings.Contains(string(orgData), "endAuction")

				mpi.DeleteFailedTX(orgTx)

				next := false
				for _, v := range actions {
					if v {
						next = true
					}
				}
				if !next {
					lerr.Println("REPEAT", "no relevant tx found", orgData)
					continue
				}

				mainTxDataStr := orgTx.Data
				mainTxData, err := base64.StdEncoding.DecodeString(mainTxDataStr)
				if err != nil {
					lerr.Println("BADERR", err.Error())
					continue
				}
				mainDataParts := strings.Split(string(mainTxData), "@")
				hexTokenId := mainDataParts[1]
				tokenId, err := hex.DecodeString(hexTokenId)
				if err != nil {
					lerr.Println("BADERR", err.Error())
					continue
				}
				hexNonce := mainDataParts[2]
				data, err := base64.StdEncoding.DecodeString(result.Data)
				if err != nil {
					lerr.Println("BADERR", err.Error())
					continue
				}
				dataStr := string(data)

				dataParts := strings.Split(dataStr, "@")

				txTimestamp := orgTx.Timestamp

				senderAdress := orgTx.Sender
				sender, err := storage.GetAccountByAddress(senderAdress)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						sender = &entities.Account{}
						sender.Name = services.RandomName()
						sender.Address = senderAdress
						err := storage.AddAccount(sender)
						if err != nil {
							lerr.Println("MAINLOOP", "couldn't add user", err.Error())
							goto mainLoop
						}
					}
				}
				err = nil
				token, err := storage.GetTokenByTokenIdAndNonceStr(string(tokenId), hexNonce)
				if token.OwnerID == 1684 {
					zlog.Info("this")
				}
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						lerr.Println("REPEAT", err.Error())
						goto txloop
					} else {
						lerr.Println("no token found", string(tokenId), hexNonce)
						tokenDetail, err := services.GetResponse(fmt.Sprintf(`%s/nfts/%s`, api, string(tokenId)+"-"+hexNonce))
						if err != nil {
							if strings.Contains(err.Error(), "404") {
								lerr.Println("BADERR", err.Error())
								continue
							}
							lerr.Println("REPEAT", err.Error())
							goto txloop
						}
						var tokenDetailObj entities.TokenBC
						err = json.Unmarshal(tokenDetail, &tokenDetailObj)
						if err != nil {
							lerr.Println("REPEAT", err.Error())
							goto txloop
						}
						col, err := storage.GetCollectionByTokenId(tokenDetailObj.Collection)
						if err != nil {
							if err == gorm.ErrRecordNotFound {
								lerr.Println("collection not found for this token!!", tokenDetailObj.Collection)
								col, err = services.CreateCollectionFromToken(tokenDetailObj, api)
								if err != nil {
									lerr.Println("CONTINUE", err.Error(), "create collection failed on market indexer", tokenDetailObj.Collection)
									continue
								}
							}
						}
						idParts := strings.Split(tokenDetailObj.Identifier, "-")
						nonceStr := idParts[len(idParts)-1]
						imageURI, metadataLink := services.GetTokenUris(tokenDetailObj)

						if strings.Contains(api, "devnet") {
							imageURI = strings.Replace(imageURI, "https://gateway.pinata.cloud/ipfs/", "https://devnet-media.elrond.com/nfts/asset/", 1)
						} else {
							imageURI = strings.Replace(imageURI, "https://gateway.pinata.cloud/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
							imageURI = strings.Replace(imageURI, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
							imageURI = strings.Replace(imageURI, "ipfs://", "https://media.elrond.com/nfts/asset/", 1)
						}

						youbeiMeta := strings.Replace(metadataLink, "https://gateway.pinata.cloud/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
						youbeiMeta = strings.Replace(youbeiMeta, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
						youbeiMeta = strings.Replace(youbeiMeta, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
						youbeiMeta = strings.Replace(youbeiMeta, "ipfs://", "https://media.elrond.com/nfts/asset/", 1)

						attrbs, err := services.GetResponse(youbeiMeta)
						if err != nil {
							lerr.Println(err.Error(), string(youbeiMeta), tokenDetailObj.Attributes, tokenDetailObj.URIs)
						}
						metadataJSON := make(map[string]interface{})
						err = json.Unmarshal(attrbs, &metadataJSON)
						var attributes datatypes.JSON
						var rarity entities.TokenRarity

						if err != nil {
							lerr.Println(err.Error(), string(youbeiMeta), tokenDetailObj.Attributes, tokenDetailObj.URIs)
						} else {
							attributesBytes, err := json.Marshal(metadataJSON["attributes"])
							if err != nil {
								lerr.Println(err.Error())
							}
							err = json.Unmarshal(attributesBytes, &attributes)
							if err != nil {
								lerr.Println(err.Error())
							}

							// Get Rarity data if exist
							if _, ok := metadataJSON["rarity"]; ok {
								rarityBody, err := json.Marshal(metadataJSON["rarity"].(map[string]interface{}))
								if err != nil {
									lerr.Println(err.Error())
								} else {
									if err := json.Unmarshal(rarityBody, &rarity); err != nil {
										lerr.Println(err.Error())
									}
								}

							}
						}
						if !utf8.Valid([]byte(youbeiMeta)) {
							token.MetadataLink = ""
						}
						token = &entities.Token{
							TokenID:      tokenDetailObj.Collection,
							MintTxHash:   "",
							OwnerID:      sender.ID,
							CollectionID: col.ID,
							Nonce:        tokenDetailObj.Nonce,
							NonceStr:     nonceStr,
							MetadataLink: youbeiMeta,
							ImageLink:    imageURI,
							TokenName:    tokenDetailObj.Name,
							Attributes:   attributes,
							OnSale:       false,
						}

						if rarity.RarityScore != 0 && !math.IsInf(rarity.RarityScore, 0) {
							token.RarityScore = rarity.RarityScore
							token.RarityScoreNorm = rarity.RarityScoreNormed
							token.RarityUsedTraitCount = uint(rarity.UsedTraitsCount)
							token.IsRarityInserted = true
						}

						err = storage.AddOrUpdateToken(token)
						if err != nil {
							if err == gorm.ErrRecordNotFound {
								storage.UpdateToken(token)
								lerr.Println("BADERR", err.Error())
								continue
							} else {
								lerr.Println("REPEAT", err.Error())
								goto txloop
							}
						}
					}
				}
				failedTx := mpi.DeleteFailedTX(orgTx)
				if failedTx {
					_, err := storage.GetLastTokenTransaction(token.ID)
					if err != nil {
						if err == gorm.ErrRecordNotFound {
							err = storage.UpdateTokenWhere(token, map[string]interface{}{
								"OnSale": false, // TODO we can't be sure if tx is messed up
							}, "token_id=? AND nonce_str=?", tokenId, hexNonce)
							if err != nil {
								lerr.Println("failed to update token when tx failed")
							}
						} else {

						}
					}
				}

				price := orgTx.Value
				bigPrice, ok := big.NewInt(0).SetString(price, 10)
				if !ok {
					lerr.Println("CRITICAL", "conversion to bigInt failed for price")
					goto txloop
				}

				fprice, err := ethconv.FromWei(bigPrice, ethconv.Ether)
				if err != nil {
					lerr.Println("CRITICAL", err.Error())
					goto txloop
				}

				toUpdate := false // we need to update token afterward  this to detect if we are on right result inside tx NEEDS REFACTOR to better detect the case
				if actions["isOnSale"] && strings.Contains(string(data), "putNftForSale") && !failedTx {
					toUpdate = true
					price, ok := big.NewInt(0).SetString(dataParts[1], 16)
					if !ok {
						lerr.Println("CRITICAL", "can not convert price", price, dataParts[1])
						goto txloop
					}
					fprice, err := ethconv.FromWei(price, ethconv.Ether)
					if err != nil {
						lerr.Println("CRITICAL", err.Error())
						goto txloop
					}

					token.OnSale = true
					token.Status = entities.ListToken
					token.OwnerID = sender.ID
					token.LastBuyPriceNominal, _ = fprice.Float64()
					token.PriceNominal, _ = fprice.Float64()
					token.PriceString = price.String()
					err = storage.AddOrUpdateTransaction(&entities.Transaction{
						PriceNominal: token.PriceNominal,
						Type:         entities.ListToken,
						Timestamp:    orgTx.Timestamp,
						SellerID:     sender.ID,
						TokenID:      token.ID,
						CollectionID: token.CollectionID,
						Hash:         orgTx.TxHash,
					})
					if err != nil {
						lerr.Println(err.Error())
					}
				} else if actions["isOffer"] && !failedTx {
					toUpdate = false
					offerStr := mainDataParts[3]
					offer, _ := big.NewInt(0).SetString(offerStr, 16)
					offerFloat, err := ethconv.FromWei(offer, ethconv.Ether)
					offerNominal, _ := offerFloat.Float64()
					if err == nil {
						fmt.Printf("%f of type %T", offerNominal, offerNominal)
					}

					offerDeadline, err := strconv.ParseUint(mainDataParts[4], 16, 64)
					if err == nil {
						fmt.Printf("%d of type %T", offerDeadline, offerDeadline)
					}
					err = storage.DeleteOfferByOfferorForTokenId(senderAdress, token.ID)
					if err != nil {
						if err != gorm.ErrRecordNotFound {
							lerr.Println("REPEAT", err.Error())
							goto txloop
						}
					}
					err = storage.AddOffer(&entities.Offer{
						AmountNominal:  offerNominal,
						AmountString:   offer.String(),
						Expire:         offerDeadline,
						OfferorAddress: senderAdress,
						Timestamp:      tx.Timestamp,
						TxHash:         orgTx.TxHash,
						TokenID:        token.ID,
					})
					if err != nil {
						lerr.Println(err.Error())
					}
				} else if actions["isAcceptOffer"] && !failedTx {
					toUpdate = true
					offerorAddrHex := mainDataParts[3]
					token.OnSale = false
					token.Status = entities.BuyToken
					offerorAddrStr, err := services.ConvertHexToBehc32(offerorAddrHex)
					if err != nil {
						lerr.Println("MAINLOOP", err.Error())
						goto mainLoop
					}
					user, err := storage.GetAccountByAddress(offerorAddrStr)
					if err != nil {
						lerr.Println("MAINLOOP", err.Error())
						goto mainLoop
					}
					token.OwnerID = user.ID

					offerStr := mainDataParts[4]
					offer, _ := big.NewInt(0).SetString(offerStr, 16)
					offerFloat, err := ethconv.FromWei(offer, ethconv.Ether)
					if err != nil {
						lerr.Println("MAINLOOP", err.Error())
						goto mainLoop
					}

					err = storage.DeleteOffersForTokenId(token.ID)
					if err != nil {
						lerr.Println(err.Error())
						if err != gorm.ErrRecordNotFound {
							lerr.Println("REPEAT", err.Error())
							goto txloop
						}
					}

					err = storage.DeleteBidsForTokenId(token.ID)
					if err != nil {
						lerr.Println(err.Error())
						if err != gorm.ErrRecordNotFound {
							lerr.Println("REPEAT", err.Error())
							goto txloop
						}
					}

					offerNominal, _ := offerFloat.Float64()
					lastBuyPriceNominal := offerNominal
					token.LastBuyPriceNominal = lastBuyPriceNominal
					token.PriceString = offer.String()
					token.PriceNominal = lastBuyPriceNominal
					err = storage.AddOrUpdateTransaction(&entities.Transaction{
						PriceNominal: token.PriceNominal,
						Type:         entities.BuyToken,
						Timestamp:    orgTx.Timestamp,
						SellerID:     sender.ID,
						BuyerID:      user.ID,
						TokenID:      token.ID,
						CollectionID: token.CollectionID,
						Hash:         orgTx.TxHash,
					})
					if err != nil {
						lerr.Println("REPEAT", err.Error())
						goto txloop
					}
				} else if actions["isCancelOffer"] && !failedTx {
					toUpdate = false
					err := storage.DeleteOfferByOfferorForTokenId(senderAdress, token.ID)
					if err != nil {
						lerr.Println("REPEAT", err.Error())
						goto txloop
					}
				} else if actions["isOnAuction"] && strings.Contains(string(data), "startAuction") && !failedTx {
					toUpdate = true
					fmt.Println("is_on_auction", dataParts)
					hexMinBid := dataParts[1]
					minBid, _ := big.NewInt(0).SetString(hexMinBid, 16)
					minBidfloat, err := ethconv.FromWei(minBid, ethconv.Ether)
					lastBuyPriceNominal, _ := minBidfloat.Float64()
					if err == nil {
						fmt.Printf("%f of type %T", lastBuyPriceNominal, lastBuyPriceNominal)
					}

					auctionDeadline, err := strconv.ParseUint(dataParts[2], 16, 64)
					if err == nil {
						fmt.Printf("%d of type %T", auctionDeadline, auctionDeadline)
					}

					auctionStartTime, err := strconv.ParseUint(dataParts[3], 16, 64)
					if err == nil {
						fmt.Printf("%d of type %T", auctionStartTime, auctionStartTime)
					}

					token.OnSale = true
					token.Status = entities.AuctionToken
					token.OwnerID = sender.ID
					token.LastBuyPriceNominal = lastBuyPriceNominal
					token.PriceString = minBid.String()
					token.PriceNominal, _ = minBidfloat.Float64()
					token.AuctionDeadline = auctionDeadline
					token.AuctionStartTime = auctionStartTime
					err = storage.AddOrUpdateTransaction(&entities.Transaction{
						PriceNominal: lastBuyPriceNominal,
						Type:         entities.AuctionToken,
						Timestamp:    orgTx.Timestamp,
						SellerID:     sender.ID,
						TokenID:      token.ID,
						CollectionID: token.CollectionID,
						Hash:         orgTx.TxHash,
					})
					if err != nil {
						lerr.Println(err.Error())

					}
				} else if actions["isWithdrawn"] && !failedTx {
					toUpdate = true
					token.OnSale = false
					token.OwnerID = sender.ID
					token.Status = entities.WithdrawToken
					err = storage.AddOrUpdateTransaction(&entities.Transaction{
						PriceNominal: token.PriceNominal,
						Type:         entities.WithdrawToken,
						Timestamp:    orgTx.Timestamp,
						SellerID:     sender.ID,
						TokenID:      token.ID,
						CollectionID: token.CollectionID,
						Hash:         orgTx.TxHash,
					})
					if err != nil {
						lerr.Println(err.Error())
					}
				} else if actions["isBuyNft"] && strings.Contains(string(data), "Seller") && !failedTx {
					toUpdate = true
					token.OnSale = false
					token.Status = entities.BuyToken
					token.OwnerID = sender.ID
					user, err := services.GetOrCreateAccount(string(tx.Receiver))
					if err != nil {
						lerr.Println("MAINLOOP", err.Error())
						goto mainLoop
					}
					err = storage.DeleteOffersForTokenId(token.ID)
					if err != nil {
						lerr.Println(err.Error())
						if err != gorm.ErrRecordNotFound {
							lerr.Println("REPEAT", err.Error())
							goto txloop
						}
					}
					err = storage.DeleteBidsForTokenId(token.ID)
					if err != nil {
						lerr.Println(err.Error())
						if err != gorm.ErrRecordNotFound {
							lerr.Println("REPEAT", err.Error())
							goto txloop
						}
					}

					token.LastBuyPriceNominal, _ = fprice.Float64()
					token.PriceString = price
					token.PriceNominal, _ = fprice.Float64()
					err = storage.AddOrUpdateTransaction(&entities.Transaction{
						PriceNominal: token.PriceNominal,
						Type:         entities.BuyToken,
						Timestamp:    orgTx.Timestamp,
						BuyerID:      sender.ID,
						SellerID:     user.ID,
						TokenID:      token.ID,
						CollectionID: token.CollectionID,
						Hash:         orgTx.TxHash,
					})
					if err != nil {
						lerr.Println(err.Error())
					}
				} else if actions["isBid"] && !failedTx {
					toUpdate = true
					bidStr := mainDataParts[3]
					bid, _ := big.NewInt(0).SetString(bidStr, 16)
					bidFloat, err := ethconv.FromWei(bid, ethconv.Ether)
					bidNominal, _ := bidFloat.Float64()
					if err == nil {
						fmt.Printf("%f of type %T", bidNominal, bidNominal)
					}
					err = storage.AddBid(&entities.Bid{
						BidAmountNominal: bidNominal,
						BidAmountString:  bid.String(),
						BidderAddress:    senderAdress,
						Timestamp:        tx.Timestamp,
						TxHash:           orgTx.TxHash,
						TokenID:          token.ID,
					})
					if err != nil {
						lerr.Println("REPEAT", err.Error())
						goto txloop
					}
				} else if actions["isEndAuction"] && strings.Contains(string(data), "ESDTNFTTransfer") && !failedTx {
					toUpdate = true
					token.OnSale = false
					token.Status = entities.BuyToken

					user, err := services.GetOrCreateAccount(string(tx.Receiver))
					if err != nil {
						lerr.Println("MAINLOOP", err.Error())
						goto mainLoop
					}
					var typeOfTx entities.TxType = entities.BuyToken
					if token.OwnerID == sender.ID {
						// auction had no winner
						typeOfTx = entities.WithdrawToken
						token.Status = entities.WithdrawToken
					}

					err = storage.DeleteBidsForTokenId(token.ID)
					if err != nil {
						if err != gorm.ErrRecordNotFound {
							lerr.Println("MAINLOOP", err.Error())
							goto mainLoop
						}
						lerr.Println(err.Error())
					}

					token.OwnerID = user.ID
					err = storage.AddOrUpdateTransaction(&entities.Transaction{
						PriceNominal: token.PriceNominal,
						Type:         typeOfTx,
						Timestamp:    orgTx.Timestamp,
						SellerID:     sender.ID,
						BuyerID:      user.ID,
						TokenID:      token.ID,
						CollectionID: token.CollectionID,
						Hash:         orgTx.TxHash,
					})
					if err != nil {
						lerr.Println(err.Error())
					}
				}
				if token.LastMarketTimestamp <= txTimestamp && toUpdate && !failedTx {
					token.LastMarketTimestamp = txTimestamp
					err = storage.UpdateTokenWhere(token, map[string]interface{}{
						"OnSale":              token.OnSale,
						"Status":              token.Status,
						"NonceStr":            hexNonce,
						"PriceString":         token.PriceString,
						"PriceNominal":        token.PriceNominal,
						"LastMarketTimestamp": token.LastMarketTimestamp,
						"OwnerID":             token.OwnerID,
						"AuctionDeadline":     token.AuctionDeadline,
						"AuctionStartTime":    token.AuctionStartTime,
					}, "token_id=? AND nonce_str=?", tokenId, hexNonce)
					if err != nil {
						if err == gorm.ErrRecordNotFound {
							lerr.Println("MAINLOOP", err.Error())
							goto mainLoop
						}
						lerr.Println(err.Error(), "MAINLOOP", "error updating token ", fmt.Sprintf("tokenID %d", token.ID))
						goto mainLoop
					}
				}
				lastProcessed = tx
			}

		}

		storage.UpdateMarketPlaceIndexerTimestamp(lastProcessed.Timestamp)
	}
}
func (mpi *MarketPlaceIndexer) DeleteFailedTX(orgTx entities.TransactionBC) bool {

	if strings.EqualFold(orgTx.Status, "fail") || strings.EqualFold(orgTx.Status, "invalid") {
		tx, err := storage.GetTransactionWhere("hash=? AND timestamp=?", orgTx.TxHash, orgTx.Timestamp)
		if err != nil {
			if err == gorm.ErrRecordNotFound {

			} else {
				storage.DeleteTransaction(tx.ID)
			}
		}
		return true
	}
	return false
}
