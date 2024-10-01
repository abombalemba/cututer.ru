package services

import (
	"log"
	"time"
	"math/rand"
	"database/sql"

	"cututer/internal/config"
	"cututer/internal/models"
)

var (
	logger *log.Logger
	db *sql.DB
)

func GenerateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	str := make([]byte, config.LengthShortUrl)

	for i := range str {
		str[i] = config.Symbols[rand.Intn(len(config.Symbols))]
	}

	return string(str)
}

func GenerateShortUrl(logger *log.Logger, db *sql.DB) string {
	var str string

	logger.Println("generateShortUrl successfully started")

	for {
		str = GenerateRandomString()

		row := db.QueryRow(models.CheckUrlSQL, str)
		var n int
		err := row.Scan(&n)
		if err != nil {
			logger.Fatalf("error or this short is existing", err)
			logger.Printf("SELECT COUNT(*) FROM urls WHERE short_url == %s\n", str)
		}

		if n == 0 {
			break
		}
	}

	logger.Println("generateShortUrl successfully executed")
	logger.Println("generateShortUrl generated new url -", str)

	return str
}

func OriginalUrlInDB(logger *log.Logger, db *sql.DB, originalUrl string) (bool, error) {
	logger.Println("originalUrlInDB successfully started")

	row := db.QueryRow("SELECT id FROM urls WHERE original_url == ? LIMIT 1", originalUrl)

	var id int

	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Println("originalUrlInDB successfully executed with false and nil")
			return false, nil
		}
		logger.Println("originalUrlInDB UNSUCCESSFULLY executed with false and", err)
		return false, err
	}

	if id == 0 {
		logger.Println("originalUrlInDB successfully executed with false and nil")
		return false, nil
	}
	logger.Println("originalUrlInDB successfully executed with true and nil")
	return true, nil
}