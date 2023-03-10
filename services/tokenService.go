package services

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/data/dtos"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/interaction"
	"github.com/ENFT-DAO/youbei-api/stats/collstats"
	"github.com/ENFT-DAO/youbei-api/storage"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/boltdb/bolt"
)

var (
	TokenSearchCacheKeyFormat = "TokenSearch:%s"
	TokenSearchExpirePeriod   = 5 * time.Minute
)

type ListTokenRequest struct {
	TxHash        string  `json:"txHash"`
	UserAddress   string  `json:"walletAddress"`
	TokenID       string  `json:"tokenName"`
	Nonce         string  `json:"tokenNonce"`
	Status        string  `json:"saleStatus"`
	OnSale        bool    `json:"saleOnSale"`
	StringPrice   string  `json:"saleStringPrice"`
	NominalPrice  float64 `json:"saleNominalPrice"`
	SaleStartDate uint64  `json:"saleStartDate"`
	SaleEndDate   uint64  `json:"saleEndDate"`
	TxConfirmed   bool    `json:"txConfirmed"`
}

type NonFungibleToken struct {
	Identifier string         `json:"identifier"`
	Collection string         `json:"collection"`
	Name       string         `json:"name"`
	Attributes string         `json:"attributes"`
	Nonce      uint64         `json:"nonce"`
	Creator    string         `json:"creator"`
	Owner      string         `json:"owner"`
	Url        string         `json:"url"`
	Hash       string         `json:"hash"`
	Royalties  float64        `json:"royalties"`
	Uris       []string       `json:"uris"`
	Metadata   datatypes.JSON `json:"metadata"`
	Ticker     string         `json:"ticker"`
}

type AvailableTokensRequest struct {
	Tokens []string `json:"tokens"`
}

type AvailableToken struct {
	Collection struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Available bool   `json:"available"`
	} `json:"collection"`
	Token struct {
		Id        string `json:"id"`
		Nonce     uint64 `json:"nonce"`
		Name      string `json:"name"`
		Available bool   `json:"available"`
	} `json:"token"`
}

