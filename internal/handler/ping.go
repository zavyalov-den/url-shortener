package handler

import (
	"context"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"net/http"
)

func Ping(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		err := db.Ping(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
