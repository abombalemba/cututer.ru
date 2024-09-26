package main

import (
	"os"
	"log"
	"sync"
	"time"
	"bytes"
	"strings"
	"net/http"
	"math/rand"
	"database/sql"
	"encoding/json"

	//"cututer/database"
	"cututer/internal/config"
	"cututer/internal/models"
	"cututer/tools"

	_ "github.com/mattn/go-sqlite3"
	//favicon "github.com/go-http-utils/favicon"
)

var (
	db *sql.DB
	mu sync.Mutex
	buffer bytes.Buffer
	logger *log.Logger
	fileLog *os.File
)

func init() {
	createLogger()
}

func main() {
	InitDB()

	defer db.Close()
	defer fileLog.Close()
	
	http.HandleFunc("/", indexUrlHandler)
	http.HandleFunc("/api", apiUrlHandler)
	http.HandleFunc("/c/", cUrlHandler)

	if err := http.ListenAndServe(":" + config.Port, nil); err != nil {
		logger.Panicln(err)
	}
}

func createLogger() {
	filename := tools.GetNow()

	fileLog, err := os.OpenFile("../../logs/" + filename + ".log", os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
		return
	}

	logger = log.New(&buffer, "", log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	logger.SetOutput(fileLog)
}

func generateShortUrl(originalUrl string) string {
	str := generateRandomString()

	checkUrlSQL := `
	SELECT COUNT(*) FROM urls WHERE short_url == ?
	`

	res, err := db.Exec(checkUrlSQL, str)
	if err != nil {
		logger.Fatalf("error or this short is existing", res, err)
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

	mu.Lock()
	defer mu.Unlock()

	is, err := originalUrlInDB(req.OriginalUrl)
	if err != nil {
		http.Error(w, "apiUrlHandler failed because SQL query got error 142", http.StatusBadRequest)
		logger.Println("apiUrlHandler failed because SQL query got error 142", err)
		return
	}

	var shortUrl string

	if is == true {
		row := db.QueryRow("SELECT short_url FROM urls WHERE original_url == ?", req.OriginalUrl)

		err := row.Scan(&shortUrl)

		if err != nil {
			http.Error(w, "apiUrlHandler failed because SQL query got error 155", http.StatusBadRequest)
			logger.Println("apiUrlHandler failed because SQL query got error 155")
			return
		}
	} else {
		shortUrl = generateShortUrl(req.OriginalUrl)

		_, err := db.Exec(
			"INSERT INTO urls (original_url, short_url) VALUES (?, ?)", req.OriginalUrl, shortUrl)

		if err != nil {
			http.Error(w, "apiUrlHandler failed because SQL query got error", http.StatusInternalServerError)
			logger.Println("apiUrlHandler failed because SQL query got error 164", err)
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