type ProxyTokenResponse struct {
	Data struct {
		Token []string `json:"token"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type NftProxyReponseToken struct {
	Balance   string   `json:"balance"`
	Name      string   `json:"name"`
	Hash      string   `json:"hash"`
	Royalties string   `json:"royalties"`
	Uris      []string `json:"uris"`
}

type NftProxyResponse struct {
	Data struct {
		TokenData NftProxyReponseToken `json:"tokenData"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type MetadataRelayRequest struct {
	Url string `json:"url"`
}

type AvailableTokensResponse struct {
	Tokens map[string]AvailableToken `json:"tokens"`
}

type TokenCacheInfo struct {
	TokenDbId uint64
	TokenName string
}

type WhitelistBuyLimitCountRequest struct {
	ContractAddress string `json:"contractAddress"`
	UserAddress     string `json:"userAddress"`
}

type ChangeTokenOwnerRequest struct {
	TokenId     string `json:"tokenId"`
	NonceHexStr string `json:"nonceHexStr"`
	NewOwner    string `json:"newOwner"`
}

const (
	minPriceUnit               = 1000
	minPercentUnit             = 1000
	minPercentRoyaltiesUnit    = 100
	minPriceDecimals           = 15
	maxPercentRoyaltiesAllowed = 1000

	maxTokenLinkResponseSize = 2048
	maxTokenNumAvailableSize = 25

	ZeroAddress           = "erd1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6gq4hu"
	NftProxyRequestFormat = "%s/address/%s/nft/%s/nonce/%d"

	UrlResponseCacheKeyFormat = "Url:%s"
	UrlResponseExpirePeriod   = 5 * time.Minute

	RefreshMetadataSetNxKeyFormat    = "Refresh:%s-%d"
	RefreshMetadataSetNxExpirePeriod = 15 * time.Minute

	ipfsProtocolURLPrefix = "ipfs://"
	ipfsDefaultGatewayURL = "https://ipfs.io/ipfs/%s"
)

var (
	TokenIdToDbIdCacheInfo = []byte("tokenToId")

	baseExp = big.NewInt(10)

	log                = logger.GetOrCreate("services")
	tooManyTokensError = errors.New("too many tokens")
)

func GetTokenBaseURIs(tokenData entities.TokenBC) (string, string) {
	var attributeUrl string
	for _, uri := range tokenData.URIs {
		attributeUrlByte, err := base64.StdEncoding.DecodeString(uri)
		if !strings.Contains(string(attributeUrlByte), ".json") {
			continue
		}
		if err != nil {
			attributeUrl = ""
		} else {
			attributeUrl = string(attributeUrlByte)

			urlParts := strings.Split(attributeUrl, "/")
			lastPart := urlParts[len(urlParts)-1]
			attributeUrl = strings.Replace(attributeUrl, lastPart, "", 1)
			if strings.Contains(attributeUrl, "metadata:") {
				attributeUrl = strings.ReplaceAll(attributeUrl, "metadata:", "")
				attributeUrl = fmt.Sprintf(ipfsDefaultGatewayURL, attributeUrl)
			}

		}
	}
	var imageURL string
	if strings.Contains(tokenData.URL, ".png") ||
		strings.Contains(tokenData.URL, ".PNG") ||
		strings.Contains(tokenData.URL, ".JPG") ||
		strings.Contains(tokenData.URL, ".JPEG") ||
		strings.Contains(tokenData.URL, ".jpeg") ||
		strings.Contains(tokenData.URL, ".mp4") ||
		strings.Contains(tokenData.URL, ".MP4") ||
		strings.Contains(tokenData.URL, ".jpg") {
		imageURL = tokenData.URL
	} else {
		urlParts := strings.Split(tokenData.URL, "/")
		lastPart := urlParts[len(urlParts)-1]
		imageURL = strings.Replace(tokenData.URL, lastPart, "", 1)

		nonce10Str := strconv.FormatUint(tokenData.Nonce, 10)

		imageURL = imageURL + nonce10Str
	}
	return imageURL, attributeUrl
}
func GetTokenUris(tokenData entities.TokenBC) (string, string) {
	var attributeUrl string
	for _, uri := range tokenData.URIs {
		if attributeUrl != "" {
			continue // pick first found json! or non png
		}
		attributeUrlByte, err := base64.StdEncoding.DecodeString(uri)
		if err != nil {
			continue
		}
		medias := tokenData.Media.([]interface{})
		if string(attributeUrlByte) != (medias[0].(map[string]interface{})["originalUrl"].(string)) {
			attributeUrl = string(attributeUrlByte)
		}
		// else {
		// 	attributeUrl = string(attributeUrlByte)
		// 	urlParts := strings.Split(attributeUrl, "/")
		// 	lastPart := urlParts[len(urlParts)-1]
		// 	attributeUrl = strings.Replace(attributeUrl, lastPart, "", 1)
		// 	stringNonce := fmt.Sprintf("%02d", tokenData.Nonce)
		// 	attributeUrl = attributeUrl + stringNonce + ".json"
		// 	break
		// }
	}
	if attributeUrl == "" {
		if tokenData.Attributes != "" {
			attributeUrlByte, err := base64.StdEncoding.DecodeString(tokenData.Attributes)
			if err != nil {

			} else {
				attributeUrl = string(attributeUrlByte)
				urlParts := strings.Split(attributeUrl, "/")
				lastPart := urlParts[len(urlParts)-1]
				attributeUrl = strings.Replace(attributeUrl, lastPart, "", 1)
				if strings.Contains(attributeUrl, "metadata:") {
					attributeUrl = strings.ReplaceAll(attributeUrl, "metadata:", "")
					attributeUrl = fmt.Sprintf(ipfsDefaultGatewayURL, attributeUrl)
					stringNonce := fmt.Sprintf("%02d", tokenData.Nonce)

					attributeUrl = attributeUrl + stringNonce + ".json"
				}
			}
		}
	}

	return tokenData.URL, attributeUrl
}
func ListTokenFromClient(request *ListTokenRequest, blockchainApi string) error {

	//get token data from blockchain
	tokenData, err := getTokenByNonce(request.TokenID, request.Nonce, blockchainApi)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	//auto create collection
	collection, err := storage.GetCollectionByTokenId(request.TokenID)
	if err != nil {
		var attributeUrl string
		for _, uri := range tokenData.Uris {
			attributeUrlByte, err := base64.StdEncoding.DecodeString(uri)
			if !strings.Contains(string(attributeUrlByte), ".json") {
				if tokenData.Attributes != "" {
					attributeUrlByte, err = base64.StdEncoding.DecodeString(tokenData.Attributes)
					if err != nil {
						continue
					}
				}
			}
			if err != nil {
				attributeUrl = ""
			} else {
				attributeUrl = string(attributeUrlByte)

				urlParts := strings.Split(attributeUrl, "/")
				lastPart := urlParts[len(urlParts)-1]
				attributeUrl = strings.Replace(attributeUrl, lastPart, "", 1)
				if strings.Contains(attributeUrl, "metadata:") {
					attributeUrl = strings.ReplaceAll(attributeUrl, "metadata:", "")
					attributeUrl = fmt.Sprintf(ipfsDefaultGatewayURL, attributeUrl)
				}

			}
		}

		urlParts := strings.Split(tokenData.Url, "/")
		lastPart := urlParts[len(urlParts)-1]
		url := strings.Replace(tokenData.Url, lastPart, "", 1)

		//if no collection auto create it
		var autoCreateCollectionRequest AutoCreateCollectionRequest
		autoCreateCollectionRequest.Name = request.TokenID
		autoCreateCollectionRequest.TokenId = request.TokenID
		autoCreateCollectionRequest.Nonce = request.Nonce
		autoCreateCollectionRequest.UserAddress = request.UserAddress
		autoCreateCollectionRequest.CreatorAddress = tokenData.Creator
		autoCreateCollectionRequest.TokenBaseURI = url
		autoCreateCollectionRequest.MetadataBaseURI = attributeUrl
		autoCreateCollectionRequest.CreatedAt = uint64(time.Now().UnixMilli())
		collection, err = AutoCreateCollection(&autoCreateCollectionRequest, blockchainApi)

		if err != nil {
			return errors.New("no collection found after autocreation")
		}
	}

	//the account was created if it did not exists in the previous step, otherwise get it
	account, err := storage.GetAccountByAddress(request.UserAddress)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	stringNonce := fmt.Sprintf("%02d", tokenData.Nonce)

	imageURI := ""
	metaDataURI := ""

	if len(tokenData.Uris) > 0 {

		byteImageURI := []byte{}
		byteMetaDataURI := []byte{}

		byteImageURI, err = base64.StdEncoding.DecodeString(tokenData.Uris[0])
		if err != nil {
			fmt.Printf("%v\n", err)
			return err
		}

		if strings.Contains(blockchainApi, "devnet") {
			imageURI = strings.Replace(string(byteImageURI), "https://gateway.pinata.cloud/ipfs/", "https://devnet-media.elrond.com/nfts/asset/", 1)
			imageURI = strings.Replace(imageURI, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
		} else {
			imageURI = strings.Replace(string(byteImageURI), "https://gateway.pinata.cloud/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
			imageURI = strings.Replace(imageURI, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
		}

		if len(tokenData.Uris) > 1 {

			byteMetaDataURI, err = base64.StdEncoding.DecodeString(tokenData.Uris[1])
			if err != nil {
				fmt.Printf("%v\n", err)
				return err
			}

			metaDataURI = strings.Replace(string(byteMetaDataURI), "https://gateway.pinata.cloud/ipfs/", "https://media.elrond.com/nfts/asset/", 1)
			metaDataURI = strings.Replace(metaDataURI, "https://ipfs.io/ipfs/", "https://media.elrond.com/nfts/asset/", 1)

		}
	}

	token := &entities.Token{
		MintTxHash:       request.TxHash,
		Nonce:            tokenData.Nonce,
		NonceStr:         stringNonce,
		OwnerID:          account.ID,
		CollectionID:     collection.ID,
		TokenID:          tokenData.Collection,
		RoyaltiesPercent: tokenData.Royalties,
		ImageLink:        imageURI,
		MetadataLink:     metaDataURI,
		CreatedAt:        uint64(time.Now().Unix()),
		Attributes:       GetAttributesFromMetadata(string(metaDataURI)),
		TokenName:        tokenData.Name,
		Hash:             tokenData.Hash,
		Status:           entities.TxType(request.Status),
		PriceString:      request.StringPrice,
		PriceNominal:     request.NominalPrice,
		AuctionStartTime: request.SaleStartDate,
		AuctionDeadline:  request.SaleEndDate,
		OnSale:           request.OnSale,
		TxConfirmed:      request.TxConfirmed,
	}

	var innerErr error

	innerErr = storage.AddToken(token)

	if innerErr == nil {
		_, cacheErr := AddTokenToCache(token.TokenID, token.Nonce, token.TokenName, token.ID)
		if cacheErr != nil {
			log.Error("could not add token to cache")
		}
	} else {
		innerErr = storage.UpdateTokenWhere(token, map[string]interface{}{
			"OnSale": token.OnSale,
			"Status": token.Status,
		}, "token_id=? AND nonce_str=?", token.TokenID, token.NonceStr)

	}

	if innerErr != nil {
		log.Debug("could not create or update token", "err", innerErr)
		return innerErr
	}

	return nil

}

func WithdrawToken(args WithdrawTokenArgs) {

	// priceNominal, err := GetPriceNominal(args.Price)
	// if err != nil {
	// 	log.Debug("could not parse price", "err", err)
	// 	return
	// }

	// ownerAccount, err := storage.GetAccountByAddress(args.OwnerAddress)
	// if err != nil {
	// 	log.Debug("could not get owner account", err)
	// 	return
	// }

	token, err := storage.GetTokenByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get token", "err", err)
		return
	}

	if args.NonceStr != "" {
		token.NonceStr = args.NonceStr
	}

	if args.TxHash != "" {
		token.MintTxHash = args.TxHash
	}

	if args.TxConfirmed {
		token.TxConfirmed = args.TxConfirmed
	}

	token.Status = entities.WithdrawToken
	token.OnSale = false

	err = storage.UpdateToken(token)
	if err != nil {
		log.Debug("could not update token", "err", err)
		return
	}

	err = storage.DeleteOffersForTokenId(token.ID)
	if err != nil {
		log.Debug("could not delete offers for token", "err", err)
	}

	err = storage.DeleteBidsForTokenId(token.ID)
	if err != nil {
		log.Debug("could not delete bids for token", "err", err)
	}
	// Indexer is doing it better TODO

	// transaction := entities.Transaction{
	// 	Hash:         args.TxHash,
	// 	Type:         entities.WithdrawToken,
	// 	PriceNominal: priceNominal,
	// 	Timestamp:    args.Timestamp,
	// 	SellerID:     0,
	// 	BuyerID:      ownerAccount.ID,
	// 	TokenID:      token.ID,
	// 	CollectionID: token.CollectionID,
	// }

	// AddTransaction(&transaction)
}

func ListToken(args ListTokenArgs, blockchainProxy string, marketplaceAddress string) {

	var priceNominal float64
	var err error

	if args.PriceNominal != "" {
		priceNominal, err = strconv.ParseFloat(args.PriceNominal, 64)
		if err != nil {
			log.Debug("could not parse nominal", "err", err)
			return
		}
	} else {
		priceNominal, err = GetPriceNominal(args.Price)
		if err != nil {
			log.Debug("could not parse price", "err", err)
			return
		}
	}

	ownerAccount, err := GetOrCreateAccount(args.OwnerAddress)
	if err != nil {
		log.Debug("could not get or create account", "err", err)
		return
	}

	collectionId := uint64(0)
	collection, err := storage.GetCollectionByTokenId(args.TokenId)
	if err == nil {
		collectionId = collection.ID
	}

	var metadataLink = ""
	if args.SecondLink != "" {
		metadataLink = args.SecondLink
		if len(metadataLink) == 0 {
			var innerErr error
			metadataLink, innerErr = TryGetMetadataLink(blockchainProxy, marketplaceAddress, args.TokenId, args.Nonce)
			if innerErr != nil {
				log.Debug("could not get metadata link", innerErr)
			}
		}
	}

	var updateToken = false
	token, err := storage.GetTokenByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		token = &entities.Token{}
	}
	priceBigFloat, err := ConvertBigFloatToFloat(args.PriceNominal)
	if err != nil {
		log.Debug("could not convert big float string to big float object", "err", err)
		return
	}
	fmt.Println(priceBigFloat.String())
	priceBigFloat, ok := TurnBigFloatoBigFloatNDec(priceBigFloat, 18)
	fmt.Println(priceBigFloat.String())
	if !ok {
		log.Debug("could not TurnBigFloatoBigFloatNDec priceBigFloat", "err", err)
		return
	}

	hexNonce := strconv.FormatInt(int64(token.Nonce), 16)
	if len(hexNonce)%2 != 0 {
		hexNonce = "0" + hexNonce
	}
	args.NonceStr = hexNonce
	args.Nonce = token.Nonce

	tokenDetail, err := GetResponse(fmt.Sprintf(`%s/nfts/%s`, blockchainProxy, string(args.TokenId)+"-"+hexNonce))
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			zlog.Error("BADERR", zap.Error(err))
		}
	}
	var tokenDetailObj entities.TokenBC
	err = json.Unmarshal(tokenDetail, &tokenDetailObj)
	if err != nil {
		zlog.Error("BADERR", zap.Error(err))

	}

	finalPriceBigInt := new(big.Int)
	priceBigFloat.Int(finalPriceBigInt)
	token.TokenID = args.TokenId
	if tokenDetailObj.Nonce != 0 {
		token.Nonce = tokenDetailObj.Nonce
		hexNonce := strconv.FormatInt(int64(token.Nonce), 16)
		if len(hexNonce)%2 != 0 {
			hexNonce = "0" + hexNonce
		}
		args.NonceStr = hexNonce
		args.Nonce = token.Nonce
	}
	token.RoyaltiesPercent = GetRoyaltiesPercentNominal(args.RoyaltiesPercent)
	token.MetadataLink = metadataLink
	token.CreatedAt = args.Timestamp
	token.Attributes = GetAttributesFromMetadata(metadataLink)
	token.ImageLink = args.FirstLink
	token.Hash = args.Hash
	token.TokenName = args.TokenName
	token.Status = entities.ListToken
	// token.PriceString = args.Price
	token.PriceString = finalPriceBigInt.String()
	token.PriceNominal = priceNominal
	token.OwnerID = ownerAccount.ID
	token.CollectionID = collectionId

	if args.NonceStr != "" {
		token.NonceStr = args.NonceStr
	}
	if args.TxHash != "" {
		token.MintTxHash = args.TxHash
	}
	if args.OnSale {
		token.OnSale = args.OnSale
	}
	if args.TxConfirmed {
		token.TxConfirmed = args.TxConfirmed
	}
	if args.AuctionStartTime > 0 {
		token.AuctionStartTime = args.AuctionStartTime
	}
	if args.AuctionDeadline > 0 {
		token.AuctionDeadline = args.AuctionDeadline
	}

	var innerErr error
	if !updateToken {
		innerErr = storage.AddToken(token)
		if innerErr == nil {
			_, cacheErr := AddTokenToCache(token.TokenID, token.Nonce, token.TokenName, token.ID)
			if cacheErr != nil {
				log.Error("could not add token to cache")
			}
			//}
		} else {
			innerErr = storage.UpdateToken(token)
		}
	}
	if innerErr != nil {
		log.Debug("could not create or update token", "err", innerErr)
		return
	}

	// Indexer is safer till later review
	// transaction := entities.Transaction{
	// 	Hash:         args.TxHash,
	// 	Type:         entities.ListToken,
	// 	PriceNominal: priceNominal,
	// 	Timestamp:    args.Timestamp,
	// 	SellerID:     ownerAccount.ID,
	// 	BuyerID:      0,
	// 	TokenID:      token.ID,
	// 	CollectionID: collectionId,
	// }

	// AddTransaction(&transaction)
}

func StakeToken(args StakeTokenArgs, blockchainProxy string, marketplaceAddress string) {

	var err error

	var updateToken = false
	token, err := storage.GetTokenByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		token = &entities.Token{}
	}

	//this logic should be expanded to other stake types

	token.Status = entities.ListToken
	token.StakeType = string(entities.DAO)
	token.OnStake = args.OnStake
	token.StakeDate = args.StakeDate

	var innerErr error
	if !updateToken {
		innerErr = storage.AddToken(token)
		if innerErr == nil {
			_, cacheErr := AddTokenToCache(token.TokenID, token.Nonce, token.TokenName, token.ID)
			if cacheErr != nil {
				log.Error("could not add token to cache")
			}
			//}
		} else {
			innerErr = storage.UpdateToken(token)
		}
	}
	if innerErr != nil {
		log.Debug("could not create or update token", "err", innerErr)
		return
	}

	// Indexer is safer till later review
	// transaction := entities.Transaction{
	// 	Hash:         args.TxHash,
	// 	Type:         entities.ListToken,
	// 	PriceNominal: priceNominal,
	// 	Timestamp:    args.Timestamp,
	// 	SellerID:     ownerAccount.ID,
	// 	BuyerID:      0,
	// 	TokenID:      token.ID,
	// 	CollectionID: collectionId,
	// }

	// AddTransaction(&transaction)
}

