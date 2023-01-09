package handlers

import (
	"fmt"
	"github.com/Khasmag06/go-url-shortener/internal/app/shorten"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {

	urlOriginal := r.FormValue("url")
	short := shorten.URLShorten()
	storage.Urls.Put(short, urlOriginal)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://localhost:8080/"+short)
}
