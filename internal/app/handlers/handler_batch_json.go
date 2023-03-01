package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Khasmag06/go-url-shortener/internal/app/middleware"
	"github.com/Khasmag06/go-url-shortener/internal/app/shorten"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"net/http"
)

type JSONBatchReq struct {
	CorID       string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type JSONBatchResp struct {
	CorID    string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

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

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(batchResp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())

}
