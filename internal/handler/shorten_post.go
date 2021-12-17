package handler

import (
	"encoding/json"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/service"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"io"
	"net/http"
)

type request struct {
	Url string `json:"url"`
}

type response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func ShortenPost(db *storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// {"url": "<some_url>"}
		// {"result": "<shorten_url>"}
		// Content-Type
		// HTTP content negotiation

		defer r.Body.Close()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "request body is not a valid json", http.StatusBadRequest)
			return
		}
		req := &request{}
		res := &response{}

		err = json.Unmarshal(data, req)
		if err != nil {
			res.Error = "request body is not a valid json"
			errData, errJson := json.Marshal(res)
			if errJson != nil {
				panic(err.Error())
			}
			http.Error(w, string(errData), http.StatusBadRequest)
			return
		}

		short := service.Shorten([]byte(req.Url))
		res.Result = config.C.BaseURL + "/" + short

		resBody, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db.Save(short, req.Url)

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(resBody)
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}
	}
}
