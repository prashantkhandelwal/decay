package db

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
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

	urlViewer := serverURL + ":" + strconv.Itoa(int(c.Server.PORT)) + "/view/" + id
	fileURL := serverURL + ":" + strconv.Itoa(int(c.Server.PORT)) + "/file/" + id
	displayURL := serverURL + ":" + strconv.Itoa(int(c.Server.PORT)) + "/display/" + id

	f.ID = id

	if f.Title == "" {
		f.Title = strings.TrimSuffix(f.Filename, filepath.Ext(f.Filename))
	}

	f.URLViewer = urlViewer
	f.URL = fileURL
	f.DisplayURL = displayURL

	f.Time = time.Now().Unix() // time the file was uploaded
	f.Filename = file.Filename

	result, err := dbase.ExecContext(ctx, query,
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
