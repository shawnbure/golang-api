package storage

import (
	"fmt"

	"gorm.io/datatypes"
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
			Where("token_id = ? AND nonce = ?", token.TokenID, token.Nonce).
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
			return err
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

func GetTokenByTokenIdAndNonce(tokenId string, nonce uint64) (*entities.Token, error) {
	var token entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Find(&token, "token_id = ? AND nonce = ?", tokenId, nonce)
	if txRead.Error != nil {
		return nil, txRead.Error
	}
	if txRead.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &token, nil
}

func GetTokensByOwnerIdWithOffsetLimit(ownerId uint64, offset int, limit int) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit).Find(&tokens, "owner_id = ?", ownerId)
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

	txRead := database.Offset(offset).Limit(limit).Find(&tokens, "owner_id = ? AND on_sale = 1", ownerId)
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

func GetTokensByCollectionIdWithOffsetLimit(
	collectionId uint64,
	offset int,
	limit int,
	attributesFilters map[string]string,
	sortRules map[string]string,
) ([]entities.Token, error) {
	var tokens []entities.Token

	database, err := GetDBOrError()
	if err != nil {
		return nil, err
	}

	txRead := database.Offset(offset).Limit(limit)
	for k, v := range attributesFilters {
		txRead.Where(datatypes.JSONQuery("attributes").Equals(v, k))
	}

	if len(sortRules) == 2 {
		query := fmt.Sprintf("%s %s", sortRules["criteria"], sortRules["mode"])
		txRead.Order(query)
	}

	txRead.Find(&tokens, "collection_id = ?", collectionId)
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
