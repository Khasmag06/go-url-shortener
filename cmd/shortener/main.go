package main

import (
	"github.com/Khasmag06/go-url-shortener/config"
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
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
	r := chi.NewRouter()
	r.Mount("/", s.Route())
	r.Mount("/debug", middleware.Profiler())

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
