package storage

import (
	"strconv"
	"testing"
	"time"

	"gorm.io/datatypes"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/stretchr/testify/require"
)

func Test_AddToken(t *testing.T) {
	connectToTestDb()

	token := defaultToken()
	err := AddToken(&token)
	require.Nil(t, err)

	var tokenRead entities.Token
	txRead := GetDB().Last(&tokenRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, tokenRead, token)
}

func Test_UpdateToken(t *testing.T) {
	connectToTestDb()

	token := defaultToken()
	err := AddToken(&token)
	require.Nil(t, err)

	token.TokenID = "new_token_id"
	err = UpdateToken(&token)

	var tokenRead entities.Token
	txRead := GetDB().Last(&tokenRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, tokenRead, token)
}

func Test_GetTokenById(t *testing.T) {
	connectToTestDb()

	token := defaultToken()
	err := AddToken(&token)
	require.Nil(t, err)

	tokenRead, err := GetTokenById(token.ID)
	require.Nil(t, err)
	require.Equal(t, tokenRead, &token)
}

func Test_GetTokenByTokenIdAndNonce(t *testing.T) {
	connectToTestDb()

	token := defaultToken()
	token.TokenID = "unique_token_id"
	token.Nonce = uint64(100)

	err := AddToken(&token)
	require.Nil(t, err)

	tokenRead, err := GetTokenByTokenIdAndNonce(token.TokenID, token.Nonce)
	require.Nil(t, err)
	require.Equal(t, tokenRead.TokenID, token.TokenID)
	require.Equal(t, tokenRead.Nonce, token.Nonce)
}

func Test_GetTokensForSaleOwnedBy(t *testing.T) {
	connectToTestDb()
	ownerId := uint64(1)

	token := defaultToken()
	err := AddToken(&token)
	require.Nil(t, err)

	otherToken := defaultToken()
	err = AddToken(&otherToken)
	require.Nil(t, err)

	tokensRead, err := GetTokensOnSaleByOwnerIdWithOffsetLimit(ownerId, 0, 100)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(tokensRead), 2)

	for _, tokenRead := range tokensRead {
		require.Equal(t, tokenRead.OwnerID, ownerId)
	}
}

func Test_GetTokensByCollectionId(t *testing.T) {
	connectToTestDb()
	collectionId := uint64(1)

	token := defaultToken()
	err := AddToken(&token)
	require.Nil(t, err)

	otherToken := defaultToken()
	err = AddToken(&otherToken)
	require.Nil(t, err)

	tokensRead, err := GetTokensByCollectionId(collectionId)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(tokensRead), 2)

	for _, tokenRead := range tokensRead {
		require.Equal(t, tokenRead.CollectionID, collectionId)
	}
}

func Test_CountListedTokensByCollectionId(t *testing.T) {
	connectToTestDb()

	token := defaultToken()
	err := AddToken(&token)
	require.Nil(t, err)

	otherToken := defaultToken()
	err = AddToken(&otherToken)
	require.Nil(t, err)

	count, err := CountListedTokensByCollectionId(1)
	require.Nil(t, err)
	require.GreaterOrEqual(t, count, uint64(2))
}

func Test_CountUniqueOwnersWithListedTokensByCollectionId(t *testing.T) {
	connectToTestDb()

	token := defaultToken()
	err := AddToken(&token)
	require.Nil(t, err)

	otherToken := defaultToken()
	err = AddToken(&otherToken)
	require.Nil(t, err)

	count, err := CountUniqueOwnersWithListedTokensByCollectionId(1)
	require.Nil(t, err)
	require.Equal(t, uint64(1), count)
}

func Test_GetTokensByCollectionIdWithOffsetLimit(t *testing.T) {
	connectToTestDb()

	coll := entities.Collection{
		Name: strconv.Itoa(int(time.Now().Unix())),
	}
	err := AddCollection(&coll)
	require.Nil(t, err)

	token1 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.ListToken,
		OwnerID:      1,
		Attributes:   datatypes.JSON(`{"hair": "red", "background": "dark"}`),
	}
	err = AddToken(&token1)
	require.Nil(t, err)

	token2 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.ListToken,
		OwnerID:      1,
		Attributes:   datatypes.JSON(`{"hair": "green", "background": "dark"}`),
	}
	err = AddToken(&token2)
	require.Nil(t, err)

	/*
		attrs := map[string]string{"background": "dark"}
		sort := map[string]string{}
		tokens, err := GetTokensByCollectionIdWithOffsetLimit(coll.ID, 0, 100, attrs, sort, true, nil)
		require.Nil(t, err)
		require.Equal(t, len(tokens), 2)

		attrs2 := map[string]string{"background": "dark", "hair": "green"}
		tokens2, err := GetTokensByCollectionIdWithOffsetLimit(coll.ID, 0, 100, attrs2, sort, true, nil)
		require.Nil(t, err)
		require.Equal(t, len(tokens2), 1)
	*/

}

func Test_GetTotalTokens(t *testing.T) {
	connectToTestDb()

	err := cleanTokenTable()
	require.Nil(t, err)

	err = insertSomeTokenRecords()
	require.Nil(t, err)

	t.Run("Get total tokens that is existed on out platform", func(t *testing.T) {
		total, err := GetTotalTokenCount()
		require.Nil(t, err)

		require.Equal(t, total, int64(3), "Total Token Count does not match")
	})
}

func Test_GetAllTokens(t *testing.T) {
	t.Run("Get tokens that is existed on the platform", func(t *testing.T) {
		lastTimeStamp := int64(0)
		currentPage := 1
		nextPage := 1
		filter := entities.QueryFilter{}
		sortOption := entities.SortOptions{}
		howMuchRows := 2

		tokens, err := GetAllTokens(lastTimeStamp, currentPage, nextPage, howMuchRows, &filter, &sortOption)
		require.Nil(t, err)

		require.Equal(t, len(tokens), int64(2), "Tokens list length is not correct")
	})

}

func defaultToken() entities.Token {
	return entities.Token{
		TokenID:      "my_token",
		Nonce:        10,
		PriceNominal: 1_000_000_000_000_000_000_000,
		Status:       entities.ListToken,
		MetadataLink: "link.com",
		OwnerID:      1,
		CollectionID: 1,
	}
}

func cleanTokenTable() error {
	_, err := DeleteAllTokens()
	if err != nil {
		return err
	}

	return nil
}

func insertSomeTokenRecords() error {
	token := defaultToken()
	token.TokenID = "my_token_1"
	token.PriceNominal = 2_000_000_000_000_000_000_000
	err := AddToken(&token)
	if err != nil {
		return err
	}

	token = defaultToken()
	token.TokenID = "my_token_2"
	token.PriceNominal = 2_300_000_000_000_000_000_000
	err = AddToken(&token)
	if err != nil {
		return err
	}

	token = defaultToken()
	token.TokenID = "my_token_3"
	token.PriceNominal = 1_453_000_000_000_000_000_000
	err = AddToken(&token)
	if err != nil {
		return err
	}

	return nil
}
