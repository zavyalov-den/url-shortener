package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zavyalov-den/url-shortener/internal/service"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"net/http"
)

func Get(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short := chi.URLParam(r, "shortUrl")

		longURL, err := db.GetURL(service.ShortToURL(short))
		if err != nil {
			fmt.Printf("long: %s, short: %s\n", longURL, short)
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Location", longURL)

		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