func getTokenBC(tokenName string, tokenHexNonce string, blockchainApi string) (entities.TokenBC, error) {
	//var resp ProxyTokenResponse
	var token entities.TokenBC

	url := fmt.Sprintf(GetNFTBaseFormat, blockchainApi, tokenName, tokenHexNonce)
	//err := HttpGet(url, &resp)
	response, err := GetResponse(url)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal([]byte(response), &token)
	if err != nil {
		return token, err
	}

	return token, nil

}
func getTokenByNonce(tokenName string, tokenNonce string, blockchainApi string) (NonFungibleToken, error) {
	//var resp ProxyTokenResponse
	var token NonFungibleToken

	intNonce, err := strconv.ParseUint(tokenNonce, 10, 64)
	if err != nil {
		return token, err
	}
	hexNonce := strconv.FormatUint(intNonce, 16)

	//Couldn't sort out padding and this quick check will work
	if len(hexNonce)%2 != 0 {
		hexNonce = "0" + hexNonce
	}

	url := fmt.Sprintf(GetNFTBaseFormat, blockchainApi, tokenName, hexNonce)

	//err := HttpGet(url, &resp)
	response, err := HttpGetRaw(url)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal([]byte(response), &token)
	if err != nil {
		return token, err
	}

	return token, nil

	//err = cache.GetCacher().Set(url, resp, HttpResponseExpirePeriod)
	//if err != nil {
	//	log.Debug("could not cache response", "err", err)
	//}

	//return response.Data.Token, nil
}

