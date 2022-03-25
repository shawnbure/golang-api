package services

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/interaction"
	"github.com/ENFT-DAO/youbei-api/stats/collstats"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/boltdb/bolt"
)

const (
	TokenIdMaxLen                  = 17
	MaxFlags                       = 10
	MaxFlagLen                     = 25
	MaxNameLen                     = 20
	MaxLinkLen                     = 100
	MaxDescLen                     = 1000
	RegisteredNFTsBaseFormat       = "%s/address/%s/registered-nfts"
	GetCollectionBaseFormat        = "%s/collections/%s"
	GetNFTBaseFormat               = "%s/nfts/%s-%s" //example https://devnet-api.elrond.com/nfts/AMIR-55a2ea-01
	HttpResponseExpirePeriod       = 10 * time.Minute
	CollectionSearchCacheKeyFormat = "CollectionSearch:%s"
	CollectionSearchExpirePeriod   = 20 * time.Minute
	MintInfoViewName               = "getMaxSupplyAndTotalSold"
	MintInfoSetNxKeyFormat         = "MintInfoNX:%s"
	MintInfoSetNxExpirePeriod      = 6 * time.Second
	MintInfoBucketName             = "MintInfo"

	CollectionVerifiedCacheKeyFormat = "CollectionVerifiedCacheKey"
	CollectionVerifiedExpirePeriod   = 5 * time.Minute

	CollectionNoteworthyCacheKeyFormat = "CollectionNoteworthyCacheKey"
	CollectionNoteworthyExpirePeriod   = 5 * time.Minute

	CollectionTrendingCacheKeyFormat = "CollectionTrendingCacheKey"
	CollectionTrendingExpirePeriod   = 5 * time.Minute
)

type MintInfo struct {
	MaxSupply uint64 `json:"maxSupply"`
	TotalSold uint64 `json:"totalSold"`
}

type CreateCollectionRequest struct {
	UserAddress             string   `json:"userAddress"`
	Name                    string   `json:"collectionName"`
	TokenId                 string   `json:"tokenId"`
	Description             string   `json:"description"`
	Website                 string   `json:"website"`
	DiscordLink             string   `json:"discordLink"`
	TwitterLink             string   `json:"twitterLink"`
	InstagramLink           string   `json:"instagramLink"`
	TelegramLink            string   `json:"telegramLink"`
	Flags                   []string `json:"flags"`
	ContractAddress         string   `json:"contractAddress"`
	MintPricePerTokenString string   `json:"mintPricePerTokenString"`
	TokenBaseURI            string   `json:"tokenBaseURI"`
	MaxSupply               uint64   `json:"maxSupply"`
	MetaDataBaseURI         string   `json:"metaDataBaseURI"`
}

type AutoCreateCollectionRequest struct {
	UserAddress     string `json:"userAddress"`
	TokenId         string `json:"tokenId"`
	Nonce           string `json:"nonce"`
	Name            string `json:"collectionName"`
	CreatorAddress  string `json:"creatorAddress"`
	TokenBaseURI    string `json:"tokenBaseURI"`
	MetadataBaseURI string `json:"metadataBaseURI"`
}

type UpdateCollectionRequest struct {
	Description   string   `json:"description"`
	Website       string   `json:"website"`
	DiscordLink   string   `json:"discordLink"`
	TwitterLink   string   `json:"twitterLink"`
	InstagramLink string   `json:"instagramLink"`
	TelegramLink  string   `json:"telegramLink"`
	Flags         []string `json:"flags"`
}

type UpdateCollectionObj struct {
	UpdateCollectionRequest
	Address string `json:"string"`
}

