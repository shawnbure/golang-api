package cdn

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ENFT-DAO/youbei-api/config"
)

var (
	UnknownImgTypeErr = errors.New("unknown image type")
)

type localUploader struct {
	baseUrl string
	rootDir string
}

func NewLocalUploader(cfg config.CDNConfig) *localUploader {
	return &localUploader{
		baseUrl: cfg.BaseUrl,
		rootDir: cfg.RootDir,
	}
}

func (lu *localUploader) UploadBase64(_ context.Context, b64Img, imgID string) (string, error) {
	imgBytes, err := Base64ToBytes(b64Img)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(lu.rootDir, imgID)
	err = ioutil.WriteFile(filePath, imgBytes, 0644)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", lu.baseUrl, imgID), nil
}

func (lu *localUploader) GetImage(fileName string) ([]byte, string, error) {
	filePath := filepath.Join(lu.rootDir, fileName)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer func(f *os.File) {
		closeErr := f.Close()
		if closeErr != nil {
			log.Warn("GetImage - failed to close file", "err", err.Error())
		}
	}(f)

	img, imgType, err := image.Decode(f)
	if err != nil {
		return nil, "", err
	}

	buf := new(bytes.Buffer)
	err = lu.encodeImageByType(imgType, img, buf)
	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), imgType, nil
}

func (lu *localUploader) encodeImageByType(imgType string, img image.Image, w io.Writer) error {
	switch imgType {
	case "jpeg", "jpg":
		return jpeg.Encode(w, img, nil)
	case "png":
		return png.Encode(w, img)
	case "gif":
		return gif.Encode(w, img, nil)
	default:
		return UnknownImgTypeErr
	}
}