func BuyToken(args BuyTokenArgs) {

	var priceNominal float64
	var err error

	if args.PriceNominal != "" {
		priceNominal, err = strconv.ParseFloat(args.PriceNominal, 64)
		if err != nil {
			log.Debug("could not parse nominal", "err", err)
			return
		}
	} else {
		priceNominal, err = GetPriceNominal(args.Price)
		if err != nil {
			log.Debug("could not parse price", "err", err)
			return
		}
	}

	ownerAccount, err := storage.GetAccountByAddress(args.OwnerAddress)
	if err != nil {
		log.Debug("could not get owner account", "err", err)
		return
	}

	// buyerAccount, err := GetOrCreateAccount(args.BuyerAddress)
	// if err != nil {
	// 	log.Debug("could not get or create account", "err", err)
	// 	return
	// }

	token, err := storage.GetTokenByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get token", "err", err)
		return
	}

	if args.NonceStr != "" {
		token.NonceStr = args.NonceStr
	}

	if args.TxHash != "" {
		token.MintTxHash = args.TxHash
	}

	if args.TxConfirmed {
		token.TxConfirmed = args.TxConfirmed
	}

	priceBigFloat, err := ConvertBigFloatToFloat(args.PriceNominal)
	if err != nil {
		log.Debug("could not convert big float string to big float object", "err", err)
		return
	}
	fmt.Println(priceBigFloat.String())
	priceBigFloat, ok := TurnBigFloatoBigFloatNDec(priceBigFloat, 18)
	fmt.Println(priceBigFloat.String())
	if !ok {
		log.Debug("could not TurnBigFloatoBigFloatNDec priceBigFloat", "err", err)
		return
	}
	finalPriceBigInt := new(big.Int)
	priceBigFloat.Int(finalPriceBigInt)

	// Owner ID was to be reset since the token will no longer be on the marketplace.
	// Could have been kept like this, but bugs may appear when querying.
	token.OwnerID = ownerAccount.ID
	token.Status = entities.BuyToken
	token.OnSale = false
	token.LastBuyPriceNominal = priceNominal
	// dec18Float := TurnIntoBigFloat18Dec((priceNominal)) // priceNominal could be float and should be treated that way
	token.PriceString = finalPriceBigInt.String()
	err = storage.UpdateToken(token)
	if err != nil {
		log.Debug("could not update token", "err", err)
		return
	}
	//indexer is safer till later review TODO

	// transaction := entities.Transaction{
	// 	Hash:         args.TxHash,
	// 	Type:         entities.BuyToken,
	// 	PriceNominal: priceNominal,
	// 	Timestamp:    args.Timestamp,
	// 	SellerID:     ownerAccount.ID,
	// 	BuyerID:      buyerAccount.ID,
	// 	TokenID:      token.ID,
	// 	CollectionID: token.CollectionID,
	// }

	// AddTransaction(&transaction)
}

