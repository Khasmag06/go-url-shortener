package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

// Config структура описывающая конфигурационные данные.
type Config struct {
	Config          string
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"server_address"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" envDefault:"false" json:"enable_https"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"" json:"file_storage_path"`
	DatabaseDsn     string `env:"DATABASE_DSN" envDefault:"" json:"database_dsn"`
	//DatabaseDsn     string `env:"DATABASE_DSN" envDefault:"postgres://localhost:5432/postgres?sslmode=disable"`
}

// NewConfig конструктор для Config
func NewConfig() (*Config, error) {
	var cfg Config
	cfg.Config = os.Getenv("CONFIG")
	flag.StringVar(&cfg.Config, "c", cfg.Config, "Config file")
	flag.Parse()

	if cfg.Config != "" {
		content, err := os.ReadFile(cfg.Config)
		if err != nil {
			return nil, fmt.Errorf("unable to open file %s: %w", cfg.Config, err)
		}
		err = json.Unmarshal(content, &cfg)
		if err != nil {
			return nil, fmt.Errorf("unable to decode contents of file %s: %w", cfg.Config, err)
		}
	}
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
