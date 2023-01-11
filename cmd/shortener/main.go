package main

import (
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	r := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.PostHandler)
		r.Get("/", handlers.HomeHandler)
		r.Get("/{id}", handlers.GetHandler)
	})
	return r
}
