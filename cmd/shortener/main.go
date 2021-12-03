package main

import (
	"crypto"
	"encoding/hex"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

func getHandler(urls map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortUrl")

		longURL := urls[shortURL]
		if longURL == "" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func postHandler(urls map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}
		url := string(data)

		md5 := crypto.MD5.New()
		md5.Write(data)
		short := hex.EncodeToString(md5.Sum(nil))[:8]

		urls[short] = url

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("http://localhost:8080/" + short))
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}
	}
}

func main() {
	urls := make(map[string]string)

	r := chi.NewRouter()

	r.Get("/{shortUrl}", getHandler(urls))
	r.Post("/", postHandler(urls))

	log.Fatal(http.ListenAndServe(":8080", r))
}
