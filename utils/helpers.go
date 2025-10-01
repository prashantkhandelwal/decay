package utils

import (
	"fmt"
	"image"
	"mime/multipart"

	"github.com/prashantkhandelwal/decay/config"
)

func ValidateFile(fconfig *config.FileSettings, file *multipart.FileHeader) bool {
	isValidType := false
	fmt.Println("Checking MIME type...")
	mimeType := file.Header.Get("Content-Type")
	for _, t := range fconfig.MimeTypes {
		if t == mimeType {
			isValidType = true
			break
		}
	}
	isValidSize := file.Size <= fconfig.MaxSize
	return isValidType && isValidSize
}

// GetImageDimensions returns width and height of an image using image.DecodeConfig.
func GetImageDimensions(fileHeader *multipart.FileHeader) (int, int, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return cfg.Width, cfg.Height, nil
}
