package handlers

import (
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	urlShort := "http://localhost:8080" + r.URL.Path
	w.Header().Add("Location", storage.Urls.Get(urlShort))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
