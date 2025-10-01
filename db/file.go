package db

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/prashantkhandelwal/decay/config"
	"github.com/prashantkhandelwal/decay/utils"
	_ "modernc.org/sqlite"
)

type File struct {
	ID         string
	Title      string
	URLViewer  string
	URL        string
	DisplayURL string
	Width      int
	Height     int
	Size       int64
	Time       int64
	Expiration int64
	Filename   string
	Mime       string
}

func InsertFile(f File, file *multipart.FileHeader) (string, error) {
	db, _ := GetDB()
	query := `INSERT INTO uploads (id, title, url_viewer, url, display_url, width, height, size, time, expiration, filename, mime)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	id, err := utils.GenerateBase62ID()
	if err != nil {
		return "", err
	}

	c := config.GetConfig()
	serverURL := c.Server.URL
	if serverURL == "" {
		serverURL = "localhost"
	}

	ctx := context.Background()

	urlViewer := serverURL + ":" + c.Server.PORT + "/view/" + id
	fileURL := serverURL + ":" + c.Server.PORT + "/file/" + id
	displayURL := serverURL + ":" + c.Server.PORT + "/display/" + id

	f.ID = id
	f.URLViewer = urlViewer
	f.URL = fileURL
	f.DisplayURL = displayURL

	f.Time = time.Now().Unix() // time the file was uploaded
	f.Expiration = 0           // in seconds, 0 means never expires
	f.Filename = file.Filename

	result, err := db.ExecContext(ctx, query,
		f.ID,
		f.Title,
		f.URLViewer,
		f.URL,
		f.DisplayURL,
		f.Width,
		f.Height,
		f.Size,
		f.Time,
		f.Expiration,
		f.Filename,
		f.Mime)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected > 0 {
		return id, nil
	}

	return "", err

}
