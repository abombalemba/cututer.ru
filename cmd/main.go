package main

import (
	//"bytes"
	//"image"
	//"image/gif"
	"log"
	"sync"
	"time"
	"strings"
	"net/http"
	"math/rand"
	"database/sql"
	"encoding/json"

	"cututer/internal/config"
	"cututer/internal/models"
	//"cututer/database"

	_ "github.com/mattn/go-sqlite3"
	favicon "github.com/go-http-utils/favicon"
)

var (
	db *sql.DB
	mu sync.Mutex
)

func main() {
	InitDB()
	defer db.Close()
	/*
	http.HandleFunc("/", indexUrlHandler)
	http.HandleFunc("/api", apiUrlHandler)
	http.HandleFunc("/c/", cUrlHandler)
	//http.HandleFunc("/favicon.ico", faviconHandler)

	if err := http.ListenAndServe(":" + config.Port, favicon.Handler()); err != nil {
		panic(err)
	}*/

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexUrlHandler)
	mux.HandleFunc("/api", apiUrlHandler)
	mux.HandleFunc("/c/", cUrlHandler)

	http.ListenAndServe(":" + config.Port, favicon.Handler(mux, "../web/favicon/32x32.ico"))
}

func generateShortUrl(originalUrl string) string {
	str := generateRandomString()

	checkUrlSQL := `
	SELECT COUNT(*) FROM urls WHERE short_url == ?
	`

	res, err := db.Exec(checkUrlSQL, str)
	if err != nil {
		log.Fatalf("error or this short is existing", res, err)
	}

	return str
}

func generateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	str := make([]byte, config.LengthShortUrl)

	for i := range str {
		str[i] = config.Letters[rand.Intn(len(config.Letters))]
	}

	return string(str)
}

func checkGeneratedShortUrl(shortUrl string) {

}

func indexUrlHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("indexUrlHandler successfully started")

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeFile(w, r, "../web/index.html")

	log.Println("indexUrlHandler successfully executed")
}

func apiUrlHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("apiUrlHandler successfully started")

	if r.Method != http.MethodPost {
		http.Error(w, "apiUrlHandler failed because another method not allowed", http.StatusMethodNotAllowed)
		log.Fatalf("apiUrlHandler failed because another method not allowed", r.Method)
		return
	}

	var req models.UrlRequest

	if err := r.ParseForm(); err != nil {
		http.Error(w, "apiUrlHandler failed because there was a parsing error", http.StatusBadRequest)
		log.Fatalf("apiUrlHandler failed because there was a parsing error", err)
		return
	}

	req.OriginalUrl = r.FormValue("original_url")

	if req.OriginalUrl == "" {
		http.Error(w, "apiUrlHandler failed because req.OriginalUrl = \"\"", http.StatusBadRequest)
		log.Fatalf("apiUrlHandler failed because req.OriginalUrl = \"\"")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	is, err := originalUrlInDB(req.OriginalUrl)
	if err != nil {
		http.Error(w, "apiUrlHandler failed because SQL query got error 142", http.StatusBadRequest)
		log.Fatalf("apiUrlHandler failed because SQL query got error 142", err)
		return
	}

	var shortUrl string

	if is == true {
		row := db.QueryRow("SELECT short_url FROM urls WHERE original_url == ?", req.OriginalUrl)

		err := row.Scan(&shortUrl)

		if err != nil {
			http.Error(w, "apiUrlHandler failed because SQL query got error 155", http.StatusBadRequest)
			log.Fatalf("apiUrlHandler failed because SQL query got error 155")
			return
		}
	} else {
		shortUrl = generateShortUrl(req.OriginalUrl)

		_, err := db.Exec(
			"INSERT INTO urls (original_url, short_url) VALUES (?, ?)", req.OriginalUrl, shortUrl)

		if err != nil {
			http.Error(w, "apiUrlHandler failed because SQL query got error", http.StatusInternalServerError)
			log.Fatalf("apiUrlHandler failed because SQL query got error 164", err)
			return
		}
	}

	shortUrl = config.Protocol + config.Host + ":" + config.Port + config.Path + shortUrl

	response := models.UrlResponse{ShortUrl: shortUrl}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Println("apiUrlHandler successfully executed")
}

func cUrlHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("cUrlHandler successfully started")

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/c/")

	row := db.QueryRow(
		"SELECT original_url FROM urls WHERE short_url == ? LIMIT 1", path)

	var originalUrl string

	err := row.Scan(&originalUrl)

	if err == nil {
		http.Redirect(w, r, originalUrl, http.StatusFound)

		log.Println("cUrlHandler successfully executed")
	} else if err == sql.ErrNoRows {
		http.ServeFile(w, r, "../web/notfound.html")

		log.Println("cUrlHandler not found url")
	} else {
		http.Error(w, "cUrlHandler failed because SQL query got error", http.StatusBadRequest)

		log.Fatalf("cUrlHandler failed because SQL query got error", err)
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

func InitDB() {
	log.Println("initDB successfully started")

	var err error

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
