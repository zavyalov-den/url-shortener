package main

import "C"
import (
	"fmt"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/handler"
	"github.com/zavyalov-den/url-shortener/internal/middlewares"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var cfg = config.C

func main() {
	st := storage.NewStorage(true)

	r := chi.NewRouter()

	r.Post("/api/shorten", handler.ShortenPost(st))
	r.Get("/{shortUrl}", handler.Get(st))
	r.Post("/", handler.Post(st))

	fmt.Println(cfg)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, middlewares.GzipHandle(r)))
}
