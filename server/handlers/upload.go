package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/prashantkhandelwal/decay/config"
)

func UploadHandler(fconfig *config.FileSettings) gin.HandlerFunc {
	fn := func(g *gin.Context) {
		// Parse the multipart form
		file, err := g.FormFile("file")
		if err != nil {
			g.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		fmt.Printf("Received file: %s\n", file.Filename)
		fmt.Printf("File size: %d bytes\n", file.Size)
		fmt.Printf("MIME header: %v\n", file.Header)

		// Validate file type and size
		fmt.Println("Validating file...")
		if !validateFile(fconfig, file) {
			g.String(http.StatusBadRequest, "Invalid file type or size exceeds limit")
			return
		}

		// Define the destination path for the uploaded file
		dst := filepath.Join(fconfig.UploadDir, file.Filename)

		// Save the uploaded file
		if err := g.SaveUploadedFile(file, dst); err != nil {
			g.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		g.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully to %s", file.Filename, dst))
	}
	return gin.HandlerFunc(fn)
}

func validateFile(fconfig *config.FileSettings, file *multipart.FileHeader) bool {
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
