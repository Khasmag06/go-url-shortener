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
	urlShort := r.URL.Path
	urlOriginal := storage.Urls.Get(urlShort)
	if urlOriginal == "" {
		http.Error(w, "url don't found", http.StatusBadRequest)
	}

	w.Header().Add("Location", urlOriginal)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
