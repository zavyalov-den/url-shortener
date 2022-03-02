package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

var Conf = parseConfig()

// Раньше здесь была вся логика по обработке переменных среды и флагам,
// но переехала в main, т.к. при запуске тестов через go test ./... возникала ошибка

func parseConfig() *config {
	cfg := &config{}

	if err := env.Parse(cfg); err != nil {
		fmt.Println("failed to parse config: " + err.Error())
	}

	serverAddress := flag.String("a", "", "server address")
	baseURL := flag.String("b", "", "base url")
	fileStoragePath := flag.String("f", "", "file storage path")

	flag.Parse()

	if serverAddress != nil && *serverAddress != "" {
		cfg.ServerAddress = *serverAddress
	}
	if baseURL != nil && *baseURL != "" {
		cfg.BaseURL = *baseURL
	}
	if fileStoragePath != nil && *fileStoragePath != "" {
		cfg.FileStoragePath = *fileStoragePath
	}

	return cfg
}

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.json"`
}
