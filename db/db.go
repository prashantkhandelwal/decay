package db

import (
	"context"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var dbase *sql.DB

func initDB(ctx context.Context, path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)

	// Pragmas (optional): WAL for concurrency
	if _, err := db.ExecContext(ctx, `PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		return nil, err
	}

	if err != nil {
		log.Fatal("ERROR:Database: Error in opening database.")
	}

	return db, nil
}

func OpenDB() error {
	ctx := context.Background()
	var err error
	dbase, err = initDB(ctx, "data\\decay.db")
	if err != nil {
		panic(err)
	}

	return nil
}
