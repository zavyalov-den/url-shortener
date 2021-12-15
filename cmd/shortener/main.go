package main

import (
	"github.com/zavyalov-den/url-shortener/internal/handler"
	"github.com/zavyalov-den/url-shortener/internal/repo"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	urls := repo.NewRepo()

	r := chi.NewRouter()

	r.Post("/api/shorten", handler.ShortenPost(urls))
	r.Get("/{shortUrl}", handler.Get(urls))
	r.Post("/", handler.Post(urls))

	log.Fatal(http.ListenAndServe(":8080", r))
}
