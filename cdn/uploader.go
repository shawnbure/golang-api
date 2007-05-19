package cdn

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"strings"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

const (
	base64Separator = ","
)

func UploadToCloudy(ctx context.Context, base64Img, imgID string) (*uploader.UploadResult, error) {
	buf, err := base64ToReader(base64Img)
	if err != nil {
		return nil, err
	}

	res, err := cloudy.Upload.Upload(ctx, buf, uploader.UploadParams{
		PublicID: imgID,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func base64ToReader(base64Img string) (io.Reader, error) {
	suffixIdx := strings.Index(base64Img, base64Separator)

	imgContent := base64Img[suffixIdx+1:]

	decoded, err := base64.StdEncoding.DecodeString(imgContent)
	if err != nil {
		return nil, err
	}

	buffReader := bytes.NewReader(decoded)
	return buffReader, nil
}
