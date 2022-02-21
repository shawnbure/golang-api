package cdn

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/ENFT-DAO/youbei-api/config"
)

type cloudyUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

func NewCloudyUploader(cfg config.CDNConfig) (*cloudyUploader, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	// newCloudy, err := cloudinary.NewFromParams(
	// 	cfg.Name,
	// 	cfg.ApiKey,
	// 	cfg.ApiSecret,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	return &cloudyUploader{
		cl:         client,
		bucketName: cfg.BucketName,
		projectID:  cfg.ProjectID,
		uploadPath: cfg.UploadPath,
		// cloudy: uploader,
	}, nil
}

func (cu *cloudyUploader) UploadBase64(ctx context.Context, b64Img, imgID string) (string, error) {
	buf, err := Base64ToReader(b64Img)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	// Upload an object with storage.Writer.
	wc := cu.cl.Bucket(cu.bucketName).Object(cu.uploadPath + imgID).NewWriter(ctx)
	if _, err := io.Copy(wc, buf); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	// res, err := cu.cloudy.Upload.Upload(ctx, buf, uploader.UploadParams{
	// 	PublicID:  imgID,
	// 	Overwrite: true,
	// })
	// if err != nil {
	// 	return "", err
	// }
	// if res.Error.Message != "" {
	// 	return "", fmt.Errorf("%s", res.Error.Message)
	// }

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s%s", cu.bucketName, cu.uploadPath, imgID), nil
}

func (cu *cloudyUploader) GetImage(_ string) ([]byte, string, error) {
	return nil, "", nil
}
