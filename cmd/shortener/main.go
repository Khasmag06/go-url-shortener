package main

import (
	"github.com/Khasmag06/go-url-shortener/config"
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	r := NewRouter()
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.PostHandler)
		r.Post("/api/shorten", handlers.PostAPIHandler)
		r.Get("/", handlers.HomeHandler)
		r.Get("/{id}", handlers.GetHandler)
	})
	return r
}
