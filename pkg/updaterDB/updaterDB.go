package main

import (
	"log"
	//"time"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func main() {
	initDB()
	defer db.Close()

	for {
		rows, err := db.Query(
			"SELECT * FROM urls"
		)

		if err != nil {
			log.Fatal("initDB failed because %v", err)
			return
		}

		defer rows.Close()
		objects := []Structre{}

		for rows.Next() {
			
		}
	}
}

func initDB() {
	var err error

	db, err := sql.Open("sqlite3", "../../database/urls.db")
	if err != nil {
		log.Fatalf("initDB failed because %v", err)
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
