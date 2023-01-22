package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/Khasmag06/go-url-shortener/internal/app/shorten"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"net/http"
)

type JSONOriginalURL struct {
	URL string `json:"url"`
}

type JSONShortURL struct {
	Result string `json:"result"`
}

func PostAPIHandler(w http.ResponseWriter, r *http.Request) {
	var u JSONOriginalURL
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	short := shorten.URLShorten()
	urlOriginal := u.URL
	shortURL := storage.ShortURL{ID: short, OriginalURL: urlOriginal}
	storage.Urls.Add(&shortURL)

	var buf bytes.Buffer
	shortJSON := JSONShortURL{Result: localhost + short}
	if err := json.NewEncoder(&buf).Encode(shortJSON); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())

}
