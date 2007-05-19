package services

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_SetAccountCoverImageLarge(t *testing.T) {
	connectToDb()

	bytes := make([]byte, 1024*1024)
	_, err := rand.Read(bytes)
	require.Nil(t, err)

	imgBase64 := base64.RawStdEncoding.EncodeToString(bytes)
	err = SetAccountCoverImage("my_addr3", &imgBase64)
	require.Nil(t, err)
}

func Test_SetAccountCoverImageTooLarge(t *testing.T) {
	connectToDb()

	bytes := make([]byte, 1024*1024+1)
	_, err := rand.Read(bytes)
	require.Nil(t, err)

	imgBase64 := base64.RawStdEncoding.EncodeToString(bytes)
	err = SetAccountCoverImage("my_addr2", &imgBase64)
	require.NotNil(t, err)
}
