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

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/services"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/emurmotol/ethconv"
	"gorm.io/datatypes"
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
	lastHash := ""
	lastHashTimestamp := uint64(0)
	lastIndex := 0
	for {
	mainLoop:
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
			lastIndex = 0
			marketStat.LastHash = lastHash
			marketStat, err = storage.UpdateMarketPlaceHash(lastHash)
			if err != nil {
				lerr.Println(err.Error())
				lerr.Println("error update marketplace index nfts ")
				continue
			}
			continue
		}
		if txResult[0].Hash == marketStat.LastHash {
			lastHashMet = true
			lastIndex = 0
			time.Sleep(time.Second * mpi.Delay)
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
			if tx.Hash == marketStat.LastHash {
				lastHashMet = true
				lastIndex = 0
				break
			} else {
				if lastHashTimestamp < tx.Timestamp {
					lastHashTimestamp = tx.Timestamp
					lastHash = tx.Hash
				}
			}
			if orgTx.Timestamp == 0 {
				continue
			}
			if strings.EqualFold(orgTx.Status, "fail") || strings.EqualFold(orgTx.Status, "invalid") {
				tx, err := storage.GetTransactionWhere("hash=? AND timestamp=?", orgTx.TxHash, orgTx.Timestamp)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						continue
					}
				}
				storage.DeleteTransaction(tx.ID)
				continue
			}
			if orgTx.TxHash == "" {
				continue
			}
			if orgTx.Status == "pending" {
				goto txloop
			}
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
			isOnAuction := strings.Contains(string(data), "startAuction")
			isBuyNft := strings.Contains(string(orgData), "buyNft")
			isOffer := strings.Contains(string(orgData), "makeOffer")
			isCancelOffer := strings.Contains(string(orgData), "cancelOffer")
			isAcceptOffer := strings.Contains(string(orgData), "acceptOffer")
			isBid := strings.Contains(string(orgData), "placeBid")
			isEndAuction := strings.Contains(string(orgData), "endAuction") && strings.Contains(string(data), "ESDTNFTTransfer")

			if !isOnSale &&
				!isOnAuction &&
				!isBuyNft &&
				!isWithdrawn &&
				!isOffer &&
				!isCancelOffer &&
				!isAcceptOffer &&
				!isBid &&
				!isEndAuction {
				continue
			}

			mainTxDataStr := orgTx.Data
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

			txTimestamp := orgTx.Timestamp

			senderAdress := orgTx.Sender
			sender, err := storage.GetAccountByAddress(senderAdress)
			if err != nil {
				if err == gorm.ErrNotImplemented {
					sender.Name = services.RandomName()
					err := storage.AddAccount(sender)
					if err != nil {
						lerr.Println("couldn't add user", err.Error())
						goto mainLoop
					}
				}
			}

			token, err := storage.GetTokenByTokenIdAndNonceStr(string(tokenId), hexNonce)
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					lerr.Println(err.Error())
					continue
				} else {
					lerr.Println("no token found", string(tokenId), hexNonce)
					tokenDetail, err := services.GetResponse(fmt.Sprintf(`%s/nfts/%s`, mpi.ElrondAPI, string(tokenId)+"-"+hexNonce))
					if err != nil {
						lerr.Println(err.Error())
						continue
					}
					var tokenDetailObj entities.TokenBC
					err = json.Unmarshal(tokenDetail, &tokenDetailObj)
					if err != nil {
						lerr.Println(err.Error())
						continue
					}
					col, err := storage.GetCollectionByTokenId(tokenDetailObj.Collection)
					if err != nil {
						if err == gorm.ErrRecordNotFound {
							lerr.Println("collection not found for this token!!", tokenDetailObj.Collection)
							col, err = services.CreateCollectionFromToken(tokenDetailObj, mpi.ElrondAPI)
							if err != nil {
								lerr.Println("create collection failed on market indexer", tokenDetailObj.Collection)
								continue
							}
						}
					}
					idParts := strings.Split(tokenDetailObj.Identifier, "-")
					nonceStr := idParts[len(idParts)-1]
					imageURI, metadataLink := services.GetTokenUris(tokenDetailObj)
					attrbs, err := services.GetResponse(metadataLink)
					metadataJSON := make(map[string]interface{})
					err = json.Unmarshal(attrbs, &metadataJSON)
					var attributes datatypes.JSON
					if err != nil {
						lerr.Println(err.Error(), string(metadataLink))
					} else {
						attributesBytes, err := json.Marshal(metadataJSON["attributes"])
						if err != nil {
							lerr.Println(err.Error())
						}
						err = json.Unmarshal(attributesBytes, &attributes)
						if err != nil {
							lerr.Println(err.Error())
							continue
						}
					}
					// owner, err := storage.GetAccountByAddress(senderAdress)
					// if err != nil {
					// 	if err == gorm.ErrRecordNotFound {
					// 		err := storage.AddAccount(&entities.Account{
					// 			Name:    tokenDetailObj.Identifier + "-Owner",
					// 			Address: tokenDetailObj.Owner,
					// 		})
					// 		if err != nil {
					// 			goto txloop
					// 		}
					// 	}
					// }
					token = &entities.Token{
						TokenID:      tokenDetailObj.Collection,
						MintTxHash:   "",
						OwnerId:      sender.ID,
						CollectionID: col.ID,
						Nonce:        tokenDetailObj.Nonce,
						NonceStr:     nonceStr,
						MetadataLink: metadataLink,
						ImageLink:    imageURI,
						TokenName:    tokenDetailObj.Name,
						Attributes:   attributes,
						OnSale:       false,
					}
					err = storage.AddToken(token)
					if err != nil {
						lerr.Println(err.Error())
						continue
					}
				}
			}

			price := orgTx.Value
			bigPrice, ok := big.NewInt(0).SetString(price, 10)
			if !ok {
				lerr.Println("conversion to bigInt failed for price")
				continue
			}
			fprice, err := ethconv.FromWei(bigPrice, ethconv.Ether)
			if err != nil {
				lerr.Println(err.Error())
				continue
			}

			if isOnSale {
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

				token.OnSale = true
				token.Status = entities.List
				token.LastBuyPriceNominal, _ = fprice.Float64()
				token.PriceNominal, _ = fprice.Float64()
				token.PriceString = price.String()
				err = storage.AddTransaction(&entities.Transaction{
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
			} else if isOffer {

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
						lerr.Println(err.Error())
						continue
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
			} else if isAcceptOffer {
				offerorAddrHex := mainDataParts[3]
				token.OnSale = false
				token.Status = entities.BuyToken
				offerorAddrStr, err := services.ConvertHexToBehc32(offerorAddrHex)
				if err != nil {
					lerr.Println(err.Error())
					goto mainLoop
				}
				user, err := storage.GetAccountByAddress(offerorAddrStr)
				if err != nil {
					lerr.Println(err.Error())
					goto mainLoop
				}
				token.OwnerId = user.ID

				offerStr := mainDataParts[4]
				offer, _ := big.NewInt(0).SetString(offerStr, 16)
				offerFloat, err := ethconv.FromWei(offer, ethconv.Ether)
				if err != nil {
					lerr.Println(err.Error())
					goto mainLoop
				}

				err = storage.DeleteOffersForTokenId(token.ID)
				if err != nil {
					lerr.Println(err.Error())
					continue
				}
				err = storage.DeleteBidsForTokenId(token.ID)
				if err != nil {
					lerr.Println(err.Error())
					continue
				}

				offerNominal, _ := offerFloat.Float64()
				lastBuyPriceNominal := offerNominal
				token.LastBuyPriceNominal = lastBuyPriceNominal
				token.PriceString = offer.String()
				token.PriceNominal = lastBuyPriceNominal
				err = storage.AddTransaction(&entities.Transaction{
					PriceNominal: token.PriceNominal,
					Type:         entities.BuyToken,
					Timestamp:    orgTx.Timestamp,
					SellerID:     sender.ID,
					TokenID:      token.ID,
					CollectionID: token.CollectionID,
					Hash:         orgTx.TxHash,
				})
				if err != nil {
					lerr.Println(err.Error())
				}

			} else if isCancelOffer {

				err := storage.DeleteOfferByOfferorForTokenId(senderAdress, token.ID)
				if err != nil {
					lerr.Println(err.Error())
					continue
				}
			} else if isOnAuction {
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
				token.LastBuyPriceNominal = lastBuyPriceNominal
				token.PriceString = minBid.String()
				token.PriceNominal, _ = minBidfloat.Float64()
				token.AuctionDeadline = auctionDeadline
				token.AuctionStartTime = auctionStartTime
				err = storage.AddTransaction(&entities.Transaction{
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
			} else if isWithdrawn {
				token.OnSale = false
				token.Status = entities.WithdrawToken
				err = storage.AddTransaction(&entities.Transaction{
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
			} else if isBuyNft {
				token.OnSale = false
				token.Status = entities.BuyToken
				token.OwnerId = sender.ID
				lastBuyPriceNominal, err := strconv.ParseFloat(dataParts[1], 64)
				if err == nil {
					fmt.Printf("%f of type %T", lastBuyPriceNominal, lastBuyPriceNominal)
				}

				err = storage.DeleteOffersForTokenId(token.ID)
				if err != nil {
					lerr.Println(err.Error())
					continue
				}
				err = storage.DeleteBidsForTokenId(token.ID)
				if err != nil {
					lerr.Println(err.Error())
					continue
				}

				token.LastBuyPriceNominal, _ = fprice.Float64()
				token.PriceString = price
				token.PriceNominal, _ = fprice.Float64()
				err = storage.AddTransaction(&entities.Transaction{
					PriceNominal: token.PriceNominal,
					Type:         entities.BuyToken,
					Timestamp:    orgTx.Timestamp,
					SellerID:     sender.ID,
					TokenID:      token.ID,
					CollectionID: token.CollectionID,
					Hash:         orgTx.TxHash,
				})
				if err != nil {
					lerr.Println(err.Error())
				}

			} else if isBid {
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
					lerr.Println(err.Error())
					continue
				}
			} else if isEndAuction {
				token.OnSale = false
				token.Status = entities.BuyToken

				user, err := storage.GetAccountByAddress(string(tx.Receiver))
				if err != nil {
					lerr.Println("CRITICAL", err.Error())
					goto mainLoop
				}
				var typeOfTx entities.TxType = entities.BuyToken
				if token.OwnerId == sender.ID {
					// auction had no winner
					typeOfTx = entities.WithdrawToken
					token.Status = entities.WithdrawToken
				}

				err = storage.DeleteBidsForTokenId(token.ID)
				if err != nil {
					lerr.Println(err.Error())
					continue
				}

				token.OwnerId = user.ID
				err = storage.AddTransaction(&entities.Transaction{
					PriceNominal: token.PriceNominal,
					Type:         typeOfTx,
					Timestamp:    orgTx.Timestamp,
					SellerID:     sender.ID,
					TokenID:      token.ID,
					CollectionID: token.CollectionID,
					Hash:         orgTx.TxHash,
				})
				if err != nil {
					lerr.Println(err.Error())
				}
			}
			if token.LastMarketTimestamp < txTimestamp {
				token.LastMarketTimestamp = txTimestamp
				err = storage.UpdateTokenWhere(token, map[string]interface{}{
					"OnSale":              token.OnSale,
					"Status":              token.Status,
					"PriceString":         token.PriceString,
					"PriceNominal":        token.PriceNominal,
					"LastMarketTimestamp": txTimestamp,
					"OwnerId":             token.OwnerId,
					"AuctionDeadline":     token.AuctionDeadline,
					"AuctionStartTime":    token.AuctionStartTime,
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
		if !lastHashMet {
			lastIndex += len(txResult)
		} else {
			marketStat, err = storage.UpdateMarketPlaceHash(lastHash)
			if err != nil {
				lerr.Println(err.Error())
				lerr.Println("error update marketplace index nfts ")
				continue
			}
		}
	}
}
