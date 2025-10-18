package config

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func InitDB() error {
	if _, err := os.Stat("data\\decay.db"); err != nil {

		log.Println("Database: decay.db does not exist....Creating")

		db, err := sql.Open("sqlite", "data\\decay.db")
		if err != nil {
			log.Fatal("Error in opening/creating the database")
			return err
		}

		defer db.Close()

		log.Println("Database: decay.db created!")

		schema := `CREATE TABLE IF NOT EXISTS refresh_tokens (
				jti         TEXT PRIMARY KEY,
				username    TEXT NOT NULL,
				expires_at  INTEGER NOT NULL,        -- unix seconds
				revoked     INTEGER NOT NULL DEFAULT 0,
				issued_at   INTEGER NOT NULL         -- unix seconds
				);
				CREATE INDEX IF NOT EXISTS idx_rt_username ON refresh_tokens(username);

				CREATE TABLE IF NOT EXISTS user_session (
				username    TEXT PRIMARY KEY,
				valid_after INTEGER NOT NULL DEFAULT 0  -- unix seconds (0 = never invalidated)
				);

				CREATE TABLE IF NOT EXISTS uploads (
				id			TEXT NOT NULL UNIQUE,
				title		TEXT NOT NULL,
				url_viewer	TEXT NOT NULL,
				url			TEXT NOT NULL,
				display_url	TEXT NOT NULL,
				width		INTEGER,
				height		INTEGER,
				size		INTEGER NOT NULL,
				time		TEXT NOT NULL,
				expiration	INTEGER NOT NULL,
				filename	TEXT NOT NULL,
				mime		TEXT,
				PRIMARY KEY(id)
			);`

		_, err = db.Exec(schema)
		if err != nil {
			db.Close()
			log.Fatalf("Database: Error in setting up database tables. - %s", err.Error())
			return err
		}

		log.Print("Database: Database setup completed successfully")

	} else {
		log.Println("Database: decay.db found!")
	}

	return nil
}
