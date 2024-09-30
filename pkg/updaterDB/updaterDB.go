package main

import (
	"log"
	//"time"
	"database/sql"

	pkg_logger "cututer/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
	logger *log.Logger
)

func init() {
	pkg_logger.CreateLogger()
	logger = pkg_logger.GetLogger()
}

func main() {
	logger.Println("updaterDB successfully started")
	initDB()
	defer db.Close()

	for {
		rows, err := db.Query(
			"SELECT * FROM urls"
		)

		if err != nil {
			logger.Println("initDB failed because %v", err)
			return
		}

		defer rows.Close()
		objects := []Structre{}

		for rows.Next() {
			
		}
	}
}