package cdn

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/ENFT-DAO/youbei-api/config"
	"github.com/stretchr/testify/require"
)

var testUploader = NewLocalUploader(config.CDNConfig{
	BaseUrl: "something://",
	RootDir: "",
})

func TestLocalUploader_UploadBase64(t *testing.T) {
	t.Parallel()

	url, err := testUploader.UploadBase64(nil, imgStr, "testLocal")
	require.Nil(t, err)
	require.True(t, strings.Contains(url, "testLocal"))
}

func TestLocalUploader_GetImage(t *testing.T) {
	t.Parallel()

	imgBytes, imgType, err := testUploader.GetImage("testLocal")
	require.Nil(t, err)
	require.True(t, imgType == "png")
	require.Equal(t, base64.StdEncoding.EncodeToString(imgBytes), stripB64Str(imgStr))
}
