package main

import (
	"flag"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/handler"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS);
//флаг -b, отвечающий за базовый адрес результирующего сокращённого URL (переменная BASE_URL);
//флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH).

func init() {

}

func main() {
	st := storage.NewStorage(true)

	flag.Parse()

	r := chi.NewRouter()

	r.Post("/api/shorten", handler.ShortenPost(st))
	r.Get("/{shortUrl}", handler.Get(st))
	r.Post("/", handler.Post(st))

	log.Fatal(http.ListenAndServe(config.C.ServerAddress, r))
}
