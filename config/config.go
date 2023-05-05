package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// Config структура описывающая конфигурационные данные.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" envDefault:"false"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	DatabaseDsn     string `env:"DATABASE_DSN" envDefault:""`
	//DatabaseDsn     string `env:"DATABASE_DSN" envDefault:"postgres://localhost:5432/postgres?sslmode=disable"`
}

// NewConfig конструктор для Config
func NewConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL")
	flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "Enable HTTPS")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "File Storage Path")
	flag.StringVar(&cfg.DatabaseDsn, "d", cfg.DatabaseDsn, "Database DSN")
	flag.Parse()
	return &cfg, err

}
