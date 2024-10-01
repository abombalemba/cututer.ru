package transport

import (
	"log"
	"strings"
	"net/http"
	"database/sql"
)

func CUrlHandler(logger *log.Logger, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println("cUrlHandler successfully started")

		path := r.URL.Path
		path = strings.TrimPrefix(path, "/c/")

		row := db.QueryRow(
			"SELECT original_url FROM urls WHERE short_url == ? LIMIT 1", path)

		var originalUrl string

		err := row.Scan(&originalUrl)

		if err == nil {
			http.Redirect(w, r, originalUrl, http.StatusFound)

			logger.Println("cUrlHandler successfully executed")
		} else if err == sql.ErrNoRows {
			http.ServeFile(w, r, "../../web/notfound.html")

			logger.Println("cUrlHandler not found url")
		} else {
			http.Error(w, "cUrlHandler failed because SQL query got error", http.StatusBadRequest)

			logger.Println("cUrlHandler failed because SQL query got error", err)
			logger.Printf("SELECT original_url FROM urls WHERE short_url == %s LIMIT 1\n", path)
		}
	}
}