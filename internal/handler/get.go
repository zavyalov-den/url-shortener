package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"net/http"
)

func Get(db *storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get handler. w: ")
		fmt.Println(w)
		shortURL := chi.URLParam(r, "shortUrl")

		longURL, err := db.Get(shortURL)
		if err != nil {
			fmt.Printf("long: %s, short: %s\n", longURL, shortURL)
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Location", longURL)

		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
