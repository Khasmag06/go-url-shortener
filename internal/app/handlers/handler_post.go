package handlers

import (
	"fmt"
	"github.com/Khasmag06/go-url-shortener/internal/app/shorten"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"io"
	"net/http"
)

func (s *Service) PostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	urlOriginal := string(body)
	short := "/" + shorten.URLShorten()
	shortURL := storage.ShortURL{ID: short, OriginalURL: urlOriginal}
	if err := s.repo.AddShortURL(&shortURL); err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, s.cfg.BaseURL+short)
}
