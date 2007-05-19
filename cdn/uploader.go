package cdn

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

type Uploader struct {
	cloudy *cloudinary.Cloudinary
}

func NewCdnUploader() *Uploader {
	cloudName := "deaezbrer"
	apiKey := "823855837497929"
	apiSecret := "9UXeyr23mESzGtVEX_1ZML54fXk"

	cloudy, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		panic(err)
	}

	return &Uploader{
		cloudy: cloudy,
	}
}

func (u *Uploader) Upload(img interface{}) {
	res, err := u.cloudy.Upload.Upload(context.Background(), img, uploader.UploadParams{
		PublicID: "penis",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func decodeBase64Img(b64 string) image.Image {
	idx := strings.Index(b64, ",")

	rawImage := b64[idx+1:]
	decoded, _ := base64.StdEncoding.DecodeString(string(rawImage))

	buf := bytes.NewReader(decoded)

	switch strings.TrimSuffix(b64[5:idx], ";base64") {
	case "image/png":
		pngI, err := png.Decode(buf)
		if err != nil {
			panic(err)
		}
		return pngI
	case "image/jpeg":
		jpgI, err := jpeg.Decode(buf)
		if err != nil {
			panic(err)
		}
		return jpgI
	default:
		return nil
	}
}

func base64ToReader(b64 string) io.Reader {
	idx := strings.Index(b64, ",")

	rawImage := b64[idx+1:]
	decoded, _ := base64.StdEncoding.DecodeString(string(rawImage))

	return bytes.NewReader(decoded)
}
