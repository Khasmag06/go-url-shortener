package handlers

import (
	"github.com/Khasmag06/go-url-shortener/config"
	myMiddleware "github.com/Khasmag06/go-url-shortener/internal/app/middleware"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Service структура, описывающая взаимодействие обработчиков с хранилищем и конфигурационными данными.
type Service struct {
	cfg  config.Config
	repo storage.Storage
}

// NewService конструктор для Service.
func NewService(cfg config.Config, repo storage.Storage) *Service {
	return &Service{cfg, repo}
}

// Route группирует запросы и возвращает роутер.
func (s *Service) Route() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(myMiddleware.GzipHandle, myMiddleware.CreateAccessToken)

	r.Route("/", func(r chi.Router) {
		r.Post("/", s.PostHandler)
		r.Get("/{id}", s.GetHandler)
		r.Get("/ping", s.PingHandler)
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", s.PostJSONHandler)
			r.Post("/shorten/batch", s.BatchHandler)
			r.Get("/user/urls", s.GetUserURLsHandler)
			r.Delete("/user/urls", s.DeleteHandler)
		})
	})
	return r
}
