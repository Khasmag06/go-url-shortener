package handlers

import (
	"github.com/Khasmag06/go-url-shortener/config"
	mw "github.com/Khasmag06/go-url-shortener/internal/app/middleware"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Service struct {
	cfg  config.Config
	repo storage.Storage
}

func NewService(cfg config.Config, repo storage.Storage) *Service {
	return &Service{cfg, repo}
}

func (s *Service) Route() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.GzipHandle, mw.CreateAccessToken)

	r.Route("/", func(r chi.Router) {
		r.Post("/", s.PostHandler)
		r.Get("/{id}", s.GetHandler)
		r.Get("/ping", s.PingHandler)
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", s.PostJSONHandler)
			r.Post("/shorten/batch", s.BatchHandler)
			r.Get("/user/urls", s.GetUserURLsHandler)
			r.Delete("/user/urls", s.DeleteHandler)
			r.Group(func(r chi.Router) {
				r.Use(func(next http.Handler) http.Handler {
					return mw.CheckTrustedSubnet(next, s.cfg.TrustedSubnet)
				})
				r.Get("/internal/stats", s.GetStats)
			})
		})
	})
	return r
}
