package storage

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/erdsea/erdsea-api/data/entities/images"
	"github.com/stretchr/testify/require"
)

func Test_AddAccountProfileImage(T *testing.T) {
	connectToTestDb()

	bytes := make([]byte, 1024*1024)
	_, err := rand.Read(bytes)
	require.Nil(T, err)

	image := images.AccountProfileImage{
		ImageBase64: base64.StdEncoding.EncodeToString(bytes),
		AccountID:   1,
	}

	err = AddOrUpdateAccountProfileImage(&image)
	require.Nil(T, err)
}

func Test_GetAccountProfileImage(T *testing.T) {
	connectToTestDb()

	bytes := make([]byte, 1024*1024)
	_, err := rand.Read(bytes)
	require.Nil(T, err)

	image := images.AccountProfileImage{
		ImageBase64: base64.StdEncoding.EncodeToString(bytes),
		AccountID:   5,
	}

	err = AddOrUpdateAccountProfileImage(&image)
	require.Nil(T, err)

	readImage, err := GetAccountProfileImageByAccountId(5)
	require.Nil(T, err)
	require.Equal(T, image.ImageBase64, readImage.ImageBase64)
}

func Test_GetCollectionProfileImage(T *testing.T) {
	connectToTestDb()

	bytes := make([]byte, 1024*1024)
	_, err := rand.Read(bytes)
	require.Nil(T, err)

	image := images.CollectionProfileImage{
		ImageBase64:  base64.StdEncoding.EncodeToString(bytes),
		CollectionID: 5,
	}

	err = AddOrUpdateCollectionProfileImage(&image)
	require.Nil(T, err)

	readImage, err := GetCollectionProfileImageByCollectionId(5)
	require.Nil(T, err)
	require.Equal(T, image.ImageBase64, readImage.ImageBase64)
}
