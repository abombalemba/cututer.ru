package database

import (
	"log"
	"os"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func InitDB() {
	var err error
	log.Println(os.Getwd())
	db, err = sql.Open("sqlite3", "../database/urls.db")
	if err != nil {
		log.Fatal(err)
		return
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_url TEXT NOT NULL,
		short_url TEXT NOT NULL UNIQUE
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("initDB successfully executed")
}
