package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Khasmag06/go-url-shortener/internal/app/middleware"
	"github.com/Khasmag06/go-url-shortener/internal/app/shorten"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
)

// JSONOriginalURL описание модели запроса оригинальной ссылки.
type JSONOriginalURL struct {
	URL string `json:"url"`
}

// JSONShortURL описание модели ответа короткой ссылки.
type JSONShortURL struct {
	Result string `json:"result"`
}

// PostJSONHandler создает короткую ссылку если оригиальная ссылка не существует.
func (s *Service) PostJSONHandler(w http.ResponseWriter, r *http.Request) {
	var u JSONOriginalURL
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(middleware.UserIDKey).(string)
	short := shorten.URLShorten()
	urlOriginal := u.URL
	shortURL := storage.ShortURL{ID: short, OriginalURL: urlOriginal}

	var buf bytes.Buffer

	err := s.repo.AddShortURL(userID, &shortURL)
	if err != nil && errors.Is(err, storage.ErrExistsURL) {
		short, err = s.repo.GetExistURL(shortURL.OriginalURL)
		if err != nil {
			log.Println(err)
		}
		shortJSON := JSONShortURL{Result: fmt.Sprintf("%s/%s", s.cfg.BaseURL, short)}
		json.NewEncoder(&buf).Encode(shortJSON)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(buf.Bytes())
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortJSON := JSONShortURL{Result: fmt.Sprintf("%s/%s", s.cfg.BaseURL, short)}
	if err := json.NewEncoder(&buf).Encode(shortJSON); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}
