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
