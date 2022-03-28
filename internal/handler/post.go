package handler

import (
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/service"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"io"
	"net/http"
)

func Post(db *storage.DB) http.HandlerFunc {
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

		db.Save(short, string(data))
		db.SaveUserUrl(userID, storage.UserURL{
			ShortURL:    config.C.BaseURL + "/" + short,
			OriginalURL: string(data),
		})

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(config.C.BaseURL + "/" + short))
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}
	}
}
