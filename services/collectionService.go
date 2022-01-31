package services

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"

	"github.com/ENFT-DAO/youbei-api/cache"
	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/ENFT-DAO/youbei-api/interaction"
	"github.com/ENFT-DAO/youbei-api/stats/collstats"
	"github.com/ENFT-DAO/youbei-api/storage"
	"github.com/boltdb/bolt"
)

const (
	TokenIdMaxLen                  = 15
	MaxFlags                       = 10
	MaxFlagLen                     = 25
	MaxNameLen                     = 20
	MaxLinkLen                     = 100
	MaxDescLen                     = 1000
	RegisteredNFTsBaseFormat       = "%s/address/%s/registered-nfts"
	HttpResponseExpirePeriod       = 10 * time.Minute
	CollectionSearchCacheKeyFormat = "CollectionSearch:%s"
	CollectionSearchExpirePeriod   = 20 * time.Minute
	MintInfoViewName               = "getMaxSupplyAndTotalSold"
	MintInfoSetNxKeyFormat         = "MintInfoNX:%s"
	MintInfoSetNxExpirePeriod      = 6 * time.Second
	MintInfoBucketName             = "MintInfo"
)

type MintInfo struct {
	MaxSupply uint64 `json:"maxSupply"`
	TotalSold uint64 `json:"totalSold"`
}

type CreateCollectionRequest struct {
	UserAddress   string   `json:"userAddress"`
	Name          string   `json:"collectionName"`
	TokenId       string   `json:"tokenId"`
	Description   string   `json:"description"`
	Website       string   `json:"website"`
	DiscordLink   string   `json:"discordLink"`
	TwitterLink   string   `json:"twitterLink"`
	InstagramLink string   `json:"instagramLink"`
	TelegramLink  string   `json:"telegramLink"`
	Flags         []string `json:"flags"`
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

type ProxyRegisteredNFTsResponse struct {
	Data struct {
		Tokens []string `json:"tokens"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
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

	collection := &entities.Collection{
		ID:            0,
		Name:          standardizedName,
		TokenID:       request.TokenId,
		Description:   request.Description,
		Website:       request.Website,
		DiscordLink:   request.DiscordLink,
		TwitterLink:   request.TwitterLink,
		InstagramLink: request.InstagramLink,
		TelegramLink:  request.TelegramLink,
		Flags:         datatypes.JSON(bytes),
		CreatorID:     account.ID,
		CreatedAt:     uint64(time.Now().Unix()),
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
