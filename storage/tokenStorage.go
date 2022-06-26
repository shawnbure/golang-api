package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"math"

	"gorm.io/gorm"

	"github.com/ENFT-DAO/youbei-api/data/entities"
)

func AddToken(token *entities.Token) error {

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	//verify the collection exists retun error if not

	collectionCount := int64(0)
	err = db.Model(&entities.Collection{}).
		Where("id = ?", token.CollectionID).
		Count(&collectionCount).
		Error

	if collectionCount > 0 {

		//if the token does not exixts in the db create it return error
		tokenCount := int64(0)
		err = db.Model(token).
			Where("token_id = ? AND nonce_str = ?", token.TokenID, token.NonceStr).
			Count(&tokenCount).
			Error

		if tokenCount == 0 {
			txCreate := database.Create(&token)
			if txCreate.Error != nil {
				return txCreate.Error
			}
			if txCreate.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		} else {
			return gorm.ErrRegistered
		}
	} else {
		return err
	}

	return nil
}

func AddOrUpdateToken(token *entities.Token) error {

	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	//verify the collection exists retun error if not

	collectionCount := int64(0)
	err = db.Model(&entities.Collection{}).
		Where("id = ?", token.CollectionID).
		Count(&collectionCount).
		Error

	if collectionCount > 0 {

		//if the token does not exixts in the db create it return error
		tokenCount := int64(0)
		err = db.Model(token).
			Where("token_id = ? AND nonce_str = ?", token.TokenID, token.NonceStr).
			Count(&tokenCount).
			Error

		if tokenCount == 0 {
			txCreate := database.Create(&token)
			if txCreate.Error != nil {
				return txCreate.Error
			}
			if txCreate.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		} else {
			txCreate := database.Model(token).Where("token_id = ? AND nonce_str = ?", token.TokenID, token.NonceStr).Updates(token)
			if txCreate.Error != nil {
				return txCreate.Error
			}
			if txCreate.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}
	} else {
		return err
	}

	return nil
}

func UpdateToken(token *entities.Token) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Save(&token)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func UpdateTokenWhere(token *entities.Token, toUpdate map[string]interface{}, where string, args ...interface{}) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Model(token).Where(where, args...).Updates(toUpdate)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func GetTokenById(id uint64) (*entities.Token, error) {
	var token entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&token, id)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &token, nil
}

func GetTokenByTokenIdAndNonceStr(tokenId string, nonce string) (*entities.Token, error) {
	var token entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&token, "token_id = ? AND nonce_str = ?", tokenId, nonce)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &token, nil
}

func GetTokenByTokenIdAndNonce(tokenId string, nonce uint64) (*entities.Token, error) {
	var token entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Preload("Owner").Find(&token, "token_id = ? AND nonce = ?", tokenId, nonce)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &token, nil
}

func GetTokensByOwnerIdWithOffsetLimit(ownerId uint64, filter entities.QueryFilter, offset int, limit int) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.
		Offset(offset).
		Limit(limit).
		Where(filter.Query, filter.Values...).
		Find(&tokens, "owner_id = ?", ownerId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func GetTokensOnSaleByOwnerIdWithOffsetLimit(ownerId uint64, offset int, limit int) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&tokens, "owner_id = ? AND on_sale = true", ownerId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func GetTokensByCollectionId(collectionId uint64) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&tokens, "collection_id = ?", collectionId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}
func GetTokensWithNoRankCount(collectionID uint64) (int64, error) {

	database, err := GetDBOrError()
	if err != nil {
		return 0, err
	}
	var count int64
	txRead := database.
		Model(&entities.Token{}).
		Where("rank = ?", 0).
		Count(&count)
	if txRead.Error != nil {
		return 0, txRead.Error
	}

	return count, nil
}
func GetTokensByCollectionIdOrderedByRarityScore(collectionId uint64, direction string) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.
		Where("rarity_score != ?", 0).
		Order(fmt.Sprintf("rarity_score %s", direction)).
		Find(&tokens, "collection_id = ?", collectionId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func GetTokensByCollectionIdNotRanked(collectionId uint64) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.
		Where("collection_id = ?", collectionId).
		Find(&tokens)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func GetLastNonceTokenByCollectionId(collectionId uint64) (entities.Token, error) {
	var token entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return token, err
	}

	txRead := database.
		Order("nonce DESC").
		First(&token, "collection_id = ?", collectionId)
	if txRead.Error != nil {
		return token, txRead.Error
	}

	return token, nil
}

