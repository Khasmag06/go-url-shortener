package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"strings"
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
	flagServerAddress := flag.String("a", cfg.ServerAddress, "Server address")
	flagBaseURL := flag.String("b", cfg.BaseURL, "Base url")
	flagFileStoragePath := flag.String("f", cfg.FileStoragePath, "File storage path")
	flag.Parse()
	if len(strings.Split(cfg.BaseURL, ":")) < 3 {
		cfg.BaseURL = "http://" + cfg.BaseURL
	}
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
