package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

var Config *config

func GetConfigInstance() *config {
	if Config == nil {
		Config = parseConfig()
	}

	return Config
}

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	AuthKey         string `env:"AUTH_KEY" envDefault:"auth"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func parseConfig() *config {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Println("failed to parse config: %w", err)
	}
	//todo: secure auth key

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base url")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database data source name")

	flag.Parse()

	return cfg
}
