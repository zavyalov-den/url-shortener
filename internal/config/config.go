package config

import (
	"flag"
	"fmt"
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

	serverAddress := flag.String("a", "", "server address")
	baseURL := flag.String("b", "", "base url")
	fileStoragePath := flag.String("f", "", "file storage path")

	//flag.Parse()

	fmt.Println(cfg)

	//if serverAddress != nil {
	if *serverAddress != "" {
		cfg.ServerAddress = *serverAddress
	}
	if *baseURL != "" {
		cfg.BaseURL = *baseURL
	}
	if *fileStoragePath != "" {
		cfg.FileStoragePath = *fileStoragePath
	}

	fmt.Println(cfg)

	return cfg
}
