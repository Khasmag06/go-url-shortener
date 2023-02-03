package handlers

import (
	"github.com/Khasmag06/go-url-shortener/config"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
)

type Service struct {
	cfg  config.Config
	repo storage.Storage
}

func NewService(cfg config.Config, repo storage.Storage) *Service {
	return &Service{cfg, repo}
}
