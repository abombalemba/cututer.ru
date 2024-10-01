package transport

import (
	"log"
	"sync"
	"strings"
	"net/http"
	"database/sql"
	"encoding/json"

	"cututer/internal/config"
	"cututer/internal/models"
	"cututer/internal/services"
	"cututer/tools"
)

var (
	mu sync.Mutex
)

func ApiUrlHandler(logger *log.Logger, db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) { 
		logger.Println("apiUrlHandler successfully started")

		if r.Method != http.MethodPost {
			http.Error(w, "apiUrlHandler failed because another method not allowed", http.StatusMethodNotAllowed)
			logger.Println("apiUrlHandler failed because another method not allowed", r.Method)
			return
		}

		var req models.UrlRequest

		if err := r.ParseForm(); err != nil {
			http.Error(w, "apiUrlHandler failed because there was a parsing error", http.StatusBadRequest)
			logger.Println("apiUrlHandler failed because there was a parsing error", err)
			return
		}

		req.OriginalUrl = r.FormValue("original_url")

		if req.OriginalUrl == "" {
			http.Error(w, "apiUrlHandler failed because req.OriginalUrl = \"\"", http.StatusBadRequest)
			logger.Println("apiUrlHandler failed because req.OriginalUrl = \"\"")
			return
		}

		if !strings.Contains(req.OriginalUrl, "://") {
			http.Error(w, "apiUrlHandler failed because req.OriginalUrl does not contains ://", http.StatusBadRequest)
			logger.Println("apiUrlHandler failed because req.OriginalUrl does not contains ://", req.OriginalUrl)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		is, err := services.OriginalUrlInDB(logger, db, req.OriginalUrl)
		if err != nil {
			http.Error(w, "apiUrlHandler failed because SQL SELECT COPY ORIGINAL URL query got error", http.StatusBadRequest)
			logger.Println("apiUrlHandler failed because SQL SELECT COPY ORIGINAL URL query got error", err)
			logger.Printf("SELECT id FROM urls WHERE original_url == %s\n", req.OriginalUrl)
			return
		}

		var shortUrl string

		if is == true {
			row := db.QueryRow("SELECT short_url FROM urls WHERE original_url == ?", req.OriginalUrl)

			err := row.Scan(&shortUrl)

			if err != nil {
				http.Error(w, "apiUrlHandler failed because SQL SELECT COPY ORIGINAL URL query got error", http.StatusBadRequest)
				logger.Println("apiUrlHandler failed because SQL SELECT COPY ORIGINAL URL query got error", err)
				logger.Printf("SELECT short_url FROM urls WHERE original_url == %s\n", req.OriginalUrl)
				return
			}
		} else {
			shortUrl = services.GenerateShortUrl(logger, db)
			now := tools.GetNow()

			_, err := db.Exec(
				"INSERT INTO urls (original_url, short_url, spawn_date) VALUES (?, ?, ?)", req.OriginalUrl, shortUrl, now)

			if err != nil {
				http.Error(w, "apiUrlHandler failed because SQL INSERT NEW URL query got error", http.StatusInternalServerError)
				logger.Println("apiUrlHandler failed because SQL INSERT NEW URL query got error", err)
				logger.Printf("INSERT INTO urls (original_url, short_url, spawn_date) VALUES (%s, %s, %s)\n", req.OriginalUrl, shortUrl, now)
				return
			}
		}

		shortUrl = config.Protocol + config.Host + ":" + config.Port + config.Path + shortUrl

		response := models.UrlResponse{ShortUrl: shortUrl}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		logger.Println("apiUrlHandler successfully executed")
	}
}