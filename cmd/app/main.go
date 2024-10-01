package main

import (
	"log"
	"sync"
	"bytes"
	"strings"
	"net/http"
	"database/sql"
	"encoding/json"
	
	"cututer/internal/config"
	intr_database "cututer/internal/database"
	"cututer/internal/models"
	"cututer/internal/services"
	"cututer/tools"
	pkg_logger "cututer/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
	//favicon "github.com/go-http-utils/favicon"
)

var (
	db *sql.DB
	mu sync.Mutex
	buffer bytes.Buffer
	logger *log.Logger
)

func init() {
	pkg_logger.InitLogger()
	logger = pkg_logger.GetLogger()

	intr_database.InitDB()
	db = intr_database.GetDB()
}

func main() {
	defer db.Close()
	defer pkg_logger.CloseLogger()
	
	http.HandleFunc("/", indexUrlHandler)
	http.HandleFunc("/api", apiUrlHandler)
	http.HandleFunc("/c/", cUrlHandler)

	if err := http.ListenAndServe(":" + config.Port, nil); err != nil {
		logger.Panicln(err)
	}
}

func generateShortUrl(originalUrl string) string {
	var str string

	for {
		str = services.GenerateRandomString()

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

	return str
}

func indexUrlHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("indexUrlHandler successfully started")

	http.ServeFile(w, r, "../../web/index.html")

	logger.Println("indexUrlHandler successfully executed")
}

func apiUrlHandler(w http.ResponseWriter, r *http.Request) {
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

	is, err := originalUrlInDB(req.OriginalUrl)
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
		shortUrl = generateShortUrl(req.OriginalUrl)
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

func cUrlHandler(w http.ResponseWriter, r *http.Request) {
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

func originalUrlInDB(originalUrl string) (bool, error) {
	row := db.QueryRow(
		"SELECT id FROM urls WHERE original_url == ?", originalUrl)

	var id int

	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if id == 0 {
		return false, nil
	}
	return true, nil
}
/*
func InitDB() {
	logger.Println("initDB successfully started")

	var err error

	db, err = sql.Open("sqlite3", "../../database/urls.db")
	if err != nil {
		logger.Fatal(err)
		return
	}

	_, err = db.Exec(models.CreateTableSQL)
	if err != nil {
		logger.Fatal(err)
		return
	}

	logger.Println("initDB successfully executed")
}*/