package storage

import (
	"strconv"
	"testing"
	"time"

	"github.com/ENFT-DAO/youbei-api/data/entities"
	"github.com/stretchr/testify/require"
)

func Test_AddCollection(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	err := AddCollection(&collection)
	require.Nil(t, err)

	var collectionRead entities.Collection
	txRead := GetDB().Last(&collectionRead)

	require.Nil(t, txRead.Error)
	require.Equal(t, collectionRead, collection)
}

func Test_GetCollectionById(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	err := AddCollection(&collection)
	require.Nil(t, err)

	collectionRead, err := GetCollectionById(collection.ID)
	require.Nil(t, err)
	require.Equal(t, collectionRead, &collection)
}

func Test_GetCollectionsCreatedById(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	err := AddCollection(&collection)
	require.Nil(t, err)

	otherCollection := defaultCollection()
	err = AddCollection(&otherCollection)
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
	err := AddCollection(&collection)
	require.Nil(t, err)

	retrievedCollection, err := GetCollectionByName(collection.Name)
	require.Nil(t, err)
	require.Equal(t, retrievedCollection.Name, collection.Name)
}

func Test_GetCollectionsWithNameAlikeWithLimit(t *testing.T) {
	connectToTestDb()

	collection := defaultCollection()
	_ = AddCollection(&collection)
	collection.ID = 0
	_ = AddCollection(&collection)

	retrievedCollection, err := GetCollectionsWithNameAlikeWithLimit("default", 5)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(retrievedCollection), 2)
	require.Equal(t, retrievedCollection[0].Name, "default")
	require.Equal(t, retrievedCollection[1].Name, "default")
}

func Test_GetCollectionsSorted(t *testing.T) {
	connectToTestDb()

	addCollectionsWithPriority()

	colls, err := GetCollectionsWithOffsetLimit(0, 2, []string{})
	require.Nil(t, err)
	require.True(t, colls[0].Priority > colls[1].Priority)
}

func Test_CollectionWithNameILike(t *testing.T) {
	connectToTestDb()

	timeNowUnix := strconv.FormatInt(time.Now().Unix(), 10)
	name := "name" + timeNowUnix
	collection := entities.Collection{
		Name: name,
	}
	err := AddCollection(&collection)
	require.Nil(t, err)

	_, err = GetCollectionWithNameILike("Name" + timeNowUnix)
	require.Nil(t, err)

	_, err = GetCollectionByName("Name" + timeNowUnix)
	require.NotNil(t, err)
}

func defaultCollection() entities.Collection {
	return entities.Collection{
		Name:      "default",
		TokenID:   "my_token",
		CreatorID: 0,
	}
}

func addCollectionsWithPriority() {
	collections := collectionsWithPriority()

	_ = AddCollection(&collections[0])
	_ = AddCollection(&collections[1])
}

func collectionsWithPriority() []entities.Collection {
	return []entities.Collection{
		{
			Name:     "first_coll",
			TokenID:  "first_token_id",
			Priority: 100,
		},
		{
			Name:     "second_coll",
			TokenID:  "second_token_id",
			Priority: 50,
		},
	}
}
