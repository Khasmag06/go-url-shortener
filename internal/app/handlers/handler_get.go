package handlers

import (
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlShort := vars["id"]

	if !validUrlParameter(`^([a-zA-Z]{6})$`, urlShort) {
		http.Error(w, "Incorrect parameters, you can only use letters", http.StatusBadRequest)
		return
	}

	urlOriginal := storage.Urls.Get(urlShort)
	if urlOriginal == "" {
		http.Error(w, "url don't found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", urlOriginal)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func validUrlParameter(re, urlParameter string) bool {
	var validPathGet = regexp.MustCompile(re)
	return validPathGet.MatchString(urlParameter)
}
