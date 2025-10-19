package utils

import (
	"fmt"
	"image"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

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

func ParseExpiration(expiration string) (int64, error) {
	var duration int64
	var err error
	if strings.HasSuffix(expiration, "s") {
		duration, err = strconv.ParseInt(strings.TrimSuffix(expiration, "s"), 10, 64)
	} else if strings.HasSuffix(expiration, "m") {
		duration, err = strconv.ParseInt(strings.TrimSuffix(expiration, "m"), 10, 64)
		duration *= 60
	} else if strings.HasSuffix(expiration, "h") {
		duration, err = strconv.ParseInt(strings.TrimSuffix(expiration, "h"), 10, 64)
		duration *= 3600
	} else {
		err = fmt.Errorf("invalid expiration format")
	}

	current_time := time.Now().Unix()
	duration = current_time + duration

	if duration < current_time {
		return 0, fmt.Errorf("expiration time is in the past")
	}

	return duration, err
}