func StartAuction(args StartAuctionArgs, blockchainProxy string, marketplaceAddress string) (*entities.Token, error) {
	amountNominal, err := GetPriceNominal(args.MinBid)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return nil, err
	}

	accountID := uint64(0)
	accountCacheInfo, err := GetOrAddAccountCacheInfo(args.OwnerAddress)
	if err != nil {
		log.Debug("could not get or add acc cache info", err)

		account, innerErr := GetOrCreateAccount(args.OwnerAddress)
		if innerErr != nil {
			log.Debug("could not get or add acc", err)
		} else {
			accountID = account.ID
		}
	} else {
		accountID = accountCacheInfo.AccountId
	}

	collectionId := uint64(0)
	collectionInfoCache, err := collstats.GetOrAddCollectionCacheInfo(args.TokenId)
	if err == nil {
		collectionId = collectionInfoCache.CollectionId
	}

	token, err := storage.GetTokenByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		metadataLink := args.SecondLink
		if len(metadataLink) == 0 {
			var innerErr error
			metadataLink, innerErr = TryGetMetadataLink(blockchainProxy, marketplaceAddress, args.TokenId, args.Nonce)
			if innerErr != nil {
				log.Debug("could not get metadata link", innerErr)
			}
		}

		token = &entities.Token{
			TokenID:          args.TokenId,
			Nonce:            args.Nonce,
			RoyaltiesPercent: GetRoyaltiesPercentNominal(args.RoyaltiesPercent),
			MetadataLink:     metadataLink,
			CreatedAt:        args.Timestamp,
			Attributes:       GetAttributesFromMetadata(metadataLink),
			ImageLink:        args.FirstLink,
			Hash:             args.Hash,
			TokenName:        args.TokenName,
		}
	}

	token.Status = entities.AuctionToken
	token.OnSale = true
	token.PriceString = args.MinBid
	token.PriceNominal = amountNominal
	token.OwnerID = accountID
	token.CollectionID = collectionId
	token.AuctionStartTime = args.StartTime
	token.AuctionDeadline = args.Deadline

	var innerErr error
	if err != nil {
		innerErr = storage.AddToken(token)
		if innerErr == nil {
			_, cacheErr := AddTokenToCache(token.TokenID, token.Nonce, token.TokenName, token.ID)
			if cacheErr != nil {
				log.Error("could not add token to cache")
			}
		}
	} else {
		innerErr = storage.UpdateToken(token)
	}

	if innerErr != nil {
		log.Debug("could not create or update token", "err", innerErr)
		return nil, err
	}

	transaction := entities.Transaction{
		Hash:         args.TxHash,
		Type:         entities.AuctionToken,
		PriceNominal: amountNominal,
		Timestamp:    args.Timestamp,
		SellerID:     accountID,
		BuyerID:      0,
		TokenID:      token.ID,
		CollectionID: collectionId,
	}

	AddTransaction(&transaction)
	return token, nil
}

func EndAuction(args EndAuctionArgs) {
	amountNominal, err := GetPriceNominal(args.Amount)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

	buyer, err := GetOrAddAccountCacheInfo(args.Winner)
	if err != nil {
		log.Debug("could not parse price", "err", err)
		return
	}

	token, err := storage.GetTokenByTokenIdAndNonce(args.TokenId, args.Nonce)
	if err != nil {
		log.Debug("could not get token", "err", err)
		return
	}
	var txType entities.TxType = entities.BuyToken
	var winner bool = true
	if token.Owner.Address == args.Caller {
		//had no winner
		txType = entities.WithdrawToken
		winner = false
	}
	sellerId := token.OwnerID
	token.OwnerID = buyer.AccountId
	token.Status = txType
	token.OnSale = false
	if winner {
		token.LastBuyPriceNominal = amountNominal
	} else {
		token.LastBuyPriceNominal = token.PriceNominal

	}
	err = storage.UpdateToken(token)
	if err != nil {
		log.Debug("could not update token", "err", err)
		return
	}

	err = storage.DeleteOffersForTokenId(token.ID)
	if err != nil {
		log.Debug("could not delete offers for token", "err", err)
	}

	err = storage.DeleteBidsForTokenId(token.ID)
	if err != nil {
		log.Debug("could not delete bids for token", "err", err)
	}

	transaction := entities.Transaction{
		Hash:         args.TxHash,
		Type:         entities.BuyToken,
		PriceNominal: amountNominal,
		Timestamp:    args.Timestamp,
		SellerID:     sellerId,
		BuyerID:      buyer.AccountId,
		TokenID:      token.ID,
		CollectionID: token.CollectionID,
	}

	AddTransaction(&transaction)
}

func GetExtendedTokenData(tokenId string, nonce uint64) (*dtos.ExtendedTokenDto, error) {
	token, err := storage.GetTokenByTokenIdAndNonce(tokenId, nonce)
	if err != nil {
		return nil, err
	}

	owner, err := storage.GetAccountById(token.OwnerID)
	if err != nil {
		return nil, err
	}

	return &dtos.ExtendedTokenDto{
		Token:              *token,
		OwnerName:          owner.Name,
		OwnerWalletAddress: owner.Address,
	}, nil
}

func GetAttributesFromMetadata(link string) datatypes.JSON {
	emptyResponse := datatypes.JSON("")
	link = strings.TrimSpace(link)
	if len(link) == 0 {
		return emptyResponse
	}

	responseRaw, err := HttpGetRaw(link)
	if err != nil {
		log.Error("could not get metadata response", "link", link, "err", err.Error())
		return emptyResponse
	}
	if len(responseRaw) > maxTokenLinkResponseSize {
		log.Error("response too long for link", "link", link)
		return emptyResponse
	}

	var response dtos.MetadataLinkResponse
	err = json.Unmarshal([]byte(responseRaw), &response)
	if err != nil {
		log.Debug("could not unmarshal", "link", link, "err", err)
		return emptyResponse
	}

	attributes := make(map[string]interface{})
	for _, key := range response.Attributes {
		attributes[key.TraitType] = key.Value
	}

	bytes, err := json.Marshal(attributes)
	if err != nil {
		log.Debug("could not marshal", "link", link, "err", err)
		return emptyResponse
	}

	return bytes
}

func GetPriceNominal(priceHex string) (float64, error) {
	priceBigUint, success := big.NewInt(0).SetString(priceHex, 16)
	if !success {
		return 0, errors.New("could not parse price")
	}

	denominatorBigUint := big.NewInt(0).Exp(baseExp, big.NewInt(minPriceDecimals), nil)
	priceNominalInt := big.NewInt(0).Div(priceBigUint, denominatorBigUint).Int64()
	priceNominal := float64(priceNominalInt) / minPercentUnit
	return priceNominal, nil
}

