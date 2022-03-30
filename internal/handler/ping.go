package handler

import (
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"net/http"
)

func Ping(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//ctx := r.Context()
		//
		//conn, err := pgx.Connect(ctx, config.Config.DatabaseDSN)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//defer conn.Close(ctx)
		//
		//err = conn.Ping(ctx)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//
		//w.WriteHeader(http.StatusOK)
	}
}
