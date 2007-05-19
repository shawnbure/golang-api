package storage

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/erdsea/erdsea-api/data/images"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AddNewAccountProfileImageByUserId(T *testing.T) {
	connectToTestDb()

	bytes := make([]byte, 1024 * 1024)
	_, err := rand.Read(bytes)
	require.Nil(T, err)

	image := images.AccountProfileImage{
		ImageBase64: base64.StdEncoding.EncodeToString(bytes),
		AccountID:   1,
	}

	err = AddNewAccountProfileImageByUserId(&image)
	require.Nil(T, err)
}