func GetTokensByCollectionIdWithOffsetLimit(
	collectionId uint64,
	offset int,
	limit int,
	attributesFilters map[string]string,
	sortRules map[string]string,
	onSaleFlag bool,
	onStakeFlag bool,
	sqlFilter entities.QueryFilter,
) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit)
	for k, v := range attributesFilters {

		txRead.Where(fmt.Sprintf(`attributes @> '[{"trait_type":"%s","value":"%s"}]'`, k, v))
		// txRead.Where(datatypes.JSONQuery("attributes").Equals(v, k))
	}

	if len(sortRules) == 2 {

		query := fmt.Sprintf("%s %s", sortRules["criteria"], sortRules["mode"])
		txRead.Order(query)
	}

	/*
		fmt.Println("sqlFilter.Query: ")
		fmt.Println(sqlFilter.Query)
		fmt.Println(sqlFilter.Values)

		if sqlFilter.Query != "" {
			txRead.Where(sqlFilter.Query, sqlFilter.Values...)
			fmt.Printf("txRead.Where: %v\n", txRead.Where.stri)
		}

		txRead.Preload("Owner").Find(&tokens, "collection_id = ? and on_sale = ? and on_stake = ?", collectionId, onSaleFlag, onStakeFlag)
	*/

	if onSaleFlag {
		txRead.Preload("Owner").Where("on_sale = True and (on_stake = False or on_stake is null)").Find(&tokens, "collection_id = ?", collectionId)
	}
	if onStakeFlag {
		txRead.Preload("Owner").Where("(on_sale = False or on_sale is null) and on_stake = True").Find(&tokens, "collection_id = ?", collectionId)
	}
	if !onStakeFlag && !onSaleFlag {
		txRead.Preload("Owner").Where("(on_sale = False or on_sale is null) and (on_stake = False or on_stake is null)").Find(&tokens, "collection_id = ?", collectionId)
	}

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func GetListedTokensByCollectionIdWithOffsetLimit(collectionId uint64, offset int, limit int) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&tokens, "collection_id = ?", collectionId)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func CountListedTokensByCollectionId(collectionId uint64) (uint64, error) {
	count := int64(0)

	database, err := GetDBOrError()
	if err != nil {
		return 0, err
	}

	txRead := database.Model(&entities.Token{}).Where("(status = 'List' OR status = 'Auction') AND collection_id = ?", collectionId)
	txRead.Count(&count)
	if txRead.Error != nil {
		return 0, txRead.Error
	}

	return uint64(count), nil
}

func CountUniqueOwnersWithListedTokensByCollectionId(collectionId uint64) (uint64, error) {
	count := int64(0)

	database, err := GetDBOrError()
	if err != nil {
		return 0, err
	}

	txRead := database.Model(&entities.Token{}).Where("(status = 'List' OR status = 'Auction') AND collection_id = ?", collectionId)
	txRead.Distinct("owner_id").Count(&count)
	if txRead.Error != nil {
		return 0, txRead.Error
	}

	return uint64(count), nil
}

func GetTokensWithOffsetLimit(
	offset int,
	limit int,
	sortRules map[string]string,
) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit)
	if len(sortRules) == 2 {
		query := fmt.Sprintf("%s %s", sortRules["criteria"], sortRules["mode"])
		txRead.Order(query)
	}

	txRead.Find(&tokens)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func GetTokensWithTokenIdAlikeWithLimit(tokenId string, limit int) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Limit(limit).Where("token_id ILIKE ?", tokenId).Find(&tokens)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func GetTokensListedWithTokenIdAlikeWithLimit(tokenId string, limit int) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Limit(limit).Where("token_id ILIKE ?", tokenId).Where("status is not NULL and status != '' and status != 'None'").Find(&tokens)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func GetTokensUnlistedWithTokenIdAlikeWithLimit(tokenId string, limit int) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Limit(limit).Where("token_id ILIKE ?", tokenId).Where("status is NULL or status = '' or status = 'None'").Find(&tokens)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func GetEndAuctionTokens() ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Where("auction_deadline <= extract(epoch from now()) ").Where("Status = ?", "Auction").Find(&tokens)
	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil

}

