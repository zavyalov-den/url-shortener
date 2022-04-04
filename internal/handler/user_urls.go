package handler

import (
	"encoding/json"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"net/http"
)

func GetUserUrls(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx := r.Context()

		userID, ok := ctx.Value(config.ContextKeyAuth).(int)
		if !ok {
			userID = 0
		}

		urls := db.GetUserURLs(userID)

		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		body, err := json.Marshal(urls)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
