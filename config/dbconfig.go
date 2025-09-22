package config

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func InitDB() error {
	if _, err := os.Stat("decay.db"); err != nil {

		log.Println("Database: decay.db does not exist....Creating")

		db, err := sql.Open("sqlite", "decay.db")
		if err != nil {
			log.Fatal("Error in opening/creating the database")
			return err
		}

		defer db.Close()

		log.Println("Database: decay.db created!")

		q := `CREATE TABLE "config" (
					"id"	INTEGER NOT NULL UNIQUE,
					"key"	TEXT NOT NULL,
					"value"	TEXT NOT NULL,
					PRIMARY KEY("id" AUTOINCREMENT)
				);
				CREATE TABLE "bookmarks" (
					"id"	INTEGER NOT NULL UNIQUE,
					"url"	TEXT NOT NULL,
					"title"	BLOB NOT NULL,
					"description"	TEXT,
					"snapshot"	REAL,
					"date_added"	DATETIME NOT NULL,
					"date_modified"	DATETIME,
					"tags"	INTEGER,
					"is_archived"	INTEGER NOT NULL,
					PRIMARY KEY("id" AUTOINCREMENT)
				);`

		_, err = db.Exec(q)
		if err != nil {
			db.Close()
			// err := os.Remove("bind.db")
			// if err != nil {
			// 	log.Fatalf("Setup failed - Cannot delete database %s", err)
			// }
			log.Fatalf("Database: Error in setting up database tables. - %s", err.Error())
			return err
		}

		log.Print("Database: Database setup completed successfully")

	} else {
		log.Println("Database: decay.db found!")
	}

	return nil
}
