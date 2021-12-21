package main

import "C"
import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/handler"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func init() {
	cfg := config.Conf
	if err := env.Parse(cfg); err != nil {
		panic("failed to parse config: " + err.Error())
	}

	serverAddress := flag.String("a", "", "server address")
	baseURL := flag.String("b", "", "base url")
	fileStoragePath := flag.String("f", "", "file storage path")

	flag.Parse()

	if *serverAddress != "" {
		cfg.ServerAddress = *serverAddress
	}
	if *baseURL != "" {
		cfg.BaseURL = *baseURL
	}
	if *fileStoragePath != "" {
		cfg.FileStoragePath = *fileStoragePath
	}
}

func main() {
	st := storage.NewStorage(false)

	flag.Parse()

	r := chi.NewRouter()

	r.Post("/api/shorten", handler.ShortenPost(st))
	r.Get("/{shortUrl}", handler.Get(st))
	r.Post("/", handler.Post(st))

	log.Fatal(http.ListenAndServe(config.Conf.ServerAddress, r))
}
