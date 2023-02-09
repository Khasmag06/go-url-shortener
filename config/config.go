package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		BaseURL:         "",
		ServerAddress:   "",
		FileStoragePath: "",
	}
	flag.StringVar(&cfg.ServerAddress, "a", "", "Server address")
	flag.StringVar(&cfg.BaseURL, "b", "", "Base url")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "File storage path")
	flag.Parse()
	cfg.BaseURL = pickFirstNonEmpty(cfg.BaseURL, os.Getenv("BASE_URL"), "http://localhost:8080")
	cfg.ServerAddress = pickFirstNonEmpty(cfg.ServerAddress, os.Getenv("SERVER_ADDRESS"), ":8080")
	cfg.FileStoragePath = pickFirstNonEmpty(cfg.FileStoragePath, os.Getenv("FILE_STORAGE_PATH"))
	return cfg, nil
}

func pickFirstNonEmpty(strings ...string) string {
	for _, str := range strings {
		if str != "" {
			return str
		}
	}
	return ""
}
