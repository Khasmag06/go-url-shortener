package handlers

import (
	"fmt"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	urlShort := r.URL.Path
	fmt.Println(urlShort)
	w.Header().Add("Location", storage.Urls.Get(urlShort))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
