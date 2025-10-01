package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prashantkhandelwal/decay/config"
	"github.com/prashantkhandelwal/decay/db"
	"github.com/prashantkhandelwal/decay/utils"
)

func UploadHandler(fconfig *config.FileSettings) gin.HandlerFunc {
	fn := func(g *gin.Context) {
		// Parse the multipart form
		file, err := g.FormFile("file")

		title := g.PostForm("title")
		var f db.File
		f.Title = title

		if err != nil {
			g.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		f.Size = file.Size
		f.Mime = file.Header.Get("Content-Type")
		if strings.HasPrefix(f.Mime, "image/") {
			f.Width, f.Height, _ = utils.GetImageDimensions(file)
		}

		// Validate file type and size
		fmt.Println("Validating file...")
		if !utils.ValidateFile(fconfig, file) {
			g.String(http.StatusBadRequest, "Invalid file type or size exceeds limit")
			return
		}

		// Define the destination path for the uploaded file
		dst := filepath.Join(fconfig.UploadDir, file.Filename)

		// Save the uploaded file
		if err := g.SaveUploadedFile(file, dst); err != nil {
			g.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		} else {
			id, err := db.InsertFile(f, file)
			if err != nil {
				g.String(http.StatusInternalServerError, fmt.Sprintf("insert file err: %s", err.Error()))
				return
			}
			f.ID = id
			//fmt.Printf("File %s uploaded successfully to %s\n", file.Filename, dst)
		}

		g.JSON(http.StatusOK, gin.H{
			"id":       f.ID,
			"filename": file.Filename,
			"title":    title,
		})
		//g.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully to %s", file.Filename, dst))
	}
	return gin.HandlerFunc(fn)
}
