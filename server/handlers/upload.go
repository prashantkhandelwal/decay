package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prashantkhandelwal/decay/config"
	"github.com/prashantkhandelwal/decay/db"
	"github.com/prashantkhandelwal/decay/server/middleware"
	"github.com/prashantkhandelwal/decay/utils"
)

func UploadHandler(fconfig *config.FileSettings) gin.HandlerFunc {
	fn := func(g *gin.Context) {
		// Parse the multipart form
		file, err := g.FormFile("file")

		// Total requests counter
		middleware.TotalFileUploadRequests.WithLabelValues("1").Inc()

		if err != nil {
			g.String(http.StatusBadRequest, fmt.Sprintf("Get form err: %s", err.Error()))
			middleware.HttpBadRequests.WithLabelValues(g.Request.Method, g.FullPath(), strconv.Itoa(g.Writer.Status())).Inc()
			return
		}

		title := g.PostForm("title")
		expiration := g.PostForm("expiration")

		var f db.File
		f.Title = title

		if strings.TrimSpace(expiration) != "" {
			expireDuration, err := utils.ParseExpiration(expiration)
			if err != nil {
				log.Println("Invalid expiration format:", err)
				g.String(http.StatusBadRequest, fmt.Sprintf("Invalid expiration format: %s", err.Error()))
				middleware.HttpBadRequests.WithLabelValues(g.Request.Method, g.FullPath(), strconv.Itoa(g.Writer.Status())).Inc()
				return
			}
			log.Printf("File will expire at Unix time: %d\n", expireDuration)
			f.Expiration = expireDuration
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
			middleware.HttpBadRequests.WithLabelValues(g.Request.Method, g.FullPath(), strconv.Itoa(g.Writer.Status())).Inc()
			return
		}

		// Define the destination path for the uploaded file
		dst := filepath.Join(fconfig.UploadDir, file.Filename)

		// Save the uploaded file
		if err := g.SaveUploadedFile(file, dst); err != nil {
			middleware.FailedFileUploadRequests.WithLabelValues("1").Inc()
			g.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			middleware.HttpBadRequests.WithLabelValues(g.Request.Method, g.FullPath(), strconv.Itoa(g.Writer.Status())).Inc()
			return
		} else {
			middleware.SuccessfulFileUploadRequests.WithLabelValues("1").Inc()
			// Insert file metadata into the database
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
