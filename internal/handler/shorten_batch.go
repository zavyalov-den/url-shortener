package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"io"
	"net/http"
)

func ShortenBatch(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		defer r.Body.Close()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "request body is not a valid json", http.StatusBadRequest)
			return
		}

		var req []storage.BatchRequest

		err = json.Unmarshal(data, &req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		result, err := db.SaveBatch(ctx, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to save batch: %s", err), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(result)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write(resp)
		if err != nil {
			return
		}
	}
}
