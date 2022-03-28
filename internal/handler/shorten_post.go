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
	URL string `json:"url"`
}

type response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func ShortenPost(db *storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		defer r.Body.Close()

		req := &request{}
		res := &response{}

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
		res.Result = config.C.BaseURL + "/" + short

		resBody, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		//fmt.Println(ctx)
		userID := ctx.Value("auth").(int)

		db.Save(short, req.URL)
		db.SaveUserUrl(userID, storage.UserURL{
			ShortURL:    short,
			OriginalURL: req.URL,
		})

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(resBody)
		if err != nil {
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
			return
		}
	}
}
