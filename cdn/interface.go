package cdn

import "context"

type ImageUploader interface {
	UploadBase64(ctx context.Context, b64Img, imgID string) (string, error)
	GetImage(fileName string) ([]byte, string, error)
}
