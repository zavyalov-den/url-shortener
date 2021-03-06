package handler

import (
	"errors"
	"fmt"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/service"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"io"
	"net/http"
)

func Shorten(db storage.Storage) http.HandlerFunc {
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

		short := service.Shorten(data)

		ctx := r.Context()

		userID, ok := ctx.Value(config.ContextKeyAuth).(int)
		if !ok {
			fmt.Println("not ok: ", userID)
			userID = 0
		}

		err = db.SaveURL(userID, storage.UserURL{
			ShortURL:    service.ShortToURL(short),
			OriginalURL: string(data),
		})
		var conflictError error

		if errors.Is(err, storage.ErrConflict) {
			conflictError = err
		} else if err != nil {
			http.Error(w, "failed to save url to database: "+err.Error(), 400)
			return
		}

		if conflictError != nil {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		_, err = w.Write([]byte(service.ShortToURL(short)))
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}
	}
}
