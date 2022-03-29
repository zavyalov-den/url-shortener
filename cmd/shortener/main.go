package main

import "C"
import (
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/handler"
	"github.com/zavyalov-den/url-shortener/internal/middlewares"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var cfg = config.Config

func main() {
	st := storage.NewStorage()

	r := chi.NewRouter()

	r.Use(middlewares.GzipHandle)
	r.Use(middlewares.Auth)

	r.Post("/api/shorten", handler.ShortenPost(st))
	r.Get("/api/user/urls", handler.GetUserUrls(st))
	r.Get("/{shortUrl}", handler.Get(st))
	r.Post("/", handler.Post(st))
	r.Get("/ping", handler.Ping(st))

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
