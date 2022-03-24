package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"net/http"
)

func Get(db *storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortUrl")

		longURL, err := db.Get(shortURL)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Location", longURL)

		fmt.Println(w.Header().Get("Location"))

		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
