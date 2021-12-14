package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Get(urls map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortUrl")

		longURL, ok := urls[shortURL]
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