func DeleteAllTokens() (int64, error) {
	database, err := GetDBOrError()
	if err != nil {
		return int64(0), err
	}

	tx := database.Where("1 = 1").Delete(&entities.Token{})
	if tx.Error != nil {
		return int64(0), err
	}

	return tx.RowsAffected, nil
}

func GetTotalTokenCount() (int64, error) {
	var total int64

	database, err := GetDBOrError()
	if err != nil {
		return int64(0), err
	}

	txRead := database.
		Table("tokens").
		Count(&total)

	if txRead.Error != nil {
		return int64(0), txRead.Error
	}

	return total, nil
}

func GetAllTokens(lastTimestamp int64, currentPage, requestedPage, pageSize int, filter *entities.QueryFilter, sortOptions *entities.SortOptions, collectionFilter *entities.QueryFilter, attributes [][]string) ([]entities.TokenExplorer, error) {
	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	tokens := []entities.TokenExplorer{}

	if sortOptions.Query == "" {
		sortOptions.Query = "tokens.last_market_timestamp %s"
		sortOptions.Values = []interface{}{"desc"}
	}
	query := ""
	offset := 0

	if lastTimestamp == 0 {
		query = filter.Query
	} else {
		query = "tokens.last_market_timestamp<?"
		if requestedPage < currentPage {
			query = "tokens.last_market_timestamp>?"

			if sortOptions.Values[0] == "asc" {
				sortOptions.Values[0] = "desc"
			} else {
				sortOptions.Values[0] = "asc"
			}
		}

		if requestedPage != currentPage {
			offset = (int(math.Abs(float64(requestedPage-currentPage))) - 1) * pageSize
		}

		if filter.Query != "" {
			query = fmt.Sprintf("(%s) and %s", filter.Query, query)
		}
		filter.Values = append(filter.Values, lastTimestamp)
	}

	//if isVerified {
	//	query += " and collections.is_verified=? "
	//	filter.Values = append(filter.Values, true)
	//}

	order := fmt.Sprintf(sortOptions.Query, sortOptions.Values...)

	//	selectQuery := `
	//tokens.token_id as token_id,
	//tokens.status as token_status,
	//tokens.token_name as token_name,
	//tokens.image_link as token_image,
	//tokens.auction_start_time as token_auction_start_time,
	//tokens.auction_deadline as token_auction_deadline,
	//tokens.created_at as token_created_at,
	//tokens.last_market_timestamp as token_last_market_timestamp,
	//tokens.last_buy_price_nominal as token_last_buy_price_nominal,
	//tokens.price_nominal as token_price_nominal,
	//accounts.address as owner_address,
	//accounts.name as owner_name,
	//collections.name as collection_name,
	//collections.token_id as collection_token_id,
	//collections.name as collection_name
	//`

	txRead := database.Table("tokens").
		Preload("Owner").
		Joins("inner join collections on collections.id=tokens.collection_id ").
		Preload("Collection").
		Order(order).
		Where(query, filter.Values...).
		Where(collectionFilter.Query, collectionFilter.Values...)

	//txRead := database.Table("tokens").Select(selectQuery).
	//	Joins("inner join collections on collections.id=tokens.collection_id ").
	//	Joins("inner join accounts on accounts.id=tokens.owner_id ").
	//	Order(order).
	//	Where(query, filter.Values...)

	for _, item := range attributes {
		txRead.Where(fmt.Sprintf(`attributes @> '[{"trait_type":"%s","value":"%s"}]'`, item[0], item[1]))
	}

	txRead.
		Offset(offset).
		Limit(pageSize).
		Find(&tokens)

	if txRead.Error != nil {
		return nil, txRead.Error
	}

	return tokens, nil
}

