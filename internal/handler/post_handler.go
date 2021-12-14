package handler

import (
	"github.com/zavyalov-den/url-shortener/internal/service"
	"io"
	"net/http"
)

func Post(urls map[string]string) http.HandlerFunc {
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

		urls[short] = string(data)

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("http://localhost:8080/" + short))
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}
	}
}