func GetPriceDenominated(price float64) *big.Int {
	priceInt := int64(price * minPriceUnit)
	if priceInt <= 0 {
		log.Error("price less than min threshold",
			"min_threshold_multiplied", "1",
			"min_threshold_nominal", 1/minPriceUnit,
			"price_int", priceInt,
		)
	}

	denominatorBigUint := big.NewInt(0).Exp(baseExp, big.NewInt(minPriceDecimals), nil)

	priceBigUint := big.NewInt(0).Mul(big.NewInt(priceInt), denominatorBigUint)
	return priceBigUint
}

func GetRoyaltiesPercentNominal(percent uint64) float64 {
	return float64(percent) / minPercentRoyaltiesUnit
}

func AddTransaction(tx *entities.Transaction) {
	err := storage.AddTransaction(tx)
	if err != nil {
		log.Debug("could not create new transaction", "err", err)
		return
	}
}

func AddTokenToCache(tokenId string, nonce uint64, tokenName string, tokenDbId uint64) (*TokenCacheInfo, error) {
	db := cache.GetBolt()
	cacheInfo := TokenCacheInfo{
		TokenDbId: tokenDbId,
		TokenName: tokenName,
	}

	entryBytes, err := json.Marshal(&cacheInfo)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, innerErr := tx.CreateBucketIfNotExists(TokenIdToDbIdCacheInfo)
		if innerErr != nil {
			return innerErr
		}

		key := fmt.Sprintf("%s-%d", tokenId, nonce)
		innerErr = bucket.Put([]byte(key), entryBytes)
		return innerErr
	})

	return &cacheInfo, nil
}

func GetTokenCacheInfo(tokenId string, nonce uint64) (*TokenCacheInfo, error) {
	db := cache.GetBolt()

	var bytes []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(TokenIdToDbIdCacheInfo)
		if bucket == nil {
			return errors.New("no bucket for token cache")
		}

		key := fmt.Sprintf("%s-%d", tokenId, nonce)
		bytes = bucket.Get([]byte(key))
		return nil
	})
	if err != nil {
		return nil, err
	}

	var cacheInfo TokenCacheInfo
	err = json.Unmarshal(bytes, &cacheInfo)
	if err != nil {
		return nil, err
	}

	return &cacheInfo, nil
}

func GetOrAddTokenCacheInfo(tokenId string, nonce uint64) (*TokenCacheInfo, error) {
	// cacheInfo, err := GetTokenCacheInfo(tokenId, nonce)
	// if err != nil {
	token, innerErr := storage.GetTokenByTokenIdAndNonce(tokenId, nonce)
	if innerErr != nil {
		return nil, innerErr
	}

	// 	cacheInfo, innerErr = AddTokenToCache(tokenId, nonce, token.TokenName, token.ID)
	// 	if innerErr != nil {
	// 		return nil, innerErr
	// 	}
	// }

	return &TokenCacheInfo{TokenDbId: token.ID, TokenName: token.TokenName}, nil
}

func GetAvailableTokens(args AvailableTokensRequest) (AvailableTokensResponse, error) {
	var response AvailableTokensResponse
	if len(args.Tokens) > maxTokenNumAvailableSize {
		return response, tooManyTokensError
	}

	response.Tokens = make(map[string]AvailableToken)
	for _, token := range args.Tokens {
		parts := strings.Split(token, "-")
		if len(parts) != 3 {
			continue
		}

		tokenId := parts[0] + "-" + parts[1]
		nonce, err := strconv.ParseUint(parts[2], 16, 64)
		if err != nil {
			continue
		}

		tokenName := ""
		tokenAvailable := false
		tokenCacheInfo, err := GetOrAddTokenCacheInfo(tokenId, nonce)
		if err == nil {
			tokenAvailable = true
			tokenName = tokenCacheInfo.TokenName
		}

		collectionName := ""
		collectionAvailable := false
		collectionCacheInfo, err := collstats.GetOrAddCollectionCacheInfo(tokenId)
		if err == nil {
			collectionAvailable = true
			collectionName = collectionCacheInfo.CollectionName
		}

		response.Tokens[token] = AvailableToken{
			Collection: struct {
				Id        string `json:"id"`
				Name      string `json:"name"`
				Available bool   `json:"available"`
			}{
				Id:        tokenId,
				Name:      collectionName,
				Available: collectionAvailable,
			},
			Token: struct {
				Id        string `json:"id"`
				Nonce     uint64 `json:"nonce"`
				Name      string `json:"name"`
				Available bool   `json:"available"`
			}{
				Id:        tokenId,
				Nonce:     nonce,
				Name:      tokenName,
				Available: tokenAvailable,
			},
		}
	}

	return response, nil
}

func TryGetMetadataLink(blockchainProxy string, address string, tokenId string, nonce uint64) (string, error) {
	proxyRequest := fmt.Sprintf(NftProxyRequestFormat, blockchainProxy, address, tokenId, nonce)

	var proxyResponse NftProxyResponse
	err := HttpGet(proxyRequest, &proxyResponse)
	if err != nil {
		log.Error("get metadata link request failed", "tokenID", tokenId, "nonce", nonce, "err", err.Error())
		return "", err
	}
	if len(proxyResponse.Data.TokenData.Uris) < 2 {
		return "", nil
	}

	link, err := base64.StdEncoding.DecodeString(proxyResponse.Data.TokenData.Uris[1])
	linkStr := string(link)

	return ParseMetadataUrl(linkStr), err
}

func ParseMetadataUrl(link string) string {
	if strings.HasPrefix(link, ipfsProtocolURLPrefix) {
		parsedUrl := fmt.Sprintf(
			ipfsDefaultGatewayURL,
			strings.Replace(link, ipfsProtocolURLPrefix, "", 1),
		)
		link = parsedUrl
	}

	return link
}

func TryGetTokenResponse(blockchainProxy string, address string, tokenId string, nonce uint64) (NftProxyReponseToken, error) {
	proxyRequest := fmt.Sprintf(NftProxyRequestFormat, blockchainProxy, address, tokenId, nonce)

	var proxyResponse NftProxyResponse
	err := HttpGet(proxyRequest, &proxyResponse)
	if err != nil {
		log.Debug("get token request failed")
		return NftProxyReponseToken{}, err
	}

	return proxyResponse.Data.TokenData, err
}

