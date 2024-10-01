package database

import (
	"log"
	"database/sql"
	
	pkg_logger "cututer/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
	logger *log.Logger
)

func InitDB() {
	logger = pkg_logger.GetLogger()

	logger.Println("initDB successfully started")

	var err error

	db, err = sql.Open("sqlite3", "../../database/urls.db")
	if err != nil {
		logger.Fatal(err)
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
		logger.Fatal(err)
		return
	}

	logger.Println("initDB successfully executed")
}

func GetDB() *sql.DB {
	return db
}