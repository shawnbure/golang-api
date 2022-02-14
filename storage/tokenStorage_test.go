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

	tokensRead, err := GetTokensForSaleByOwnerIdWithOffsetLimit(ownerId, 0, 100)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(tokensRead), 2)

	for _, tokenRead := range tokensRead {
		require.Equal(t, tokenRead.OwnerId, ownerId)
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
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "red", "background": "dark"}`),
	}
	err = AddToken(&token1)
	require.Nil(t, err)

	token2 := entities.Token{
		CollectionID: coll.ID,
		Status:       entities.List,
		OwnerId:      1,
		Attributes:   datatypes.JSON(`{"hair": "green", "background": "dark"}`),
	}
	err = AddToken(&token2)
	require.Nil(t, err)

	attrs := map[string]string{"background": "dark"}
	sort := map[string]string{}
	tokens, err := GetTokensByCollectionIdWithOffsetLimit(coll.ID, 0, 100, attrs, sort)
	require.Nil(t, err)
	require.Equal(t, len(tokens), 2)

	attrs2 := map[string]string{"background": "dark", "hair": "green"}
	tokens2, err := GetTokensByCollectionIdWithOffsetLimit(coll.ID, 0, 100, attrs2, sort)
	require.Nil(t, err)
	require.Equal(t, len(tokens2), 1)
}

func defaultToken() entities.Token {
	return entities.Token{
		TokenID:      "my_token",
		Nonce:        10,
		PriceNominal: 1_000_000_000_000_000_000_000,
		Status:       entities.List,
		MetadataLink: "link.com",
		OwnerId:      1,
		CollectionID: 1,
	}
}