func ConstructOwnedTokensFromTokens(tokens []entities.Token) []dtos.OwnedTokenDto {
	tokenIds := make(map[string]bool)
	for _, token := range tokens {
		tokenIds[token.TokenID] = true
	}

	collections := make(map[string]dtos.CollectionCacheInfo)
	for tokenId := range tokenIds {
		info, innerErr := collstats.GetOrAddCollectionCacheInfo(tokenId)
		if innerErr == nil {
			collections[tokenId] = *info
		}
	}

	ownedTokens := make([]dtos.OwnedTokenDto, len(tokens))
	for index, token := range tokens {
		ownedToken := dtos.OwnedTokenDto{
			Token:               token,
			CollectionCacheInfo: collections[token.TokenID],
		}
		ownedTokens[index] = ownedToken
	}

	return ownedTokens
}

func ClearResponseCached(url string) error {
	redis := cache.GetRedis()
	redisCtx := cache.GetContext()
	url = strings.TrimSpace(url)

	key := fmt.Sprintf(UrlResponseCacheKeyFormat, url)
	cmd := redis.Del(redisCtx, key)
	_, err := cmd.Result()
	if err != nil {
		return nil
	}
	return err
}
func TryGetResponseCached(url string) (string, error) {
	redis := cache.GetRedis()
	redisCtx := cache.GetContext()
	url = strings.TrimSpace(url)

	key := fmt.Sprintf(UrlResponseCacheKeyFormat, url)
	metadataBytes, err := redis.Get(redisCtx, key).Result()
	if err == nil {
		return metadataBytes, nil
	}

	metadataBytes, err = HttpGetRaw(url)
	if err != nil {
		log.Debug("http get returned error", err)
	}
	if len(metadataBytes) > maxTokenLinkResponseSize {
		metadataBytes = ""
	}

	err = redis.Set(redisCtx, key, metadataBytes, UrlResponseExpirePeriod).Err()
	if err != nil {
		log.Debug("could not set to redis", err)
	}

	return metadataBytes, nil
}

func AddOrRefreshToken(
	tokenId string,
	nonce uint64,
	collectionId uint64,
	userAddress string,
	blockchainProxy string,
	marketplaceAddress string,
) (datatypes.JSON, error) {
	redisClient := cache.GetRedis()
	redisContext := cache.GetContext()

	tokenIsInDb := false
	attributes := datatypes.JSON("")
	emptyAttributes := datatypes.JSON("")
	token, err := storage.GetTokenByTokenIdAndNonce(tokenId, nonce)
	if err == nil {
		tokenIsInDb = true
		attributes = token.Attributes
	}

	refreshKey := fmt.Sprintf(RefreshMetadataSetNxKeyFormat, tokenId, nonce)
	ok, err := redisClient.SetNX(redisContext, refreshKey, true, RefreshMetadataSetNxExpirePeriod).Result()
	if err != nil {
		log.Debug("set nx resulted in error", "err", err.Error())
	}

	shouldTry := ok == true && err == nil
	if !shouldTry {
		return JsonOrEmpty(attributes), nil
	}

	if !tokenIsInDb {
		token = &entities.Token{
			TokenID:    tokenId,
			Nonce:      nonce,
			CreatedAt:  uint64(time.Now().Unix()),
			Status:     entities.None,
			Attributes: datatypes.JSON(""),
		}
	}

	tokenProxyResponse, err := TryGetTokenResponse(blockchainProxy, userAddress, tokenId, nonce)
	if err != nil {
		var innerErr error
		tokenProxyResponse, innerErr = TryGetTokenResponse(blockchainProxy, marketplaceAddress, tokenId, nonce)
		if innerErr != nil {
			return emptyAttributes, innerErr
		}
	}

	if tokenProxyResponse.Balance != "1" {
		return emptyAttributes, errors.New("balance not 1")
	}

	if len(tokenProxyResponse.Uris) == 0 {
		return emptyAttributes, errors.New("no uris")
	}

	metadataLink := ""
	if len(tokenProxyResponse.Uris) >= 2 {
		link, innerErr := base64.StdEncoding.DecodeString(tokenProxyResponse.Uris[1])
		if innerErr != nil {
			return emptyAttributes, innerErr
		}

		metadataLink = ParseMetadataUrl(string(link))
	}

	if len(tokenProxyResponse.Royalties) == 0 {
		tokenProxyResponse.Royalties = "0"
	}

	royalties, err := strconv.Atoi(tokenProxyResponse.Royalties)
	if err != nil {
		return emptyAttributes, nil
	}

	royaltiesNominal := GetRoyaltiesPercentNominal(uint64(royalties))
	if royalties > maxPercentRoyaltiesAllowed {
		return emptyAttributes, nil
	}

	link, err := base64.StdEncoding.DecodeString(tokenProxyResponse.Uris[0])
	if err != nil {
		return emptyAttributes, err
	}

	token.ImageLink = string(link)
	token.CollectionID = collectionId
	token.RoyaltiesPercent = royaltiesNominal
	token.MetadataLink = metadataLink
	token.TokenName = tokenProxyResponse.Name
	token.Hash = tokenProxyResponse.Hash

	newAttributes := GetAttributesFromMetadata(metadataLink)
	if len(newAttributes) != 0 {
		token.Attributes = newAttributes
	}

	var innerErr error
	if tokenIsInDb {
		innerErr = storage.UpdateToken(token)
		if innerErr != nil {
			log.Debug("could not update token")
		}
	} else {
		innerErr = storage.AddToken(token)
		if innerErr != nil {
			log.Debug("could not add token")
		}
	}
	if innerErr != nil {
		return emptyAttributes, nil
	}

	return JsonOrEmpty(token.Attributes), nil
}

func JsonOrEmpty(value datatypes.JSON) datatypes.JSON {
	if value != nil {
		return value
	}

	return datatypes.JSON("")
}

func GetTokensWithTokenIdAlike(tokenId string, limit int) ([]entities.Token, error) {
	var byteArray []byte
	var tokenArray []entities.Token

	cacheKey := fmt.Sprintf(TokenSearchCacheKeyFormat, tokenId)
	err := cache.GetCacher().Get(cacheKey, &byteArray)
	if err == nil {
		err = json.Unmarshal(byteArray, &tokenArray)
		return tokenArray, err
	}

	searchName := "%" + tokenId + "%"
	tokenArray, err = storage.GetTokensWithTokenIdAlikeWithLimit(searchName, limit)
	if err != nil {
		return nil, err
	}

	byteArray, err = json.Marshal(tokenArray)
	if err == nil {
		err = cache.GetCacher().Set(cacheKey, byteArray, TokenSearchExpirePeriod)
		if err != nil {
			log.Debug("could not set cache", "err", err)
		}
	}

	return tokenArray, nil
}

