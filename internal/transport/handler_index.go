package transport

import (
	"log"
	"net/http"
)

func IndexUrlHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println("indexUrlHandler successfully started")

		http.ServeFile(w, r, "../../web/index.html")

		logger.Println("indexUrlHandler successfully executed")
	}
}