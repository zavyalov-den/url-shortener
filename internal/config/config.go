package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var C = parseConfig()

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.json"`
}

func parseConfig() *config {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		panic("failed to parse config")
	}

	serverAddress := flag.String("a", ":8080", "server address")
	baseURL := flag.String("b", "http://localhost:8080", "base url")
	fileStoragePath := flag.String("f", "./storage.json", "file storage path")

	flag.Parse()

	if serverAddress != nil {
		cfg.ServerAddress = *serverAddress
	}
	if baseURL != nil {
		cfg.BaseURL = *baseURL
	}
	if fileStoragePath != nil {
		cfg.FileStoragePath = *fileStoragePath
	}

	return cfg
}
