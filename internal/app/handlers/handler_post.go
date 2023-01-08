package handlers

import (
	"fmt"
	"github.com/Khasmag06/go-url-shortener/internal/app/shorten"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	urlOriginal := r.FormValue("url")
	urlShort := shorten.URLShorten()
	storage.Urls.Put(urlShort, urlOriginal)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://localhost:8080/"+urlShort)
}
