package storage

import (
	"github.com/erdsea/erdsea-api/data"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AddNewCollection(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	err := AddNewCollection(&collection)
	require.Nil(t, err)

	var collectionRead data.Collection
	txRead := GetDB().Last(&collectionRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, collectionRead, collection)
}

func Test_GetCollectionById(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	err := AddNewCollection(&collection)
	require.Nil(t, err)

	collectionRead, err := GetCollectionById(collection.ID)
	require.Nil(t, err)
	require.Equal(t, collectionRead, &collection)
}

func Test_GetCollectionsCreatedById(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	err := AddNewCollection(&collection)
	require.Nil(t, err)

	otherCollection := defaultCollection()
	err = AddNewCollection(&otherCollection)
	require.Nil(t, err)

	collectionsRead, err := GetCollectionsCreatedBy(collection.CreatorID)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(collectionsRead), 2)

	for _, collectionRead := range collectionsRead {
		require.Equal(t, collectionRead.CreatorID, collection.CreatorID)
	}
}

func Test_GetCollectionByName(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	collection.Name = "insane_unique_name"
	err := AddNewCollection(&collection)
	require.Nil(t, err)

	retrievedCollection, err := GetCollectionByName(collection.Name)
	require.Nil(t, err)
	require.Equal(t, retrievedCollection.Name, collection.Name)
}

func Test_GetCollectionsWithNameAlikeWithLimit(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	_ = AddNewCollection(&collection)
	collection.ID = 0
	_ = AddNewCollection(&collection)

	retrievedCollection, err := GetCollectionsWithNameAlikeWithLimit("default", 5)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(retrievedCollection), 2)
	require.Equal(t, retrievedCollection[0].Name, "default")
	require.Equal(t, retrievedCollection[1].Name, "default")
}

func defaultCollection() data.Collection {
	return data.Collection{
		Name:      "default",
		TokenID:   "my_token",
		CreatorID: 0,
	}
}
