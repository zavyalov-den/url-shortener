package handler

import (
	"encoding/json"
	"fmt"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"net/http"
)

func GetUserUrls(db *storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID := ctx.Value("auth").(int)

		urls := db.GetUserUrls(userID)

		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
		}

		body, err := json.Marshal(urls)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, err = w.Write(body)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
