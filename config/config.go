package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	flagServerAddress := flag.String("a", "localhost:8080", "Server address")
	flagBaseURL := flag.String("b", "http://localhost:8080", "Base url")
	flagFileStoragePath := flag.String("f", "", "File storage path")
	flag.Parse()
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = *flagServerAddress
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = *flagBaseURL
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *flagFileStoragePath
	}
	return &cfg, err

}
