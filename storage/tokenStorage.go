package storage

import (
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/erdsea/erdsea-api/data/entities"
)

func AddToken(token *entities.Token) error {
	database, err := GetDBOrError()
	if err != nil {
		return err
	}

	txCreate := database.Create(&token)
	if txCreate.Error != nil {
		return txCreate.Error
	}
	if txCreate.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
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

	txRead := database.Offset(offset).Limit(limit).Find(&tokens, "listed = true AND collection_id = ?", collectionId)
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

	txRead := database.Model(&entities.Token{}).Where("listed = true AND collection_id = ?", collectionId).Count(&count)
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

	txRead := database.Model(&entities.Token{}).Where("listed = true AND collection_id = ?", collectionId).Distinct("owner_id").Count(&count)
	if txRead.Error != nil {
		return 0, txRead.Error
	}

	return uint64(count), nil
}
