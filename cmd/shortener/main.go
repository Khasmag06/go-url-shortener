package main

import (
	"github.com/Khasmag06/go-url-shortener/config"
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	myMiddlewere "github.com/Khasmag06/go-url-shortener/internal/app/middleware"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	repo := storage.NewMemoryStorage()
	if fp := cfg.FileStoragePath; fp != "" {
		repo = storage.NewFileStorage(fp)
	}
	s := handlers.NewService(*cfg, repo)

	r := NewRouter(s)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}

func NewRouter(s *handlers.Service) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Use(myMiddlewere.GzipHandle)
		r.Post("/", s.PostHandler)
		r.Post("/api/shorten", s.PostAPIHandler)
		r.Get("/", handlers.HomeHandler)
		r.Get("/{id}", s.GetHandler)
	})
	return r
}
