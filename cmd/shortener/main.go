// Сервис по сокращению ориганальной ссылки.
package main

import (
	"context"
	"fmt"
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers/gRPC"
	"log"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Khasmag06/go-url-shortener/config"
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	printInfo()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := getStorage(cfg)
	if err != nil {
		log.Fatalf("enable to create database or file storage: %v", err)
	}

	grpcServer := gRPC.NewShortenerServer(*cfg, repo)
	grpcServer.Run()

	s := handlers.NewService(*cfg, repo)
	r := chi.NewRouter()
	r.Mount("/", s.Route())
	r.Mount("/debug", middleware.Profiler())

	srv := http.Server{Addr: cfg.ServerAddress, Handler: r}
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if cfg.EnableHTTPS {
		srv.Addr = ":443"
		if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ServeTLS: %v", err)
		}
	} else {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}

	<-idleConnsClosed

	fmt.Println("Server Shutdown gracefully")
}

func printInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
func getStorage(cfg *config.Config) (storage.Storage, error) {
	if dsn := cfg.DatabaseDsn; dsn != "" {
		return storage.NewDB(dsn)
	} else if fp := cfg.FileStoragePath; fp != "" {
		return storage.NewFileStorage(fp)
	}
	return storage.NewMemoryStorage(), nil
}