func GetTokensCountWithCriteria(filter *entities.QueryFilter, collectionFilter *entities.QueryFilter, attributes [][]string) (int64, error) {
	database, err := GetDBOrError()
	if err != nil {
		return 0, err
	}

	var total int64

	//if isVerified {
	//	filter.Query += " and collections.is_verified=? "
	//	filter.Values = append(filter.Values, true)
	//}

	var txRead *gorm.DB
	//if isVerified {
	txRead = database.Table("tokens").
		Joins("inner join collections on tokens.collection_id=collections.id").
		Where(filter.Query, filter.Values...).
		Where(collectionFilter.Query, collectionFilter.Values...)

	for _, item := range attributes {
		txRead.Where(fmt.Sprintf(`attributes @> '[{"trait_type":"%s","value":"%s"}]'`, item[0], item[1]))
	}

	txRead.Count(&total)
	//} else {
	//	txRead = database.Table("tokens").
	//		Where(filter.Query, filter.Values...)
	//
	//	for _, item := range attributes {
	//		txRead.Where(fmt.Sprintf(`attributes @> '[{"trait_type":"%s","value":"%s"}]'`, item[0], item[1]))
	//	}
	//
	//	txRead.Count(&total)
	//}

	if txRead.Error != nil {
		return 0, txRead.Error
	}

	return total, nil
}

func GetTokensPriceBoundary(filter *entities.QueryFilter, collectionFilter *entities.QueryFilter, attributes [][]string) (float64, float64, error) {
	database, err := GetDBOrError()
	if err != nil {
		return 0, 0, err
	}

	type tempStruct struct {
		Min sql.NullFloat64 `json:"min"`
		Max sql.NullFloat64 `json:"max"`
	}

	p := tempStruct{}

	txRead := database.Table("tokens").
		Joins("inner join collections on tokens.collection_id=collections.id").
		Select("min(tokens.price_nominal) as min, max(tokens.price_nominal) as max").
		Where(filter.Query, filter.Values...).
		Where(collectionFilter.Query, collectionFilter.Values...)

	for _, item := range attributes {
		txRead.Where(fmt.Sprintf(`attributes @> '[{"trait_type":"%s","value":"%s"}]'`, item[0], item[1]))
	}

	txRead.Scan(&p)

	if txRead.Error != nil {
		return 0, 0, txRead.Error
	}

	if p.Min.Valid && p.Max.Valid {
		return p.Min.Float64, p.Max.Float64, nil
	}

	return 0, 0, errors.New("Cannot get values from database")
}

//func GetVerifiedTokensPriceBoundary(filter *entities.QueryFilter, collectionFilter *entities.QueryFilter, attributes [][]string) (float64, float64, error) {
//	database, err := GetDBOrError()
//	if err != nil {
//		return 0, 0, err
//	}
//
//	type tempStruct struct {
//		Min sql.NullFloat64 `json:"min"`
//		Max sql.NullFloat64 `json:"max"`
//	}
//	p := tempStruct{}
//
//	txRead := database.Table("tokens").
//		Select("min(tokens.price_nominal) as min, max(tokens.price_nominal) as max").
//		Joins("inner join collections on collections.id=tokens.collection_id").
//		Where("collections.is_verified=?", true).
//		Where(filter.Query, filter.Values...)
//
//	for _, item := range attributes {
//		txRead.Where(fmt.Sprintf(`attributes @> '[{"trait_type":"%s","value":"%s"}]'`, item[0], item[1]))
//	}
//
//	txRead.Scan(&p)
//
//	if txRead.Error != nil {
//		return 0, 0, txRead.Error
//	}
//
//	if p.Min.Valid && p.Max.Valid {
//		return p.Min.Float64, p.Max.Float64, nil
//	}
//
//	return 0, 0, errors.New("Cannot get values from database")
//}

func GetOldTokenWithZeroRarity() (entities.Token, error) {
	var tokenInstance entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return tokenInstance, err
	}

	txRead := database.
		Where("is_rarity_inserted=?", false).
		Order("rarity_last_updated ASC").
		First(&tokenInstance)
	if txRead.Error != nil {
		return tokenInstance, txRead.Error
	}

	return tokenInstance, nil

}
