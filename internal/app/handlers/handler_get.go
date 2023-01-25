package handlers

import (
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"regexp"
)

var ShortIDValid = regexp.MustCompile(`^/([a-zA-Z]{6})$`)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	shortID := "/" + chi.URLParam(r, "id")

	if !ShortIDValid.MatchString(shortID) {
		http.Error(w, "Incorrect parameters, you can only use letters", http.StatusBadRequest)
		return
	}

	url, err := storage.Urls.Get(shortID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home page"))
}
