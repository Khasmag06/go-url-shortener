package handlers

import (
	"encoding/json"
	"github.com/Khasmag06/go-url-shortener/internal/app/models"
	"log"
	"net/http"
)

// GetStats возвращает количество сокращенных ссылок и пользователей.
func (s *Service) GetStats(w http.ResponseWriter, r *http.Request) {
	var stats models.InternalStats
	if err := s.repo.GetShortAndUserCount(&stats); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
