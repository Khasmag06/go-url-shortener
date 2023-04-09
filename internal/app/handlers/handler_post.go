package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Khasmag06/go-url-shortener/internal/app/middleware"
	"github.com/Khasmag06/go-url-shortener/internal/app/shorten"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
)

func (s *Service) PostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	userID := r.Context().Value(middleware.UserIDKey).(string)
	urlOriginal := string(body)
	short := shorten.URLShorten()
	shortURL := storage.ShortURL{ID: short, OriginalURL: urlOriginal}
	err = s.repo.AddShortURL(userID, &shortURL)
	if err != nil && errors.Is(err, storage.ErrExistsURL) {
		short, err = s.repo.GetExistURL(string(body))
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "%s/%s", s.cfg.BaseURL, short)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", s.cfg.BaseURL, short)
}