type ProxyRegisteredNFTsResponse struct {
	Data struct {
		Tokens []string `json:"tokens"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type CollectionDetailBCResponse struct {
	Collection      string                   `json:"collection"`
	Type            string                   `json:"type"`
	Name            string                   `json:"name"`
	Ticker          string                   `json:"ticker"`
	Owner           string                   `json:"owner"`
	Timestamp       uint64                   `json:"timestamp"`
	CanFreeze       bool                     `json:"canFreeze"`
	CanWipe         bool                     `json:"canWipe"`
	CanTransferRole bool                     `json:"canTransferRole"`
	Roles           []map[string]interface{} `json:"roles"`
}

func CreateCollection(request *CreateCollectionRequest, blockchainProxy string) (*entities.Collection, error) {

	err := checkValidInputOnCreate(request)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(request.Flags)
	if err != nil {
		return nil, err
	}

	tokensRegisteredByUser, err := getTokensRegisteredByUser(request.UserAddress, blockchainProxy)
	if err != nil {
		return nil, err
	}
	if !contains(tokensRegisteredByUser, request.TokenId) {
		return nil, errors.New("token not owner by user")
	}

	standardizedName := standardizeName(request.Name)
	_, err = storage.GetCollectionWithNameILike(standardizedName)
	if err == nil {
		return nil, errors.New("collection name already taken")
	}

	_, err = storage.GetCollectionByTokenId(request.TokenId)
	if err == nil {
		return nil, errors.New("token id already has a collection associated")
	}

	account, err := GetOrCreateAccount(request.UserAddress)
	if err != nil {
		return nil, err
	}

	// ContractAddress
	// nominal price

	/*
		mintPricePerTokenNominalrequest, err := strconv.ParseFloat(request.MintPricePerTokenString, 64)

		if err != nil {
			mintPricePerTokenNominalrequest = 0.1
		}
	*/

	//strMintPricePerToken = request.MintPricePerTokenString
	//const fMintPricePerTokenNominal = 0.0
	floatPrice, ok := big.NewFloat(0).SetString(request.MintPricePerTokenString)
	if !ok {
		return nil, errors.New("couldn't convert price string to blockchain unit")
	}
	multiplier := new(big.Float)
	multiplier.SetInt(big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
	floatPrice.Mul(multiplier, floatPrice)
	fmt.Println(floatPrice.String())
	priceBig, _ := floatPrice.Int(nil)
	// priceBig, ok := big.NewInt(0).Setint(uintPrice)
	if !ok {
		return nil, errors.New("couldn't convert price string to blockchain unit")
	}
	mintPricePerTokenNominalrequest, err := strconv.ParseFloat(request.MintPricePerTokenString, 64)

	if err != nil {
		mintPricePerTokenNominalrequest = 0.1
	}

	collection := &entities.Collection{
		ID:                       0,
		Name:                     request.Name,
		TokenID:                  request.TokenId,
		Description:              "",
		Website:                  "",
		DiscordLink:              "",
		TwitterLink:              "",
		InstagramLink:            "",
		TelegramLink:             "",
		Flags:                    datatypes.JSON(bytes),
		ContractAddress:          request.ContractAddress,
		MintPricePerTokenString:  priceBig.String(),
		MintPricePerTokenNominal: mintPricePerTokenNominalrequest,
		TokenBaseURI:             request.TokenBaseURI,
		MetaDataBaseURI:          request.MetaDataBaseURI,
		MaxSupply:                request.MaxSupply,
		CreatorID:                account.ID,
		CreatedAt:                uint64(time.Now().Unix()),
	}

	err = storage.AddCollection(collection)
	if err != nil {
		return nil, err
	}

	_, err = collstats.AddCollectionToCache(collection.ID, collection.Name, collection.Flags, collection.TokenID)
	if err != nil {
		log.Debug("could not add to coll stats")
	}

	return collection, nil
}
func CreateCollectionFromToken(token entities.TokenBC, blockchainApi string) (*entities.Collection, error) {

	// ========== STEP: CHECK TO SEE IF COLLECTION EXIST ==========
	collection, errGetCollection := storage.GetCollectionByTokenId(token.Collection)

	if errGetCollection != nil {

		if errGetCollection == gorm.ErrRecordNotFound {
			//no collection found so don't exit out of the process and
			//auto create collection
		} else {
			return nil, errGetCollection
		}
	} else {
		//collection Found - so no need to auto create
		return collection, nil
	}

	//set the token creator adress
	tokenCreatorAddress := token.Owner

	//
	// ========== STEP: GET CREATOR ID FROM ACCOUNT BY ADDRESS   ==========
	//get account to get the "creator id"
	account, err := storage.GetAccountByAddress(tokenCreatorAddress)

	if err != nil {
		return nil, err
	}

	creatorID := account.ID

	// ========== STEP: AUTO CREATE COLLECTION ==========
	collection = &entities.Collection{
		Name:      token.Collection,
		TokenID:   token.Identifier,
		CreatorID: uint64(creatorID), //set the creator id
		CreatedAt: uint64(time.Now().Unix()),
	}

	// ========== GET COLLECTION DETAIL FROM BC ========
	colDetail, err := GetCollectionDetailBC(token.Collection, blockchainApi)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		collection.Name = colDetail.Name
		var address string
		for _, role := range colDetail.Roles {
			rolesStr, ok := role["roles"].([]string)
			if ok {
				for _, roleStr := range rolesStr {
					if strings.EqualFold(roleStr, "ESDTRoleNFTCreate") {
						address = role["address"].(string)
					}
				}
			}
		}
		collection.ContractAddress = address
		// collection.ContractAddress = colDetail.Roles[]
	}
	//
	errCollection := storage.AddCollection(collection)
	if errCollection != nil {
		fmt.Printf("%n", errCollection)
		return nil, errCollection
	}

	_, err = collstats.AddCollectionToCache(collection.ID, collection.Name, nil, collection.TokenID)
	if err != nil {
		log.Debug("could not add to coll stats")
	}

	return collection, nil
}
func AutoCreateCollection(request *AutoCreateCollectionRequest, blockchainApi string) (*entities.Collection, error) {

	// ========== STEP: CHECK TO SEE IF COLLECTION EXIST ==========
	collection, errGetCollection := storage.GetCollectionByTokenId(request.TokenId)

	if errGetCollection != nil {

		if errGetCollection == gorm.ErrRecordNotFound {
			//no collection found so don't exit out of the process and
			//auto create collection
		} else {
			return nil, errGetCollection
		}
	} else {
		//collection Found - so no need to auto create
		return collection, nil
	}

	//set the token creator adress
	tokenCreatorAddress := request.CreatorAddress

	// ========== STEP: CHECK IF USER IS THE CREATOR OF TOKEN  ==========
	// if user wallet not the creator of token, then check if the original
	// creator is register - if not register, then auto register
	if request.UserAddress != tokenCreatorAddress {
		//user is NOT the creator of the token, so auto creator the creator account

		//Check if address is not already in account
		//(for cases userWallet in account but haven't register account)
		_, err := storage.GetAccountByAddress(tokenCreatorAddress)
		if err != nil && err == gorm.ErrRecordNotFound {
			//account doesn't exist, so auto register
			accountTokenCreator := &entities.Account{
				Name:      request.TokenId + " Creator",
				Address:   tokenCreatorAddress,
				CreatedAt: uint64(time.Now().Unix()),
			}
			errAddAccount := storage.AddAccount(accountTokenCreator)
			if errAddAccount != nil {
				if !strings.Contains(errAddAccount.Error(), "duplicate") {
					return nil, err
				} else {
					err = storage.UpdateAccountProfileWhereName(accountTokenCreator.Name, *accountTokenCreator)
					if err != nil {
						if err != gorm.ErrRecordNotFound {
							return nil, err
						}
					}
				}
			}
		}
	}

	//
	// ========== STEP: GET CREATOR ID FROM ACCOUNT BY ADDRESS   ==========
	//get account to get the "creator id"
	account, err := storage.GetAccountByAddress(tokenCreatorAddress)

	if err != nil {
		return nil, err
	}

	creatorID := account.ID

	// ========== STEP: AUTO CREATE COLLECTION ==========
	collection = &entities.Collection{
		Name:      request.Name,
		TokenID:   request.TokenId,
		CreatorID: uint64(creatorID), //set the creator id
		CreatedAt: uint64(time.Now().Unix()),
	}

	// ========== GET COLLECTION DETAIL FROM BC ========
	colDetail, err := GetCollectionDetailBC(request.TokenId, blockchainApi)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		collection.Name = colDetail.Name
		var address string
		for _, role := range colDetail.Roles {
			rolesStr, ok := role["roles"].([]string)
			if ok {
				for _, roleStr := range rolesStr {
					if strings.EqualFold(roleStr, "ESDTRoleNFTCreate") {
						address = role["address"].(string)
					}
				}
			}
		}
		collection.ContractAddress = address
		// collection.ContractAddress = colDetail.Roles[]
	}
	//
	errCollection := storage.AddCollection(collection)
	if errCollection != nil {
		fmt.Printf("%n", errCollection)
		return nil, errCollection
	}

	_, err = collstats.AddCollectionToCache(collection.ID, collection.Name, nil, collection.TokenID)
	if err != nil {
		log.Debug("could not add to coll stats")
	}

	return collection, nil
}

func UpdateCollection(collection *entities.Collection, request *UpdateCollectionRequest) error {
	err := checkValidInputOnUpdate(request)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(request.Flags)
	if err != nil {
		return err
	}

	collection.Description = request.Description
	collection.Website = request.Website
	collection.DiscordLink = request.DiscordLink
	collection.TwitterLink = request.TwitterLink
	collection.InstagramLink = request.InstagramLink
	collection.TelegramLink = request.TelegramLink
	collection.Flags = bytes

	err = storage.UpdateCollection(collection)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCollectionWithAddress(collection *entities.Collection, request map[string]interface{}) error {

	collection.ContractAddress = request["ContractAddress"].(string)
	collection.Name = request["Name"].(string)

	err := storage.UpdateCollection(collection)
	if err != nil {
		return err
	}

	return nil
}

func GetAllCollections() ([]entities.Collection, error) {
	//var collectionArray []entities.Collection

	collectionArray, err := storage.GetAllCollections()
	if err != nil {
		return nil, err
	}

	return collectionArray, nil
}

func GetCollectionsWithNameAlike(name string, limit int) ([]entities.Collection, error) {
	var byteArray []byte
	var collectionArray []entities.Collection

	cacheKey := fmt.Sprintf(CollectionSearchCacheKeyFormat, name)
	err := cache.GetCacher().Get(cacheKey, &byteArray)
	if err == nil {
		err = json.Unmarshal(byteArray, &collectionArray)
		return collectionArray, err
	}

	searchName := "%" + name + "%"
	collectionArray, err = storage.GetCollectionsWithNameAlikeWithLimit(searchName, limit)
	if err != nil {
		return nil, err
	}

	byteArray, err = json.Marshal(collectionArray)
	if err == nil {
		err = cache.GetCacher().Set(cacheKey, byteArray, CollectionSearchExpirePeriod)
		if err != nil {
			log.Debug("could not set cache", "err", err)
		}
	}

	return collectionArray, nil
}

func GetCollectionsVerified(limit int) ([]entities.Collection, error) {
	var byteArray []byte
	var collectionArray []entities.Collection

	cacheKey := CollectionVerifiedCacheKeyFormat
	err := cache.GetCacher().Get(cacheKey, &byteArray)
	if err == nil {
		err = json.Unmarshal(byteArray, &collectionArray)
		return collectionArray, err
	}

	collectionArray, err = storage.GetCollectionsVerified(limit)
	if err != nil {
		return nil, err
	}

	byteArray, err = json.Marshal(collectionArray)
	if err == nil {
		err = cache.GetCacher().Set(cacheKey, byteArray, CollectionVerifiedExpirePeriod)
		if err != nil {
			log.Debug("could not set cache", "err", err)
		}
	}

	return collectionArray, nil
}

func GetCollectionsNoteworthy(limit int) ([]entities.Collection, error) {
	var byteArray []byte
	var collectionArray []entities.Collection

	cacheKey := CollectionNoteworthyCacheKeyFormat
	err := cache.GetCacher().Get(cacheKey, &byteArray)
	if err == nil {
		err = json.Unmarshal(byteArray, &collectionArray)
		return collectionArray, err
	}

	collectionArray, err = storage.GetCollectionsNoteworthy(limit)
	if err != nil {
		return nil, err
	}

	byteArray, err = json.Marshal(collectionArray)
	if err == nil {
		err = cache.GetCacher().Set(cacheKey, byteArray, CollectionNoteworthyExpirePeriod)
		if err != nil {
			log.Debug("could not set cache", "err", err)
		}
	}

	return collectionArray, nil
}

func GetCollectionsTrending(limit int) ([]entities.Collection, error) {
	var byteArray []byte
	var collectionArray []entities.Collection

	cacheKey := CollectionTrendingCacheKeyFormat
	err := cache.GetCacher().Get(cacheKey, &byteArray)
	if err == nil {
		err = json.Unmarshal(byteArray, &collectionArray)
		return collectionArray, err
	}

	collectionArray, err = storage.GetCollectionsTrending(limit)
	if err != nil {
		return nil, err
	}

	byteArray, err = json.Marshal(collectionArray)
	if err == nil {
		err = cache.GetCacher().Set(cacheKey, byteArray, CollectionTrendingExpirePeriod)
		if err != nil {
			log.Debug("could not set cache", "err", err)
		}
	}

	return collectionArray, nil
}

func GetMintInfoForContract(contractAddress string) (*MintInfo, error) {
	redisClient := cache.GetRedis()
	redisContext := cache.GetContext()

	mintInfoSetNxKey := fmt.Sprintf(MintInfoSetNxKeyFormat, contractAddress)
	ok, err := redisClient.SetNX(redisContext, mintInfoSetNxKey, true, MintInfoSetNxExpirePeriod).Result()
	if err != nil {
		log.Debug("set nx resulted in error", err)
	}

	shouldDoQuery := ok == true && err == nil
	if shouldDoQuery {
		mintInfo, innerErr := setMintInfoCache(contractAddress)
		if innerErr != nil {
			_, _ = redisClient.Del(redisContext, mintInfoSetNxKey).Result()
			return nil, innerErr
		}
		return mintInfo, nil
	} else {
		return getMintInfoCache(contractAddress)
	}
}

func setMintInfoCache(contractAddress string) (*MintInfo, error) {
	db := cache.GetBolt()

	bi := interaction.GetBlockchainInteractor()
	if bi == nil {
		return nil, errors.New("no blockchain interactor")
	}

	result, err := bi.DoVmQuery(contractAddress, MintInfoViewName, []string{})
	if err != nil {
		return nil, err
	}
	if len(result) != 2 {
		return nil, errors.New("unknown result len")
	}

	maxSupplyHex := hex.EncodeToString(result[0])
	maxSupply, err := strconv.ParseUint(maxSupplyHex, 16, 64)
	if err != nil {
		return nil, err
	}

	totalSoldHex := hex.EncodeToString(result[1])
	if len(result[1]) == 0 {
		totalSoldHex = "0"
	}
	totalSold, err := strconv.ParseUint(totalSoldHex, 16, 64)
	if err != nil {
		return nil, err
	}

	mintInfo := MintInfo{
		MaxSupply: maxSupply,
		TotalSold: totalSold,
	}

	entryBytes, err := json.Marshal(&mintInfo)
	if err != nil {
		log.Debug("could not marshal", err.Error())
	} else {
		innerErr := db.Update(func(tx *bolt.Tx) error {
			bucket, innerErr := tx.CreateBucketIfNotExists([]byte(MintInfoBucketName))
			if innerErr != nil {
				return innerErr
			}

			innerErr = bucket.Put([]byte(contractAddress), entryBytes)
			return innerErr
		})

		if innerErr != nil {
			log.Debug("could not set to bolt db")
		}
	}

	return &mintInfo, nil
}

func getMintInfoCache(contractAddress string) (*MintInfo, error) {
	db := cache.GetBolt()

	var bytes []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(MintInfoBucketName))
		if bucket == nil {
			return errors.New("no bucket for collection cache")
		}

		bytes = bucket.Get([]byte(contractAddress))
		return nil
	})
	if err != nil {
		return nil, err
	}

	var mintInfo MintInfo
	err = json.Unmarshal(bytes, &mintInfo)
	if err != nil {
		return nil, err
	}

	return &mintInfo, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getTokensRegisteredByUser(userAddress string, blockchainProxy string) ([]string, error) {
	var resp ProxyRegisteredNFTsResponse

	url := fmt.Sprintf(RegisteredNFTsBaseFormat, blockchainProxy, userAddress)
	// err := cache.GetCacher().Get(url, &resp)
	// if err == nil {
	// 	return resp.Data.Tokens, nil
	// }

	err := HttpGet(url, &resp)
	if err != nil {
		return nil, err
	}

	err = cache.GetCacher().Set(url, resp, HttpResponseExpirePeriod)
	if err != nil {
		log.Debug("could not cache response", "err", err)
	}

	return resp.Data.Tokens, nil
}

func GetCollectionDetailBC(collectionName string, blockchainApi string) (CollectionDetailBCResponse, error) {
	var resp CollectionDetailBCResponse

	url := fmt.Sprintf(GetCollectionBaseFormat, blockchainApi, collectionName)
	// err := cache.GetCacher().Get(url, &resp)
	// if err == nil {
	// 	return resp.Data.Tokens, nil
	// }

	err := HttpGet(url, &resp)
	if err != nil {
		return resp, err
	}

	err = cache.GetCacher().Set(url, resp, HttpResponseExpirePeriod)
	if err != nil {
		log.Debug("could not cache response", "err", err)
	}

	return resp, nil
}

func checkValidInputOnCreate(request *CreateCollectionRequest) error {
	if len(request.TokenId) == 0 {
		return errors.New("empty token id")
	}

	if len(request.TokenId) > TokenIdMaxLen {
		return errors.New("token id too long")
	}

	if len(request.Name) == 0 {
		return errors.New("empty name")
	}

	if len(request.Name) > MaxNameLen {
		return errors.New("name too long")
	}

	if len(request.Description) > MaxDescLen {
		return errors.New("description too long")
	}

	if len(request.Website) > MaxLinkLen {
		return errors.New("website too long")
	}

	if len(request.DiscordLink) > MaxLinkLen {
		return errors.New("discord link too long")
	}

	if len(request.TwitterLink) > MaxLinkLen {
		return errors.New("twitter link too long")
	}

	if len(request.InstagramLink) > MaxLinkLen {
		return errors.New("instagram link too long")
	}

	if len(request.TelegramLink) > MaxLinkLen {
		return errors.New("telegram link too long")
	}

	if len(request.Flags) > MaxFlags {
		return errors.New("too many flags")
	}

	for _, flag := range request.Flags {
		if len(flag) > MaxFlagLen {
			return errors.New("flag too long")
		}
	}

	return nil
}

func checkValidInputOnUpdate(request *UpdateCollectionRequest) error {
	if len(request.Description) > MaxDescLen {
		return errors.New("description too long")
	}

	if len(request.Website) > MaxLinkLen {
		return errors.New("website too long")
	}

	if len(request.DiscordLink) > MaxLinkLen {
		return errors.New("discord link too long")
	}

	if len(request.TwitterLink) > MaxLinkLen {
		return errors.New("twitter link too long")
	}

	if len(request.InstagramLink) > MaxLinkLen {
		return errors.New("instagram link too long")
	}

	if len(request.TelegramLink) > MaxLinkLen {
		return errors.New("telegram link too long")
	}

	if len(request.Flags) > MaxFlags {
		return errors.New("too many flags")
	}

	for _, flag := range request.Flags {
		if len(flag) > MaxFlagLen {
			return errors.New("flag too long")
		}
	}

	return nil
}

func standardizeName(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func CheckValidFlags(flags []string) error {
	for _, flag := range flags {
		if !hasOnlyLettersAndWhitespaces(flag) {
			return errors.New("invalid flag")
		}
	}

	return nil
}

func hasOnlyLettersAndWhitespaces(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r != ' ') {
			return false
		}
	}

	return true
}
