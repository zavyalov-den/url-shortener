package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/handler"
	"github.com/zavyalov-den/url-shortener/internal/middlewares"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfigInstance()
	st := storage.NewStorage()

	r := chi.NewRouter()

	r.Use(middlewares.GzipHandle)
	r.Use(middlewares.Auth)

	r.Post("/api/shorten/batch", handler.ShortenBatch(st))
	r.Post("/api/shorten", handler.ShortenJSON(st))
	r.Get("/api/user/urls", handler.GetUserUrls(st))
	r.Get("/{shortUrl}", handler.GetFullURL(st))
	r.Post("/", handler.Shorten(st))
	r.Get("/ping", handler.PingDB(st))

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
