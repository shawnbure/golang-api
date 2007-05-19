package cdn

import (
	"bytes"
	"encoding/base64"
	"io"
	"strings"
)

const (
	Base64Separator = ","
)

func Base64ToReader(base64Img string) (io.Reader, error) {
	imgContent := stripB64Str(base64Img)

	decoded, err := base64.StdEncoding.DecodeString(imgContent)
	if err != nil {
		return nil, err
	}

	buffReader := bytes.NewReader(decoded)
	return buffReader, nil
}

func Base64ToBytes(base64Img string) ([]byte, error) {
	imgContent := stripB64Str(base64Img)

	decoded, err := base64.StdEncoding.DecodeString(imgContent)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func stripB64Str(s string) string {
	suffixIdx := strings.Index(s, Base64Separator)

	return s[suffixIdx+1:]
}
