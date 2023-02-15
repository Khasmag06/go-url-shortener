package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"net/http"
)

func (s *Service) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	var userShorts []*storage.ShortURL
	for _, el := range s.repo.GetAllShortURL(userID) {
		el.ID = fmt.Sprintf("%s/%s", s.cfg.BaseURL, el.ID)
		userShorts = append(userShorts, el)

	}
	if userShorts == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(userShorts); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}