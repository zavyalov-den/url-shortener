package handler

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/zavyalov-den/url-shortener/internal/service"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"net/http"
)

func GetFullURL(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short := chi.URLParam(r, "shortUrl")

		longURL, err := db.GetURL(service.ShortToURL(short))
		if err != nil {
			if errors.Is(err, storage.ErrRowDeleted) {
				w.WriteHeader(http.StatusGone)
				return
			}
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
