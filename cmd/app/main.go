package main

import (
	"log"
	"net/http"
	"database/sql"

	"cututer/internal/config"
	"cututer/internal/database"
	"cututer/internal/transport"
	pkg_logger "cututer/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
	//favicon "github.com/go-http-utils/favicon"
)

var (
	db *sql.DB
	logger *log.Logger
)

func init() {
	pkg_logger.InitLogger()
	logger = pkg_logger.GetLogger()

	database.InitDB()
	db = database.GetDB()
}

func main() {
	defer pkg_logger.CloseLogger()
	defer database.CloseDB()
	
	http.HandleFunc("/", transport.IndexUrlHandler(logger))
	http.HandleFunc("/api", transport.ApiUrlHandler(logger, db))
	http.HandleFunc("/c/", transport.CUrlHandler(logger, db))

	if err := http.ListenAndServe(":" + config.Port, nil); err != nil {
		logger.Panicln(err)
	}
}