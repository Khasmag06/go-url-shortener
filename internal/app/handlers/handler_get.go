package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

// ShortIDValid валидация параметра короткой ссылки.
var ShortIDValid = regexp.MustCompile(`^([a-zA-Z]{6})$`)

// GetHandler возвращеет короткую ссылку
func (s *Service) GetHandler(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "id")
	if !ShortIDValid.MatchString(shortID) {
		http.Error(w, "Incorrect parameters, you can only use letters", http.StatusBadRequest)
		return
	}
	url, err := s.repo.GetShortURL(shortID)

	if err != nil && errors.Is(err, storage.ErrNotAvailable) {
		http.Error(w, err.Error(), http.StatusGone)
		return
	}
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, storage.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	w.Header().Set("Location", url.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
