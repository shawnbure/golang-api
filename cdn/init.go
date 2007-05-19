package cdn

import (
	"errors"
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/erdsea/erdsea-api/config"
)

var (
	once        sync.Once
	imgUploader ImageUploader

	log = logger.GetOrCreate("cdn")
)

const (
	local     = "local"
	cloudyCDN = "cloudy"
)

func InitUploader(cfg config.CDNConfig) {
	once.Do(func() {
		upl, err := makeUploader(cfg)
		if err != nil {
			panic(err)
		}

		imgUploader = upl
	})
}

func GetImageUploaderOrErr() (ImageUploader, error) {
	if imgUploader == nil {
		return nil, errors.New("no uploader initialized")
	}

	return imgUploader, nil
}

func makeUploader(cfg config.CDNConfig) (ImageUploader, error) {
	switch cfg.Selector {
	case local:
		return NewLocalUploader(cfg), nil
	case cloudyCDN:
		return NewCloudyUploader(cfg)
	default:
		return nil, errors.New("unknown selector provided")
	}
}
