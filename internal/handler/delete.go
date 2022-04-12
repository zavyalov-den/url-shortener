package handler

import (
	"encoding/json"
	"fmt"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"io"
	"net/http"
)

func Delete(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(data) == 0 {
			http.Error(w, "request body must not be empty", 400)
			return
		}

		var arr []string

		err = json.Unmarshal(data, &arr)
		if err != nil {
			http.Error(w, "request body data is not valid", 400)
			return
		}

		ctx := r.Context()

		userID, ok := ctx.Value(config.ContextKeyAuth).(int)
		if !ok {
			userID = 0
		}

		err = db.DeleteBatch(ctx, userID, arr)
		if err != nil {
			fmt.Printf("failed to delete batch: %s", err.Error())

		}

		w.WriteHeader(http.StatusAccepted)
	}
}
