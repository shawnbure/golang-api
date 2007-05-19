package cdn

import (
	"context"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/erdsea/erdsea-api/config"
)

type cloudyUploader struct {
	cloudy *cloudinary.Cloudinary
}

func NewCloudyUploader(cfg config.CDNConfig) (*cloudyUploader, error) {
	newCloudy, err := cloudinary.NewFromParams(
		cfg.Name,
		cfg.ApiKey,
		cfg.ApiSecret,
	)
	if err != nil {
		return nil, err
	}

	return &cloudyUploader{
		cloudy: newCloudy,
	}, nil
}

func (cu *cloudyUploader) UploadBase64(ctx context.Context, b64Img, imgID string) (string, error) {
	buf, err := Base64ToReader(b64Img)
	if err != nil {
		return "", err
	}
	res, err := cu.cloudy.Upload.Upload(ctx, buf, uploader.UploadParams{
		PublicID:  imgID,
		Overwrite: true,
	})
	if err != nil {
		return "", err
	}

	return res.SecureURL, nil
}

func (cu *cloudyUploader) GetImage(fileName string) ([]byte, string, error) {
	return nil, "", nil
}
