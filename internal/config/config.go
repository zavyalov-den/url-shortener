package config

import (
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

	return cfg
}
