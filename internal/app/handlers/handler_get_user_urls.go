package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Khasmag06/go-url-shortener/internal/app/middleware"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
)

// GetUserURLsHandler возвращает список всех коротких ссылок пользователя.
func (s *Service) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	var userShorts []storage.ShortURL
	shorts, _ := s.repo.GetAllShortURL(userID)
	for _, el := range shorts {
		var short = *el
		short.ID = fmt.Sprintf("%s/%s", s.cfg.BaseURL, el.ID)
		userShorts = append(userShorts, short)

	}
	if userShorts == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userShorts); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