func GetTokensListedWithTokenIdAlikeWithStatus(tokenId string, limit int) ([]entities.Token, error) {
	var byteArray []byte
	var tokenArray []entities.Token

	cacheKey := fmt.Sprintf(TokenSearchCacheKeyFormat, tokenId)
	err := cache.GetCacher().Get(cacheKey, &byteArray)
	if err == nil {
		err = json.Unmarshal(byteArray, &tokenArray)
		return tokenArray, err
	}

	searchName := "%" + tokenId + "%"
	tokenArray, err = storage.GetTokensListedWithTokenIdAlikeWithLimit(searchName, limit)
	if err != nil {
		return nil, err
	}

	byteArray, err = json.Marshal(tokenArray)
	if err == nil {
		err = cache.GetCacher().Set(cacheKey, byteArray, TokenSearchExpirePeriod)
		if err != nil {
			log.Debug("could not set cache", "err", err)
		}
	}

	return tokenArray, nil
}

func GetTokensUnlistedWithTokenIdAlikeWithStatus(tokenId string, limit int) ([]entities.Token, error) {
	var byteArray []byte
	var tokenArray []entities.Token

	cacheKey := fmt.Sprintf(TokenSearchCacheKeyFormat, tokenId)
	err := cache.GetCacher().Get(cacheKey, &byteArray)
	if err == nil {
		err = json.Unmarshal(byteArray, &tokenArray)
		return tokenArray, err
	}

	searchName := "%" + tokenId + "%"
	tokenArray, err = storage.GetTokensUnlistedWithTokenIdAlikeWithLimit(searchName, limit)
	if err != nil {
		return nil, err
	}

	byteArray, err = json.Marshal(tokenArray)
	if err == nil {
		err = cache.GetCacher().Set(cacheKey, byteArray, TokenSearchExpirePeriod)
		if err != nil {
			log.Debug("could not set cache", "err", err)
		}
	}

	return tokenArray, nil
}

func GetBuyerWhiteListCheck(contractAddress string) (string, error) {

	strWhiteListCheck, err := DoGetBuyerWhitelistCheckVmQuery(contractAddress)
	if err != nil {
		return "", err
	}

	return strWhiteListCheck, nil //return "ON" | "OFF"
}

func DoGetBuyerWhitelistCheckVmQuery(contractAddress string) (string, error) {

	bi := interaction.GetBlockchainInteractor()

	whiteListCheck, errWhiteListCheck := bi.DoVmQuery(contractAddress, "getBuyerWhiteListCheck", []string{})
	if errWhiteListCheck != nil {
		return "", errWhiteListCheck
	}

	whiteListReturn := "ON"

	if len(whiteListCheck) != 0 {
		whiteListCheckInt := big.NewInt(0).SetBytes(whiteListCheck[0])
		if whiteListCheckInt.String() == "1" {
			whiteListReturn = "ON"
		} else {
			whiteListReturn = "OFF"
		}

	}

	return whiteListReturn, nil
}

func GetWhitelistBuyCountLimit(contractAddress string, userAddress string) (string, error) {
	localCacher := cache.GetLocalCacher()
	key := fmt.Sprintf(GetBuyerLimitCount, userAddress, contractAddress)

	count, errRead := localCacher.Get(key)
	if errRead == nil {
		return strconv.FormatInt(count.(int64), 10), nil
	}

	strBuyCountLimit, err := DoGetWhitelistBuyCountLimitVmQuery(contractAddress, userAddress)
	if err != nil {
		return "", err
	}

	return strBuyCountLimit, nil
}

func DoGetWhitelistBuyCountLimitVmQuery(contractAddress string, userAddress string) (string, error) {

	bi := interaction.GetBlockchainInteractor()

	//get the user address prep (decoded & Hexed)
	userAddressDecoded, errUserAddress := data.NewAddressFromBech32String(userAddress)
	if errUserAddress != nil {
		return "", errUserAddress
	}
	userAddressHex := hex.EncodeToString(userAddressDecoded.AddressBytes())

	whiteListCheck, errWhiteListCheck := bi.DoVmQuery(contractAddress, "getBuyerWhiteListCheck", []string{})
	if errWhiteListCheck != nil {
		return "", errWhiteListCheck
	}

	if len(whiteListCheck) != 0 {
		buyCount := big.NewInt(0).SetBytes(whiteListCheck[0])
		if buyCount.String() != "1" {
			strBuyCountLimit := "-1" + "," + "-1"
			return strBuyCountLimit, nil
		}
	}
	//get the "buy_count" from SC
	resultBuyCount, errBuyCount := bi.DoVmQuery(contractAddress, "getBuyCount", []string{userAddressHex})
	if errBuyCount != nil {
		return "", errBuyCount
	}

	strBuyCount := "0"
	if len(resultBuyCount) != 0 {
		buyCount := big.NewInt(0).SetBytes(resultBuyCount[0])
		strBuyCount = buyCount.String()
	}

	//get the "buy_limit" from SC
	resultBuyLimit, errBuyLimit := bi.DoVmQuery(contractAddress, "getBuyLimit", []string{userAddressHex})

	if errBuyLimit != nil {
		return "", errBuyLimit
	}

	strBuyLimit := "0"
	if len(resultBuyLimit) != 0 {
		buyLimit := big.NewInt(0).SetBytes(resultBuyLimit[0])
		strBuyLimit = buyLimit.String()
	}

	strBuyCountLimit := strBuyCount + "," + strBuyLimit
	return strBuyCountLimit, nil
}

func ChangeTokenOwner(tokenIdentifier string, nonceStr string, newOwner string, blockchainAPI string) error {
	token, err := storage.GetTokenByTokenIdAndNonceStr(tokenIdentifier, nonceStr)
	if err != nil {
		return err
	}
	if token.OnSale == true {
		return nil
	}
	if token.Owner.Address == newOwner {
		return nil
	}
	reqURL := fmt.Sprintf("%s/nfts/%s", blockchainAPI, tokenIdentifier+"-"+nonceStr)
	resByte, err := GetResponse(reqURL)
	if err != nil {
		return err
	}
	var tokenBC entities.TokenBC
	err = json.Unmarshal(resByte, &tokenBC)
	if err != nil {
		return err
	}
	if tokenBC.Owner == newOwner {
		acc, err := GetOrCreateAccount(tokenBC.Owner)
		if err != nil {
			return err
		}

		token.OwnerID = acc.ID
		err = storage.UpdateToken(token)
		if err != nil {
			return err
		}
	}
	return nil
}
