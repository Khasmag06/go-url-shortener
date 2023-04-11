package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Khasmag06/go-url-shortener/internal/app/middleware"
)

type userShort struct {
	userID  string
	shortID string
}

// DeleteHandler удаляет ссылки пользователя.
func (s *Service) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var shortIDs []string
	userID := r.Context().Value(middleware.UserIDKey).(string)

	err := json.NewDecoder(r.Body).Decode(&shortIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ch := make(chan userShort, len(shortIDs))

	for i := 0; i <= 100; i++ {
		go func(ch <-chan userShort) {
			for el := range ch {
				err := s.repo.DeleteShortURL(el.userID, el.shortID)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}(ch)
	}

	for _, shortID := range shortIDs {
		newUserShort := userShort{userID: userID, shortID: shortID}
		ch <- newUserShort
	}
	close(ch)
	w.WriteHeader(http.StatusAccepted)

}
