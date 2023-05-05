// Сервис по сокращению ориганальной ссылки.
package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/Khasmag06/go-url-shortener/config"
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

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

	if cfg.EnableHTTPS {
		log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", r))
	}
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
