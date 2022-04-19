package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/service"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"io"
	"net/http"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func ShortenJSON(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		defer r.Body.Close()

		req := &shortenRequest{}
		res := &shortenResponse{}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "request body is not a valid json", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(data, req)
		if err != nil {
			res.Error = "request body is not a valid json"
			errData, errJSON := json.Marshal(res)
			if errJSON != nil {
				panic(err.Error())
			}
			http.Error(w, string(errData), http.StatusBadRequest)
			return
		}

		short := service.Shorten([]byte(req.URL))
		res.Result = service.ShortToURL(short)

		resBody, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		userID, ok := ctx.Value(config.ContextKeyAuth).(int)
		if !ok {
			fmt.Println("not ok (shorten json): ", userID)
			userID = 0
		}

		fmt.Println("current user ID: ", userID)

		err = db.SaveURL(userID, storage.UserURL{
			ShortURL:    res.Result,
			OriginalURL: req.URL,
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
		_, err = w.Write(resBody)
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}
	}
}
