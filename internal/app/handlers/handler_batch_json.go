package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Khasmag06/go-url-shortener/internal/app/middleware"
	"github.com/Khasmag06/go-url-shortener/internal/app/shorten"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
)

// JSONBatchReq описание модели запроса оригинальной ссылки.
type JSONBatchReq struct {
	CorID       string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

// JSONBatchResp описание модели ответа короткой ссылки.
type JSONBatchResp struct {
	CorID    string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// BatchHandler создает список коротких ссылок.
func (s *Service) BatchHandler(w http.ResponseWriter, r *http.Request) {
	var batchReq []JSONBatchReq
	err := json.NewDecoder(r.Body).Decode(&batchReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(middleware.UserIDKey).(string)
	batchResp := make([]JSONBatchResp, len(batchReq))
	for i, el := range batchReq {
		short := shorten.URLShorten()
		shortURL := storage.ShortURL{ID: short, OriginalURL: el.OriginalURL}
		if err := s.repo.AddShortURL(userID, &shortURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		batchResp[i].CorID = el.CorID
		batchResp[i].ShortURL = fmt.Sprintf("%s/%s", s.cfg.BaseURL, short)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(batchResp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
