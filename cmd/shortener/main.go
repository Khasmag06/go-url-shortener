package main

import (
	"flag"
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
	flagServerAddress := flag.String("a", "localhost:8080", "Domain name")
	flagBaseURL := flag.String("b", "http://localhost:8080", "Net address")
	flagFileStoragePath := flag.String("f", "", "File name")
	flag.Parse()
	if config.Cfg.ServerAddress == "" {
		config.Cfg.ServerAddress = *flagServerAddress
	}
	if config.Cfg.BaseURL == "" {
		config.Cfg.BaseURL = *flagBaseURL
	}
	if config.Cfg.FileStoragePath == "" {
		config.Cfg.FileStoragePath = *flagFileStoragePath
	}
	storage.Urls = storage.NewStorage()
	r := NewRouter()
	log.Fatal(http.ListenAndServe(config.Cfg.ServerAddress, r))
}

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Use(myMiddlewere.GzipHandle)
		r.Post("/", handlers.PostHandler)
		r.Post("/api/shorten", handlers.PostAPIHandler)
		r.Get("/", handlers.HomeHandler)
		r.Get("/{id}", handlers.GetHandler)
	})
	return r
}
