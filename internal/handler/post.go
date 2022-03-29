package handler

import (
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/service"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"io"
	"net/http"
)

func Post(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}

		if len(data) == 0 {
			http.Error(w, "request body must not be empty", 400)
			return
		}

		short := service.Shorten(data)

		ctx := r.Context()
		userID := ctx.Value("auth").(int)

		db.SaveURL(short, string(data))
		db.SaveUserURL(userID, storage.UserURL{
			ShortURL:    config.Config.BaseURL + "/" + short,
			OriginalURL: string(data),
		})

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(config.Config.BaseURL + "/" + short))
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}
	}
}
