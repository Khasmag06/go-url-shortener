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
	userID := r.Context().Value("userID").(string)
	urlOriginal := string(body)
	short := shorten.URLShorten()
	shortURL := storage.ShortURL{ID: short, OriginalURL: urlOriginal}
	if err := s.repo.AddShortURL(userID, &shortURL); err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", s.cfg.BaseURL, short)
}
