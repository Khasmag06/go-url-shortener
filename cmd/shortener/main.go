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
	if dsn := cfg.DatabaseDsn; dsn != "" {
		repo, err = storage.NewDB(dsn)
		if err != nil {
			log.Fatalf("unable to create database storage: %v", err)
		}

	} else if fp := cfg.FileStoragePath; fp != "" {
		repo, err = storage.NewFileStorage(fp)
		if err != nil {
			log.Fatalf("unable to create file storage: %v", err)
		}
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
	r.Use(myMiddlewere.GzipHandle, myMiddlewere.CreateAccessToken)

	r.Route("/", func(r chi.Router) {
		r.Post("/", s.PostHandler)
		r.Get("/", handlers.HomeHandler)
		r.Get("/{id}", s.GetHandler)
		r.Get("/ping", s.PingHandler)
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", s.PostJSONHandler)
			r.Post("/shorten/batch", s.BatchHandler)
			r.Get("/user/urls", s.GetUserURLsHandler)
		})
	})
	return r
}
