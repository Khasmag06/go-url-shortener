package handlers

import (
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"regexp"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "id")

	if !validURLParameter(`^([a-zA-Z]{6})$`, shortID) {
		http.Error(w, "Incorrect parameters, you can only use letters", http.StatusBadRequest)
		return
	}

	urlOriginal := storage.Urls.Get(shortID)
	if urlOriginal == "" {
		http.Error(w, "Link don't found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", urlOriginal)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home page"))
}

func validURLParameter(re, urlParameter string) bool {
	var validPathGet = regexp.MustCompile(re)
	return validPathGet.MatchString(urlParameter)
}
