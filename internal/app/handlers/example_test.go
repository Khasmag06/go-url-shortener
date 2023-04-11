package handlers

import (
	"log"
	"net/http"

	"github.com/Khasmag06/go-url-shortener/config"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

func ExampleService_Route() {
	r := chi.NewRouter()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	repo := storage.NewMemoryStorage()
	s := NewService(*cfg, repo)
	s.Route()
	http.ListenAndServe(cfg.ServerAddress, r)
}
